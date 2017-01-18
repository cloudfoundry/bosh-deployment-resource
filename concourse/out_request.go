package concourse

import (
	"bytes"
	"encoding/json"
	"fmt"
)

type OutRequest struct {
	Params OutParams `json:"params"`
	Source Source    `json:"source"`
}

func NewOutRequest(request []byte) (OutRequest, error) {
	var outRequest OutRequest
	if err := json.NewDecoder(bytes.NewReader(request)).Decode(&outRequest); err != nil {
		return OutRequest{}, fmt.Errorf("Invalid parameters: %s\n", err)
	}

	dynamicSource, err := NewDynamicSource(request)
	if err != nil {
		return OutRequest{}, err
	}

	outRequest.Source = dynamicSource

	return outRequest, nil
}
