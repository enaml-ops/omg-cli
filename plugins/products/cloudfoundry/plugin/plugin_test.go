package cloudfoundry_test

import (
	. "github.com/enaml-ops/omg-cli/plugins/products/cloudfoundry/plugin"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Cloud Foundry Plugin", func() {
	Describe("given a cf with Go-Routers", func() {
		XContext("when the plugin is called by a operator", func() {
			It("then it should allow the user to configure the IPs", func() {
				Ω(Plugin{}).Should(BeNil())
			})

			It("then it should allow the user to configure the AZs", func() {
				Ω(Plugin{}).Should(BeNil())
			})

			It("then it should allow the user to configure vm-type", func() {
				Ω(Plugin{}).Should(BeNil())
			})

			It("then it should allow the user to configure the used stemcell", func() {
				Ω(Plugin{}).Should(BeNil())
			})

			It("then it should allow the user to configure the cert & key used", func() {
				Ω(Plugin{}).Should(BeNil())
			})

			It("then it should allow the user to configure if we enable ssl", func() {
				Ω(Plugin{}).Should(BeNil())
			})

			It("then it should allow the user to configure the nats pool to use", func() {
				Ω(Plugin{}).Should(BeNil())
			})

			It("then it should allow the user to configure the loggregator pool to use", func() {
				Ω(Plugin{}).Should(BeNil())
			})
		})
	})
})
