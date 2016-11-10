package photoncli_test

import (
	photoncli "github.com/enaml-ops/omg-cli/photon-cli"
	"github.com/enaml-ops/pluginlib/pluginutil"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gopkg.in/urfave/cli.v2"
)

var _ = Describe("given the photon cli", func() {
	Context("when called with a complete set of flags", func() {
		It("then it should not return an error", func() {
			action := photoncli.GetAction(func(s string) {})
			var ctx *cli.Context
			ctx = pluginutil.NewContext([]string{"someapp",
				"--photon-target", "some",
				"--photon-project-id", "stuff",
				"--photon-user", "to",
				"--photon-password", "do",
				"--photon-network-id", "92895-35-2975340-34346346",
				"--bosh-private-ip", "10.0.0.3",
				"--gateway", "10.0.0.254",
				"--cidr", "10.0.0.1/24",
				"--dns", "10.0.0.2",
				"--ntp-server", "10.0.0.2",
				"--bosh-private-ip", "10.0.10.2",
			}, pluginutil.ToCliFlagArray(photoncli.GetFlags()))
			err := action(ctx)
			Ω(err).ShouldNot(HaveOccurred())
		})
	})

	Context("when called with an incomplete set of flags", func() {
		It("then it should return an error", func() {
			action := photoncli.GetAction(func(s string) {})
			var ctx *cli.Context
			ctx = pluginutil.NewContext([]string{"someapp"}, pluginutil.ToCliFlagArray(photoncli.GetFlags()))
			err := action(ctx)
			Ω(err).Should(HaveOccurred())
		})
	})

	Context("when called with a username but without a password", func() {
		It("then it should return an error", func() {
			action := photoncli.GetAction(func(s string) {})
			var ctx *cli.Context
			ctx = pluginutil.NewContext([]string{"someapp",
				"--photon-target", "some",
				"--photon-project-id", "stuff",
				"--photon-user", "to",
				"--photon-network-id", "92895-35-2975340-34346346",
				"--bosh-private-ip", "10.0.0.3",
				"--gateway", "10.0.0.254",
				"--cidr", "10.0.0.1/24",
			}, pluginutil.ToCliFlagArray(photoncli.GetFlags()))
			err := action(ctx)
			Ω(err).Should(HaveOccurred())
		})
	})

	Context("when called with a password but without a username", func() {
		It("then it should return an error", func() {
			action := photoncli.GetAction(func(s string) {})
			var ctx *cli.Context
			ctx = pluginutil.NewContext([]string{"someapp",
				"--photon-target", "some",
				"--photon-project-id", "stuff",
				"--photon-password", "do",
				"--photon-network-id", "92895-35-2975340-34346346",
				"--bosh-private-ip", "10.0.0.3",
				"--gateway", "10.0.0.254",
				"--cidr", "10.0.0.1/24",
			}, pluginutil.ToCliFlagArray(photoncli.GetFlags()))
			err := action(ctx)
			Ω(err).Should(HaveOccurred())
		})
	})
})
