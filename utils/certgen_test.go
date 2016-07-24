package utils_test

import (
	"crypto/x509"
	"encoding/pem"

	. "github.com/enaml-ops/omg-cli/utils"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("certgen", func() {

	Describe("given GenerateCert", func() {
		BeforeEach(func() {

		})
		Context("when called hosts", func() {
			It("then it should return valid cacert, cert and key", func() {
				caCert, cert, key, err := GenerateCert([]string{"*.sys.test.com", "*.apps.test.com", "sys.test.com"})
				Ω(err).Should(BeNil())
				Ω(len(caCert) > 0).Should(BeTrue())
				Ω(len(cert) > 0).Should(BeTrue())
				Ω(len(key) > 0).Should(BeTrue())

				roots := x509.NewCertPool()
				ok := roots.AppendCertsFromPEM([]byte(caCert))
				if !ok {
					panic("failed to parse root certificate")
				}

				block, _ := pem.Decode([]byte(cert))
				if block == nil {
					panic("failed to parse certificate PEM")
				}
				certificate, theError := x509.ParseCertificate(block.Bytes)
				if theError != nil {
					panic("failed to parse certificate: " + theError.Error())
				}

				opts := x509.VerifyOptions{
					DNSName: "sys.test.com",
					Roots:   roots,
				}

				if _, theError := certificate.Verify(opts); theError != nil {
					panic("failed to verify certificate: " + theError.Error())
				}
			})
		})
	})
})
