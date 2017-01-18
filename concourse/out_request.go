package concourse

import (
	"encoding/json"
	"io"
	"io/ioutil"
)

type OutRequest struct {
	Params OutParams `json:"params"`
	Source Source    `json:"source"`
}

func NewOutRequest(r io.Reader) (OutRequest, error) {
	var outRequest OutRequest
	if err := json.NewDecoder(r).Decode(&outRequest); err != nil {
		return OutRequest{}, err
	}

	if outRequest.Params.TargetFile != "" {
		target, err := ioutil.ReadFile(outRequest.Params.TargetFile)
		if err != nil {
			return OutRequest{}, err
		}

		outRequest.Source.Target = string(target)
	}

	return outRequest, nil
}
