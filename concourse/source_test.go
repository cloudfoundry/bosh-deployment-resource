package concourse_test

import (
	"fmt"
	"github.com/cloudfoundry/bosh-deployment-resource/concourse"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io/ioutil"
	"path/filepath"
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

	Context("when the config has a target_file", func() {
		var (
			sourcesDir          string
			targetFileName      string
			requestJsonTemplate string = `{
				"params": {
					"target_file": "%s"
				},
				"source": {
					"deployment": "mydeployment",
					"target": "director.example.com",
					"client": "foo",
					"client_secret": "foobar"
				}
			}`
		)

		BeforeEach(func() {
			targetFile, _ := ioutil.TempFile("", "")
			targetFile.WriteString("director.example.net")
			targetFile.Close()

			sourcesDir = filepath.Dir(targetFile.Name())
			targetFileName = filepath.Base(targetFile.Name())
		})

		It("uses the contents of that file instead of the target parameter", func() {
			config := []byte(fmt.Sprintf(requestJsonTemplate, targetFileName))

			source, err := concourse.NewDynamicSource(config, sourcesDir)
			Expect(err).NotTo(HaveOccurred())

			Expect(source.Target).To(Equal("director.example.net"))
		})

		Context("when the target_file cannot be read", func() {
			BeforeEach(func() {
				targetFileName = "not-a-real-file"
			})

			It("errors", func() {
				config := []byte(fmt.Sprintf(requestJsonTemplate, targetFileName))

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
