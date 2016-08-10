package photoncli_test

import (
	"github.com/codegangsta/cli"
	photoncli "github.com/enaml-ops/omg-cli/photon-cli"
	"github.com/enaml-ops/pluginlib/util"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("given the photon cli", func() {
	Context("when called with a complete set of flags", func() {
		It("then it should NOT panic", func() {
			action := photoncli.GetAction(func(s string) {})
			var ctx *cli.Context
			ctx = pluginutil.NewContext([]string{"someapp",
				"--photon-target", "some",
				"--photon-project-id", "stuff",
				"--photon-user", "to",
				"--photon-password", "do",
			}, photoncli.GetFlags())
			Ω(func() { action(ctx) }).ShouldNot(Panic())
		})
	})
	Context("when called with an incomplete set of flags", func() {
		It("then it should panic and exit", func() {
			action := photoncli.GetAction(func(s string) {})
			var ctx *cli.Context
			ctx = pluginutil.NewContext([]string{"someapp"}, photoncli.GetFlags())
			Ω(func() { action(ctx) }).Should(Panic())
		})
	})
})
