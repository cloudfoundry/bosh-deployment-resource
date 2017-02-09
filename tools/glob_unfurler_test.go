package tools_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"os"
	"io/ioutil"
	"github.com/cloudfoundry/bosh-deployment-resource/tools"
	"fmt"
)

var _ = Describe("GlobUnfurler", func() {
	var (
		releaseOne, releaseTwo, releaseThree *os.File
		primaryReleaseDir string
	)

	BeforeEach(func() {
		primaryReleaseDir, _ = ioutil.TempDir("", "primary-releases")

		releaseOne, _ = ioutil.TempFile(primaryReleaseDir, "release-one")
		releaseOne.Close()

		releaseTwo, _ = ioutil.TempFile(primaryReleaseDir, "release-two")
		releaseTwo.Close()

		secondaryReleaseDir, _ := ioutil.TempDir("", "secondary-releases")

		releaseThree, _ = ioutil.TempFile(secondaryReleaseDir, "release-three")
		releaseThree.Close()
	})

	It("returns all filepaths matching the globs", func() {
		filepaths, err := tools.UnfurlGlobs(
			fmt.Sprintf("%s/release-*", primaryReleaseDir),
			releaseThree.Name(),
		)

		Expect(err).ToNot(HaveOccurred())
		Expect(filepaths).To(ConsistOf(
			releaseOne.Name(),
			releaseTwo.Name(),
			releaseThree.Name(),
		))
	})

	Context("when some globs unfurl to the same file", func() {
		It("removes duplicate filepaths", func() {
			filepaths, err := tools.UnfurlGlobs(
				fmt.Sprintf("%s/release-*", primaryReleaseDir),
				fmt.Sprintf("%s/rel*se-*", primaryReleaseDir),
			)

			Expect(err).ToNot(HaveOccurred())
			Expect(filepaths).To(ConsistOf(
				releaseOne.Name(),
				releaseTwo.Name(),
			))
		})
	})

	Context("when a bad glob is passed", func() {
		It("returns an error", func() {
			_, err := tools.UnfurlGlobs("/[")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("/[ is not a valid file glob"))
		})
	})
})