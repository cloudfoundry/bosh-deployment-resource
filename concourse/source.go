package concourse

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
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

	return sourceRequest.Source, nil
}
