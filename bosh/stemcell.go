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
	FilePath        string
}

func NewStemcells(basePath string, stemcellPathGlobs []string) ([]Stemcell, error) {
	stemcellPaths, err := tools.UnfurlGlobs(basePath, stemcellPathGlobs)
	if err != nil {
		return nil, fmt.Errorf("Invalid stemcell name: %s", err)
	}

	stemcells := []Stemcell{}
	for _, stemcellPath := range stemcellPaths {
		stemcell, err := newStemcell(stemcellPath)
		if err != nil {
			return nil, err
		}
		stemcells = append(stemcells, stemcell)
	}

	return stemcells, nil
}

func newStemcell(filePath string) (Stemcell, error) {
	stemcell := Stemcell{}

	stemcellFileContents, err := tools.ReadTgzFile(filePath, "stemcell.MF")
	if err != nil {
		return Stemcell{}, fmt.Errorf("Could not read stemcell: %s", err)
	}

	err = yaml.Unmarshal(stemcellFileContents, &stemcell)
	if err != nil {
		return Stemcell{}, fmt.Errorf("Stemcell %s is not a valid stemcell", filePath)
	}

	stemcell.FilePath = filePath

	return stemcell, nil
}