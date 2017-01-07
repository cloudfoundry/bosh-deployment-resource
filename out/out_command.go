package out

import (
	"github.com/cloudfoundry/bosh-deployment-resource/concourse"
	"github.com/cloudfoundry/bosh-deployment-resource/bosh"
	"fmt"
	"crypto/sha1"
)

type OutResponse struct {
	Version Version `json:"version"`
}

type OutCommand struct {
	director bosh.Director
}

type Version struct {
	ManifestSha1 string `json:"manifest_sha1"`
	Target       string `json:"target"`
}

func NewOutCommand(director bosh.Director) OutCommand {
	return OutCommand{
		director: director,
	}
}

func (c OutCommand) Run(outRequest concourse.OutRequest) (OutResponse, error){
	if err := c.director.Deploy(outRequest.Params.Manifest); err != nil {
		return OutResponse{}, err
	}

	manifest, err := c.director.DownloadManifest()
	if err != nil {
		return OutResponse{}, err
	}

	byteSum := sha1.Sum(manifest)

	concourseOutput := OutResponse{
		Version: Version{
			ManifestSha1: fmt.Sprintf("%x", byteSum),
			Target: outRequest.Source.Target,
		},
	}

	return concourseOutput, nil
}
