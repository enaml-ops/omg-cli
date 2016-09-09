package utils_test

import (
	"os"
	"path"

	. "github.com/enaml-ops/omg-cli/utils"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gopkg.in/urfave/cli.v2"
)

var _ = Describe("utils", func() {

	Describe("given a GetCloudConfigCommands func", func() {
		Context("when called with a valid plugin dir", func() {
			var commands []*cli.Command
			BeforeEach(func() {
				commands = GetCloudConfigCommands("../pluginlib/registry/fixtures/cloudconfig")
			})
			It("then it should return a set of commands for the plugins in the dir", func() {
				Ω(len(commands)).Should(Equal(1))
				Ω(commands[0].Name).Should(ContainSubstring("testplugin-"))
				Ω(commands[0].Action).ShouldNot(BeNil())
			})
		})
	})

	Describe("given a GetProductCommands func", func() {
		Context("when called with a valid plugin dir", func() {
			var commands []*cli.Command
			BeforeEach(func() {
				commands = GetProductCommands("../pluginlib/registry/fixtures/product")
			})
			It("then it should return a set of commands for the plugins in the dir", func() {
				Ω(len(commands)).Should(Equal(1))
				Ω(commands[0].Name).Should(ContainSubstring("testproductplugin-"))
				Ω(commands[0].Action).ShouldNot(BeNil())
			})
		})
	})

	Describe("given a deployyaml method", func() {
		Context("when called with valid arguments", func() {
			var stringSpy string

			AfterEach(func() {
				stringSpy = ""
			})

			It("then it should create a temporary file in the pwd", func() {
				DeployYaml("something", func(s string) { stringSpy = s })
				matchstr, _ := os.Getwd()
				matchstr = path.Join(matchstr, "omg-bosh.*")
				ismatch, err := path.Match(matchstr, stringSpy)
				Ω(ismatch).Should(BeTrue())
				Ω(err).ShouldNot(HaveOccurred())
			})
			It("then it should cleanup its tempfile", func() {
				DeployYaml("something", func(s string) { stringSpy = s })
				_, errnofile := os.Stat(stringSpy)
				Ω(errnofile).Should(HaveOccurred())
			})
		})
	})
})
