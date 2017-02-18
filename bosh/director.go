package bosh

import (
	"io"
	"fmt"

	"github.com/cloudfoundry/bosh-deployment-resource/concourse"

	boshcmd "github.com/cloudfoundry/bosh-cli/cmd"
	boshtpl "github.com/cloudfoundry/bosh-cli/director/template"
	boshsys "github.com/cloudfoundry/bosh-utils/system"
	boshlog "github.com/cloudfoundry/bosh-utils/logger"
	"io/ioutil"
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
	UploadRelease(releaseURL string) error
	UploadStemcell(stemcellURL string) error
}

type BoshDirector struct {
	source             concourse.Source
	commandRunner      Runner
	out                io.Writer
}

func NewBoshDirector(source concourse.Source, commandRunner Runner, out io.Writer) BoshDirector {
	return BoshDirector{
		source:             source,
		commandRunner:      commandRunner,
		out:                out,
	}
}

func (d BoshDirector) Deploy(manifestBytes []byte, deployParams DeployParams) error {
	if deployParams.Cleanup {
		d.commandRunner.Execute(&boshcmd.CleanUpOpts{})
	}

	boshVarsFiles, err := parsedVarsFiles(deployParams.VarsFiles)
	if err != nil {
		return err
	}

	boshOpsFiles, err := parsedOpsFiles(deployParams.OpsFiles)
	if err != nil {
		return err
	}

	err = d.commandRunner.Execute(&boshcmd.DeployOpts{
		Args:     boshcmd.DeployArgs{Manifest: boshcmd.FileBytesArg{Bytes: manifestBytes}},
		NoRedact: deployParams.NoRedact,
		VarFlags: boshcmd.VarFlags{
			VarKVs: varKVsFromVars(deployParams.Vars),
			VarsFiles: boshVarsFiles,
		},
		OpsFlags: boshcmd.OpsFlags{
			OpsFiles: boshOpsFiles,
		},
	})
	if err != nil {
		return fmt.Errorf("Could not deploy: %s\n", err)
	}

	return nil
}

func (d BoshDirector) DownloadManifest() ([]byte, error) {
	bytes, err := d.commandRunner.GetResult(&boshcmd.ManifestOpts{})

	if err != nil {
		return nil, fmt.Errorf("Could not get deployment manifest: %s\n", err)
	}

	return bytes, nil
}

func (d BoshDirector) UploadRelease(URL string) error {
	err := d.commandRunner.Execute(&boshcmd.UploadReleaseOpts{
		Args:     boshcmd.UploadReleaseArgs{URL: boshcmd.URLArg(URL)},
	})

	if err != nil {
		return fmt.Errorf("Could not upload release %s: %s\n", URL, err)
	}

	return nil
}

func (d BoshDirector) UploadStemcell(URL string) error {
	err := d.commandRunner.Execute(&boshcmd.UploadStemcellOpts{
		Args:     boshcmd.UploadStemcellArgs{URL: boshcmd.URLArg(URL)},
	})

	if err != nil {
		return fmt.Errorf("Could not upload stemcell %s: %s\n", URL, err)
	}

	return nil
}

func varKVsFromVars(vars map[string]interface{}) []boshtpl.VarKV {
	varKVs := []boshtpl.VarKV{}
	for k, v := range vars {
		varKVs = append(varKVs, boshtpl.VarKV{Name: k, Value: v})
	}
	return varKVs
}

func parsedVarsFiles(varsFiles []string) ([]boshtpl.VarsFileArg, error) {
	nullLogger := boshlog.NewWriterLogger(boshlog.LevelInfo, ioutil.Discard, ioutil.Discard)
	boshFS := boshsys.NewOsFileSystemWithStrictTempRoot(nullLogger)

	varsFileArgs := []boshtpl.VarsFileArg{}
	for _, varsFile := range varsFiles {
		varsFileArg := boshtpl.VarsFileArg{FS: boshFS}
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
