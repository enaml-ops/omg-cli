package boshinit_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("NewBoshDeployment", func() {
	XDescribe("given the function", func() {
		Context("when called w/ valid parameters", func() {
			It("then it should generate a valid base bosh init object", func() {
				Ω(true).Should(Equal(false))
			})
		})
	})
})
