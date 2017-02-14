package bosh

import (
	"os/exec"
	"fmt"
	"gopkg.in/yaml.v2"
)

type Stemcell struct {
	Name            string
	OperatingSystem string `yaml:"operating_system"`
	Version         string
}

func NewStemcell(filePath string) (Stemcell, error) {
	stemcell := Stemcell{}

	readStemcellCommand := exec.Command("tar", "-Oxzf", filePath, "./stemcell.MF")
	stemcellFileContents, err := readStemcellCommand.Output()
	if err != nil {
		return Stemcell{}, fmt.Errorf("Could not read stemcell %s", filePath)
	}

	err = yaml.Unmarshal(stemcellFileContents, &stemcell)
	if err != nil {
		return Stemcell{}, fmt.Errorf("Stemcell %s is not a valid stemcell", filePath)
	}

	return stemcell, nil
}