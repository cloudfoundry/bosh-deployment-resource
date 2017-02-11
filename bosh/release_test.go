package bosh_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/cloudfoundry/bosh-deployment-resource/bosh"
)

var _ = Describe("Release", func() {
	Describe("NewRelease", func() {
		It("parses the release tar", func() {
			release, err := bosh.NewRelease("fixtures/small-release.tgz")
			Expect(err).ToNot(HaveOccurred())

			Expect(release).To(Equal(bosh.Release{
				Name: "small-release",
				Version: "53",
			}))
		})

		Context("when the release doesn't exist", func() {
			It("returns an error", func() {
				_, err := bosh.NewRelease("fixtures/not-a-file.tgz")
				Expect(err).To(HaveOccurred())

				Expect(err.Error()).To(ContainSubstring("Could not read release fixtures/not-a-file.tgz"))
			})
		})

		Context("when the release is malformed exist", func() {
			It("returns an error", func() {
				_, err := bosh.NewRelease("fixtures/malformed-release.tgz")
				Expect(err).To(HaveOccurred())

				Expect(err.Error()).To(ContainSubstring("Release fixtures/malformed-release.tgz is not a valid release"))
			})
		})
	})
})