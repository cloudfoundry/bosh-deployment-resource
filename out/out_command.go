package out

import (
	"github.com/cloudfoundry/bosh-deployment-resource/bosh"
	"github.com/cloudfoundry/bosh-deployment-resource/concourse"
	"github.com/cloudfoundry/bosh-deployment-resource/tools"
	"fmt"
	"io/ioutil"
	"path"
)

type OutResponse struct {
	Version concourse.Version `json:"version"`
}

type OutCommand struct {
	director           bosh.Director
	resourcesDirectory string
}

func NewOutCommand(director bosh.Director, resourcesDirectory string) OutCommand {
	return OutCommand{
		director: director,
		resourcesDirectory: resourcesDirectory,
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

	stemcellPaths, err := tools.UnfurlGlobs(outRequest.Params.Stemcells...)
	if err != nil {
		return OutResponse{}, fmt.Errorf("Invalid stemcell name: %s", err)
	}
	for _, stemcellPath := range stemcellPaths {
		c.director.UploadStemcell(stemcellPath)
	}

	manifestBytes, err := ioutil.ReadFile(path.Join(c.resourcesDirectory, outRequest.Params.Manifest))
	if err != nil {
		return OutResponse{}, err
	}

	manifest, err := bosh.NewDeploymentManifest(manifestBytes)
	if err != nil {
		return OutResponse{}, err
	}

	err = updateReleasesInManifest(&manifest, releasePaths)
	if err != nil {
		return OutResponse{}, err
	}

	err = updateStemcellsInManifest(&manifest, stemcellPaths)
	if err != nil {
		return OutResponse{}, err
	}

	err = c.director.Deploy(manifest.Manifest(), bosh.DeployParams{
		NoRedact: outRequest.Params.NoRedact,
		Cleanup:  outRequest.Params.Cleanup,
	})
	if err != nil {
		return OutResponse{}, err
	}

	uploadedManifest, err := c.director.DownloadManifest()
	if err != nil {
		return OutResponse{}, err
	}

	concourseOutput := OutResponse{
		Version: concourse.NewVersion(uploadedManifest, outRequest.Source.Target),
	}

	return concourseOutput, nil
}

type ReleaseFile struct {
	Name string
	Version string
}

func updateReleasesInManifest(manifest *bosh.DeploymentManifest, releasePaths []string) error {
	for _, releasePath := range releasePaths {
		release, err := bosh.NewRelease(releasePath)
		if err != nil {
			return err
		}

		if err = manifest.UseReleaseVersion(release.Name, release.Version); err != nil {
			return err
		}
	}

	return nil
}

func updateStemcellsInManifest(manifest *bosh.DeploymentManifest, stemcellPaths []string) error {
	for _, stemcellPath := range stemcellPaths {
		stemcell, err := bosh.NewStemcell(stemcellPath)
		if err != nil {
			return err
		}

		if err = manifest.UseStemcellVersion(stemcell.Name, stemcell.OperatingSystem, stemcell.Version); err != nil {
			return err
		}
	}

	return nil
}