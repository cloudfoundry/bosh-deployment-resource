package bosh

import (
	"os/exec"
	"fmt"
	"gopkg.in/yaml.v2"
)

type Release struct {
	Name    string
	Version string
}

func NewRelease(filePath string) (Release, error) {
	release := Release{}

	readReleaseCommand := exec.Command("tar", "-Oxzf", filePath, "release.MF")
	releaseFileContents, err := readReleaseCommand.Output()
	if err != nil {
		return Release{}, fmt.Errorf("Could not read release %s", filePath)
	}

	err = yaml.Unmarshal(releaseFileContents, &release)
	if err != nil {
		return Release{}, fmt.Errorf("Release %s is not a valid release", filePath)
	}

	return release, nil
}