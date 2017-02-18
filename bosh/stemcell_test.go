package bosh_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/cloudfoundry/bosh-deployment-resource/bosh"
)

var _ = Describe("Stemcell", func() {
	Describe("NewStemcells", func() {
		It("parses the stemcell tar", func() {
			stemcell, err := bosh.NewStemcells("fixtures", []string{"small-stemcell.tgz"})
			Expect(err).ToNot(HaveOccurred())

			Expect(stemcell).To(Equal([]bosh.Stemcell{
				{
					Name: "small-stemcell",
					OperatingSystem: "ubuntu-trusty",
					Version: "8675309",
					FilePath: "fixtures/small-stemcell.tgz",
				},
			}))
		})

		Context("when the tgz is not a stemcell", func() {
			It("returns an error", func() {
				_, err := bosh.NewStemcells("fixtures", []string{"small-release.tgz"})
				Expect(err).To(HaveOccurred())

				Expect(err.Error()).To(ContainSubstring("Could not read stemcell:"))
			})
		})

		Context("when the stemcell is malformed", func() {
			It("returns an error", func() {
				_, err := bosh.NewStemcells("fixtures", []string{"malformed-stemcell.tgz"})
				Expect(err).To(HaveOccurred())

				Expect(err.Error()).To(ContainSubstring("Stemcell fixtures/malformed-stemcell.tgz is not a valid stemcell"))
			})
		})

		Context("when a stemcell glob is bad", func() {
			It("gives a useful error", func() {
				_, err := bosh.NewStemcells("fixtures", []string{"/["})
				Expect(err).To(HaveOccurred())

				Expect(err.Error()).To(ContainSubstring("Invalid stemcell name: /["))
			})
		})
	})
})