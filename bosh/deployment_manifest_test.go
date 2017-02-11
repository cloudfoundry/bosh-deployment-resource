package bosh_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/cloudfoundry/bosh-deployment-resource/bosh"
	"strings"
)

var _ = Describe("DeploymentManifest", func() {
	Describe("NewDeploymentManifest", func() {
		It("returns an error if parsing invalid yaml", func() {
			_, err := bosh.NewDeploymentManifest([]byte("&&&"))
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Failed to unmarshal manifest: "))
		})
	})

	Describe("UseRelease", func() {
		It("updates the request release version to match the provided release", func() {
			d, _ := bosh.NewDeploymentManifest(properYaml(`
				releases:
				- name: cool-release
			`))

			err := d.UseReleaseVersion("cool-release", "6")
			Expect(err).ToNot(HaveOccurred())

			Expect(d.Manifest()).To(MatchYAML(properYaml(`
				releases:
				- name: cool-release
				  version: "6"
			`)))
		})

		Context("when the release is not found", func() {
			It("returns an error", func() {
				d, _ := bosh.NewDeploymentManifest(properYaml(`
					releases:
					- cool-release:
						version: 5
				`))

				err := d.UseReleaseVersion("unknown-release", "6")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Release unknown-release not defined"))
			})
		})

		Context("when there is no releases section", func() {
			It("returns an error", func() {
				d, _ := bosh.NewDeploymentManifest([]byte(`
					jobs:
					- my_job: 5
				`))

				err := d.UseReleaseVersion("unknown-release", "6")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("No releases section"))
			})
		})
	})
})

func properYaml(improperYaml string) []byte {
	return []byte(strings.Replace(improperYaml, "\t", "  ", -1))
}