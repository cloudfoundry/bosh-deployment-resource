package in_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"strings"
	"testing"
)

func TestOut(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "In Suite")
}

func properYaml(improperYaml string) []byte {
	return []byte(strings.Replace(improperYaml, "\t", "  ", -1))
}
