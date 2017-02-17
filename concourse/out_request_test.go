package concourse_test

import (
	"io/ioutil"
	"fmt"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/cloudfoundry/bosh-deployment-resource/concourse"
)

var _ = Describe("NewOutRequest", func() {
	It("converts the config into an OutRequest", func() {
		targetFile, _ := ioutil.TempFile("", "")
		targetFile.WriteString("director.example.net")
		targetFile.Close()

		configTemplate := `{
			"params": {
				"target_file": "%s",
				"manifest": "path/to/manifest.yml",
				"vars": {
					"foo": "bar",
					"slice": [1, "two"]
				},
				"vars_files": [
					"path/to/file",
					"second/path/to/file"
				],
				"ops_files": [
					"ops-file1",
					"path/to/ops-file2"
				]
			},
			"source": {
				"deployment": "mydeployment",
				"target": "director.example.com",
				"client": "foo",
				"client_secret": "foobar"
			}
		}`
		config := []byte(fmt.Sprintf(configTemplate, filepath.Base(targetFile.Name())))

		source, err := concourse.NewOutRequest(config, filepath.Dir(targetFile.Name()))
		Expect(err).NotTo(HaveOccurred())

		Expect(source).To(Equal(concourse.OutRequest{
			Source: concourse.Source{
				Deployment: "mydeployment",
				Target: "director.example.net",
				Client: "foo",
				ClientSecret: "foobar",
			},
			Params: concourse.OutParams {
				Manifest: "path/to/manifest.yml",
				Vars: map[string]interface{} {
					"foo": "bar",
					"slice": []interface{}{float64(1), "two"},
				},
				VarsFiles: []string {
					"path/to/file",
					"second/path/to/file",
				},
				OpsFiles: []string {
					"ops-file1",
					"path/to/ops-file2",
				},
			},
		}))
	})

	Context("when decoding fails", func() {
		It("errors", func() {
			config := []byte("not-json")

			_, err := concourse.NewOutRequest(config, "")
			Expect(err).To(HaveOccurred())
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
