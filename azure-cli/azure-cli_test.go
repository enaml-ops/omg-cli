package azurecli_test

import (
	"gopkg.in/urfave/cli.v2"
	azurecli "github.com/enaml-ops/omg-cli/azure-cli"
	"github.com/enaml-ops/pluginlib/util"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("given the photon cli", func() {
	Context("when called with a complete set of flags", func() {
		It("then it should NOT panic", func() {
			action := azurecli.GetAction(func(s string) {})
			var ctx *cli.Context
			ctx = pluginutil.NewContext([]string{"someapp",
				"--azure-vnet", "some",
				"--azure-subnet", "stuff",
				"--azure-subscription-id", "to",
				"--azure-tenant-id", "do",
				"--azure-client-id", "92895-35-2975340-34346346",
				"--azure-client-secret", "asdfdf",
				"--azure-resource-group", "asdfdf",
				"--azure-security-group", "asdfdf",
				"--azure-storage-account", "asdfdf",
				"--azure-client-secret", "asdfdf",
				"--azure-ssh-pub-key-path", "./fixtures/fake.key",
				"--azure-ssh-user", "asdfasf",
				"--azure-private-key-path", "asfasdffas",
				"--gateway", "10.0.0.254",
				"--cidr", "10.0.0.1/24",
			}, azurecli.GetFlags())
			err := action(ctx)
			Ω(err).ShouldNot(HaveOccurred())
		})
	})

	Context("when called with an incomplete set of flags", func() {
		It("then it should panic and exit", func() {
			action := azurecli.GetAction(func(s string) {})
			var ctx *cli.Context
			ctx = pluginutil.NewContext([]string{"someapp"}, azurecli.GetFlags())
			err := action(ctx)
			Ω(err).Should(HaveOccurred())
		})
	})
})
