package cloudfoundry_test

import (
	. "github.com/enaml-ops/omg-cli/plugins/products/cloudfoundry/plugin"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Cloud Foundry Deployment", func() {
	XContext("", func() {
		It("", func() {
			Ω(Plugin{}).Should(BeNil())
		})
	})
})
