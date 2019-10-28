package out_test

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"errors"

	boshdir "github.com/cloudfoundry/bosh-cli/director"
	"github.com/cloudfoundry/bosh-deployment-resource/bosh"
	"github.com/cloudfoundry/bosh-deployment-resource/bosh/boshfakes"
	"github.com/cloudfoundry/bosh-deployment-resource/concourse"
	"github.com/cloudfoundry/bosh-deployment-resource/out"
	"github.com/cloudfoundry/bosh-deployment-resource/storage/storagefakes"
)

var _ = Describe("OutCommand", func() {
	var (
		outCommand   out.OutCommand
		director     *boshfakes.FakeDirector
		boshIOClient *boshfakes.FakeBoshIO
		resourcesDir string
		manifestYaml []byte
	)

	BeforeEach(func() {
		director = new(boshfakes.FakeDirector)
		boshIOClient = new(boshfakes.FakeBoshIO)
		resourcesDir, _ = ioutil.TempDir("", "resources-dir")
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
		Expect(ioutil.WriteFile(filepath.Join(resourcesDir, "manifest"), manifestYaml, 0600)).To(Succeed())
		director.InterpolateReturns(manifestYaml, nil)
		outCommand = out.NewOutCommand(director, boshIOClient, nil, resourcesDir)
	})

	AfterEach(func() {
		Expect(os.RemoveAll(resourcesDir)).To(Succeed())
	})

	Describe("Run", func() {
		var outRequest concourse.OutRequest

		BeforeEach(func() {
			outRequest = concourse.OutRequest{
				Source: concourse.Source{
					Target: "director.example.com",
				},
				Params: concourse.OutParams{
					Manifest: "manifest",
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
			Expect(director.DeleteCallCount()).To(Equal(0))

			_, actualInterpolateParams := director.InterpolateArgsForCall(0)
			Expect(actualInterpolateParams.Vars).To(Equal(
				map[string]interface{}{
					"foo": "bar",
				},
			))

			Expect(director.DeployCallCount()).To(Equal(1))
			actualManifestYaml, actualDeployParams := director.DeployArgsForCall(0)
			Expect(actualManifestYaml).To(MatchYAML(manifestYaml))
			Expect(actualDeployParams).To(Equal(bosh.DeployParams{
				NoRedact: true,
				VarFiles: map[string]string{},
			}))
		})

		It("deploys with recreate", func() {
			outRequest.Params.Recreate = true

			_, err := outCommand.Run(outRequest)
			Expect(err).ToNot(HaveOccurred())

			_, actualInterpolateParams := director.InterpolateArgsForCall(0)
			Expect(actualInterpolateParams.Vars).To(Equal(
				map[string]interface{}{
					"foo": "bar",
				},
			))

			Expect(director.DeployCallCount()).To(Equal(1))
			actualManifestYaml, actualDeployParams := director.DeployArgsForCall(0)
			Expect(actualManifestYaml).To(MatchYAML(manifestYaml))
			Expect(actualDeployParams).To(Equal(bosh.DeployParams{
				NoRedact: true,
				Recreate: true,
				VarFiles: map[string]string{},
			}))
		})

		It("deploys with fix", func() {
			outRequest.Params.Fix = true

			_, err := outCommand.Run(outRequest)
			Expect(err).ToNot(HaveOccurred())

			_, actualInterpolateParams := director.InterpolateArgsForCall(0)
			Expect(actualInterpolateParams.Vars).To(Equal(
				map[string]interface{}{
					"foo": "bar",
				},
			))

			Expect(director.DeployCallCount()).To(Equal(1))
			actualManifestYaml, actualDeployParams := director.DeployArgsForCall(0)
			Expect(actualManifestYaml).To(MatchYAML(manifestYaml))
			Expect(actualDeployParams).To(Equal(bosh.DeployParams{
				NoRedact: true,
				Fix:      true,
				VarFiles: map[string]string{},
			}))
		})

		It("dryrun deploys", func() {
			outRequest.Params.DryRun = true

			_, err := outCommand.Run(outRequest)
			Expect(err).ToNot(HaveOccurred())

			_, actualInterpolateParams := director.InterpolateArgsForCall(0)
			Expect(actualInterpolateParams.Vars).To(Equal(
				map[string]interface{}{
					"foo": "bar",
				},
			))

			Expect(director.DeployCallCount()).To(Equal(1))
			actualManifestYaml, actualDeployParams := director.DeployArgsForCall(0)
			Expect(actualManifestYaml).To(MatchYAML(manifestYaml))
			Expect(actualDeployParams).To(Equal(bosh.DeployParams{
				NoRedact: true,
				DryRun:   true,
				VarFiles: map[string]string{},
			}))
		})

		It("deploys with max in flight", func() {
			outRequest.Params.MaxInFlight = 5

			_, err := outCommand.Run(outRequest)
			Expect(err).ToNot(HaveOccurred())

			_, actualInterpolateParams := director.InterpolateArgsForCall(0)
			Expect(actualInterpolateParams.Vars).To(Equal(
				map[string]interface{}{
					"foo": "bar",
				},
			))

			Expect(director.DeployCallCount()).To(Equal(1))
			actualManifestYaml, actualDeployParams := director.DeployArgsForCall(0)
			Expect(actualManifestYaml).To(MatchYAML(manifestYaml))
			Expect(actualDeployParams).To(Equal(bosh.DeployParams{
				NoRedact:    true,
				MaxInFlight: 5,
				VarFiles:    map[string]string{},
			}))
		})

		It("deploys with skip drain", func() {
			outRequest.Params.SkipDrain = []string{"all"}

			_, err := outCommand.Run(outRequest)
			Expect(err).ToNot(HaveOccurred())

			_, actualInterpolateParams := director.InterpolateArgsForCall(0)
			Expect(actualInterpolateParams.Vars).To(Equal(
				map[string]interface{}{
					"foo": "bar",
				},
			))

			Expect(director.DeployCallCount()).To(Equal(1))
			actualManifestYaml, actualDeployParams := director.DeployArgsForCall(0)
			Expect(actualManifestYaml).To(MatchYAML(manifestYaml))
			Expect(actualDeployParams).To(Equal(bosh.DeployParams{
				NoRedact:  true,
				SkipDrain: []string{"all"},
				VarFiles:  map[string]string{},
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
				Metadata: []concourse.Metadata{},
			}))
		})

		It("waits for locks on the deployment", func() {
			_, err := outCommand.Run(outRequest)
			Expect(err).ToNot(HaveOccurred())
			Expect(director.WaitForDeployLockCallCount()).To(Equal(1))
		})

		Context("when varsFiles are provided", func() {
			var (
				varsFileOne, varsFileTwo, varsFileThree *os.File
				varsFiles                               []string
			)

			BeforeEach(func() {
				// Update varFile generation to yield expected bosh varFile format
				primaryVarFileDir, _ := ioutil.TempDir(resourcesDir, "")

				varsFileOne, _ = ioutil.TempFile(primaryVarFileDir, "varFile-one")
				varsFileOne.Close()

				varsFileTwo, _ = ioutil.TempFile(primaryVarFileDir, "varFile-two")
				varsFileTwo.Close()

				varsFileThree, _ = ioutil.TempFile(resourcesDir, "varFile-three")
				varsFileThree.Close()

				primaryDirWithoutResourcesDir, _ := filepath.Rel(resourcesDir, fmt.Sprintf("%s/varFile-*", primaryVarFileDir))
				varsFiles = []string{
					filepath.Base(varsFileThree.Name()),
					primaryDirWithoutResourcesDir,
				}
				outRequest.Params.VarsFiles = varsFiles
			})

			It("interpolates the varsFiles into the manifest but does not delpoy with them", func() {
				_, err := outCommand.Run(outRequest)
				Expect(err).ToNot(HaveOccurred())

				_, actualInterpolateParams := director.InterpolateArgsForCall(0)
				Expect(actualInterpolateParams.VarsFiles).To(ConsistOf(
					varsFileThree.Name(),
					varsFileOne.Name(),
					varsFileTwo.Name(),
				))

				_, actualDeployParams := director.DeployArgsForCall(0)
				Expect(actualDeployParams.VarsFiles).To(BeEmpty())
			})

			Context("when a varsFile glob is bad", func() {
				It("gives a useful error", func() {
					outRequest.Params.VarsFiles = []string{"/["}
					_, err := outCommand.Run(outRequest)

					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring("Invalid vars_file name: /["))
				})
			})
		})

		Context("when varFiles are provided", func() {
			It("prepends the paths with the resources directory and passes them to deploy params", func() {
				varFiles := map[string]string{"awesome-var": "relative/path/to/file-with/value"}
				outRequest.Params.VarFiles = varFiles

				_, err := outCommand.Run(outRequest)
				Expect(err).ToNot(HaveOccurred())

				_, actualDeployParams := director.DeployArgsForCall(0)
				Expect(actualDeployParams.VarFiles).To(Equal(
					map[string]string{"awesome-var": filepath.Join(resourcesDir, "relative/path/to/file-with/value")},
				))
			})
		})

		Context("when opsFiles are provided", func() {
			var (
				opsFileOne, opsFileTwo, opsFileThree *os.File
				opsFiles                             []string
			)

			BeforeEach(func() {
				// Update opsFile generation to yield expected bosh opsFile format
				primaryopsFileDir, _ := ioutil.TempDir(resourcesDir, "")

				opsFileOne, _ = ioutil.TempFile(primaryopsFileDir, "opsFile-one")
				opsFileOne.Close()

				opsFileTwo, _ = ioutil.TempFile(primaryopsFileDir, "opsFile-two")
				opsFileTwo.Close()

				opsFileThree, _ = ioutil.TempFile(resourcesDir, "opsFile-three")
				opsFileThree.Close()

				primaryDirWithoutResourcesDir, _ := filepath.Rel(resourcesDir, fmt.Sprintf("%s/opsFile-*", primaryopsFileDir))
				opsFiles = []string{
					filepath.Base(opsFileThree.Name()),
					primaryDirWithoutResourcesDir,
				}
				outRequest.Params.OpsFiles = opsFiles
			})

			It("interpolates the opsfiles into the manifest but does not delpoy with them", func() {
				_, err := outCommand.Run(outRequest)
				Expect(err).ToNot(HaveOccurred())

				_, actualInterpolateParams := director.InterpolateArgsForCall(0)
				Expect(actualInterpolateParams.OpsFiles).To(ConsistOf(
					opsFileThree.Name(),
					opsFileOne.Name(),
					opsFileTwo.Name(),
				))

				_, actualDeployParams := director.DeployArgsForCall(0)
				Expect(actualDeployParams.OpsFiles).To(BeEmpty())
			})

			Context("when a opsFile glob is bad", func() {
				It("gives a useful error", func() {
					outRequest.Params.OpsFiles = []string{"/["}
					_, err := outCommand.Run(outRequest)

					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring("Invalid ops_file name: /["))
				})
			})

		})

		Context("when releases are provided", func() {
			var (
				releaseOne, releaseTwo, releaseThree *os.File
			)

			BeforeEach(func() {
				// Update release generation to yield expected bosh release format
				primaryReleaseDir, _ := ioutil.TempDir(resourcesDir, "")

				smallRelease, _ := ioutil.ReadFile("fixtures/small-release.tgz")

				releaseOne, _ = ioutil.TempFile(primaryReleaseDir, "release-one")
				io.Copy(releaseOne, bytes.NewReader(smallRelease))
				releaseOne.Close()

				releaseTwo, _ = ioutil.TempFile(primaryReleaseDir, "release-two")
				io.Copy(releaseTwo, bytes.NewReader(smallRelease))
				releaseTwo.Close()

				releaseThree, _ = ioutil.TempFile(resourcesDir, "release-three")
				io.Copy(releaseThree, bytes.NewReader(smallRelease))
				releaseThree.Close()

				primaryDirWithoutResourcesDir, _ := filepath.Rel(resourcesDir, fmt.Sprintf("%s/release-*", primaryReleaseDir))
				outRequest.Params.Releases = []string{
					primaryDirWithoutResourcesDir,
					filepath.Base(releaseThree.Name()),
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
						Name:  "release",
						Value: "small-release v53",
					},
					{
						Name:  "release",
						Value: "small-release v53",
					},
					{
						Name:  "release",
						Value: "small-release v53",
					},
				}))
			})
		})

		Context("when stemcells are provided", func() {
			var (
				stemcellOne, stemcellTwo, stemcellThree *os.File
				interpolatedManifest                    []byte
			)

			BeforeEach(func() {
				primaryStemcellDir, _ := ioutil.TempDir(resourcesDir, "")

				smallStemcell, _ := ioutil.ReadFile("fixtures/small-stemcell.tgz")

				stemcellOne, _ = ioutil.TempFile(primaryStemcellDir, "stemcell-one")
				io.Copy(stemcellOne, bytes.NewReader(smallStemcell))
				stemcellOne.Close()

				stemcellTwo, _ = ioutil.TempFile(primaryStemcellDir, "stemcell-two")
				io.Copy(stemcellTwo, bytes.NewReader(smallStemcell))
				stemcellTwo.Close()

				stemcellThree, _ = ioutil.TempFile(resourcesDir, "stemcell-three")
				io.Copy(stemcellThree, bytes.NewReader(smallStemcell))
				stemcellThree.Close()

				primaryDirWithoutResourcesDir, _ := filepath.Rel(resourcesDir, fmt.Sprintf("%s/stemcell-*", primaryStemcellDir))
				outRequest.Params.Stemcells = []string{
					primaryDirWithoutResourcesDir,
					filepath.Base(stemcellThree.Name()),
				}

				interpolatedManifest = properYaml(`
					releases:
						- name: small-release
						  version: latest
						  url: file://release.tgz
						  sha1: SHA1FORMAT
					stemcells:
						- alias: super-awesome-stemcell
						  name: small-stemcell
						  version: "8675309"
				`)
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
				_, err := outCommand.Run(outRequest)
				Expect(err).ToNot(HaveOccurred())

				updatedManifest, _ := director.DeployArgsForCall(0)

				Expect(updatedManifest).To(MatchYAML(interpolatedManifest))
			})

			It("includes the provided stemcells in the metadata", func() {
				director.InterpolateReturns(interpolatedManifest, nil)
				outResponse, err := outCommand.Run(outRequest)
				Expect(err).ToNot(HaveOccurred())

				Expect(outResponse.Metadata).To(Equal([]concourse.Metadata{
					{
						Name:  "stemcell",
						Value: "small-stemcell v8675309",
					},
					{
						Name:  "stemcell",
						Value: "small-stemcell v8675309",
					},
					{
						Name:  "stemcell",
						Value: "small-stemcell v8675309",
					},
				}))
			})
		})

		Context("when bosh_io_stemcell_type is provided", func() {
			var (
				interpolatedManifest []byte
			)

			BeforeEach(func() {
				interpolatedManifest = properYaml(`
					stemcells:
					- alias: default
					  os: ubuntu-xenial
					  version: "456.40"
				`)

				stemcells := []byte(`[{
                                  "name":"bosh-google-kvm-ubuntu-xenial-go_agent",
                                  "version":"456.40",
                                  "regular":{
                                    "url":"https://example.com/bosh-stemcell-456.40-google-kvm-ubuntu-xenial-go_agent.tgz",
                                    "sha1":"e3fe3b2fa7e5f0111bfc0b22f30eb5658eba89c5"
                                  }}]`)
				outRequest.Params.BoshIOStemcellType = "regular"

				director.InterpolateReturns(interpolatedManifest, nil)
				director.InfoReturns(boshdir.Info{CPI: "google_cpi"}, nil)
				boshIOClient.StemcellsReturns(stemcells, nil)
			})

			It("uploads all of the stemcells", func() {
				_, err := outCommand.Run(outRequest)
				Expect(err).ToNot(HaveOccurred())

				Expect(director.UploadRemoteStemcellCallCount()).To(Equal(1))

				url, name, version, sha := director.UploadRemoteStemcellArgsForCall(0)
				Expect(url).To(HaveSuffix("bosh-stemcell-456.40-google-kvm-ubuntu-xenial-go_agent.tgz"))
				Expect(name).To(Equal("bosh-google-kvm-ubuntu-xenial-go_agent"))
				Expect(version).To(Equal("456.40"))
				Expect(sha).To(Equal("e3fe3b2fa7e5f0111bfc0b22f30eb5658eba89c5"))
			})

			It("includes the provided stemcells in the metadata", func() {
				outResponse, err := outCommand.Run(outRequest)
				Expect(err).ToNot(HaveOccurred())

				Expect(outResponse.Metadata).To(Equal([]concourse.Metadata{
					{
						Name:  "stemcell",
						Value: "bosh-google-kvm-ubuntu-xenial-go_agent v456.40",
					},
				}))
			})
		})

		Context("when a vars store config is provided", func() {
			var (
				fakeStorageClient *storagefakes.FakeStorageClient
			)

			It("downloads the vars store, uses it, and uploads it", func() {
				director = new(boshfakes.FakeDirector)
				fakeStorageClient = new(storagefakes.FakeStorageClient)
				outCommand = out.NewOutCommand(director, boshIOClient, fakeStorageClient, resourcesDir)
				_, err := outCommand.Run(outRequest)
				Expect(err).ToNot(HaveOccurred())

				Expect(fakeStorageClient.DownloadCallCount()).To(Equal(1))
				filePath := fakeStorageClient.DownloadArgsForCall(0)

				Expect(fakeStorageClient.UploadCallCount()).To(Equal(1))
				Expect(fakeStorageClient.UploadArgsForCall(0)).To(Equal(filePath))

				_, actualDeployParams := director.DeployArgsForCall(0)
				Expect(actualDeployParams.VarsStore).To(Equal(filePath))
			})

			Describe("when the download fails", func() {
				It("returns an error", func() {
					director = new(boshfakes.FakeDirector)
					fakeStorageClient = new(storagefakes.FakeStorageClient)
					fakeStorageClient.DownloadReturns(errors.New("Failed to download"))

					outCommand = out.NewOutCommand(director, boshIOClient, fakeStorageClient, resourcesDir)
					_, err := outCommand.Run(outRequest)
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(Equal("Failed to download"))
				})
			})

			Describe("when the upload fails", func() {
				It("returns an error", func() {
					director = new(boshfakes.FakeDirector)
					fakeStorageClient = new(storagefakes.FakeStorageClient)
					fakeStorageClient.UploadReturns(errors.New("Failed to upload"))

					outCommand = out.NewOutCommand(director, boshIOClient, fakeStorageClient, resourcesDir)
					_, err := outCommand.Run(outRequest)
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(Equal("Failed to upload"))
				})
			})
		})

		Context("when the requested operation is a delete", func() {
			BeforeEach(func() {
				outRequest.Params = concourse.OutParams{
					Delete: concourse.DeleteParams{
						Enabled: true,
						Force:   true,
					},
				}
			})

			It("deletes the deployment", func() {
				response, err := outCommand.Run(outRequest)
				Expect(err).ToNot(HaveOccurred())

				Expect(director.DeployCallCount()).To(Equal(0))
				Expect(director.DeleteCallCount()).To(Equal(1))
				Expect(director.DeleteArgsForCall(0)).To(Equal(true))

				Expect(response).To(Equal(out.OutResponse{}))
			})

			Context("when the delete errors", func() {
				BeforeEach(func() {
					director.DeleteReturns(fmt.Errorf("Delete failed!"))
				})

				It("returns an error", func() {
					response, err := outCommand.Run(outRequest)

					Expect(err).To(HaveOccurred())
					Expect(response).To(Equal(out.OutResponse{}))
				})
			})

		})
	})
})
