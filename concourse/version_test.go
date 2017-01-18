package concourse_test

import (
	"github.com/cloudfoundry/bosh-deployment-resource/concourse"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Version", func() {
	sillyBytes := []byte{0xFE, 0xED, 0xDE, 0xAD, 0xBE, 0xEF}
	sillyBytesSha1 := "33bf00cb7a45258748f833a47230124fcc8fa3a4"

	It("presents the SHA1 as a string", func() {
		Expect(concourse.NewVersion(sillyBytes, "director.example.com")).To(Equal(concourse.Version{
			ManifestSha1: sillyBytesSha1,
			Target:       "director.example.com",
		}))
	})
})
