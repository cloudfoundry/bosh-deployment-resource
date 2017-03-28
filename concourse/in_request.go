package concourse

import (
	"bytes"
	"encoding/json"
	"fmt"
)

const MissingTarget = "MISSING-TARGET-SHORTCIRCUIT.example.com"

type InRequest struct {
	Source  Source   `json:"source"`
	Version Version  `json:"version"`
	Params  InParams `json:"params"`
}

func NewInRequest(request []byte) (InRequest, error) {
	var inRequest InRequest
	if err := json.NewDecoder(bytes.NewReader(request)).Decode(&inRequest); err != nil {
		return InRequest{}, fmt.Errorf("Invalid parameters: %s\n", err)
	}

	if inRequest.Source.Target == "" {
		inRequest.Source.Target = MissingTarget
	}

	return inRequest, nil
}
