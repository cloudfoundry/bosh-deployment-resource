package out_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"strings"
	"testing"
)

func TestOut(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Out Suite")
}

func properYaml(improperYaml string) []byte {
	return []byte(strings.Replace(improperYaml, "\t", "  ", -1))
}
