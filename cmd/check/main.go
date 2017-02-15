package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/cloudfoundry/bosh-deployment-resource/bosh"
	"github.com/cloudfoundry/bosh-deployment-resource/check"
	"github.com/cloudfoundry/bosh-deployment-resource/concourse"
	"io/ioutil"
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

	stdin, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Cannot read configuration: %s\n", err)
		os.Exit(1)
	}

	checkRequest, err := concourse.NewCheckRequest(stdin)
	if err != nil {
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

	fmt.Fprintf(realStdout, "%s", concourseOutputFormatted)
}
