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
	Version  concourse.Version    `json:"version"`
	Metadata []concourse.Metadata `json:"metadata"`
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
	releasePaths, err := tools.UnfurlGlobs(c.resourcesDirectory, outRequest.Params.Releases)
	if err != nil {
		return OutResponse{}, fmt.Errorf("Invalid release name: %s", err)
	}
	for _, releasePath := range releasePaths {
		c.director.UploadRelease(releasePath)
	}

	stemcellPaths, err := tools.UnfurlGlobs(c.resourcesDirectory, outRequest.Params.Stemcells)
	if err != nil {
		return OutResponse{}, fmt.Errorf("Invalid stemcell name: %s", err)
	}
	for _, stemcellPath := range stemcellPaths {
		c.director.UploadStemcell(stemcellPath)
	}

	varFilePaths, err := tools.UnfurlGlobs(c.resourcesDirectory, outRequest.Params.VarsFiles)
	if err != nil {
		return OutResponse{}, fmt.Errorf("Invalid var_file name: %s", err)
	}

	opsFilePaths, err := tools.UnfurlGlobs(c.resourcesDirectory, outRequest.Params.OpsFiles)
	if err != nil {
		return OutResponse{}, fmt.Errorf("Invalid ops_file name: %s", err)
	}
	
	manifestBytes, err := ioutil.ReadFile(path.Join(c.resourcesDirectory, outRequest.Params.Manifest))
	if err != nil {
		return OutResponse{}, err
	}

	manifest, err := bosh.NewDeploymentManifest(manifestBytes)
	if err != nil {
		return OutResponse{}, err
	}

	releases, err := globbedReleases(releasePaths)
	if err != nil {
		return OutResponse{}, err
	}

	err = updateReleasesInManifest(&manifest, releases)
	if err != nil {
		return OutResponse{}, err
	}

	stemcells, err := globbedStemcells(stemcellPaths)
	if err != nil {
		return OutResponse{}, err
	}

	err = updateStemcellsInManifest(&manifest, stemcells)
	if err != nil {
		return OutResponse{}, err
	}

	err = c.director.Deploy(manifest.Manifest(), bosh.DeployParams{
		NoRedact: outRequest.Params.NoRedact,
		Cleanup:  outRequest.Params.Cleanup,
		Vars:  outRequest.Params.Vars,
		VarsFiles:  varFilePaths,
		OpsFiles:  opsFilePaths,
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

	for _, release := range releases {
		metadata := concourse.Metadata{
			Name: "release",
			Value: fmt.Sprintf("%s v%s", release.Name, release.Version),
		}
		concourseOutput.Metadata = append(concourseOutput.Metadata, metadata)
	}

	for _, stemcell := range stemcells {
		metadata := concourse.Metadata{
			Name: "stemcell",
			Value: fmt.Sprintf("%s v%s", stemcell.Name, stemcell.Version),
		}
		concourseOutput.Metadata = append(concourseOutput.Metadata, metadata)
	}

	return concourseOutput, nil
}

type ReleaseFile struct {
	Name string
	Version string
}

func globbedReleases(releasePaths []string) ([]bosh.Release, error) {
	releases := []bosh.Release{}
	for _, releasePath := range releasePaths {
		release, err := bosh.NewRelease(releasePath)
		if err != nil {
			return nil, err
		}
		releases = append(releases, release)
	}

	return releases, nil
}

func globbedStemcells(stemcellPaths []string) ([]bosh.Stemcell, error) {
	stemcells := []bosh.Stemcell{}
	for _, stemcellPath := range stemcellPaths {
		stemcell, err := bosh.NewStemcell(stemcellPath)
		if err != nil {
			return nil, err
		}
		stemcells = append(stemcells, stemcell)
	}

	return stemcells, nil
}

func updateReleasesInManifest(manifest *bosh.DeploymentManifest, releases []bosh.Release) error {
	for _, release := range releases {
		if err := manifest.UseReleaseVersion(release.Name, release.Version); err != nil {
			return err
		}
	}

	return nil
}

func updateStemcellsInManifest(manifest *bosh.DeploymentManifest, stemcells []bosh.Stemcell) error {
	for _, stemcell := range stemcells {
		if err := manifest.UseStemcellVersion(stemcell.Name, stemcell.OperatingSystem, stemcell.Version); err != nil {
			return err
		}
	}

	return nil
}