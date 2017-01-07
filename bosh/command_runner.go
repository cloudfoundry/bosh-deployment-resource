package bosh

import (
	boshcmd "github.com/cloudfoundry/bosh-cli/cmd"
	boshui "github.com/cloudfoundry/bosh-cli/ui"
	boshlog "github.com/cloudfoundry/bosh-utils/logger"

	goflags "github.com/jessevdk/go-flags"
	"github.com/cloudfoundry/bosh-deployment-resource/concourse"
	"io"
	"io/ioutil"
	"bytes"
)

type Runner interface {
	Execute(commandOpts interface{}) error
	GetResult(commandOpts interface{}) ([]byte, error)
}

type CommandRunner struct {
	source             concourse.Source
	out io.Writer
}

func NewCommandRunner(source concourse.Source, out io.Writer) CommandRunner {
	return CommandRunner{
		source: source,
		out: out,
	}
}

func (c CommandRunner) Execute(commandOpts interface{}) error {
	return c.internalExecute(commandOpts, c.streamingBasicDeps())
}

func (c CommandRunner) GetResult(commandOpts interface{}) ([]byte, error) {
	b := bytes.NewBufferString("")
	err := c.internalExecute(commandOpts, c.capturedBasicDeps(b))

	if err != nil {
		return nil, err
	}

	return b.Bytes(), nil
}

func (c CommandRunner) internalExecute(commandOpts interface{}, deps boshcmd.BasicDeps) error {
	globalOpts := &boshcmd.BoshOpts{
		NonInteractiveOpt: true,
		CACertOpt: boshcmd.CACertArg{Content: c.source.CACert},
		ClientOpt: c.source.Client,
		ClientSecretOpt: c.source.ClientSecret,
		EnvironmentOpt: c.source.Target,
		DeploymentOpt: c.source.Deployment,
	}

	setDefaults(globalOpts)
	setDefaults(commandOpts)

	cmd := boshcmd.NewCmd(*globalOpts, commandOpts, deps)

	return cmd.Execute()
}

func setDefaults(obj interface{}) {
	parser := goflags.NewParser(obj, goflags.None)

	// Intentionally ignoring errors. We are not parsing user-passed options,
	// we are just doing this to populate defaults.
	parser.ParseArgs([]string{})
}

func (c CommandRunner) streamingBasicDeps() boshcmd.BasicDeps {
	logger := nullLogger()

	parentUI := boshui.NewPaddingUI(boshui.NewWriterUI(c.out, c.out, logger))

	ui := boshui.NewWrappingConfUI(parentUI, logger)
	return boshcmd.NewBasicDeps(ui, logger)
}

func (c CommandRunner) capturedBasicDeps(ourByteWriter io.Writer) boshcmd.BasicDeps {
	logger := nullLogger()

	parentUI := boshui.NewNonTTYUI(boshui.NewWriterUI(ourByteWriter, c.out, logger))

	ui := boshui.NewWrappingConfUI(parentUI, logger)
	return boshcmd.NewBasicDeps(ui, logger)
}

func nullLogger() boshlog.Logger {
	return boshlog.NewWriterLogger(boshlog.LevelInfo, ioutil.Discard, ioutil.Discard)
}