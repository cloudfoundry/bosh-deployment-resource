package main

import (
	"fmt"
	"os"
	"encoding/json"

	"github.com/cloudfoundry/bosh-deployment-resource/concourse"
	"github.com/cloudfoundry/bosh-deployment-resource/check"
	"github.com/cloudfoundry/bosh-deployment-resource/bosh"
)

func main() {
	var checkRequest concourse.CheckRequest
	if err := json.NewDecoder(os.Stdin).Decode(&checkRequest); err != nil {
		fmt.Fprintf(os.Stderr, "Invalid parameters: %s\n", err)
		os.Exit(1)
	}

	sourcesDir := os.Args[1]

	commandRunner := bosh.NewCommandRunner(checkRequest.Source, os.Stderr)
	director := bosh.NewBoshDirector(checkRequest.Source, commandRunner, sourcesDir, os.Stderr)

	checkCommand := check.NewCheckCommand(director)
	checkResponse, err := checkCommand.Run(checkRequest)
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)

	}

	concourseOutputFormatted, err := json.MarshalIndent(checkResponse, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not generate version: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("%s", concourseOutputFormatted)
}