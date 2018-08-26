package bosh

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"time"

	"github.com/cloudfoundry/bosh-deployment-resource/concourse"

	boshcmd "github.com/cloudfoundry/bosh-cli/cmd"
	boshdir "github.com/cloudfoundry/bosh-cli/director"
	boshtpl "github.com/cloudfoundry/bosh-cli/director/template"
	boshlog "github.com/cloudfoundry/bosh-utils/logger"
	boshsys "github.com/cloudfoundry/bosh-utils/system"
)

type DeployParams struct {
	Vars      map[string]interface{}
	VarFiles  map[string]string
	VarsFiles []string
	OpsFiles  []string
	NoRedact  bool
	DryRun    bool
	Recreate  bool
	SkipDrain []string
	Cleanup   bool
	VarsStore string
}

type InterpolateParams struct {
	Vars      map[string]interface{}
	VarsFiles []string
	OpsFiles  []string
	VarsStore string
}

type ReleaseSpec struct {
	Name string
	Jobs []string
}

//go:generate counterfeiter . Director
type Director interface {
	Delete(force bool) error
	Deploy(manifestBytes []byte, deployParams DeployParams) error
	Interpolate(manifestBytes []byte, interpolateParams InterpolateParams) ([]byte, error)
	DownloadManifest() ([]byte, error)
	ExportReleases(targetDirectory string, releases []ReleaseSpec) error
	UploadRelease(releaseURL string) error
	UploadStemcell(stemcellURL string) error
	WaitForDeployLock() error
}

type BoshDirector struct {
	source        concourse.Source
	commandRunner Runner
	cliDirector   boshdir.Director
	writer        io.Writer
}

func NewBoshDirector(source concourse.Source, commandRunner Runner, cliDirector boshdir.Director, writer io.Writer) BoshDirector {
	return BoshDirector{
		source:        source,
		commandRunner: commandRunner,
		cliDirector:   cliDirector,
		writer:        writer,
	}
}

func (d BoshDirector) Delete(force bool) error {
	return d.commandRunner.Execute(&boshcmd.DeleteDeploymentOpts{Force: force})
}

func (d BoshDirector) Deploy(manifestBytes []byte, deployParams DeployParams) error {
	boshVarsFiles, err := parsedVarsFiles(deployParams.VarsFiles)
	if err != nil {
		return err
	}

	boshVarFiles, err := parsedVarFiles(deployParams.VarFiles)
	if err != nil {
		return err
	}

	boshOpsFiles, err := parsedOpsFiles(deployParams.OpsFiles)
	if err != nil {
		return err
	}

	skipDrains, err := parsedSkipDrains(deployParams.SkipDrain)
	if err != nil {
		return err
	}

	deployOpts := boshcmd.DeployOpts{
		Args:      boshcmd.DeployArgs{Manifest: boshcmd.FileBytesArg{Bytes: manifestBytes}},
		NoRedact:  deployParams.NoRedact,
		DryRun:    deployParams.DryRun,
		Recreate:  deployParams.Recreate,
		SkipDrain: skipDrains,
		VarFlags: boshcmd.VarFlags{
			VarKVs:    varKVsFromVars(deployParams.Vars),
			VarsFiles: boshVarsFiles,
			VarFiles:  boshVarFiles,
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

func (d BoshDirector) Interpolate(manifestBytes []byte, interpolateParams InterpolateParams) ([]byte, error) {
	boshVarsFiles, err := parsedVarsFiles(interpolateParams.VarsFiles)
	if err != nil {
		return nil, err
	}

	boshOpsFiles, err := parsedOpsFiles(interpolateParams.OpsFiles)
	if err != nil {
		return nil, err
	}

	interpolateOpts := boshcmd.InterpolateOpts{
		Args: boshcmd.InterpolateArgs{Manifest: boshcmd.FileBytesArg{Bytes: manifestBytes}},
		VarFlags: boshcmd.VarFlags{
			VarKVs:    varKVsFromVars(interpolateParams.Vars),
			VarsFiles: boshVarsFiles,
		},
		OpsFlags: boshcmd.OpsFlags{
			OpsFiles: boshOpsFiles,
		},
	}

	writer := new(bytes.Buffer)
	err = d.commandRunner.ExecuteWithWriter(&interpolateOpts, writer)
	if err != nil {
		return nil, fmt.Errorf("Could not interpolate: %s\n", err)
	}

	return writer.Bytes(), nil
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

func (d BoshDirector) WaitForDeployLock() error {
	fmt.Fprint(d.writer, "Waiting for deployment lock")

	locked, err := d.deploymentIsLocked()
	if err != nil {
		return err
	} else if locked {
		for locked {
			time.Sleep(3 * time.Second)
			locked, err = d.deploymentIsLocked()
			if err != nil {
				return err
			}
		}
	}
	fmt.Fprintln(d.writer, " Done")
	return nil
}

func (d BoshDirector) deploymentIsLocked() (bool, error) {
	locks, err := d.cliDirector.Locks()
	if err != nil {
		return true, fmt.Errorf("Could not check if deployment was locked: %s\n", err)
	}

	for _, lock := range locks {
		resources := lock.Resource
		for _, resource := range resources {
			if resource == d.source.Deployment {
				fmt.Fprint(d.writer, ".")
				return true, nil
			}
		}
	}
	return false, nil
}

func (d BoshDirector) ExportReleases(targetDirectory string, releases []ReleaseSpec) error {
	deploymentReleases, stemcell, err := d.releasesAndStemcell()
	if err != nil {
		return fmt.Errorf("could not export releases: %s", err)
	}

	releasesToDownload := []boshdir.Release{}

	for _, release := range releases {
		foundRelease := false
		for _, deploymentRelease := range deploymentReleases {
			if deploymentRelease.Name() == release.Name {
				releasesToDownload = append(releasesToDownload, deploymentRelease)
				foundRelease = true
			}
		}

		if !foundRelease {
			return fmt.Errorf("could not find release %s in deployment", release.Name)
		}
	}

	for i, deploymentRelease := range releasesToDownload {
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
			Jobs:      releases[i].Jobs,
			Directory: directory,
		}, directoryFixFunction, nil)
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

func parsedVarFiles(varFiles map[string]string) ([]boshtpl.VarFileArg, error) {
	varFileArgs := []boshtpl.VarFileArg{}
	for varKey, varFile := range varFiles {
		varFileArg := boshtpl.VarFileArg{FS: boshFileSystem()}
		if err := varFileArg.UnmarshalFlag(fmt.Sprintf("%s=%s", varKey, varFile)); err != nil {
			return nil, err
		}
		varFileArgs = append(varFileArgs, varFileArg)
	}
	return varFileArgs, nil
}

func parsedOpsFiles(opsFiles []string) ([]boshcmd.OpsFileArg, error) {
	nullLogger := boshlog.NewWriterLogger(boshlog.LevelInfo, ioutil.Discard)
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

func parsedSkipDrains(drains []string) ([]boshdir.SkipDrain, error) {
	parsedDrains := make([]boshdir.SkipDrain, len(drains))
	for idx, drain := range drains {
		parsedDrain := boshdir.SkipDrain{}
		if err := parsedDrain.UnmarshalFlag(drain); err != nil {
			return parsedDrains, err
		}
		parsedDrains[idx] = parsedDrain
	}
	return parsedDrains, nil
}

func boshFileSystem() boshsys.FileSystem {
	nullLogger := boshlog.NewWriterLogger(boshlog.LevelInfo, ioutil.Discard)
	return boshsys.NewOsFileSystemWithStrictTempRoot(nullLogger)
}
