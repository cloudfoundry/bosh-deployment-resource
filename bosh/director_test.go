package bosh_test

import (
	"bytes"
	"io"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/cloudfoundry/bosh-deployment-resource/bosh"
	"github.com/cloudfoundry/bosh-deployment-resource/bosh/boshfakes"
	"github.com/cloudfoundry/bosh-deployment-resource/concourse"

	"errors"
	boshcmd "github.com/cloudfoundry/bosh-cli/cmd"
	"io/ioutil"
	"path/filepath"
)

var _ = Describe("BoshDirector", func() {
	var (
		director      bosh.BoshDirector
		out           io.Writer
		commandRunner *boshfakes.FakeRunner
		tempDir       string
		sillyBytes    = []byte{0xFE, 0xED, 0xDE, 0xAD, 0xBE, 0xEF}
	)

	BeforeEach(func() {
		commandRunner = new(boshfakes.FakeRunner)
		out = bytes.NewBufferString("")
		tempDir, _ = ioutil.TempDir("", "")
		director = bosh.NewBoshDirector(concourse.Source{}, commandRunner, tempDir, out)
	})

	Describe("Deploy", func() {
		var (
			manifestPath string
		)

		BeforeEach(func() {
			manifest, _ := ioutil.TempFile(tempDir, "")
			manifest.Write(sillyBytes)
			manifest.Close()

			manifestPath = filepath.Base(manifest.Name())
		})

		It("tells BOSH to deploy the given manifest and parameters", func() {
			noRedact := true
			err := director.Deploy(manifestPath, noRedact)
			Expect(err).ToNot(HaveOccurred())

			Expect(commandRunner.ExecuteCallCount()).To(Equal(1))

			deployOpts := commandRunner.ExecuteArgsForCall(0).(*boshcmd.DeployOpts)
			Expect(deployOpts.Args.Manifest.Bytes).To(Equal(sillyBytes))
			Expect(deployOpts.NoRedact).To(Equal(noRedact))
		})

		Context("when deploying fails", func() {
			It("returns an error", func() {
				commandRunner.ExecuteReturns(errors.New("Your deploy failed"))

				err := director.Deploy(manifestPath, false)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Your deploy failed"))
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
})
