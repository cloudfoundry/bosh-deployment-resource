package bosh

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/cloudfoundry/bosh-deployment-resource/bosh/proxy"
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
	proxy  proxy.Proxy
}

func NewCLICoordinator(source concourse.Source, out io.Writer, proxy proxy.Proxy) CLICoordinator {
	return CLICoordinator{
		source: source,
		out:    out,
		proxy:  proxy,
	}
}

func (c CLICoordinator) GlobalOpts(proxyAddr string) boshcmd.BoshOpts {
	globalOpts := &boshcmd.BoshOpts{
		NonInteractiveOpt: true,
		CACertOpt:         boshcmd.CACertArg{Content: c.source.CACert},
		ClientOpt:         c.source.Client,
		ClientSecretOpt:   c.source.ClientSecret,
		EnvironmentOpt:    c.source.Target,
		DeploymentOpt:     c.source.Deployment,
	}

	if proxyAddr != "" {
		proxyAddr = fmt.Sprintf("socks5://%s", proxyAddr)
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
	addr, err := c.StartProxy()
	if err != nil {
		return nil, fmt.Errorf("start proxy: %s", err) // untested
	}
	globalOpts := c.GlobalOpts(addr)
	deps := c.BasicDeps(bytes.NewBufferString(""))
	config, _ := cmdconf.NewFSConfigFromPath(globalOpts.ConfigPathOpt, deps.FS)
	session := boshcmd.NewSessionFromOpts(globalOpts, config, deps.UI, true, true, deps.FS, deps.Logger)

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

		err = c.proxy.Start(c.source.JumpboxSSHKey, c.source.JumpboxURL)
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
	goflags.FactoryFunc = func(val interface{}) {}

	parser := goflags.NewParser(obj, goflags.None)

	// Intentionally ignoring errors. We are not parsing user-passed options,
	// we are just doing this to populate defaults.
	parser.ParseArgs([]string{})
}
