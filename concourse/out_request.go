package concourse

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

type OutRequest struct {
	Params OutParams `json:"params"`
	Source Source    `json:"source"`
}

func NewOutRequest(request []byte, sourcesDir string) (OutRequest, error) {
	var outRequest OutRequest
	if err := json.NewDecoder(bytes.NewReader(request)).Decode(&outRequest); err != nil {
		return OutRequest{}, fmt.Errorf("Invalid parameters: %s\n", err)
	}

	dynamicSource, err := NewDynamicSource(request, sourcesDir)
	if err != nil {
		return OutRequest{}, err
	}

	outRequest.Source = dynamicSource

	if err := checkRequiredOutParameters(outRequest.Params); err != nil {
		return OutRequest{}, err
	}

	return outRequest, nil
}

func checkRequiredOutParameters(params OutParams) error {
	missingParameters := []string{}

	if params.Manifest == "" && !params.Delete.Enabled {
		missingParameters = append(missingParameters, "manifest")
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
