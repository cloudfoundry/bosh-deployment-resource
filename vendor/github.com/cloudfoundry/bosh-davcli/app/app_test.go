package app_test

import (
	"errors"
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/cloudfoundry/bosh-davcli/app"
	davconf "github.com/cloudfoundry/bosh-davcli/config"
)

type FakeRunner struct {
	Config       davconf.Config
	SetConfigErr error
	RunArgs      []string
	RunErr       error
}

func (r *FakeRunner) SetConfig(newConfig davconf.Config) (err error) {
	r.Config = newConfig
	return r.SetConfigErr
}

func (r *FakeRunner) Run(cmdArgs []string) (err error) {
	r.RunArgs = cmdArgs
	return r.RunErr
}

func pathToFixture(file string) string {
	pwd, err := os.Getwd()
	Expect(err).ToNot(HaveOccurred())

	fixturePath := filepath.Join(pwd, "../test_assets", file)

	absPath, err := filepath.Abs(fixturePath)
	Expect(err).ToNot(HaveOccurred())

	return absPath
}

var _ = Describe("App", func() {
	It("reads the CA cert from config", func() {
		runner := &FakeRunner{}

		app := New(runner)
		err := app.Run([]string{"dav-cli", "-c", pathToFixture("dav-cli-config-with-ca.json"), "put", "localFile", "remoteFile"})
		Expect(err).ToNot(HaveOccurred())

		expectedConfig := davconf.Config{
			User:     "some user",
			Password: "some pwd",
			Endpoint: "https://example.com/some/endpoint",
			Secret:   "77D47E3A0B0F590B73CF3EBD9BB6761E244F90FA6F28BB39F941B0905789863FBE2861FDFD8195ADC81B72BB5310BC18969BEBBF4656366E7ACD3F0E4186FDDA",
			TLS: davconf.TLS{
				Cert: davconf.Cert{
					CA: "ca-cert",
				},
			},
		}

		Expect(runner.Config).To(Equal(expectedConfig))
		Expect(runner.Config.TLS.Cert.CA).ToNot(BeNil())
	})

	It("returns error if CA Cert is invalid", func() {
		runner := &FakeRunner{
			SetConfigErr: errors.New("invalid cert"),
		}

		app := New(runner)
		err := app.Run([]string{"dav-cli", "-c", pathToFixture("dav-cli-config-with-ca.json"), "put", "localFile", "remoteFile"})
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("Invalid CA Certificate: invalid cert"))

	})

	It("runs the put command", func() {
		runner := &FakeRunner{}

		app := New(runner)
		err := app.Run([]string{"dav-cli", "-c", pathToFixture("dav-cli-config.json"), "put", "localFile", "remoteFile"})
		Expect(err).ToNot(HaveOccurred())

		expectedConfig := davconf.Config{
			User:     "some user",
			Password: "some pwd",
			Endpoint: "http://example.com/some/endpoint",
			Secret:   "77D47E3A0B0F590B73CF3EBD9BB6761E244F90FA6F28BB39F941B0905789863FBE2861FDFD8195ADC81B72BB5310BC18969BEBBF4656366E7ACD3F0E4186FDDA",
		}

		Expect(runner.Config).To(Equal(expectedConfig))
		Expect(runner.Config.TLS.Cert.CA).To(BeEmpty())
		Expect(runner.RunArgs).To(Equal([]string{"put", "localFile", "remoteFile"}))
	})

	It("returns error with no config argument", func() {
		runner := &FakeRunner{}

		app := New(runner)
		err := app.Run([]string{"put", "localFile", "remoteFile"})
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("Config file arg `-c` is missing"))
	})
	It("prints the version info with the -v flag", func() {
		runner := &FakeRunner{}
		app := New(runner)
		err := app.Run([]string{"dav-cli", "-v"})
		Expect(err).ToNot(HaveOccurred())
	})

	It("returns error from the cmd runner", func() {
		runner := &FakeRunner{
			RunErr: errors.New("fake-run-error"),
		}

		app := New(runner)
		err := app.Run([]string{"dav-cli", "-c", pathToFixture("dav-cli-config.json"), "put", "localFile", "remoteFile"})
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("fake-run-error"))
	})
})
