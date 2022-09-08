package bosh

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/cloudfoundry/bosh-utils/httpclient"
	"io"
	"io/ioutil"
	"log"
	"os"

	"github.com/cloudfoundry/bosh-deployment-resource/concourse"

	boshcmd "github.com/cloudfoundry/bosh-cli/v7/cmd"
	cmdconf "github.com/cloudfoundry/bosh-cli/v7/cmd/config"
	boshcmdopts "github.com/cloudfoundry/bosh-cli/v7/cmd/opts"
	boshdir "github.com/cloudfoundry/bosh-cli/v7/director"
	boshui "github.com/cloudfoundry/bosh-cli/v7/ui"
	boshlog "github.com/cloudfoundry/bosh-utils/logger"
	goflags "github.com/jessevdk/go-flags"
)

//go:generate counterfeiter . Proxy
type Proxy interface {
	Start(string, string, string) error
	Addr() (string, error)
}

type CLICoordinator struct {
	source concourse.Source
	out    io.Writer
	proxy  Proxy
}

func NewCLICoordinator(source concourse.Source, out io.Writer, proxy Proxy) CLICoordinator {
	return CLICoordinator{
		source: source,
		out:    out,
		proxy:  proxy,
	}
}

func (c CLICoordinator) GlobalOpts() boshcmdopts.BoshOpts {
	globalOpts := &boshcmdopts.BoshOpts{
		NonInteractiveOpt: true,
		CACertOpt:         boshcmdopts.CACertArg{Content: c.source.CACert},
		ClientOpt:         c.source.Client,
		ClientSecretOpt:   c.source.ClientSecret,
		EnvironmentOpt:    c.source.Target,
		DeploymentOpt:     c.source.Deployment,
	}

	if c.source.JumpboxSSHKey != "" && c.source.JumpboxURL != "" {
		sshKeyFile, err := ioutil.TempFile("", "ssh-key-file")
		if err != nil {
			log.Fatal(err)
		}
		sshKeyFile.Write([]byte(c.source.JumpboxSSHKey))
		proxyAddr := fmt.Sprintf("ssh+socks5://%s@%s?private-key=%s", c.source.JumpboxUsername, c.source.JumpboxURL, sshKeyFile.Name())
		os.Setenv("BOSH_ALL_PROXY", proxyAddr)
		globalOpts.SSH.GatewayFlags.SOCKS5Proxy = proxyAddr
		globalOpts.SCP.GatewayFlags.SOCKS5Proxy = proxyAddr
		globalOpts.Logs.GatewayFlags.SOCKS5Proxy = proxyAddr
	}

	setDefaults(globalOpts)

	return *globalOpts
}

func (c CLICoordinator) BasicDeps(writer io.Writer) boshcmd.BasicDeps {
	logger := nullLogger()

	if writer == nil {
		writer = c.out
	}
	parentUI := boshui.NewPaddingUI(boshui.NewWriterUI(writer, writer, logger))

	ui := boshui.NewWrappingConfUI(parentUI, logger)
	return boshcmd.NewBasicDeps(ui, logger)
}

func (c CLICoordinator) Director() (boshdir.Director, error) {
	globalOpts := c.GlobalOpts()
	deps := c.BasicDeps(bytes.NewBufferString(""))
	config, _ := cmdconf.NewFSConfigFromPath(globalOpts.ConfigPathOpt, deps.FS)
	session := boshcmd.NewSessionFromOpts(globalOpts, config, deps.UI, true, true, deps.FS, deps.Logger)
	httpclient.ResetDialerContext()
	return session.Director()
}

func (c CLICoordinator) StartProxy() (string, error) {
	if c.source.JumpboxSSHKey == "" && c.source.JumpboxURL == "" {
		return "", nil
	}

	if c.source.JumpboxSSHKey != "" && c.source.JumpboxURL != "" {
		addr, err := c.proxy.Addr()
		if err == nil {
			return addr, nil
		}

		err = c.proxy.Start(c.source.JumpboxUsername, c.source.JumpboxSSHKey, c.source.JumpboxURL)
		if err != nil {
			panic(err)
		}

		addr, err = c.proxy.Addr()
		if err != nil {
			panic(err)
		}
		return addr, nil
	}

	return "", errors.New("Jumpbox URL and Jumpbox SSH Key are both required to use a jumpbox")
}

func nullLogger() boshlog.Logger {
	return boshlog.NewWriterLogger(boshlog.LevelInfo, ioutil.Discard)
}

func setDefaults(obj interface{}) {
	parser := goflags.NewParser(obj, goflags.None)

	// Intentionally ignoring errors. We are not parsing user-passed options,
	// we are just doing this to populate defaults.
	parser.ParseArgs([]string{})
}
