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

	commandRunner := bosh.NewCommandRunner(checkRequest.Source, os.Stderr)
	director := bosh.NewBoshDirector(checkRequest.Source, commandRunner, os.Stderr)

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
