package in_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"crypto/sha1"
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/cloudfoundry/bosh-deployment-resource/bosh/boshfakes"
	"github.com/cloudfoundry/bosh-deployment-resource/concourse"
	"github.com/cloudfoundry/bosh-deployment-resource/in"
)

var _ = Describe("InCommand", func() {
	var (
		inCommand in.InCommand
		director  *boshfakes.FakeDirector
	)

	BeforeEach(func() {
		director = new(boshfakes.FakeDirector)
		inCommand = in.NewInCommand(director)
	})

	Describe("Run", func() {
		var inRequest concourse.InRequest
		var targetDir string
		sillyBytes := []byte{0xFE, 0xED, 0xDE, 0xAD, 0xBE, 0xEF}
		sillyBytesSha1 := "33bf00cb7a45258748f833a47230124fcc8fa3a4"
		wrongBytes := []byte{0x0F, 0xFE, 0xEF, 0xBE, 0xCF, 0xF0}

		BeforeEach(func() {
			inRequest = concourse.InRequest{
				Source: concourse.Source{
					Target: "director.example.com",
				},
				Version: concourse.Version{
					ManifestSha1: sillyBytesSha1,
					Target:       "director.example.com",
				},
			}

			var err error
			targetDir, err = ioutil.TempDir("", "")
			Expect(err).ToNot(HaveOccurred())

			director.DownloadManifestReturns(sillyBytes, nil)
		})

		It("writes the manifest and target to disk and returns the version as a response", func() {
			inResponse, err := inCommand.Run(inRequest, targetDir)
			Expect(err).ToNot(HaveOccurred())

			manifestBytes, err := ioutil.ReadFile(filepath.Join(targetDir, "manifest.yml"))
			Expect(err).ToNot(HaveOccurred())
			Expect(manifestBytes).To(Equal(sillyBytes))
			Expect(director.DownloadManifestCallCount()).To(Equal(1))

			targetBytes, err := ioutil.ReadFile(filepath.Join(targetDir, "target"))
			Expect(err).ToNot(HaveOccurred())
			Expect(string(targetBytes)).To(Equal("director.example.com"))

			Expect(inResponse).To(Equal(in.InResponse{
				Version: concourse.Version{
					ManifestSha1: sillyBytesSha1,
					Target:       "director.example.com",
				},
			}))
		})

		Context("when no target is provided", func() {
			BeforeEach(func() {
				inRequest.Source.Target = ""
			})

			It("no-ops assuming this is an implicit get after using a dynamic source file", func() {
				inResponse, err := inCommand.Run(inRequest, targetDir)
				Expect(err).ToNot(HaveOccurred())

				Expect(director.DownloadManifestCallCount()).To(Equal(0))
				_, err = ioutil.ReadFile(filepath.Join(targetDir, "manifest.yml"))
				Expect(err).To(HaveOccurred())

				_, err = ioutil.ReadFile(filepath.Join(targetDir, "target"))
				Expect(err).To(HaveOccurred())

				Expect(inResponse).To(Equal(in.InResponse{
					Version: inRequest.Version,
				}))
			})
		})

		Context("when the manifest download fails", func() {
			BeforeEach(func() {
				director.DownloadManifestReturns(nil, errors.New("could not download manifest"))
			})

			It("returns an error", func() {
				_, err := inCommand.Run(inRequest, targetDir)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("could not download manifest"))
			})
		})

		Context("when downloaded manifest does not match the requested version", func() {
			BeforeEach(func() {
				director.DownloadManifestReturns(wrongBytes, nil)
			})

			It("returns an error", func() {
				_, err := inCommand.Run(inRequest, targetDir)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("Requested deployment version is not available"))
			})
		})

		Context("when director target does not match the requested version", func() {
			BeforeEach(func() {
				inRequest.Source.Target = "weird.example.com"
			})

			It("returns an error", func() {
				_, err := inCommand.Run(inRequest, targetDir)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("Requested deployment director is different than configured source"))
			})
		})

		Context("when requesting compiled_releases", func() {
			manifest := properYaml(`
				releases:
				- name: real-one
				  version: "1"
				- name: real-two
				  version: "2.2"
				stemcells:
				- alias: default
				  os: ubuntu-trusty
				  version: "3309.8"
			`)

			BeforeEach(func() {
				director.DownloadManifestReturns(manifest, nil)
				inRequest.Version.ManifestSha1 = fmt.Sprintf("%x", sha1.Sum(manifest))
				inRequest.Params.CompiledReleases = []concourse.CompiledRelease{
					{Name: "real-one"},
					{Name: "real-two"},
				}
			})

			It("downloads each release with the version specified in the manifest", func() {
				_, err := inCommand.Run(inRequest, targetDir)
				Expect(err).ToNot(HaveOccurred())

				Expect(director.ExportReleasesCallCount()).To(Equal(1))

				targetDir, releases := director.ExportReleasesArgsForCall(0)
				Expect(targetDir).To(Equal(targetDir))
				Expect(releases).To(Equal([]string{"real-one", "real-two"}))
			})

			Context("when exporting releases fails", func() {
				It("errors", func() {
					director.ExportReleasesReturns(errors.New("could not export"))
					_, err := inCommand.Run(inRequest, targetDir)
					Expect(err).To(MatchError(ContainSubstring("could not export")))
				})
			})
		})
	})
})
