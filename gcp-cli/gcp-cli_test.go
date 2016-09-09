package gcpcli_test

import (
	"gopkg.in/urfave/cli.v2"
	gcpcli "github.com/enaml-ops/omg-cli/gcp-cli"
	"github.com/enaml-ops/pluginlib/util"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("given the photon cli", func() {
	Context("when called with a complete set of flags", func() {
		It("then it should NOT panic", func() {
			action := gcpcli.GetAction(func(s string) {})
			var ctx *cli.Context
			ctx = pluginutil.NewContext([]string{"someapp",
				"--gcp-network-name", "some",
				"--gcp-subnetwork-name", "stuff",
				"--gcp-default-zone", "to",
				"--gcp-project", "do",
				"--gateway", "10.0.0.254",
				"--cidr", "10.0.0.1/24",
			}, gcpcli.GetFlags())
			err := action(ctx)
			Ω(err).ShouldNot(HaveOccurred())
		})
	})

	Context("when called with an incomplete set of flags", func() {
		It("then it should panic and exit", func() {
			action := gcpcli.GetAction(func(s string) {})
			var ctx *cli.Context
			ctx = pluginutil.NewContext([]string{"someapp"}, gcpcli.GetFlags())
			err := action(ctx)
			Ω(err).Should(HaveOccurred())
		})
	})
})
