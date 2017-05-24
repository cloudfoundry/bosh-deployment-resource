package bosh

import (
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/cloudfoundry/bosh-deployment-resource/concourse"

	boshcmd "github.com/cloudfoundry/bosh-cli/cmd"
	boshdir "github.com/cloudfoundry/bosh-cli/director"
	boshtpl "github.com/cloudfoundry/bosh-cli/director/template"
	boshlog "github.com/cloudfoundry/bosh-utils/logger"
	boshsys "github.com/cloudfoundry/bosh-utils/system"
)

type DeployParams struct {
	Vars      map[string]interface{}
	VarsFiles []string
	OpsFiles  []string
	NoRedact  bool
	Cleanup   bool
	VarsStore string
}

type Director interface {
	Deploy(manifestBytes []byte, deployParams DeployParams) error
	DownloadManifest() ([]byte, error)
	ExportReleases(targetDirectory string, releases []string) error
	UploadRelease(releaseURL string) error
	UploadStemcell(stemcellURL string) error
}

type BoshDirector struct {
	source        concourse.Source
	commandRunner Runner
	cliDirector   boshdir.Director
}

func NewBoshDirector(source concourse.Source, commandRunner Runner, cliDirector boshdir.Director) BoshDirector {
	return BoshDirector{
		source:        source,
		commandRunner: commandRunner,
		cliDirector:   cliDirector,
	}
}

func (d BoshDirector) Deploy(manifestBytes []byte, deployParams DeployParams) error {
	boshVarsFiles, err := parsedVarsFiles(deployParams.VarsFiles)
	if err != nil {
		return err
	}

	boshOpsFiles, err := parsedOpsFiles(deployParams.OpsFiles)
	if err != nil {
		return err
	}

	deployOpts := boshcmd.DeployOpts{
		Args:     boshcmd.DeployArgs{Manifest: boshcmd.FileBytesArg{Bytes: manifestBytes}},
		NoRedact: deployParams.NoRedact,
		VarFlags: boshcmd.VarFlags{
			VarKVs:    varKVsFromVars(deployParams.Vars),
			VarsFiles: boshVarsFiles,
		},
		OpsFlags: boshcmd.OpsFlags{
			OpsFiles: boshOpsFiles,
		},
	}

	if deployParams.VarsStore != "" {
		varsFSStore := boshcmd.VarsFSStore{}
		varsFSStore.FS = boshFileSystem()
		varsFSStore.UnmarshalFlag(deployParams.VarsStore)
		deployOpts.VarsFSStore = varsFSStore
	}

	err = d.commandRunner.Execute(&deployOpts)
	if err != nil {
		return fmt.Errorf("Could not deploy: %s\n", err)
	}

	if deployParams.Cleanup {
		d.commandRunner.Execute(&boshcmd.CleanUpOpts{})
	}

	return nil
}

func (d BoshDirector) DownloadManifest() ([]byte, error) {
	desiredDeployment, err := d.cliDirector.FindDeployment(d.source.Deployment)
	if err != nil {
		return nil, fmt.Errorf("Could not get deployment manifest: %s\n", err)
	}

	manifest, err := desiredDeployment.Manifest()
	return []byte(manifest), err
}

func (d BoshDirector) UploadRelease(URL string) error {
	err := d.commandRunner.Execute(&boshcmd.UploadReleaseOpts{
		Args: boshcmd.UploadReleaseArgs{URL: boshcmd.URLArg(URL)},
	})

	if err != nil {
		return fmt.Errorf("Could not upload release %s: %s\n", URL, err)
	}

	return nil
}

func (d BoshDirector) ExportReleases(targetDirectory string, releases []string) error {
	deploymentReleases, stemcell, err := d.releasesAndStemcell()
	if err != nil {
		return fmt.Errorf("could not export releases: %s", err)
	}

	releasesToDownload := []boshdir.Release{}

	for _, release := range releases {
		foundRelease := false
		for _, deploymentRelease := range deploymentReleases {
			if deploymentRelease.Name() == release {
				releasesToDownload = append(releasesToDownload, deploymentRelease)
				foundRelease = true
			}
		}

		if !foundRelease {
			return fmt.Errorf("could not find release %s in deployment", release)
		}
	}

	for _, deploymentRelease := range releasesToDownload {
		releaseSlug := boshdir.NewReleaseSlug(deploymentRelease.Name(), deploymentRelease.Version().AsString())
		osVersionSlug := boshdir.NewOSVersionSlug(stemcell.OSName(), stemcell.Version().AsString())

		directory := boshcmd.DirOrCWDArg{}
		directoryFixFunction := func(defaultedOps interface{}) (interface{}, error) {
			switch v := defaultedOps.(type) {
			case (*boshcmd.ExportReleaseOpts):
				v.Directory.Path = targetDirectory
			default:
				panic("todo")
			}
			return defaultedOps, nil
		}
		err = d.commandRunner.ExecuteWithDefaultOverride(&boshcmd.ExportReleaseOpts{
			Args:      boshcmd.ExportReleaseArgs{ReleaseSlug: releaseSlug, OSVersionSlug: osVersionSlug},
			Directory: directory,
		}, directoryFixFunction)
		if err != nil {
			return fmt.Errorf("could not export release %s: %s", deploymentRelease.Name(), err)
		}
	}

	return nil
}

func (d BoshDirector) UploadStemcell(URL string) error {
	err := d.commandRunner.Execute(&boshcmd.UploadStemcellOpts{
		Args: boshcmd.UploadStemcellArgs{URL: boshcmd.URLArg(URL)},
	})

	if err != nil {
		return fmt.Errorf("Could not upload stemcell %s: %s\n", URL, err)
	}

	return nil
}

func (d BoshDirector) deployment() (boshdir.Deployment, error) {
	deployment, err := d.cliDirector.FindDeployment(d.source.Deployment)
	if err != nil {
		return nil, fmt.Errorf("could not fetch deployment %s: %s", d.source.Deployment, err)
	}

	return deployment, nil
}

func (d BoshDirector) releasesAndStemcell() ([]boshdir.Release, boshdir.Stemcell, error) {
	deployment, err := d.deployment()
	if err != nil {
		return []boshdir.Release{}, nil, err
	}

	releases, err := deployment.Releases()
	if err != nil {
		return []boshdir.Release{}, nil, fmt.Errorf("could not fetch releases: %s", err)
	}

	deploymentStemcells, err := deployment.Stemcells()
	if err != nil {
		return []boshdir.Release{}, nil, fmt.Errorf("could not fetch stemcells: %s", err)
	}
	if len(deploymentStemcells) > 1 {
		return []boshdir.Release{}, nil, errors.New("exporting releases from a deployment with multiple stemcells is unsupported")
	}
	directorStemcells, err := d.cliDirector.Stemcells()
	if err != nil {
		return []boshdir.Release{}, nil, fmt.Errorf("could not fetch stemcells: %s", err)
	}

	var stemcell boshdir.Stemcell
	for _, directorStemcell := range directorStemcells {
		if directorStemcell.Name() == deploymentStemcells[0].Name() && directorStemcell.Version().IsEq(deploymentStemcells[0].Version()) {
			stemcell = directorStemcell
			break
		}
	}

	return releases, stemcell, nil
}

func varKVsFromVars(vars map[string]interface{}) []boshtpl.VarKV {
	varKVs := []boshtpl.VarKV{}
	for k, v := range vars {
		varKVs = append(varKVs, boshtpl.VarKV{Name: k, Value: v})
	}
	return varKVs
}

func parsedVarsFiles(varsFiles []string) ([]boshtpl.VarsFileArg, error) {
	varsFileArgs := []boshtpl.VarsFileArg{}
	for _, varsFile := range varsFiles {
		varsFileArg := boshtpl.VarsFileArg{FS: boshFileSystem()}
		if err := varsFileArg.UnmarshalFlag(varsFile); err != nil {
			return nil, err
		}
		varsFileArgs = append(varsFileArgs, varsFileArg)
	}
	return varsFileArgs, nil
}

func parsedOpsFiles(opsFiles []string) ([]boshcmd.OpsFileArg, error) {
	nullLogger := boshlog.NewWriterLogger(boshlog.LevelInfo, ioutil.Discard, ioutil.Discard)
	boshFS := boshsys.NewOsFileSystemWithStrictTempRoot(nullLogger)

	opsFileArgs := []boshcmd.OpsFileArg{}
	for _, opsFile := range opsFiles {
		opsFileArg := boshcmd.OpsFileArg{FS: boshFS}
		if err := opsFileArg.UnmarshalFlag(opsFile); err != nil {
			return nil, err
		}
		opsFileArgs = append(opsFileArgs, opsFileArg)
	}

	return opsFileArgs, nil
}

func boshFileSystem() boshsys.FileSystem {
	nullLogger := boshlog.NewWriterLogger(boshlog.LevelInfo, ioutil.Discard, ioutil.Discard)
	return boshsys.NewOsFileSystemWithStrictTempRoot(nullLogger)
}
