package bosh_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/cloudfoundry/bosh-deployment-resource/bosh"
	"github.com/cloudfoundry/bosh-deployment-resource/bosh/boshfakes"
)

var _ = Describe("LookupBoshIOStemcell", func() {
	var (
		fakeBoshIOClient *boshfakes.FakeBoshIO
		cpi, os, version string
		light            bool
		stemcells        []byte
		stemcell         BoshIOStemcell
		err              error
	)

	BeforeEach(func() {
		fakeBoshIOClient = new(boshfakes.FakeBoshIO)
		version = "456.40"
		cpi = "google_cpi"
		os = "ubuntu-xenial"
		light = true
		stemcells = []byte(`[{
                                  "name":"bosh-google-kvm-ubuntu-xenial-go_agent",
                                  "version":"456.40",
                                  "regular":{
                                    "url":"https://example.com/bosh-stemcell-456.40-google-kvm-ubuntu-xenial-go_agent.tgz",
                                    "sha1":"regular-sha"
                                  },
                                  "light":{
                                    "url":"https://example.com/light-bosh-stemcell-456.40-google-kvm-ubuntu-xenial-go_agent.tgz",
                                    "sha1":"light-sha"
                                  }}]`)

	})

	JustBeforeEach(func() {
		fakeBoshIOClient.StemcellsReturns(stemcells, nil)
		stemcell, err = LookupBoshIOStemcell(
			fakeBoshIOClient, cpi, os, version, light)
	})

	Context("when using latests stemcell", func() {
		BeforeEach(func() {
			version = "latest"
		})

		It("raises en error", func() {
			Expect(err).To(
				MatchError("Auto upload of \"latest\" stemcell is not support, please use bosh-io-stemcell-resource"))
		})
	})

	Context("when stemcell version not found", func() {
		BeforeEach(func() {
			version = "non-existing-version"
		})

		It("raises en error", func() {
			Expect(err).To(
				MatchError("did not find a suitable stemcell with version: non-existing-version"))
		})
	})

	Context("when using light stemcells", func() {
		It("returns stemcell", func() {
			Expect(err).ToNot(HaveOccurred())
			Expect(stemcell).To(Equal(BoshIOStemcell{
				Name:    "bosh-google-kvm-ubuntu-xenial-go_agent",
				Version: "456.40",
				URL:     "https://example.com/light-bosh-stemcell-456.40-google-kvm-ubuntu-xenial-go_agent.tgz",
				Sha1:    "light-sha",
			}))
		})
	})

	Context("when using regular stemcells", func() {
		BeforeEach(func() {
			light = false
		})
		It("returns stemcell", func() {
			Expect(err).ToNot(HaveOccurred())
			Expect(stemcell).To(Equal(BoshIOStemcell{
				Name:    "bosh-google-kvm-ubuntu-xenial-go_agent",
				Version: "456.40",
				URL:     "https://example.com/bosh-stemcell-456.40-google-kvm-ubuntu-xenial-go_agent.tgz",
				Sha1:    "regular-sha",
			}))
		})
	})
})
