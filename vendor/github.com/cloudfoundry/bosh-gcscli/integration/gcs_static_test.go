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
	"github.com/cloudfoundry/bosh-gcscli/config"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"net/http"
	"strings"
)

var _ = Describe("Integration", func() {
	Context("static credentials configuration with a regional bucket", func() {
		var (
			ctx AssertContext
			cfg *config.GCSCli
		)
		BeforeEach(func() {
			cfg = getRegionalConfig()
			ctx = NewAssertContext(AsStaticCredentials)
			ctx.AddConfig(cfg)
		})
		AfterEach(func() {
			ctx.Cleanup()
		})

		It("can perform blobstore lifecycle", func() {
			AssertLifecycleWorks(gcsCLIPath, ctx)
		})

		It("can generate a signed url for a given object and action", func() {
			session, err := RunGCSCLI(gcsCLIPath, ctx.ConfigPath,
				"sign", ctx.GCSFileName, "PUT", "1h")

			Expect(err).ToNot(HaveOccurred())
			Expect(session.ExitCode()).To(Equal(0))
			url := string(session.Out.Contents())
			Expect(url).To(MatchRegexp("https://"))

			body := strings.NewReader(`bar`)
			req, err := http.NewRequest("PUT", url, body)
			Expect(err).ToNot(HaveOccurred())

			resp, err := http.DefaultClient.Do(req)
			Expect(err).ToNot(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(200))
			defer resp.Body.Close()
		})
	})
})
