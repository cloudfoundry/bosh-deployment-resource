package cmd

import (
	"errors"
	"fmt"
	davclient "github.com/cloudfoundry/bosh-davcli/client"
	"time"
)

type SignCmd struct {
	client davclient.Client
}

func newSignCmd(client davclient.Client) (cmd SignCmd) {
	cmd.client = client
	return
}

func (cmd SignCmd) Run(args []string) (err error) {
	if len(args) != 3 {
		err = errors.New("incorrect usage, sign requires: <object_id> <action> <duration>")
		return
	}

	objectID, action := args[0], args[1]

	expiration, err := time.ParseDuration(args[2])
	if err != nil {
		err = fmt.Errorf("expiration should be a duration value eg: 45s or 1h43m. Got: %s", args[2])
		return
	}

	signedURL, err := cmd.client.Sign(objectID, action, expiration)
	if err != nil {
		return err
	}

	fmt.Print(signedURL)
	return
}
