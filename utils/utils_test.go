package utils_test

import (
	"github.com/codegangsta/cli"
	. "github.com/enaml-ops/omg-cli/utils"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("utils", func() {

	Describe("given a GetCloudConfigCommands func", func() {
		Context("when called with a valid plugin dir", func() {
			var commands []cli.Command
			BeforeEach(func() {
				commands = GetCloudConfigCommands("../pluginlib/registry/fixtures/cloudconfig")
			})
			It("then it should return a set of commands for the plugins in the dir", func() {
				Ω(len(commands)).Should(Equal(1))
				Ω(commands[0].Name).Should(Equal("testplugin-linux"))
				Ω(commands[0].Action).ShouldNot(BeNil())
			})
		})
	})

	Describe("given a GetProductCommands func", func() {
		Context("when called with a valid plugin dir", func() {
			var commands []cli.Command
			BeforeEach(func() {
				commands = GetProductCommands("../pluginlib/registry/fixtures/product")
			})
			It("then it should return a set of commands for the plugins in the dir", func() {
				Ω(len(commands)).Should(Equal(1))
				Ω(commands[0].Name).Should(Equal("testproductplugin-linux"))
				Ω(commands[0].Action).ShouldNot(BeNil())
			})
		})
	})

	Describe("given ClearDefaultStringSliceValue", func() {
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
})
