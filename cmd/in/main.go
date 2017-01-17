package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/cloudfoundry/bosh-deployment-resource/bosh"
	"github.com/cloudfoundry/bosh-deployment-resource/concourse"
	"github.com/cloudfoundry/bosh-deployment-resource/in"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr,
			"not enough args - usage: %s <sources directory>\n",
			os.Args[0],
		)
		os.Exit(1)
	}

	var inRequest concourse.InRequest
	if err := json.NewDecoder(os.Stdin).Decode(&inRequest); err != nil {
		fmt.Fprintf(os.Stderr, "Invalid parameters: %s\n", err)
		os.Exit(1)
	}

	sourcesDir := os.Args[1]

	commandRunner := bosh.NewCommandRunner(inRequest.Source, os.Stderr)
	director := bosh.NewBoshDirector(inRequest.Source, commandRunner, sourcesDir, os.Stderr)

	inCommand := in.NewInCommand(director)
	inResponse, err := inCommand.Run(inRequest, sourcesDir)
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)

	}

	concourseInputFormatted, err := json.MarshalIndent(inResponse, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not generate version: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("%s", concourseInputFormatted)
}
