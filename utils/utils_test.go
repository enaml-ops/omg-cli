package utils_test

import (
	. "github.com/enaml-ops/omg-cli/utils"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("given ClearDefaultStringSliceValue", func() {
	Context("when called on a stringslice containing a default & added values", func() {
		It("then it should clear the default value from the list", func() {
			stringSlice := []string{"default", "useradded1", "useradded2"}
			clearSlice := ClearDefaultStringSliceValue(stringSlice...)
			Ω(clearSlice).Should(ConsistOf("useradded1", "useradded2"))
			Ω(clearSlice).ShouldNot(ContainElement("default"))
		})
	})

	Context("when called on a stringslice only containing a default value", func() {
		It("then it should simply pass through the default value", func() {
			stringSlice := []string{"default"}
			Ω(ClearDefaultStringSliceValue(stringSlice...)).Should(Equal(stringSlice))
		})
	})
})
