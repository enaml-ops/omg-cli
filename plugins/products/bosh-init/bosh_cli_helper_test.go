package boshinit_test

import (
	. "github.com/enaml-ops/omg-cli/plugins/products/bosh-init"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("BOSH CLI helpers", func() {
	Describe("given a set of bosh required flags", func() {
		Context("when used to check for use of required flags", func() {
			It("then it should only allow fields to be defined when they are valid bosh flags", func() {
				var validFlags []string
				for _, f := range BoshFlags(NewPhotonBoshBase(new(BoshBase))) {
					validFlags = append(validFlags, f.Names()[0])
				}
				for _, required := range RequiredBoshFlags {
					Ω(validFlags).Should(ContainElement(required))
				}
			})
		})
	})
})
