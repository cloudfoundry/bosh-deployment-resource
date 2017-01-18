package concourse_test

import (
	"fmt"
	"github.com/cloudfoundry/bosh-deployment-resource/concourse"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io/ioutil"
	"strings"
)

var _ = Describe("NewOutRequest", func() {
	It("converts the reader into an OutRequest", func() {
		reader := strings.NewReader(`{
  "params": {
    "manifest": "/tmp/manifest.yml"
  },
  "source": {
    "deployment": "mydeployment",
    "target": "director.example.com"
  }
}`)

		outRequest, err := concourse.NewOutRequest(reader)
		Expect(err).NotTo(HaveOccurred())

		Expect(outRequest).To(Equal(concourse.OutRequest{
			Params: concourse.OutParams{
				Manifest: "/tmp/manifest.yml",
			},
			Source: concourse.Source{
				Deployment: "mydeployment",
				Target:     "director.example.com",
			},
		}))
	})

	Context("when the reader has a target_file", func() {
		var (
			targetFilePath      string
			requestJsonTemplate string = `{
  "params": {
    "manifest": "/tmp/manifest.yml",
    "target_file": "%s"
  },
  "source": {
    "deployment": "mydeployment",
    "target": "director.example.com"
  }
}`
		)

		BeforeEach(func() {
			targetFile, _ := ioutil.TempFile("", "")
			targetFile.WriteString("director.example.net")
			targetFile.Close()

			targetFilePath = targetFile.Name()
		})

		It("uses the contents of that file instead of the target parameter", func() {
			reader := strings.NewReader(fmt.Sprintf(requestJsonTemplate, targetFilePath))

			outRequest, err := concourse.NewOutRequest(reader)
			Expect(err).NotTo(HaveOccurred())

			Expect(outRequest.Source.Target).To(Equal("director.example.net"))
		})

		Context("when the targetfile cannot be read", func() {
			BeforeEach(func() {
				targetFilePath = "not-a-real-file"
			})

			It("errors", func() {
				reader := strings.NewReader(fmt.Sprintf(requestJsonTemplate, targetFilePath))

				_, err := concourse.NewOutRequest(reader)
				Expect(err).To(HaveOccurred())
			})
		})
	})

	Context("when decoding fails", func() {
		It("errors", func() {
			reader := strings.NewReader(`not-json`)

			_, err := concourse.NewOutRequest(reader)
			Expect(err).To(HaveOccurred())
		})
	})
})
