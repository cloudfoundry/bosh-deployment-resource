package bosh

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"github.com/cloudfoundry/bosh-deployment-resource/tools"
)

type Release struct {
	Name    string
	Version string
}

func NewRelease(filePath string) (Release, error) {
	release := Release{}

	releaseFileContents, err := tools.ReadTgzFile(filePath, "release.MF")
	if err != nil {
		return Release{}, fmt.Errorf("Could not read release: %s", err)
	}

	err = yaml.Unmarshal(releaseFileContents, &release)
	if err != nil {
		return Release{}, fmt.Errorf("Release %s is not a valid release", filePath)
	}

	return release, nil
}