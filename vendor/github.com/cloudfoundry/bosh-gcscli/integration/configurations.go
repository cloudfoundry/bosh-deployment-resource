/*
 * Copyright 2017 Google Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package integration

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"

	"cloud.google.com/go/storage"

	"github.com/cloudfoundry/bosh-gcscli/config"

	. "github.com/onsi/ginkgo/extensions/table"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2/jwt"
	"google.golang.org/api/option"
)

const regionalBucketEnv = "REGIONAL_BUCKET_NAME"
const multiRegionalBucketEnv = "MULTIREGIONAL_BUCKET_NAME"
const publicBucketEnv = "PUBLIC_BUCKET_NAME"

// noBucketMsg is the template used when a BucketEnv's environment variable
// has not been populated.
const noBucketMsg = "environment variable %s expected to contain a valid Google Cloud Storage bucket but was empty"

const getConfigErrMsg = "creating %s configs: %v"

func readBucketEnv(env string) (string, error) {
	bucket := os.Getenv(env)
	if len(bucket) == 0 {
		return "", fmt.Errorf(noBucketMsg, env)
	}
	return bucket, nil
}

func getRegionalConfig() *config.GCSCli {
	var regional string
	var err error

	if regional, err = readBucketEnv(regionalBucketEnv); err != nil {
		panic(fmt.Errorf(getConfigErrMsg, "base", err))
	}

	return &config.GCSCli{BucketName: regional}
}

func getMultiRegionConfig() *config.GCSCli {
	var multiRegional string
	var err error

	if multiRegional, err = readBucketEnv(multiRegionalBucketEnv); err != nil {
		panic(fmt.Errorf(getConfigErrMsg, "base", err))
	}

	return &config.GCSCli{BucketName: multiRegional}
}

func getBaseConfigs() []TableEntry {
	regional := getRegionalConfig()
	multiRegion := getMultiRegionConfig()

	return []TableEntry{
		Entry("Regional bucket, default StorageClass", regional),
		Entry("MultiRegion bucket, default StorageClass", multiRegion),
	}
}

func getPublicConfig() *config.GCSCli {
	public, err := readBucketEnv(publicBucketEnv)
	if err != nil {
		panic(fmt.Errorf(getConfigErrMsg, "public", err))
	}

	return &config.GCSCli{
		BucketName: public,
	}
}

func getInvalidStorageClassConfigs() []TableEntry {
	regional := getRegionalConfig()
	multiRegion := getMultiRegionConfig()

	multiRegion.StorageClass = "REGIONAL"
	regional.StorageClass = "MULTI_REGIONAL"

	return []TableEntry{
		Entry("Multi-Region bucket, regional StorageClass", regional),
		Entry("Regional bucket, Multi-Region StorageClass", multiRegion),
	}
}

// newSDK builds the GCS SDK Client from a valid config.GCSCli
// TODO: Simplify and remove this. Tests should expect a single config and use it.
func newSDK(ctx context.Context, c config.GCSCli) (*storage.Client, error) {
	var client *storage.Client
	var err error
	var opt option.ClientOption
	switch c.CredentialsSource {
	case config.DefaultCredentialsSource:
		var tokenSource oauth2.TokenSource
		tokenSource, err = google.DefaultTokenSource(ctx, storage.ScopeFullControl)
		if err == nil {
			opt = option.WithTokenSource(tokenSource)
		}
	case config.NoneCredentialsSource:
		opt = option.WithHTTPClient(http.DefaultClient)
	case config.ServiceAccountFileCredentialsSource:
		var token *jwt.Config
		token, err = google.JWTConfigFromJSON([]byte(c.ServiceAccountFile), storage.ScopeFullControl)
		if err == nil {
			tokenSource := token.TokenSource(ctx)
			opt = option.WithTokenSource(tokenSource)
		}
	default:
		err = errors.New("unknown credentials_source in configuration")
	}
	if err != nil {
		return client, err
	}

	return storage.NewClient(ctx, opt)
}
