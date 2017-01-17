package check_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/cloudfoundry/bosh-deployment-resource/check"
	"github.com/cloudfoundry/bosh-deployment-resource/bosh/boshfakes"
	"github.com/cloudfoundry/bosh-deployment-resource/concourse"
	"errors"
)

var _ = Describe("CheckCommand", func() {
	var (
		checkCommand check.CheckCommand
		director *boshfakes.FakeDirector
	)

	BeforeEach(func() {
		director = new(boshfakes.FakeDirector)
		checkCommand = check.NewCheckCommand(director)
	})

	Describe("Run", func() {
		var checkRequest concourse.CheckRequest

		BeforeEach(func() {
			manifestContents := []byte{0xFE, 0xED, 0xDE, 0xAD, 0xBE, 0xEF}
			director.DownloadManifestReturns(manifestContents, nil)
		})

		Context("When the manifest sha is a mismatch with the version provided", func() {
			BeforeEach(func() {
				checkRequest = concourse.CheckRequest{
					Source: concourse.Source{
						Target: "director.example.com",
					},
					Version: concourse.Version{},
				}
			})
			It("returns the sha1 of the manifest", func() {
				checkResponse, err := checkCommand.Run(checkRequest)
				Expect(err).ToNot(HaveOccurred())

				Expect(director.DownloadManifestCallCount()).To(Equal(1))
				Expect(checkResponse).To(Equal([]concourse.Version{
					{
						ManifestSha1: "33bf00cb7a45258748f833a47230124fcc8fa3a4",
						Target: "director.example.com",
					},
				}))
			})
		})

		Context("When the manifest sha matches the version provided to the check", func() {
			BeforeEach(func() {
				checkRequest = concourse.CheckRequest{
					Source: concourse.Source{
						Target: "director.example.com",
					},
					Version: concourse.Version{
						ManifestSha1: "33bf00cb7a45258748f833a47230124fcc8fa3a4",
						Target: "director.example.com",
					},
				}
			})
			It("an empty versions array", func() {
				checkResponse, err := checkCommand.Run(checkRequest)
				Expect(err).ToNot(HaveOccurred())

				Expect(director.DownloadManifestCallCount()).To(Equal(1))
				Expect(checkResponse).To(Equal([]concourse.Version{}))
			})
		})
		Context("When the there is an error downloading the manifest", func() {
			BeforeEach(func() {
				director.DownloadManifestReturns([]byte{}, errors.New("No manifest for you"))
			})

			It("returns the error", func() {
				_, err := checkCommand.Run(checkRequest)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("No manifest for you"))
			})
		})
	})
})