package storage

import (
	"github.com/cloudfoundry/bosh-deployment-resource/concourse"
	"github.com/cloudfoundry/bosh-deployment-resource/gcp"
)

type StorageClient interface {
	Download(filePath string) error
	Upload(filePath string) error
}

func NewStorageClient(source concourse.Source) (StorageClient, error) {
	if source.VarsStore.Provider == "gcs" {
		return gcp.NewStorage(
			source.VarsStore.Config["json_key"].(string),
			source.VarsStore.Config["bucket"].(string),
			source.VarsStore.Config["file_name"].(string),
		)
	}

	return nil, nil
}
