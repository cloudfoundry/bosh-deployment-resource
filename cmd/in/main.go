package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/cloudfoundry/bosh-deployment-resource/bosh"
	"github.com/cloudfoundry/bosh-deployment-resource/concourse"
	"github.com/cloudfoundry/bosh-deployment-resource/in"
	"github.com/cloudfoundry/socks5-proxy"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr,
			"not enough args - usage: %s <target directory>\n",
			os.Args[0],
		)
		os.Exit(1)
	}

	stdin, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Cannot read configuration: %s\n", err)
		os.Exit(1)
	}

	inRequest, err := concourse.NewInRequest(stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Invalid parameters: %s\n", err)
		os.Exit(1)
	}

	targetDir := os.Args[1]

	if inRequest.Source.Target == concourse.MissingTarget {
		inResponse := in.InResponse{Version: inRequest.Version}
		printResponse(inResponse)
		os.Exit(0)
	}

	hostKeyGetter := proxy.NewHostKey()
	socks5Proxy := proxy.NewSocks5Proxy(hostKeyGetter, log.New(ioutil.Discard, "", log.LstdFlags), 1*time.Minute)
	cliCoordinator := bosh.NewCLICoordinator(inRequest.Source, os.Stderr, socks5Proxy)
	commandRunner := bosh.NewCommandRunner(cliCoordinator)
	cliDirector, err := cliCoordinator.Director()
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}
	director := bosh.NewBoshDirector(
		inRequest.Source,
		commandRunner,
		cliDirector,
		os.Stderr,
	)

	inCommand := in.NewInCommand(director)
	inResponse, err := inCommand.Run(inRequest, targetDir)
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}

	printResponse(inResponse)
}

func printResponse(inResponse in.InResponse) {
	concourseInputFormatted, err := json.MarshalIndent(inResponse, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not generate version: %s\n", err)
		os.Exit(1)
	}

	fmt.Fprintf(os.Stdout, "%s", concourseInputFormatted)
}
