package bosh

import (
	"bytes"
	"io"
	"io/ioutil"

	"github.com/cloudfoundry/bosh-deployment-resource/concourse"

	boshcmd "github.com/cloudfoundry/bosh-cli/cmd"
	cmdconf "github.com/cloudfoundry/bosh-cli/cmd/config"
	boshdir "github.com/cloudfoundry/bosh-cli/director"
	boshui "github.com/cloudfoundry/bosh-cli/ui"
	boshlog "github.com/cloudfoundry/bosh-utils/logger"
	goflags "github.com/jessevdk/go-flags"
)

type CLICoordinator struct {
	source concourse.Source
	out    io.Writer
}

func NewCLICoordinator(source concourse.Source, out io.Writer) CLICoordinator {
	return CLICoordinator{
		source: source,
		out:    out,
	}
}

func (c CLICoordinator) GlobalOpts() boshcmd.BoshOpts {
	globalOpts := &boshcmd.BoshOpts{
		NonInteractiveOpt: true,
		CACertOpt:         boshcmd.CACertArg{Content: c.source.CACert},
		ClientOpt:         c.source.Client,
		ClientSecretOpt:   c.source.ClientSecret,
		EnvironmentOpt:    c.source.Target,
		DeploymentOpt:     c.source.Deployment,
	}

	setDefaults(globalOpts)

	return *globalOpts
}

func (c CLICoordinator) StreamingBasicDeps() boshcmd.BasicDeps {
	logger := nullLogger()

	parentUI := boshui.NewPaddingUI(boshui.NewWriterUI(c.out, c.out, logger))

	ui := boshui.NewWrappingConfUI(parentUI, logger)
	return boshcmd.NewBasicDeps(ui, logger)
}

func (c CLICoordinator) CapturedBasicDeps() boshcmd.BasicDeps {
	byteWriter := bytes.NewBufferString("")
	logger := nullLogger()

	parentUI := boshui.NewNonTTYUI(boshui.NewWriterUI(byteWriter, c.out, logger))

	ui := boshui.NewWrappingConfUI(parentUI, logger)
	return boshcmd.NewBasicDeps(ui, logger)
}

func (c CLICoordinator) Director() (boshdir.Director, error) {
	globalOpts := c.GlobalOpts()
	deps := c.CapturedBasicDeps()
	config, _ := cmdconf.NewFSConfigFromPath(globalOpts.ConfigPathOpt, deps.FS)
	session := boshcmd.NewSessionFromOpts(globalOpts, config, deps.UI, true, true, deps.FS, deps.Logger)

	return session.Director()
}

func nullLogger() boshlog.Logger {
	return boshlog.NewWriterLogger(boshlog.LevelInfo, ioutil.Discard, ioutil.Discard)
}

func setDefaults(obj interface{}) {
	parser := goflags.NewParser(obj, goflags.None)

	// Intentionally ignoring errors. We are not parsing user-passed options,
	// we are just doing this to populate defaults.
	parser.ParseArgs([]string{})
}
