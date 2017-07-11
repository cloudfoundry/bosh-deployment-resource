package in

import (
	"errors"
	"io/ioutil"
	"path/filepath"

	"github.com/cloudfoundry/bosh-deployment-resource/bosh"
	"github.com/cloudfoundry/bosh-deployment-resource/concourse"
)

type InCommand struct {
	director bosh.Director
}

type InResponse struct {
	Version concourse.Version `json:"version"`
}

func NewInCommand(director bosh.Director) InCommand {
	return InCommand{
		director: director,
	}
}

func (c InCommand) Run(inRequest concourse.InRequest, targetDir string) (InResponse, error) {
	manifest, err := c.director.DownloadManifest()

	if err != nil {
		return InResponse{}, err
	}
	actualVersion := concourse.NewVersion(manifest, inRequest.Source.Target)

	if actualVersion.Target != inRequest.Version.Target {
		return InResponse{}, errors.New("Requested deployment director is different than configured source")
	}

	if actualVersion.ManifestSha1 != inRequest.Version.ManifestSha1 {
		return InResponse{}, errors.New("Requested deployment version is not available")
	}

	if len(inRequest.Params.CompiledReleases) > 0 {
		releases := []string{}
		for _, compiledRelease := range inRequest.Params.CompiledReleases {
			releases = append(releases, compiledRelease.Name)
		}
		if err := c.director.ExportReleases(targetDir, releases); err != nil {
			return InResponse{}, err
		}
	}

	err = ioutil.WriteFile(filepath.Join(targetDir, "manifest.yml"), manifest, 0644)
	if err != nil {
		return InResponse{}, err
	}

	err = ioutil.WriteFile(filepath.Join(targetDir, "target"), []byte(actualVersion.Target), 0644)
	if err != nil {
		return InResponse{}, err
	}

	return InResponse{Version: actualVersion}, nil
}
