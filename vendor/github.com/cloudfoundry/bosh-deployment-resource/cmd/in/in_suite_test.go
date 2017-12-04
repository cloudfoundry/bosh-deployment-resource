package main_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

// There are no tests, we just want to make sure it builds

func TestIn(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "In Suite")
}
