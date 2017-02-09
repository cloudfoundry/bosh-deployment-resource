package out

import (
	"github.com/cloudfoundry/bosh-deployment-resource/bosh"
	"github.com/cloudfoundry/bosh-deployment-resource/concourse"
	"github.com/cloudfoundry/bosh-deployment-resource/tools"
	"fmt"
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
	releasePaths, err := tools.UnfurlGlobs(outRequest.Params.Releases...)
	if err != nil {
		return OutResponse{}, fmt.Errorf("Invalid release name: %s", err)
	}
	for _, releasePath := range releasePaths {
		c.director.UploadRelease(releasePath)
	}

	err = c.director.Deploy(outRequest.Params.Manifest, bosh.DeployParams{
		NoRedact: outRequest.Params.NoRedact,
		Cleanup:  outRequest.Params.Cleanup,
	})
	if err != nil {
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
