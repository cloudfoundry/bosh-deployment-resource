package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	proxy "github.com/cloudfoundry/socks5-proxy"

	"github.com/cloudfoundry/bosh-deployment-resource/bosh"
	"github.com/cloudfoundry/bosh-deployment-resource/check"
	"github.com/cloudfoundry/bosh-deployment-resource/concourse"
)

func main() {
	stdin, err := io.ReadAll(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Cannot read configuration: %s\n", err) //nolint:errcheck
		os.Exit(1)
	}

	checkRequest, err := concourse.NewCheckRequest(stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Invalid parameters: %s\n", err) //nolint:errcheck
		os.Exit(1)
	}

	var checkResponse []concourse.Version

	if checkRequest.Source.SkipCheck {
		checkResponse = []concourse.Version{}
	} else {

		hostKeyGetter := proxy.NewHostKey()
		socks5Proxy := proxy.NewSocks5Proxy(hostKeyGetter, log.New(io.Discard, "", log.LstdFlags), 1*time.Minute)
		cliCoordinator := bosh.NewCLICoordinator(checkRequest.Source, os.Stderr, socks5Proxy)
		commandRunner := bosh.NewCommandRunner(cliCoordinator)
		cliDirector, err := cliCoordinator.Director()
		if err != nil {
			fmt.Fprint(os.Stderr, err) //nolint:errcheck
			os.Exit(1)
		}

		director := bosh.NewBoshDirector(
			checkRequest.Source,
			commandRunner,
			cliDirector,
			os.Stderr,
		)

		checkCommand := check.NewCheckCommand(director)
		checkResponse, err = checkCommand.Run(checkRequest)
		if err != nil {
			fmt.Fprint(os.Stderr, err) //nolint:errcheck
			os.Exit(1)
		}
	}

	concourseOutputFormatted, err := json.MarshalIndent(checkResponse, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not generate version: %s\n", err) //nolint:errcheck
		os.Exit(1)
	}

	fmt.Fprintf(os.Stdout, "%s", concourseOutputFormatted) //nolint:errcheck
}
