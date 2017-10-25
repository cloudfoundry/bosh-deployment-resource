package cmd

import (
	"fmt"

	"crypto/x509"
	davclient "github.com/cloudfoundry/bosh-davcli/client"
	davconf "github.com/cloudfoundry/bosh-davcli/config"
	boshhttpclient "github.com/cloudfoundry/bosh-utils/httpclient"
	boshlog "github.com/cloudfoundry/bosh-utils/logger"
	boshcrypto "github.com/cloudfoundry/bosh-utils/crypto"
)

type Factory interface {
	Create(name string) (cmd Cmd, err error)
	SetConfig(config davconf.Config) (err error)
}

func NewFactory(logger boshlog.Logger) Factory {
	return &factory{
		cmds:   make(map[string]Cmd),
		logger: logger,
	}
}

type factory struct {
	config davconf.Config
	cmds   map[string]Cmd
	logger boshlog.Logger
}

func (f *factory) Create(name string) (cmd Cmd, err error) {
	cmd, found := f.cmds[name]
	if !found {
		err = fmt.Errorf("Could not find command with name %s", name)
	}
	return
}

func (f *factory) SetConfig(config davconf.Config) (err error) {
	var httpClient boshhttpclient.Client
	var certPool *x509.CertPool

	if len(config.CACert) != 0 {
		certPool, err = boshcrypto.CertPoolFromPEM([]byte(config.CACert))
	}

	httpClient = boshhttpclient.CreateDefaultClient(certPool)

	client := davclient.NewClient(config, httpClient, f.logger)

	f.cmds = map[string]Cmd{
		"put":    newPutCmd(client),
		"get":    newGetCmd(client),
		"exists": newExistsCmd(client),
		"delete": newDeleteCmd(client),
	}

	return
}
