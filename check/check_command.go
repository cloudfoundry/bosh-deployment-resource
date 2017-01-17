package check

import (
	"github.com/cloudfoundry/bosh-deployment-resource/bosh"
	"github.com/cloudfoundry/bosh-deployment-resource/concourse"
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

	version := concourse.NewVersion(manifest, checkRequest.Source.Target)

	var concourseOutput = []concourse.Version{}
	if version != checkRequest.Version {
		concourseOutput = append(concourseOutput, version)
	}

	return concourseOutput, nil
}
