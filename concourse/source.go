package concourse

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/cloudfoundry/cli/cf/errors"
	"io/ioutil"
	"strings"
)

type Source struct {
	Deployment   string `json:"deployment"`
	Client       string `json:"client"`
	ClientSecret string `json:"client_secret"`
	Target       string `json:"target"`
	CACert       string `json:"ca_cert"`
}

type sourceRequest struct {
	Source Source              `json:"source"`
	Params dynamicSourceParams `json:"params"`
}

type dynamicSourceParams struct {
	TargetFile string `json:"target_file,omitempty"`
}

func NewDynamicSource(config []byte) (Source, error) {
	var sourceRequest sourceRequest
	if err := json.NewDecoder(bytes.NewReader(config)).Decode(&sourceRequest); err != nil {
		return Source{}, fmt.Errorf("Invalid dynamic source config: %s", err)
	}

	if sourceRequest.Params.TargetFile != "" {
		target, err := ioutil.ReadFile(sourceRequest.Params.TargetFile)
		if err != nil {
			return Source{}, fmt.Errorf("Invalid dynamic source config: %s", err)
		}

		sourceRequest.Source.Target = string(target)
	}

	if err := checkRequiredSourceParameters(sourceRequest.Source); err != nil {
		return Source{}, err
	}

	return sourceRequest.Source, nil
}

func checkRequiredSourceParameters(source Source) error {
	missingParameters := []string{}

	if source.Deployment == "" {
		missingParameters = append(missingParameters, "deployment")
	}
	if source.Target == "" {
		missingParameters = append(missingParameters, "target")
	}
	if source.Client == "" {
		missingParameters = append(missingParameters, "client")
	}
	if source.ClientSecret == "" {
		missingParameters = append(missingParameters, "client_secret")
	}

	if len(missingParameters) > 0 {
		parametersString := "parameters"
		if len(missingParameters) > 2 {
			parametersString = parametersString + "s"
		}
		errorMessage := fmt.Sprintf("Missing required source %s: %s", parametersString, strings.Join(missingParameters, ", "))
		return errors.New(errorMessage)
	}

	return nil
}
