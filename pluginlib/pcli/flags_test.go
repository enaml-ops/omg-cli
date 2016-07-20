package pcli

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("CLI flags", func() {
	Context("when creating new CLI flags", func() {
		const controlType = BoolFlag
		const controlName = "my-cool-flag"
		const controlUsage = "usage for my cool flag"

		It("should set the type, name, and usage correctly", func() {
			flag := NewFlag(controlType, controlName, controlUsage, "")
			Ω(flag.Typ).Should(Equal(controlType))
			Ω(flag.Name).Should(Equal(controlName))
			Ω(flag.Usage).Should(Equal(controlUsage))
		})

		It("should create the env variable automatically", func() {
			flag := NewFlag(controlType, controlName, controlUsage, "")
			Ω(flag.EnvVar).Should(Equal("MY_COOL_FLAG"))
		})
	})
})
