package bosh_test

import (
	"bytes"
	"io"
	"io/ioutil"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/cloudfoundry/bosh-deployment-resource/bosh"
	"github.com/cloudfoundry/bosh-deployment-resource/bosh/boshfakes"
	"github.com/cloudfoundry/bosh-deployment-resource/concourse"

	"errors"
	boshcmd "github.com/cloudfoundry/bosh-cli/cmd"
	boshtpl "github.com/cloudfoundry/bosh-cli/director/template"
)

var _ = Describe("BoshDirector", func() {
	var (
		director      bosh.BoshDirector
		out           io.Writer
		commandRunner *boshfakes.FakeRunner
		sillyBytes    = []byte{0xFE, 0xED, 0xDE, 0xAD, 0xBE, 0xEF}
	)

	BeforeEach(func() {
		commandRunner = new(boshfakes.FakeRunner)
		out = bytes.NewBufferString("")
		director = bosh.NewBoshDirector(concourse.Source{}, commandRunner, out)
	})

	Describe("Deploy", func() {
		It("tells BOSH to deploy the given manifest and parameters", func() {
			vars := map[string]interface{}{"foo": "bar"}
			varKVs := []boshtpl.VarKV{
				{
					Name: "foo",
					Value: "bar",
				},
			}
			varFileContents := properYaml(`
				baz: "best-bar"
			`)
			varFile, _ := ioutil.TempFile("", "var-file-1")
			varFile.Write(varFileContents)

			noRedact := true
			err := director.Deploy(sillyBytes, bosh.DeployParams{
				NoRedact: noRedact,
				Vars: vars,
				VarsFiles: []string{varFile.Name()},
			})
			Expect(err).ToNot(HaveOccurred())

			Expect(commandRunner.ExecuteCallCount()).To(Equal(1))

			deployOpts := commandRunner.ExecuteArgsForCall(0).(*boshcmd.DeployOpts)
			Expect(deployOpts.Args.Manifest.Bytes).To(Equal(sillyBytes))
			Expect(deployOpts.NoRedact).To(Equal(noRedact))
			Expect(deployOpts.VarKVs).To(Equal(varKVs))
			Expect(len(deployOpts.VarsFiles)).To(Equal(1))
			Expect(deployOpts.VarsFiles[0].Vars).To(Equal(boshtpl.StaticVariables{
				"baz": "best-bar",
			}))
		})

		Context("when deploying fails", func() {
			It("returns an error", func() {
				commandRunner.ExecuteReturns(errors.New("Your deploy failed"))

				err := director.Deploy(sillyBytes, bosh.DeployParams{})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Your deploy failed"))
			})
		})

		Context("when cleanup is specified", func() {
			It("runs a cleanup before the deploy", func() {
				err := director.Deploy(sillyBytes, bosh.DeployParams{Cleanup: true})
				Expect(err).ToNot(HaveOccurred())

				Expect(commandRunner.ExecuteCallCount()).To(Equal(2))

				cleanUpOpts := commandRunner.ExecuteArgsForCall(0).(*boshcmd.CleanUpOpts)
				Expect(cleanUpOpts.All).To(Equal(false))

				deployOpts := commandRunner.ExecuteArgsForCall(1).(*boshcmd.DeployOpts)
				Expect(deployOpts.Args.Manifest.Bytes).To(Equal(sillyBytes))
			})
		})
	})

	Describe("DownloadManifest", func() {
		It("gets the deployment manifest", func() {
			commandRunner.GetResultReturns(sillyBytes, nil)

			manifestBytes, err := director.DownloadManifest()
			Expect(err).ToNot(HaveOccurred())

			Expect(manifestBytes).To(Equal(sillyBytes))
		})

		Context("when getting the manifest fails", func() {
			It("returns an error", func() {
				commandRunner.GetResultReturns(nil, errors.New("Your manifest is missing"))

				_, err := director.DownloadManifest()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Your manifest is missing"))
			})
		})
	})

	Describe("UploadRelease", func() {
		It("uploads the given release", func() {
			err := director.UploadRelease("my-cool-release")
			Expect(err).ToNot(HaveOccurred())

			Expect(commandRunner.ExecuteCallCount()).To(Equal(1))

			uploadReleaseOpts := commandRunner.ExecuteArgsForCall(0).(*boshcmd.UploadReleaseOpts)
			Expect(string(uploadReleaseOpts.Args.URL)).To(Equal("my-cool-release"))
		})

		Context("when uploading the release fails", func() {
			It("returns an error", func() {
				commandRunner.ExecuteReturns(errors.New("failed communicating with director"))

				err := director.UploadRelease("my-cool-release")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Could not upload release my-cool-release: failed communicating with director"))
			})
		})
	})

	Describe("UploadStemcell", func() {
		It("uploads the given stemcell", func() {
			err := director.UploadStemcell("my-cool-stemcell")
			Expect(err).ToNot(HaveOccurred())

			Expect(commandRunner.ExecuteCallCount()).To(Equal(1))

			uploadStemcellOpts := commandRunner.ExecuteArgsForCall(0).(*boshcmd.UploadStemcellOpts)
			Expect(string(uploadStemcellOpts.Args.URL)).To(Equal("my-cool-stemcell"))
		})

		Context("when uploading the stemcell fails", func() {
			It("returns an error", func() {
				commandRunner.ExecuteReturns(errors.New("failed communicating with director"))

				err := director.UploadStemcell("my-cool-stemcell")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Could not upload stemcell my-cool-stemcell: failed communicating with director"))
			})
		})
	})
})
