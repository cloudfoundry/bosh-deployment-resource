package bosh_test

import (
	"github.com/cloudfoundry/bosh-deployment-resource/bosh"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Release", func() {
	Describe("NewRelease", func() {
		It("parses the release tar", func() {
			release, err := bosh.NewReleases("fixtures", []string{"small-release.tgz"})
			Expect(err).ToNot(HaveOccurred())

			Expect(release).To(Equal([]bosh.Release{
				{
					Name:     "small-release",
					Version:  "53",
					FilePath: "fixtures/small-release.tgz",
				},
			}))
		})

		Context("when the tgz is not a release", func() {
			It("returns an error", func() {
				_, err := bosh.NewReleases("fixtures", []string{"small-stemcell.tgz"})
				Expect(err).To(HaveOccurred())

				Expect(err.Error()).To(ContainSubstring("Could not read release:"))
			})
		})

		Context("when the release is malformed", func() {
			It("returns an error", func() {
				_, err := bosh.NewReleases("fixtures", []string{"malformed-release.tgz"})
				Expect(err).To(HaveOccurred())

				Expect(err.Error()).To(ContainSubstring("Release fixtures/malformed-release.tgz is not a valid release"))
			})
		})

		Context("when a release glob is bad", func() {
			It("gives a useful error", func() {
				_, err := bosh.NewReleases("fixtures", []string{"/["})
				Expect(err).To(HaveOccurred())

				Expect(err.Error()).To(ContainSubstring("Invalid release name: /["))
			})
		})
	})
})
