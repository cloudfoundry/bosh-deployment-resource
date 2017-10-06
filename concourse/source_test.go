package concourse_test

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/cloudfoundry/bosh-deployment-resource/concourse"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("NewDynamicSource", func() {
	It("converts the config into a Source", func() {
		config := []byte(`{
			"source": {
				"deployment": "mydeployment",
				"target": "director.example.com",
				"client": "foo",
				"client_secret": "foobar"
			}
		}`)

		source, err := concourse.NewDynamicSource(config, "")
		Expect(err).NotTo(HaveOccurred())

		Expect(source).To(Equal(concourse.Source{
			Deployment:   "mydeployment",
			Target:       "director.example.com",
			Client:       "foo",
			ClientSecret: "foobar",
		}))
	})

	Context("when source_file param is passed", func() {
		var (
			sourcesDir          string
			sourceFileName      string
			requestJsonTemplate string = `{
				"params": {
					"source_file": "%s"
				},
				"source": {
					"deployment": "mydeployment",
					"target": "director.example.com",
					"client": "original_client",
					"client_secret": "foobar",
					"jumpbox_ssh_key": "some-ssh-key",
					"jumpbox_url": "jumpbox.example.com",
					"vars_store": {
						"provider": "gcs",
						"config": {
							"some": "dynamic",
							"keys": "per-provider"
						}
					}
				}
			}`
		)

		BeforeEach(func() {
			sourceFile, _ := ioutil.TempFile("", "")
			sourceFile.Write(properYaml(`
				deployment: fileDeployment
				target: fileDirector.com
				client_secret: fileSecret
				vars_store:
					provider: fileProvider
					config:
						file: vars
						keys: dynamic-keys
			`))
			sourceFile.Close()

			sourcesDir = filepath.Dir(sourceFile.Name())
			sourceFileName = filepath.Base(sourceFile.Name())
		})

		It("overrides source with the values in the source_file", func() {
			config := []byte(fmt.Sprintf(
				requestJsonTemplate,
				filepath.Base(sourceFileName),
			))

			source, err := concourse.NewDynamicSource(config, sourcesDir)
			Expect(err).NotTo(HaveOccurred())

			Expect(source).To(Equal(concourse.Source{
				Deployment:    "fileDeployment",
				Target:        "fileDirector.com",
				Client:        "original_client",
				ClientSecret:  "fileSecret",
				JumpboxSSHKey: "some-ssh-key",
				JumpboxURL:    "jumpbox.example.com",
				VarsStore: concourse.VarsStore{
					Provider: "fileProvider",
					Config: map[string]interface{}{
						"file": "vars",
						"some": "dynamic",
						"keys": "dynamic-keys",
					},
				},
			}))
		})

		Context("when the target_file cannot be read", func() {
			BeforeEach(func() {
				sourceFileName = "not-a-real-file"
			})

			It("errors", func() {
				config := []byte(fmt.Sprintf(requestJsonTemplate, sourceFileName))

				_, err := concourse.NewDynamicSource(config, sourcesDir)
				Expect(err).To(HaveOccurred())
			})
		})
	})

	Context("when decoding fails", func() {
		It("errors", func() {
			reader := []byte("not-json")

			_, err := concourse.NewDynamicSource(reader, "")
			Expect(err).To(HaveOccurred())
		})
	})

	Context("when a required parameter is missing", func() {
		It("returns an error with each missing parameter", func() {
			config := []byte("{}")

			_, err := concourse.NewDynamicSource(config, "")
			Expect(err).To(HaveOccurred())

			Expect(err.Error()).To(ContainSubstring("deployment"))
			Expect(err.Error()).To(ContainSubstring("target"))
			Expect(err.Error()).To(ContainSubstring("client"))
			Expect(err.Error()).To(ContainSubstring("client_secret"))
		})
	})
})
