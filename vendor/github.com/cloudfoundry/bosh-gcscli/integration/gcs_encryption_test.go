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
	"bytes"
	"crypto/sha256"

	"github.com/cloudfoundry/bosh-gcscli/client"
	"github.com/cloudfoundry/bosh-gcscli/config"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

// encryptionKeyBytes are used as the key in tests requiring encryption.
var encryptionKeyBytes = []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31}

// encryptionKeyBytesHash is the has of the encryptionKeyBytes
//
// Typical usage is ensuring the encryption key is actually used by GCS.
var encryptionKeyBytesHash = sha256.Sum256(encryptionKeyBytes)

var _ = Describe("Integration", func() {
	Context("general (Default Applicaton Credentials) configuration", func() {
		var (
			env AssertContext
			cfg *config.GCSCli
		)
		BeforeEach(func() {
			cfg = getMultiRegionConfig()
			cfg.EncryptionKey = encryptionKeyBytes

			env = NewAssertContext(AsDefaultCredentials)
			env.AddConfig(cfg)
		})
		AfterEach(func() {
			env.Cleanup()
		})

		// tests that a blob uploaded with a specified encryption_key can be downloaded again.
		It("can perform encrypted lifecycle", func() {
			AssertLifecycleWorks(gcsCLIPath, env)
		})

		// tests that uploading a blob with encryption
		// results in failure to download when the key is changed.
		It("fails to get with the wrong encryption_key", func() {
			Expect(env.Config.EncryptionKey).ToNot(BeNil(),
				"Need encryption key for test")

			session, err := RunGCSCLI(gcsCLIPath, env.ConfigPath,
				"put", env.ContentFile, env.GCSFileName)
			Expect(err).ToNot(HaveOccurred())
			Expect(session.ExitCode()).To(BeZero())

			blobstoreClient, err := client.New(env.ctx, env.Config)
			Expect(err).ToNot(HaveOccurred())

			env.Config.EncryptionKey[0]++

			var target bytes.Buffer
			err = blobstoreClient.Get(env.GCSFileName, &target)
			Expect(err).To(HaveOccurred())

			session, err = RunGCSCLI(gcsCLIPath, env.ConfigPath, "delete", env.GCSFileName)
			Expect(err).ToNot(HaveOccurred())
			Expect(session.ExitCode()).To(BeZero())
		})

		// tests that uploading a blob with encryption
		// results in failure to download without encryption.
		It("fails to get with no encryption_key", func() {
			Expect(env.Config.EncryptionKey).ToNot(BeNil(),
				"Need encryption key for test")

			session, err := RunGCSCLI(gcsCLIPath, env.ConfigPath, "put", env.ContentFile, env.GCSFileName)
			Expect(err).ToNot(HaveOccurred())
			Expect(session.ExitCode()).To(BeZero())

			blobstoreClient, err := client.New(env.ctx, env.Config)
			Expect(err).ToNot(HaveOccurred())

			env.Config.EncryptionKey = nil

			var target bytes.Buffer
			err = blobstoreClient.Get(env.GCSFileName, &target)
			Expect(err).To(HaveOccurred())

			session, err = RunGCSCLI(gcsCLIPath, env.ConfigPath, "delete", env.GCSFileName)
			Expect(err).ToNot(HaveOccurred())
			Expect(session.ExitCode()).To(BeZero())
		})
	})
})
