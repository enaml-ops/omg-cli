package cloudfoundry_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/enaml-ops/omg-cli/plugins/products/cloudfoundry/plugin"
)

var _ = Describe("Cloud Foundry Plugin", func() {
	Describe("given the GetProduct Method", func() {
		Context("when called w/ a `vault-active` flag set to TRUE and an INCOMPLETE set of vault values", func() {
			var plugin *Plugin

			BeforeEach(func() {
				plugin = new(Plugin)
			})

			It("then it should panic", func() {
				Ω(func() {
					plugin.GetProduct([]string{
						"my-app",
						"--vault-active",
					},
						[]byte(``),
					)
				}).Should(Panic())
			})
		})
	})
})
