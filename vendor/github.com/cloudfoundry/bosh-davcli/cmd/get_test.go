package cmd_test

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/cloudfoundry/bosh-davcli/cmd"
	testcmd "github.com/cloudfoundry/bosh-davcli/cmd/testing"
	davconf "github.com/cloudfoundry/bosh-davcli/config"
	boshlog "github.com/cloudfoundry/bosh-utils/logger"
)

func runGet(config davconf.Config, args []string) error {
	logger := boshlog.NewLogger(boshlog.LevelNone)
	factory := NewFactory(logger)
	factory.SetConfig(config)

	cmd, err := factory.Create("get")
	Expect(err).ToNot(HaveOccurred())

	return cmd.Run(args)
}

func getFileContent(path string) string {
	file, err := os.Open(path)
	Expect(err).ToNot(HaveOccurred())

	fileBytes, err := ioutil.ReadAll(file)
	Expect(err).ToNot(HaveOccurred())

	return string(fileBytes)
}

var _ = Describe("GetCmd", func() {
	var (
		handler        func(http.ResponseWriter, *http.Request)
		targetFilePath string
		requestedBlob  string
		ts             *httptest.Server
		config         davconf.Config
	)

	BeforeEach(func() {
		requestedBlob = "0ca907f2-dde8-4413-a304-9076c9d0978b"
		targetFilePath = filepath.Join(os.TempDir(), "testRunGetCommand.txt")

		handler = func(w http.ResponseWriter, r *http.Request) {
			req := testcmd.NewHTTPRequest(r)

			username, password, err := req.ExtractBasicAuth()
			Expect(err).ToNot(HaveOccurred())
			Expect(req.URL.Path).To(Equal("/0d/" + requestedBlob))
			Expect(req.Method).To(Equal("GET"))
			Expect(username).To(Equal("some user"))
			Expect(password).To(Equal("some pwd"))

			w.Write([]byte("this is your blob"))
		}

	})

	AfterEach(func() {
		os.RemoveAll(targetFilePath)
		ts.Close()
	})

	AssertGetBehavior := func() {
		It("get run with valid args", func() {
			err := runGet(config, []string{requestedBlob, targetFilePath})
			Expect(err).ToNot(HaveOccurred())
			Expect(getFileContent(targetFilePath)).To(Equal("this is your blob"))
		})

		It("get run with incorrect arg count", func() {
			err := runGet(davconf.Config{}, []string{})
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Incorrect usage"))
		})
	}

	Context("with http endpoint", func() {
		BeforeEach(func() {
			ts = httptest.NewServer(http.HandlerFunc(handler))

			config = davconf.Config{
				User:     "some user",
				Password: "some pwd",
				Endpoint: ts.URL,
			}
		})

		AssertGetBehavior()
	})

	Context("with https endpoint", func() {
		BeforeEach(func() {
			ts = httptest.NewTLSServer(http.HandlerFunc(handler))

			rootCa, err := testcmd.ExtractRootCa(ts)
			Expect(err).ToNot(HaveOccurred())

			config = davconf.Config{
				User:     "some user",
				Password: "some pwd",
				Endpoint: ts.URL,
				TLS: davconf.TLS{
					Cert: davconf.Cert{
						CA: rootCa,
					},
				},
			}
		})

		AssertGetBehavior()
	})
})
