package concourse_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
	"strings"
)

func TestOut(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Concourse Suite")
}

func properYaml(improperYaml string) []byte {
	return []byte(strings.Replace(improperYaml, "\t", "  ", -1))
}
