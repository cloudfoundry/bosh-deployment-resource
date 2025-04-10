package tools_test

import (
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/cloudfoundry/bosh-deployment-resource/tools"
)

var _ = Describe("GlobUnfurler", func() {
	var (
		releaseOne, releaseTwo, releaseThree *os.File
		releaseDir                           string
	)

	BeforeEach(func() {
		releaseDir, _ = os.MkdirTemp("", "primary-releases") //nolint:errcheck

		releaseOne, _ = os.CreateTemp(releaseDir, "release-one") //nolint:errcheck
		releaseOne.Close()                                       //nolint:errcheck

		releaseTwo, _ = os.CreateTemp(releaseDir, "release-two") //nolint:errcheck
		releaseTwo.Close()                                       //nolint:errcheck

		releaseThree, _ = os.CreateTemp(releaseDir, "coolio-three") //nolint:errcheck
		releaseThree.Close()                                        //nolint:errcheck
	})

	It("returns all filepaths matching the globs", func() {
		filepaths, err := tools.UnfurlGlobs(
			releaseDir, []string{
				"release-*",
				"coolio-three*",
			},
		)

		Expect(err).ToNot(HaveOccurred())
		Expect(filepaths).To(ConsistOf(
			releaseOne.Name(),
			releaseTwo.Name(),
			releaseThree.Name(),
		))
	})

	It("returns all filepaths in order", func() {
		filepaths, err := tools.UnfurlGlobs(
			releaseDir, []string{
				"release-two*",
				"coolio-three*",
				"release-one*",
			},
		)

		Expect(err).ToNot(HaveOccurred())
		Expect(filepaths).To(Equal([]string{
			releaseTwo.Name(),
			releaseThree.Name(),
			releaseOne.Name(),
		}))
	})

	Context("when some globs unfurl to the same file", func() {
		It("removes duplicate filepaths", func() {
			filepaths, err := tools.UnfurlGlobs(
				releaseDir, []string{
					"release-*",
					"rel*se-*",
				},
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
			_, err := tools.UnfurlGlobs(releaseDir, []string{"/["})
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("/[ is not a valid file glob"))
		})
	})

	Context("when a glob matches no files", func() {
		It("returns an error", func() {
			_, err := tools.UnfurlGlobs(releaseDir, []string{"zzzzz"})
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("zzzzz does not match any files"))
		})
	})
})
