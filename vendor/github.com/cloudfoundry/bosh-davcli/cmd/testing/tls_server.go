package testing

import (
	"bytes"
	"crypto/x509"
	"encoding/pem"
	"net/http/httptest"
)

func ExtractRootCa(server *httptest.Server) (rootCaStr string, err error) {
	rootCa := new(bytes.Buffer)

	cert, err := x509.ParseCertificate(server.TLS.Certificates[0].Certificate[0])
	if err != nil {
		panic(err.Error())
	}
	// TODO: Replace above with following on Go 1.9
	//cert := server.Certificate()

	block := &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: cert.Raw,
	}

	err = pem.Encode(rootCa, block)
	if err != nil {
		return "", err
	}

	return rootCa.String(), nil
}
