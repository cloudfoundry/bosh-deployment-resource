package out

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/cloudfoundry/bosh-deployment-resource/bosh"
	"github.com/cloudfoundry/bosh-deployment-resource/concourse"
	"github.com/cloudfoundry/bosh-deployment-resource/storage"
	"github.com/cloudfoundry/bosh-deployment-resource/tools"
)

type OutResponse struct {
	Version  concourse.Version    `json:"version"`
	Metadata []concourse.Metadata `json:"metadata"`
}

type OutCommand struct {
	director           bosh.Director
	storageClient      storage.StorageClient
	resourcesDirectory string
}

func NewOutCommand(director bosh.Director, storageClient storage.StorageClient, resourcesDirectory string) OutCommand {
	return OutCommand{
		director:           director,
		storageClient:      storageClient,
		resourcesDirectory: resourcesDirectory,
	}
}

func (c OutCommand) Run(outRequest concourse.OutRequest) (OutResponse, error) {
	if err := c.director.WaitForDeployLock(); err != nil {
		return OutResponse{}, err
	}

	if outRequest.Params.Delete.Enabled {
		return OutResponse{}, c.director.Delete(outRequest.Params.Delete.Force)
	} else {
		return c.deploy(outRequest)
	}
}

func (c OutCommand) deploy(outRequest concourse.OutRequest) (OutResponse, error) {
	manifestBytes, err := ioutil.ReadFile(path.Join(c.resourcesDirectory, outRequest.Params.Manifest))
	if err != nil {
		return OutResponse{}, err
	}

	varsFilePaths, err := tools.UnfurlGlobs(c.resourcesDirectory, outRequest.Params.VarsFiles)
	if err != nil {
		return OutResponse{}, fmt.Errorf("Invalid var_file name: %s", err)
	}

	opsFilePaths, err := tools.UnfurlGlobs(c.resourcesDirectory, outRequest.Params.OpsFiles)
	if err != nil {
		return OutResponse{}, fmt.Errorf("Invalid ops_file name: %s", err)
	}

	interpolateParams := bosh.InterpolateParams{
		Vars:      outRequest.Params.Vars,
		VarsFiles: varsFilePaths,
		OpsFiles:  opsFilePaths,
	}

	manifestBytes, err = c.director.Interpolate(manifestBytes, interpolateParams)
	if err != nil {
		return OutResponse{}, err
	}

	manifest, err := bosh.NewDeploymentManifest(manifestBytes)
	if err != nil {
		return OutResponse{}, err
	}

	releaseMetadata, err := c.consumeReleases(manifest, outRequest.Params.Releases)
	if err != nil {
		return OutResponse{}, err
	}

	stemcellMetadata, err := c.consumeStemcells(manifest, outRequest.Params.Stemcells)
	if err != nil {
		return OutResponse{}, err
	}

	deployParams := bosh.DeployParams{
		NoRedact: outRequest.Params.NoRedact,
		DryRun:   outRequest.Params.DryRun,
		Recreate: outRequest.Params.Recreate,
		Cleanup:  outRequest.Params.Cleanup,
	}

	var varsStoreFile *os.File
	if c.storageClient != nil {
		varsStoreFile, err = ioutil.TempFile("", "vars-store")
		if err != nil {
			return OutResponse{}, err
		}
		defer varsStoreFile.Close()

		if err = c.storageClient.Download(varsStoreFile.Name()); err != nil {
			return OutResponse{}, err
		}

		deployParams.VarsStore = varsStoreFile.Name()
	}

	if err := c.director.Deploy(manifest.Manifest(), deployParams); err != nil {
		return OutResponse{}, err
	}

	if c.storageClient != nil {
		if err := c.storageClient.Upload(varsStoreFile.Name()); err != nil {
			return OutResponse{}, err
		}
	}

	uploadedManifest, err := c.director.DownloadManifest()
	if err != nil {
		return OutResponse{}, err
	}

	concourseOutput := OutResponse{
		Version:  concourse.NewVersion(uploadedManifest, outRequest.Source.Target),
		Metadata: append(releaseMetadata, stemcellMetadata...),
	}

	return concourseOutput, nil
}

func (c OutCommand) consumeReleases(manifest bosh.DeploymentManifest, releaseGlobs []string) ([]concourse.Metadata, error) {
	releases, err := bosh.NewReleases(c.resourcesDirectory, releaseGlobs)
	if err != nil {
		return nil, err
	}

	metadata := []concourse.Metadata{}

	for _, release := range releases {
		if err := c.director.UploadRelease(release.FilePath); err != nil {
			return nil, err
		}

		if err := manifest.UseReleaseVersion(release.Name, release.Version); err != nil {
			return nil, err
		}

		metadata = append(metadata, concourse.Metadata{
			Name:  "release",
			Value: fmt.Sprintf("%s v%s", release.Name, release.Version),
		})
	}

	return metadata, nil
}

func (c OutCommand) consumeStemcells(manifest bosh.DeploymentManifest, stemcellGlobs []string) ([]concourse.Metadata, error) {
	stemcells, err := bosh.NewStemcells(c.resourcesDirectory, stemcellGlobs)
	if err != nil {
		return nil, err
	}

	metadata := []concourse.Metadata{}

	for _, stemcell := range stemcells {
		if err := c.director.UploadStemcell(stemcell.FilePath); err != nil {
			return nil, err
		}

		if err := manifest.UseStemcellVersion(stemcell.Name, stemcell.OperatingSystem, stemcell.Version); err != nil {
			return nil, err
		}

		metadata = append(metadata, concourse.Metadata{
			Name:  "stemcell",
			Value: fmt.Sprintf("%s v%s", stemcell.Name, stemcell.Version),
		})
	}

	return metadata, nil
}
