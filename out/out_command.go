package out

import (
	"github.com/cloudfoundry/bosh-deployment-resource/bosh"
	"github.com/cloudfoundry/bosh-deployment-resource/concourse"
)

type OutResponse struct {
	Version concourse.Version `json:"version"`
}

type OutCommand struct {
	director bosh.Director
}

func NewOutCommand(director bosh.Director) OutCommand {
	return OutCommand{
		director: director,
	}
}

func (c OutCommand) Run(outRequest concourse.OutRequest) (OutResponse, error) {
	if err := c.director.Deploy(outRequest.Params.Manifest); err != nil {
		return OutResponse{}, err
	}

	manifest, err := c.director.DownloadManifest()
	if err != nil {
		return OutResponse{}, err
	}

	concourseOutput := OutResponse{
		Version: concourse.NewVersion(manifest, outRequest.Source.Target),
	}

	return concourseOutput, nil
}
