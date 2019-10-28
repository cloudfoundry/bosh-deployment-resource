package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"io/ioutil"

	"github.com/cloudfoundry/bosh-deployment-resource/bosh"
	"github.com/cloudfoundry/bosh-deployment-resource/concourse"
	"github.com/cloudfoundry/bosh-deployment-resource/out"
	"github.com/cloudfoundry/bosh-deployment-resource/storage"
	proxy "github.com/cloudfoundry/socks5-proxy"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr,
			"not enough args - usage: %s <sources directory>\n",
			os.Args[0],
		)
		os.Exit(1)
	}

	sourcesDir := os.Args[1]

	stdin, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Cannot read configuration: %s\n", err)
		os.Exit(1)
	}

	outRequest, err := concourse.NewOutRequest(stdin, sourcesDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Invalid parameters: %s\n", err)
		os.Exit(1)
	}

	hostKeyGetter := proxy.NewHostKey()
	socks5Proxy := proxy.NewSocks5Proxy(hostKeyGetter, log.New(ioutil.Discard, "", log.LstdFlags), 1*time.Minute)
	cliCoordinator := bosh.NewCLICoordinator(outRequest.Source, os.Stderr, socks5Proxy)
	commandRunner := bosh.NewCommandRunner(cliCoordinator)
	cliDirector, err := cliCoordinator.Director()
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}
	director := bosh.NewBoshDirector(
		outRequest.Source,
		commandRunner,
		cliDirector,
		os.Stderr,
	)

	storageClient, err := storage.NewStorageClient(outRequest.Source)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Invalid vars store: %s\n", err)
		os.Exit(1)
	}

	outCommand := out.NewOutCommand(director, bosh.BoshIOClient{}, storageClient, sourcesDir)
	outResponse, err := outCommand.Run(outRequest)
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}

	concourseOutputFormatted, err := json.MarshalIndent(outResponse, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not generate version: %s\n", err)
		os.Exit(1)
	}

	fmt.Fprintf(os.Stdout, "%s", concourseOutputFormatted)
}
