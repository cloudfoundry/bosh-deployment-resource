package tools_test

import (
	"compress/gzip"
	"fmt"
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/cloudfoundry/bosh-deployment-resource/tools"
)

var _ = Describe("ReadTgzFile", func() {
	It("returns the contents of a file in a gzip tar archive", func() {
		contents, err := tools.ReadTgzFile("fixtures/small-release.tgz", "release.MF")
		Expect(err).ToNot(HaveOccurred())

		Expect(string(contents)).To(Equal("---\nname: small-release\nversion: 53"))
	})

	Describe("when the archive does not exist", func() {
		It("returns an error", func() {
			_, err := tools.ReadTgzFile("fixtures/does-not-exist.nope", "release.MF")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Could not read archive fixtures/does-not-exist.nope"))
		})
	})

	Describe("when the archive is not a valid gzip", func() {
		It("returns an error", func() {
			notArchive, _ := os.CreateTemp("", "release-one") //nolint:errcheck
			notArchive.Close()                                //nolint:errcheck

			_, err := tools.ReadTgzFile(notArchive.Name(), "release.MF")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(
				ContainSubstring(fmt.Sprintf("%s is not a valid gzip archive", notArchive.Name())),
			)
		})
	})

	Describe("when the gzip archive does not contain a valid tar", func() {
		It("returns an error", func() {
			notArchive, _ := os.CreateTemp("", "release-one") //nolint:errcheck
			gzipWriter := gzip.NewWriter(notArchive)
			gzipWriter.Write([]byte("hello")) //nolint:errcheck
			gzipWriter.Close()                //nolint:errcheck
			notArchive.Close()                //nolint:errcheck

			_, err := tools.ReadTgzFile(notArchive.Name(), "release.MF")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(
				ContainSubstring(fmt.Sprintf("%s is not a valid tar", notArchive.Name())),
			)
		})
	})

	Describe("when file is not in the archive", func() {
		It("returns an error", func() {
			_, err := tools.ReadTgzFile("fixtures/small-release.tgz", "not-a-file.nope")
			Expect(err).To(HaveOccurred())

			Expect(err.Error()).To(
				ContainSubstring("fixtures/small-release.tgz does not contain file not-a-file.nope"),
			)
		})
	})
})
