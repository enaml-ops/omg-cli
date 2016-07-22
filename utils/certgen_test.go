package utils_test

import (
	. "github.com/enaml-ops/omg-cli/utils"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("certgen", func() {

	Describe("given GenerateCert", func() {
		Context("when called hosts", func() {
			It("then it should return valid cacert, cert and key", func() {
				caCert, cert, key, err := GenerateCert([]string{"*.sys.test.com", "*.apps.test.com"})
				Ω(err).Should(BeNil())
				Ω(len(caCert) > 0).Should(BeTrue())
				Ω(len(cert) > 0).Should(BeTrue())
				Ω(len(key) > 0).Should(BeTrue())
			})
		})
		Context("when called empty hosts and ca", func() {
			It("then it should return valid cacert, cert and key", func() {
				caCert, cert, key, err := GenerateCert(nil)
				Ω(err).Should(BeNil())
				Ω(len(caCert) > 0).Should(BeTrue())
				Ω(len(cert) > 0).Should(BeTrue())
				Ω(len(key) > 0).Should(BeTrue())
			})
		})
	})
})
