package bosh_test

import (
	"errors"

	"github.com/cloudfoundry/bosh-deployment-resource/bosh"
	"github.com/cloudfoundry/bosh-deployment-resource/bosh/boshfakes"
	"github.com/cloudfoundry/bosh-deployment-resource/concourse"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("CLI coordinator", func() {
	var (
		cliCoordinator bosh.CLICoordinator
		fakeProxy      *boshfakes.FakeProxy
		source         concourse.Source
	)

	BeforeEach(func() {
		fakeProxy = &boshfakes.FakeProxy{}
		fakeProxy.AddrReturns("some-proxy-addr", nil)
		fakeProxy.AddrReturnsOnCall(0, "", errors.New("proxy is not running"))
		source = concourse.Source{JumpboxSSHKey: "some-key", JumpboxURL: "some-url"}
		cliCoordinator = bosh.NewCLICoordinator(source, GinkgoWriter, fakeProxy)
	})

	Describe("StartProxy", func() {
		It("starts a proxy server and returns the proxy address", func() {
			addr, err := cliCoordinator.StartProxy()
			Expect(err).NotTo(HaveOccurred())

			Expect(fakeProxy.StartCallCount()).To(Equal(1))
			key, url := fakeProxy.StartArgsForCall(0)
			Expect(key).To(Equal("some-key"))
			Expect(url).To(Equal("some-url"))

			Expect(fakeProxy.AddrCallCount()).To(Equal(2))

			Expect(addr).To(Equal("some-proxy-addr"))
		})

		Context("when the proxy is already running", func() {
			BeforeEach(func() {
				_, err := cliCoordinator.StartProxy()
				Expect(err).NotTo(HaveOccurred())
			})

			It("returns the address for the existing proxy server", func() {
				addr, err := cliCoordinator.StartProxy()
				Expect(err).NotTo(HaveOccurred())

				Expect(fakeProxy.StartCallCount()).To(Equal(1))
				Expect(fakeProxy.AddrCallCount()).To(Equal(3))

				Expect(addr).To(Equal("some-proxy-addr"))
			})
		})

		Context("when the jumpbox url and the jumpbox ssh key are not set", func() {
			BeforeEach(func() {
				source = concourse.Source{}
				cliCoordinator = bosh.NewCLICoordinator(source, GinkgoWriter, fakeProxy)
			})

			It("does not start a proxy", func() {
				addr, err := cliCoordinator.StartProxy()
				Expect(err).NotTo(HaveOccurred())
				Expect(addr).To(Equal(""))

				Expect(fakeProxy.StartCallCount()).To(Equal(0))
				Expect(fakeProxy.AddrCallCount()).To(Equal(0))
			})
		})

		Context("when the jumpbox url is set and the ssh key is missing", func() {
			BeforeEach(func() {
				source = concourse.Source{JumpboxURL: "some-url"}
				cliCoordinator = bosh.NewCLICoordinator(source, GinkgoWriter, fakeProxy)
			})

			It("returns an error", func() {
				_, err := cliCoordinator.StartProxy()
				Expect(err).To(MatchError("Jumpbox URL and Jumpbox SSH Key are both required to use a jumpbox"))
			})
		})

		Context("when the jumpbox ssh key is set and the jumpbox url is missing", func() {
			BeforeEach(func() {
				source = concourse.Source{JumpboxSSHKey: "some-key"}
				cliCoordinator = bosh.NewCLICoordinator(source, GinkgoWriter, fakeProxy)
			})

			It("returns an error", func() {
				_, err := cliCoordinator.StartProxy()
				Expect(err).To(MatchError("Jumpbox URL and Jumpbox SSH Key are both required to use a jumpbox"))
			})
		})
	})
})
