package concourse

import (
	"bytes"
	"encoding/json"
	"fmt"
)

type CheckRequest struct {
	Source  Source  `json:"source"`
	Version Version `json:"version"`
}

func NewCheckRequest(request []byte) (CheckRequest, error) {
	var checkRequest CheckRequest
	if err := json.NewDecoder(bytes.NewReader(request)).Decode(&checkRequest); err != nil {
		return CheckRequest{}, fmt.Errorf("Invalid parameters: %s\n", err)
	}

	return checkRequest, nil
}
