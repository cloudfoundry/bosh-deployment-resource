package bosh_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/cloudfoundry/bosh-deployment-resource/bosh"
)

var _ = Describe("Stemcell", func() {
	Describe("NewStemcell", func() {
		It("parses the stemcell tar", func() {
			stemcell, err := bosh.NewStemcell("fixtures/small-stemcell.tgz")
			Expect(err).ToNot(HaveOccurred())

			Expect(stemcell).To(Equal(bosh.Stemcell{
				Name: "small-stemcell",
				OperatingSystem: "ubuntu-trusty",
				Version: "8675309",
			}))
		})

		Context("when the stemcell doesn't exist", func() {
			It("returns an error", func() {
				_, err := bosh.NewStemcell("fixtures/not-a-file.tgz")
				Expect(err).To(HaveOccurred())

				Expect(err.Error()).To(ContainSubstring("Could not read stemcell fixtures/not-a-file.tgz"))
			})
		})

		Context("when the stemcell is malformed exist", func() {
			It("returns an error", func() {
				_, err := bosh.NewStemcell("fixtures/malformed-stemcell.tgz")
				Expect(err).To(HaveOccurred())

				Expect(err.Error()).To(ContainSubstring("Stemcell fixtures/malformed-stemcell.tgz is not a valid stemcell"))
			})
		})
	})
})