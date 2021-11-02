package cmd_test

import (
	"reflect"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/cloudfoundry/bosh-davcli/cmd"
	davconf "github.com/cloudfoundry/bosh-davcli/config"
	boshlog "github.com/cloudfoundry/bosh-utils/logger"
)

func buildFactory() (factory Factory) {
	config := davconf.Config{User: "some user"}
	logger := boshlog.NewLogger(boshlog.LevelNone)
	factory = NewFactory(logger)
	factory.SetConfig(config)
	return
}

var _ = Describe("Factory", func() {
	Describe("Create", func() {
		It("factory create a put command", func() {
			factory := buildFactory()
			cmd, err := factory.Create("put")

			Expect(err).ToNot(HaveOccurred())
			Expect(reflect.TypeOf(cmd)).To(Equal(reflect.TypeOf(PutCmd{})))
		})

		It("factory create a get command", func() {
			factory := buildFactory()
			cmd, err := factory.Create("get")

			Expect(err).ToNot(HaveOccurred())
			Expect(reflect.TypeOf(cmd)).To(Equal(reflect.TypeOf(GetCmd{})))
		})

		It("factory create a delete command", func() {
			factory := buildFactory()
			cmd, err := factory.Create("delete")

			Expect(err).ToNot(HaveOccurred())
			Expect(reflect.TypeOf(cmd)).To(Equal(reflect.TypeOf(DeleteCmd{})))
		})

		It("factory create when cmd is unknown", func() {
			factory := buildFactory()
			_, err := factory.Create("some unknown cmd")

			Expect(err).To(HaveOccurred())
		})
	})

	Describe("SetConfig", func() {
		It("returns an error if CaCert is given but invalid", func() {
			factory := buildFactory()
			config := davconf.Config{
				TLS: davconf.TLS{
					Cert: davconf.Cert{
						CA: "--- INVALID CERTIFICATE ---",
					},
				},
			}

			err := factory.SetConfig(config)
			Expect(err).To(HaveOccurred())
		})
		It("does not return an error if CaCert is valid", func() {
			factory := buildFactory()
			cert := `-----BEGIN CERTIFICATE-----
MIICEzCCAXygAwIBAgIQMIMChMLGrR+QvmQvpwAU6zANBgkqhkiG9w0BAQsFADAS
MRAwDgYDVQQKEwdBY21lIENvMCAXDTcwMDEwMTAwMDAwMFoYDzIwODQwMTI5MTYw
MDAwWjASMRAwDgYDVQQKEwdBY21lIENvMIGfMA0GCSqGSIb3DQEBAQUAA4GNADCB
iQKBgQDuLnQAI3mDgey3VBzWnB2L39JUU4txjeVE6myuDqkM/uGlfjb9SjY1bIw4
iA5sBBZzHi3z0h1YV8QPuxEbi4nW91IJm2gsvvZhIrCHS3l6afab4pZBl2+XsDul
rKBxKKtD1rGxlG4LjncdabFn9gvLZad2bSysqz/qTAUStTvqJQIDAQABo2gwZjAO
BgNVHQ8BAf8EBAMCAqQwEwYDVR0lBAwwCgYIKwYBBQUHAwEwDwYDVR0TAQH/BAUw
AwEB/zAuBgNVHREEJzAlggtleGFtcGxlLmNvbYcEfwAAAYcQAAAAAAAAAAAAAAAA
AAAAATANBgkqhkiG9w0BAQsFAAOBgQCEcetwO59EWk7WiJsG4x8SY+UIAA+flUI9
tyC4lNhbcF2Idq9greZwbYCqTTTr2XiRNSMLCOjKyI7ukPoPjo16ocHj+P3vZGfs
h1fIw3cSS2OolhloGw/XM6RWPWtPAlGykKLciQrBru5NAPvCMsb/I1DAceTiotQM
fblo6RBxUQ==
-----END CERTIFICATE-----`
			config := davconf.Config{
				TLS: davconf.TLS{
					Cert: davconf.Cert{
						CA: cert,
					},
				},
			}

			err := factory.SetConfig(config)
			Expect(err).ToNot(HaveOccurred())
		})
		It("does not return an error if CaCert is not provided", func() {
			factory := buildFactory()
			config := davconf.Config{
				TLS: davconf.TLS{
					Cert: davconf.Cert{
						CA: "",
					},
				},
			}

			err := factory.SetConfig(config)
			Expect(err).ToNot(HaveOccurred())
		})
	})
})
