package bosh

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"github.com/cloudfoundry/bosh-deployment-resource/tools"
)

type Release struct {
	Name     string

	Version  string
	FilePath string
}

func NewReleases(basePath string, releasePathGlobs []string) ([]Release, error) {
	releasePaths, err := tools.UnfurlGlobs(basePath, releasePathGlobs)
	if err != nil {
		return nil, fmt.Errorf("Invalid release name: %s", err)
	}

	releases := []Release{}
	for _, releasePath := range releasePaths {
		release, err := newRelease(releasePath)
		if err != nil {
			return nil, err
		}
		releases = append(releases, release)
	}

	return releases, nil
}

func newRelease(filePath string) (Release, error) {
	release := Release{}

	releaseFileContents, err := tools.ReadTgzFile(filePath, "release.MF")
	if err != nil {
		return Release{}, fmt.Errorf("Could not read release: %s", err)
	}

	err = yaml.Unmarshal(releaseFileContents, &release)
	if err != nil {
		return Release{}, fmt.Errorf("Release %s is not a valid release", filePath)
	}

	release.FilePath = filePath

	return release, nil
}