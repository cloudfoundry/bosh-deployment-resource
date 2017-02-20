package gcp

import (
	"golang.org/x/oauth2"
	oauthgoogle "golang.org/x/oauth2/google"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/storage/v1"
	"net/http"
	"os"

	"io/ioutil"
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
	storageClient = storageJwtConf.Client(oauth2.NoContext)

	storageService, err := storage.New(storageClient)
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
		switch err.(type) {
		case *googleapi.Error:
			if err.(*googleapi.Error).Code == 404 {
				return nil
			}
		}

		return err
	}

	response, err := getCall.Download()
	if err != nil {
		return err
	}
	defer response.Body.Close()

	responseBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(filePath, responseBytes, 0600)
	if err != nil {
		return err
	}

	return nil
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
		return err
	}

	return nil
}
