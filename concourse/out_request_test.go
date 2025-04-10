package concourse_test

import (
	"fmt"
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/cloudfoundry/bosh-deployment-resource/concourse"
)

var _ = Describe("NewOutRequest", func() {
	It("converts the config into an OutRequest", func() {
		config := []byte(`{
			"params": {
				"manifest": "path/to/manifest.yml",
				"vars": {
					"foo": "bar",
					"slice": [1, "two"]
				},
				"vars_files": [
					"path/to/file",
					"second/path/to/file"
				],
				"var_files": {
					"var_from_file": "/path/to/file/with/var/content",
					"another_var_from_file": "/more/path/to/content"
				},
				"ops_files": [
					"ops-file1",
					"path/to/ops-file2"
				]
			},
			"source": {
				"deployment": "mydeployment",
				"target": "director.example.com",
				"client": "foo",
				"client_secret": "foobar",
				"vars_store": {
					"provider": "gcs",
					"config": {
						"some": "dynamic",
						"keys": "per-provider"
					}
				}
			}
		}`)

		source, err := concourse.NewOutRequest(config, "")
		Expect(err).NotTo(HaveOccurred())

		Expect(source).To(Equal(concourse.OutRequest{
			Source: concourse.Source{
				Deployment:   "mydeployment",
				Target:       "director.example.com",
				Client:       "foo",
				ClientSecret: "foobar",
				VarsStore: concourse.VarsStore{
					Provider: "gcs",
					Config: map[string]interface{}{
						"some": "dynamic",
						"keys": "per-provider",
					},
				},
			},
			Params: concourse.OutParams{
				Manifest: "path/to/manifest.yml",
				Vars: map[string]interface{}{
					"foo":   "bar",
					"slice": []interface{}{float64(1), "two"},
				},
				VarsFiles: []string{
					"path/to/file",
					"second/path/to/file",
				},
				VarFiles: map[string]string{
					"var_from_file":         "/path/to/file/with/var/content",
					"another_var_from_file": "/more/path/to/content",
				},
				OpsFiles: []string{
					"ops-file1",
					"path/to/ops-file2",
				},
			},
		}))
	})

	Context("when the dry run flag is true", func() {
		It("set dryrun to true in OutParams", func() {
			config := []byte(`{
				"params": {
					"manifest": "path/to/manifest.yml",
					"dry_run": true
				},
				"source": {
					"deployment": "mydeployment",
					"target": "director.example.com",
					"client": "foo",
					"client_secret": "foobar"
				}
			}`)

			source, err := concourse.NewOutRequest(config, "")
			Expect(err).NotTo(HaveOccurred())

			Expect(source).To(Equal(concourse.OutRequest{
				Source: concourse.Source{
					Deployment:   "mydeployment",
					Target:       "director.example.com",
					Client:       "foo",
					ClientSecret: "foobar",
				},
				Params: concourse.OutParams{
					Manifest: "path/to/manifest.yml",
					DryRun:   true,
				},
			}))
		})
	})

	Context("when bosh_io_stemcell_type is set", func() {
		It("set BoshIOStemcellType in OutParams", func() {
			config := []byte(`{
				"params": {
					"manifest": "path/to/manifest.yml",
					"bosh_io_stemcell_type": "regular"
				},
				"source": {
					"deployment": "mydeployment",
					"target": "director.example.com",
					"client": "foo",
					"client_secret": "foobar"
				}
			}`)

			source, err := concourse.NewOutRequest(config, "")
			Expect(err).NotTo(HaveOccurred())

			Expect(source).To(Equal(concourse.OutRequest{
				Source: concourse.Source{
					Deployment:   "mydeployment",
					Target:       "director.example.com",
					Client:       "foo",
					ClientSecret: "foobar",
				},
				Params: concourse.OutParams{
					Manifest:           "path/to/manifest.yml",
					BoshIOStemcellType: "regular",
				},
			}))
		})

		It("raises an error when invalid type is passed", func() {
			config := []byte(`{
				"params": {
					"manifest": "path/to/manifest.yml",
					"bosh_io_stemcell_type": "non-existing-type"
				},
				"source": {
					"deployment": "mydeployment",
					"target": "director.example.com",
					"client": "foo",
					"client_secret": "foobar"
				}
			}`)

			_, err := concourse.NewOutRequest(config, "")
			Expect(err).To(MatchError(
				"bosh_io_stemcell_type only supports 'light' or 'regular' got: non-existing-type"))
		})
	})

	Context("when source_file param is passed", func() {
		It("overrides source with the values in the source_file", func() {
			sourceFile, _ := os.CreateTemp("", "") //nolint:errcheck
			sourceFile.WriteString(                //nolint:errcheck
				`{
				"deployment": "fileDeployment",
				"target": "fileDirector.com",
				"client_secret": "fileSecret",
				"vars_store": {
					"provider": "fileProvider",
					"config": {
						"file": "vars"
					}
				}
			}`)
			sourceFile.Close() //nolint:errcheck

			configTemplate := `{
				"params": {
					"manifest": "path/to/manifest.yml",
					"source_file": "%s"
				},
				"source": {
					"deployment": "mydeployment",
					"target": "director.example.com",
					"client": "original_client",
					"client_secret": "foobar",
					"vars_store": {
						"provider": "gcs",
						"config": {
							"some": "dynamic",
							"keys": "per-provider"
						}
					}
				}
			}`
			config := []byte(fmt.Sprintf(
				configTemplate,
				filepath.Base(sourceFile.Name()),
			))

			source, err := concourse.NewOutRequest(config, filepath.Dir(sourceFile.Name()))
			Expect(err).NotTo(HaveOccurred())

			Expect(source).To(Equal(concourse.OutRequest{
				Source: concourse.Source{
					Deployment:   "fileDeployment",
					Target:       "fileDirector.com",
					Client:       "original_client",
					ClientSecret: "fileSecret",
					VarsStore: concourse.VarsStore{
						Provider: "fileProvider",
						Config: map[string]interface{}{
							"file": "vars",
							"some": "dynamic",
							"keys": "per-provider",
						},
					},
				},
				Params: concourse.OutParams{
					Manifest: "path/to/manifest.yml",
				},
			}))
		})
	})

	Context("when decoding fails", func() {
		It("errors", func() {
			config := []byte("not-json")

			_, err := concourse.NewOutRequest(config, "")
			Expect(err).To(HaveOccurred())
		})
	})

	Context("when delete is specified", func() {
		It("does not require the manifest parameter", func() {
			config := []byte(`{
				"source": {
					"deployment": "mydeployment",
					"target": "director.example.com",
					"client": "foo",
					"client_secret": "foobar"
				},
				"params": {
					"delete": {
						"enabled": true
					}
				}
			}`)

			_, err := concourse.NewOutRequest(config, "")
			Expect(err).NotTo(HaveOccurred())
		})
	})

	Context("when a required parameter is missing", func() {
		It("returns an error with each missing parameter", func() {
			config := []byte(`{
				"source": {
					"deployment": "mydeployment",
					"target": "director.example.com",
					"client": "foo",
					"client_secret": "foobar"
				}
			}`)

			_, err := concourse.NewOutRequest(config, "")
			Expect(err).To(HaveOccurred())

			Expect(err.Error()).To(ContainSubstring("manifest"))
		})
	})
})
