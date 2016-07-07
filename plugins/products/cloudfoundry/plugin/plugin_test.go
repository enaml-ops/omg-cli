package cloudfoundry_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Cloud Foundry Plugin", func() {
	XDescribe("given a fully initialized cf plugin", func() {
		Context("when calling getproduct", func() {
			It("then it should return a deployment with all required releases", func() {
				Ω(true).Should(BeFalse())
			})
		})
	})
})
