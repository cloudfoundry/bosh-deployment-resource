package check

import (
	"github.com/cloudfoundry/bosh-deployment-resource/bosh"
	"github.com/cloudfoundry/bosh-deployment-resource/concourse"
	"crypto/sha1"
	"fmt"
)

type CheckCommand struct {
	director bosh.Director
}

func NewCheckCommand(director bosh.Director) CheckCommand {
	return CheckCommand{
		director: director,
	}
}

func (c CheckCommand) Run(checkRequest concourse.CheckRequest) ([]concourse.Version, error) {
	manifest, err := c.director.DownloadManifest()
	if err != nil {
		return []concourse.Version{}, err
	}

	byteSum := sha1.Sum(manifest)

	version := concourse.Version{
		ManifestSha1: fmt.Sprintf("%x", byteSum),
		Target: checkRequest.Source.Target,
	}

	var concourseOutput = []concourse.Version{}
	if version != checkRequest.Version {
		 concourseOutput = append(concourseOutput, version)
	}

	return concourseOutput, nil
}
