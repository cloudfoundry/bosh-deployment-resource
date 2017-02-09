package out_test

import (
	"fmt"
	"os"
	"io/ioutil"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/cloudfoundry/bosh-deployment-resource/bosh"
	"github.com/cloudfoundry/bosh-deployment-resource/bosh/boshfakes"
	"github.com/cloudfoundry/bosh-deployment-resource/concourse"
	"github.com/cloudfoundry/bosh-deployment-resource/out"
)

var _ = Describe("OutCommand", func() {
	var (
		outCommand out.OutCommand
		director   *boshfakes.FakeDirector
	)

	BeforeEach(func() {
		director = new(boshfakes.FakeDirector)
		outCommand = out.NewOutCommand(director)
	})

	Describe("Run", func() {
		var outRequest concourse.OutRequest

		BeforeEach(func() {
			outRequest = concourse.OutRequest{
				Source: concourse.Source{
					Target: "director.example.com",
				},
				Params: concourse.OutParams{
					Manifest: "path/to/manifest.yml",
					NoRedact: true,
				},
			}
		})

		It("deploys", func() {
			_, err := outCommand.Run(outRequest)
			Expect(err).ToNot(HaveOccurred())

			Expect(director.DeployCallCount()).To(Equal(1))
			actualManifestPath, actualDeployParams := director.DeployArgsForCall(0)
			Expect(actualManifestPath).To(Equal("path/to/manifest.yml"))
			Expect(actualDeployParams).To(Equal(bosh.DeployParams{NoRedact: true}))
		})

		It("returns the new version", func() {
			sillyBytes := []byte{0xFE, 0xED, 0xDE, 0xAD, 0xBE, 0xEF}
			director.DownloadManifestReturns(sillyBytes, nil)

			outResponse, err := outCommand.Run(outRequest)
			Expect(err).ToNot(HaveOccurred())

			Expect(outResponse).To(Equal(out.OutResponse{
				Version: concourse.Version{
					ManifestSha1: "33bf00cb7a45258748f833a47230124fcc8fa3a4",
					Target:       "director.example.com",
				},
			}))
		})

		Context("when releases are provided", func() {
			var (
				releaseOne, releaseTwo, releaseThree *os.File
			)

			BeforeEach(func() {
				primaryReleaseDir, _ := ioutil.TempDir("", "")

				releaseOne, _ = ioutil.TempFile(primaryReleaseDir, "release-one")
				releaseOne.Close()

				releaseTwo, _ = ioutil.TempFile(primaryReleaseDir, "release-two")
				releaseTwo.Close()

				secondaryReleaseDir, _ := ioutil.TempDir("", "")

				releaseThree, _ = ioutil.TempFile(secondaryReleaseDir, "release-three")
				releaseThree.Close()

				outRequest.Params.Releases = []string{
					fmt.Sprintf("%s/release-*", primaryReleaseDir),
					releaseThree.Name(),
				}
			})

			It("uploads all of the releases", func() {
				_, err := outCommand.Run(outRequest)
				Expect(err).ToNot(HaveOccurred())

				Expect(director.UploadReleaseCallCount()).To(Equal(3))

				uploadedReleases := []string{}
				uploadedReleases = append(uploadedReleases,
					director.UploadReleaseArgsForCall(0),
					director.UploadReleaseArgsForCall(1),
					director.UploadReleaseArgsForCall(2),
				)
				Expect(uploadedReleases).To(ContainElement(releaseOne.Name()))
				Expect(uploadedReleases).To(ContainElement(releaseTwo.Name()))
				Expect(uploadedReleases).To(ContainElement(releaseThree.Name()))
			})

			Context("when a release glob is bad", func() {
				It("gives a useful error", func() {
					outRequest.Params.Releases = []string{"/["}
					_, err := outCommand.Run(outRequest)

					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring("Invalid release name: /["))
				})
			})
		})

		Context("when stemcells are provided", func() {
			var (
				stemcellOne, stemcellTwo, stemcellThree *os.File
			)

			BeforeEach(func() {
				primaryStemcellDir, _ := ioutil.TempDir("", "")

				stemcellOne, _ = ioutil.TempFile(primaryStemcellDir, "stemcell-one")
				stemcellOne.Close()

				stemcellTwo, _ = ioutil.TempFile(primaryStemcellDir, "stemcell-two")
				stemcellTwo.Close()

				secondaryStemcellDir, _ := ioutil.TempDir("", "")

				stemcellThree, _ = ioutil.TempFile(secondaryStemcellDir, "stemcell-three")
				stemcellThree.Close()

				outRequest.Params.Stemcells = []string{
					fmt.Sprintf("%s/stemcell-*", primaryStemcellDir),
					stemcellThree.Name(),
				}
			})

			It("uploads all of the stemcells", func() {
				_, err := outCommand.Run(outRequest)
				Expect(err).ToNot(HaveOccurred())

				Expect(director.UploadStemcellCallCount()).To(Equal(3))

				uploadedStemcells := []string{}
				uploadedStemcells = append(uploadedStemcells,
					director.UploadStemcellArgsForCall(0),
					director.UploadStemcellArgsForCall(1),
					director.UploadStemcellArgsForCall(2),
				)
				Expect(uploadedStemcells).To(ContainElement(stemcellOne.Name()))
				Expect(uploadedStemcells).To(ContainElement(stemcellTwo.Name()))
				Expect(uploadedStemcells).To(ContainElement(stemcellThree.Name()))
			})

			Context("when a stemcell glob is bad", func() {
				It("gives a useful error", func() {
					outRequest.Params.Stemcells = []string{"/["}
					_, err := outCommand.Run(outRequest)

					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring("Invalid stemcell name: /["))
				})
			})
		})
	})
})
