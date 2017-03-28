package storage

import (
	"encoding/json"

	"github.com/cloudfoundry/bosh-deployment-resource/concourse"
	"github.com/cloudfoundry/bosh-deployment-resource/gcp"
)

type GCSConfig struct {
	FileName string `json:"file_name"`
	Bucket   string `json:"bucket"`
	JSONKey  string `json:"json_key"`
}

type StorageClient interface {
	Download(filePath string) error
	Upload(filePath string) error
}

func NewStorageClient(source concourse.Source) (StorageClient, error) {
	if source.VarsStore.Provider == "gcs" {
		gcsConfigJson, err := json.Marshal(source.VarsStore.Config)
		if err != nil {
			return nil, err
		}

		gcsConfig := GCSConfig{}
		if err := json.Unmarshal(gcsConfigJson, &gcsConfig); err != nil {
			return nil, err
		}

		return gcp.NewStorage(
			gcsConfig.JSONKey,
			gcsConfig.Bucket,
			gcsConfig.FileName,
		)
	}

	return nil, nil
}
