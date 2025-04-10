package gcp

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"golang.org/x/oauth2"
	oauthgoogle "golang.org/x/oauth2/google"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/storage/v1"
)

type Storage struct {
	bucket         string
	objectPath     string
	storageService *storage.Service
}

func NewStorage(jsonKey, bucket, objectPath string) (Storage, error) {
	var err error
	var storageClient *http.Client
	var userAgent = "bosh-deployment-resource"

	storageJwtConf, err := oauthgoogle.JWTConfigFromJSON([]byte(jsonKey), storage.DevstorageFullControlScope)
	if err != nil {
		return Storage{}, err
	}
	storageClient = storageJwtConf.Client(oauth2.NoContext) //nolint:staticcheck

	storageService, err := storage.New(storageClient) //nolint:staticcheck
	if err != nil {
		return Storage{}, err
	}
	storageService.UserAgent = userAgent

	return Storage{
		bucket:         bucket,
		objectPath:     objectPath,
		storageService: storageService,
	}, nil
}

func (s Storage) Download(filePath string) error {
	getCall := s.storageService.Objects.Get(s.bucket, s.objectPath)

	_, err := getCall.Do()
	if err != nil {
		switch err.(type) { //nolint:staticcheck
		case *googleapi.Error:
			if err.(*googleapi.Error).Code == 404 { //nolint:staticcheck
				return s.Upload(filePath)
			}
		}

		return err
	}

	response, err := getCall.Download()
	if err != nil {
		return err
	}
	defer response.Body.Close() //nolint:errcheck

	responseBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}

	err = os.WriteFile(filePath, responseBytes, 0600)
	if err != nil {
		return err
	}

	// Check that we can not only read the file, but can also write it
	return s.Upload(filePath)
}

func (s Storage) Upload(filePath string) error {
	f, err := os.Open(filePath)
	if err != nil {
		return err
	}

	object := &storage.Object{
		Name: s.objectPath,
	}

	if _, err = s.storageService.Objects.Insert(s.bucket, object).Media(f).Do(); err != nil {
		return fmt.Errorf("Can not write to %s in bucket %s", s.objectPath, s.bucket) //nolint:staticcheck
	}

	return nil
}
