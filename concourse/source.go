package concourse

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"
)

type Source struct {
	Deployment      string    `json:"deployment,omitempty" yaml:"deployment"`
	Client          string    `json:"client,omitempty" yaml:"client"`
	ClientSecret    string    `json:"client_secret,omitempty" yaml:"client_secret"`
	Target          string    `json:"target,omitempty" yaml:"target"`
	CACert          string    `json:"ca_cert,omitempty" yaml:"ca_cert"`
	JumpboxSSHKey   string    `json:"jumpbox_ssh_key,omitempty" yaml:"jumpbox_ssh_key"`
	JumpboxURL      string    `json:"jumpbox_url,omitempty" yaml:"jumpbox_url"`
	JumpboxUsername string    `json:"jumpbox_username,omitempty" yaml:"jumpbox_username"`
	VarsStore       VarsStore `json:"vars_store,omitempty" yaml:"vars_store"`
	SkipCheck       bool      `json:"skip_check,omitempty" yaml:"skip_check"`
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
		return Source{}, fmt.Errorf("Invalid dynamic source config: %s", err) //nolint:staticcheck
	}

	if sourceRequest.Params.SourceFile != "" {
		source, err := os.ReadFile(filepath.Join(sourcesDir, sourceRequest.Params.SourceFile))
		if err != nil {
			return Source{}, fmt.Errorf("Invalid dynamic source config: %s", err) //nolint:staticcheck
		}

		var tempSource Source
		err = yaml.Unmarshal(source, &tempSource)
		if err != nil {
			return Source{}, fmt.Errorf("Invalid dynamic source config: %s", err) //nolint:staticcheck
		}

		jsonBytes, err := json.Marshal(tempSource)
		if err != nil {
			return Source{}, fmt.Errorf("Invalid dynamic source config: %s", err) //nolint:staticcheck
		}

		if err := json.Unmarshal(jsonBytes, &sourceRequest.Source); err != nil {
			return Source{}, fmt.Errorf("Invalid dynamic source config: %s", err) //nolint:staticcheck
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
