package concourse_test

import (
	"github.com/cloudfoundry/bosh-deployment-resource/concourse"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("InRequest", func() {
	Describe("NewInRequest", func() {
		Context("when the target is empty", func() {
			It("It sets a placeholder so newing up the director does not fail validation", func() {
				request := []byte(`{}`)

				inRequest, err := concourse.NewInRequest(request)
				Expect(err).ToNot(HaveOccurred())

				Expect(inRequest.Source.Target).To(Equal(concourse.MissingTarget))
			})
		})
	})
})
