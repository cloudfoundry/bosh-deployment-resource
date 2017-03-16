package main

import (
	"encoding/json"
	"fmt"
	"os"

	"io/ioutil"

	"github.com/cloudfoundry/bosh-deployment-resource/bosh"
	"github.com/cloudfoundry/bosh-deployment-resource/concourse"
	"github.com/cloudfoundry/bosh-deployment-resource/in"
)

func main() {
	//remove when https://github.com/cloudfoundry/bosh-cli/pull/135
	realStdout := os.Stdout
	os.Stdout = os.Stderr

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

	cliCoordinator := bosh.NewCLICoordinator(inRequest.Source, os.Stderr)
	commandRunner := bosh.NewCommandRunner(cliCoordinator)
	cliDirector, err := cliCoordinator.Director()
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}
	director := bosh.NewBoshDirector(inRequest.Source, commandRunner, cliDirector)

	inCommand := in.NewInCommand(director)
	inResponse, err := inCommand.Run(inRequest, targetDir)
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}

	concourseInputFormatted, err := json.MarshalIndent(inResponse, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not generate version: %s\n", err)
		os.Exit(1)
	}

	fmt.Fprintf(realStdout, "%s", concourseInputFormatted)
}
