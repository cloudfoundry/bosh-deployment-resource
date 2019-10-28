package bosh_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/cloudfoundry/bosh-deployment-resource/bosh"
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
		It("updates the requested release version to match the provided release", func() {
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

	Describe("UseStemcell", func() {
		It("updates the requested stemcell version to match the provided stemcell", func() {
			d, _ := bosh.NewDeploymentManifest(properYaml(`
				stemcells:
				- name: bosh-best-iaas-light-stemcell
				  version: latest
				- os: super-virtual-os
				  version: latest
			`))

			err := d.UseStemcellVersion("bosh-best-iaas-light-stemcell", "lame-os", "9002")
			Expect(err).ToNot(HaveOccurred())
			err = d.UseStemcellVersion("bosh-ok-iaas-heavy-stemcell", "super-virtual-os", "1002")
			Expect(err).ToNot(HaveOccurred())

			Expect(d.Manifest()).To(MatchYAML(properYaml(`
				stemcells:
				- name: bosh-best-iaas-light-stemcell
				  version: "9002"
				- os: super-virtual-os
				  version: "1002"
			`)))
		})

		It("does not update stemcells when the version is not latest", func() {
			d, _ := bosh.NewDeploymentManifest(properYaml(`
				stemcells:
				- name: bosh-best-iaas-light-stemcell
				  version: 1
			`))

			err := d.UseStemcellVersion("bosh-best-iaas-light-stemcell", "lame-os", "9002")
			Expect(err).ToNot(HaveOccurred())

			Expect(d.Manifest()).To(MatchYAML(properYaml(`
				stemcells:
				- name: bosh-best-iaas-light-stemcell
				  version: 1
			`)))
		})

		Context("when the stemcell is not found", func() {
			It("returns an error", func() {
				d, _ := bosh.NewDeploymentManifest(properYaml(`
					stemcells:
					- name: bosh-best-iaas-light-stemcell
					  version: 1
				`))

				err := d.UseStemcellVersion("bosh-unknown-light-stemcell", "lame-os", "9002")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Stemcell bosh-unknown-light-stemcell not defined in deployment manifest"))
			})

			Context("when there is no stemcells section", func() {
				It("returns an error", func() {
					d, _ := bosh.NewDeploymentManifest([]byte(`
						jobs:
						- my_job: 5
					`))

					err := d.UseStemcellVersion("uknown-stemcell", "lame-os", "6")
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring("No stemcells section"))
				})
			})
		})

		Context("when more than one stemcell matches", func() {
			It("returns an error", func() {
				d, _ := bosh.NewDeploymentManifest(properYaml(`
					stemcells:
					- name: bosh-best-iaas-light-stemcell
					  version: latest
					- name: bosh-best-iaas-light-stemcell
					  version: latest
				`))

				err := d.UseStemcellVersion("bosh-best-iaas-light-stemcell", "os", "6")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Multiple matches for stemcell bosh-best-iaas-light-stemcell"))
			})
		})
	})
	Describe("Stemcells", func() {
		It("returns version and os for defined stemcells", func() {
			d, err := bosh.NewDeploymentManifest(properYaml(`
				stemcells:
				- name: bosh-best-iaas-light-stemcell
				  os: ubuntu-xenial
				  version: 1.2.3
				- os: ubuntu-trusty
				  version: latest
			`))
			Expect(err).ToNot(HaveOccurred())

			stemcells, err := d.Stemcells()
			Expect(err).ToNot(HaveOccurred())
			Expect(stemcells[0].Version).To(Equal("1.2.3"))
			Expect(stemcells[0].OperatingSystem).To(Equal("ubuntu-xenial"))
			Expect(stemcells[1].Version).To(Equal("latest"))
			Expect(stemcells[1].OperatingSystem).To(Equal("ubuntu-trusty"))
		})
	})
})
