package cmd_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/cloudfoundry/bosh-davcli/cmd"

	davconf "github.com/cloudfoundry/bosh-davcli/config"
	boshlog "github.com/cloudfoundry/bosh-utils/logger"
)

func runSign(config davconf.Config, args []string) error {
	logger := boshlog.NewLogger(boshlog.LevelNone)
	factory := NewFactory(logger)
	factory.SetConfig(config)

	cmd, err := factory.Create("sign")
	Expect(err).ToNot(HaveOccurred())

	return cmd.Run(args)
}

var _ = Describe("SignCmd", func() {
	var (
		objectID = "0ca907f2-dde8-4413-a304-9076c9d0978b"
		config   davconf.Config
	)

	It("with valid args", func() {
		err := runSign(config, []string{objectID, "get", "15m"})
		Expect(err).ToNot(HaveOccurred())
	})

	It("returns err with incorrect arg count", func() {
		err := runSign(davconf.Config{}, []string{})
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("incorrect usage"))
	})

	It("returns err with non-implemented action", func() {
		err := runSign(davconf.Config{}, []string{objectID, "delete", "15m"})
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("action not implemented"))
	})

	It("returns err with incorrect duration", func() {
		err := runSign(davconf.Config{}, []string{objectID, "put", "15"})
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("expiration should be a duration value"))
	})
})
