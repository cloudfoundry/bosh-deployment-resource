package main

import (
	"encoding/json"
	"fmt"
	"os"

	"io/ioutil"

	"github.com/cloudfoundry/bosh-deployment-resource/bosh"
	"github.com/cloudfoundry/bosh-deployment-resource/concourse"
	"github.com/cloudfoundry/bosh-deployment-resource/out"
	"github.com/cloudfoundry/bosh-deployment-resource/storage"
)

func main() {
	//remove when https://github.com/cloudfoundry/bosh-cli/pull/135
	realStdout := os.Stdout
	os.Stdout = os.Stderr

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

	cliCoordinator := bosh.NewCLICoordinator(outRequest.Source, os.Stderr)
	commandRunner := bosh.NewCommandRunner(cliCoordinator)
	cliDirector, err := cliCoordinator.Director()
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}
	director := bosh.NewBoshDirector(outRequest.Source, commandRunner, cliDirector)

	storageClient, err := storage.NewStorageClient(outRequest.Source)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Invalid vars store: %s\n", err)
		os.Exit(1)
	}

	outCommand := out.NewOutCommand(director, storageClient, sourcesDir)
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

	fmt.Fprintf(realStdout, "%s", concourseOutputFormatted)
}
