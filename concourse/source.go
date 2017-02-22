package concourse

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
)

type Source struct {
	Deployment   string    `json:"deployment"`
	Client       string    `json:"client"`
	ClientSecret string    `json:"client_secret"`
	Target       string    `json:"target"`
	CACert       string    `json:"ca_cert"`
	VarsStore    VarsStore `json:"vars_store"`
}

type sourceRequest struct {
	Source Source              `json:"source"`
	Params dynamicSourceParams `json:"params"`
}

type dynamicSourceParams struct {
	SourceFile string `json:"source_file,omitempty"`
}

func NewDynamicSource(config []byte, sourcesDir string) (Source, error) {
	var sourceRequest sourceRequest
	if err := json.NewDecoder(bytes.NewReader(config)).Decode(&sourceRequest); err != nil {
		return Source{}, fmt.Errorf("Invalid dynamic source config: %s", err)
	}

	if sourceRequest.Params.SourceFile != "" {
		source, err := ioutil.ReadFile(filepath.Join(sourcesDir, sourceRequest.Params.SourceFile))
		if err != nil {
			return Source{}, fmt.Errorf("Invalid dynamic source config: %s", err)
		}

		if err := json.Unmarshal(source, &sourceRequest.Source); err != nil {
			return Source{}, fmt.Errorf("Invalid dynamic source config: %s", err)
		}
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
		parametersString := "parameter"
		if len(missingParameters) > 2 {
			parametersString = parametersString + "s"
		}
		errorMessage := fmt.Sprintf("Missing required source %s: %s", parametersString, strings.Join(missingParameters, ", "))
		return errors.New(errorMessage)
	}

	return nil
}
