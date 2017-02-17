package out_test

import (
	"fmt"
	"os"
	"io"
	"bytes"
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
		outCommand   out.OutCommand
		director     *boshfakes.FakeDirector
		manifest     *os.File
		manifestYaml []byte
	)

	BeforeEach(func() {
		director = new(boshfakes.FakeDirector)
		outCommand = out.NewOutCommand(director, "")
		manifest, _ = ioutil.TempFile("", "manifest")
		manifestYaml = properYaml(`
			releases:
			- name: small-release
			  version: latest
			  url: file://release.tgz
			  sha1: SHA1FORMAT
			stemcells:
			- name: small-stemcell
			  alias: super-awesome-stemcell
			  version: latest
		`)
		manifest.Write(manifestYaml)
		manifest.Close()
	})

	Describe("Run", func() {
		var outRequest concourse.OutRequest

		BeforeEach(func() {
			outRequest = concourse.OutRequest{
				Source: concourse.Source{
					Target: "director.example.com",
				},
				Params: concourse.OutParams{
					Manifest: manifest.Name(),
					NoRedact: true,
					Vars: map[string]interface{}{
						"foo": "bar",
					},
				},
			}
		})

		It("deploys", func() {
			_, err := outCommand.Run(outRequest)
			Expect(err).ToNot(HaveOccurred())

			Expect(director.DeployCallCount()).To(Equal(1))
			actualManifestYaml, actualDeployParams := director.DeployArgsForCall(0)
			Expect(actualManifestYaml).To(MatchYAML(manifestYaml))
			Expect(actualDeployParams).To(Equal(bosh.DeployParams{
				NoRedact: true,
				VarsFiles: []string{},
				Vars: map[string]interface{}{
					"foo": "bar",
				},
			}))
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

		Context("when varFiles are provided", func() {
			var (
				varFileOne, varFileTwo, varFileThree *os.File
				varFiles []string
			)

			BeforeEach(func() {
				// Update varFile generation to yield expected bosh varFile format
				primaryVarFileDir, _ := ioutil.TempDir("", "")

				varFileOne, _ = ioutil.TempFile(primaryVarFileDir, "varFile-one")
				varFileOne.Close()

				varFileTwo, _ = ioutil.TempFile(primaryVarFileDir, "varFile-two")
				varFileTwo.Close()

				secondaryVarFileDir, _ := ioutil.TempDir("", "")

				varFileThree, _ = ioutil.TempFile(secondaryVarFileDir, "varFile-three")
				varFileThree.Close()

				varFiles =[]string{
					varFileThree.Name(),
					fmt.Sprintf("%s/varFile-*", primaryVarFileDir),
				}
				outRequest.Params.VarsFiles = varFiles
			})

			It("specifies the varFiles in the deploy command", func() {
				_, err := outCommand.Run(outRequest)
				Expect(err).ToNot(HaveOccurred())

				_, actualDeployParams := director.DeployArgsForCall(0)
				Expect(actualDeployParams.VarsFiles).To(ConsistOf(
					varFileThree.Name(),
					varFileOne.Name(),
					varFileTwo.Name(),
				))
			})

			Context("when a varFile glob is bad", func() {
				It("gives a useful error", func() {
					outRequest.Params.VarsFiles = []string{"/["}
					_, err := outCommand.Run(outRequest)

					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring("Invalid var_file name: /["))
				})
			})
		})
		
		Context("when releases are provided", func() {
			var (
				releaseOne, releaseTwo, releaseThree *os.File
			)

			BeforeEach(func() {
				// Update release generation to yield expected bosh release format
				primaryReleaseDir, _ := ioutil.TempDir("", "")

				smallRelease, _ := ioutil.ReadFile("fixtures/small-release.tgz")

				releaseOne, _ = ioutil.TempFile(primaryReleaseDir, "release-one")
				io.Copy(releaseOne, bytes.NewReader(smallRelease))
				releaseOne.Close()

				releaseTwo, _ = ioutil.TempFile(primaryReleaseDir, "release-two")
				io.Copy(releaseTwo, bytes.NewReader(smallRelease))
				releaseTwo.Close()

				secondaryReleaseDir, _ := ioutil.TempDir("", "")

				releaseThree, _ = ioutil.TempFile(secondaryReleaseDir, "release-three")
				io.Copy(releaseThree, bytes.NewReader(smallRelease))
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

			It("updates the version information in the manifest", func() {
				outRequest.Params.Releases = []string{"fixtures/small-release.tgz"}
				_, err := outCommand.Run(outRequest)
				Expect(err).ToNot(HaveOccurred())

				updatedManifest, _ := director.DeployArgsForCall(0)

				Expect(updatedManifest).To(MatchYAML(properYaml(`
					releases:
						- name: small-release
						  version: "53"
						  url: file://release.tgz
						  sha1: SHA1FORMAT
					stemcells:
						- name: small-stemcell
						  alias: super-awesome-stemcell
						  version: latest
				`)))
			})

			It("includes the provided releases in the metadata", func() {
				outResponse, err := outCommand.Run(outRequest)
				Expect(err).ToNot(HaveOccurred())

				Expect(outResponse.Metadata).To(Equal([]concourse.Metadata{
					{
						Name: "release",
						Value: "small-release v53",
					},
					{
						Name: "release",
						Value: "small-release v53",
					},
					{
						Name: "release",
						Value: "small-release v53",
					},
				}))
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

				smallStemcell, _ := ioutil.ReadFile("fixtures/small-stemcell.tgz")

				stemcellOne, _ = ioutil.TempFile(primaryStemcellDir, "stemcell-one")
				io.Copy(stemcellOne, bytes.NewReader(smallStemcell))
				stemcellOne.Close()

				stemcellTwo, _ = ioutil.TempFile(primaryStemcellDir, "stemcell-two")
				io.Copy(stemcellTwo, bytes.NewReader(smallStemcell))
				stemcellTwo.Close()

				secondaryStemcellDir, _ := ioutil.TempDir("", "")

				stemcellThree, _ = ioutil.TempFile(secondaryStemcellDir, "stemcell-three")
				io.Copy(stemcellThree, bytes.NewReader(smallStemcell))
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

			It("updates the version information in the manifest", func() {
				outRequest.Params.Stemcells = []string{"fixtures/small-stemcell.tgz"}
				_, err := outCommand.Run(outRequest)
				Expect(err).ToNot(HaveOccurred())

				updatedManifest, _ := director.DeployArgsForCall(0)

				Expect(updatedManifest).To(MatchYAML(properYaml(`
					releases:
						- name: small-release
						  version: latest
						  url: file://release.tgz
						  sha1: SHA1FORMAT
					stemcells:
						- alias: super-awesome-stemcell
						  name: small-stemcell
						  version: "8675309"
				`)))
			})

			It("includes the provided stemcells in the metadata", func() {
				outResponse, err := outCommand.Run(outRequest)
				Expect(err).ToNot(HaveOccurred())

				Expect(outResponse.Metadata).To(Equal([]concourse.Metadata{
					{
						Name: "stemcell",
						Value: "small-stemcell v8675309",
					},
					{
						Name: "stemcell",
						Value: "small-stemcell v8675309",
					},
					{
						Name: "stemcell",
						Value: "small-stemcell v8675309",
					},
				}))
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
