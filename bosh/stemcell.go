package bosh

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"github.com/cloudfoundry/bosh-deployment-resource/tools"
)

type Stemcell struct {
	Name            string
	OperatingSystem string `yaml:"operating_system"`
	Version         string
}

func NewStemcell(filePath string) (Stemcell, error) {
	stemcell := Stemcell{}

	stemcellFileContents, err := tools.ReadTgzFile(filePath, "stemcell.MF")
	if err != nil {
		return Stemcell{}, fmt.Errorf("Could not read stemcell: %s", err)
	}

	err = yaml.Unmarshal(stemcellFileContents, &stemcell)
	if err != nil {
		return Stemcell{}, fmt.Errorf("Stemcell %s is not a valid stemcell", filePath)
	}

	return stemcell, nil
}