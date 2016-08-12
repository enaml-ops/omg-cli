package cloudconfigs_test

import (
	"github.com/codegangsta/cli"
	"github.com/enaml-ops/omg-cli/plugins/cloudconfigs"
	"github.com/enaml-ops/pluginlib/util"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("cloud config", func() {
	Context("when creating a cloud config with a single AZ", func() {
		var c *cli.Context

		It("can create a network with more than one DNS address", func() {
			nop := func(f []cli.Flag, i int) []cli.Flag { return f }
			c = pluginutil.NewContext([]string{
				"foo",
				"--az", "az1",
				"--network-name-1", "private",
				"--network-az-1", "az1",
				"--network-cidr-1", "10.180.132.0/22",
				"--network-gateway-1", "10.180.132.254",
				"--network-reserved-1", "10.180.134.0-10.180.135.250",
				"--network-static-1", "MyStaticNetwork",
				"--network-dns-1", "10.148.20.5",
				"--network-dns-1", "10.148.20.6",
			}, cloudconfigs.CreateNetworkFlags([]cli.Flag{
				cli.StringFlag{Name: "az"},
			}, nop))

			validateCP := func(i, j int) error {
				return nil
			}
			cp := func(i, j int) interface{} {
				return nil
			}

			_, err := cloudconfigs.CreateNetworks(c, validateCP, cp)
			Ω(err).ShouldNot(HaveOccurred())
		})

		It("can create a network with more than one reserved address", func() {
			nop := func(f []cli.Flag, i int) []cli.Flag { return f }
			c = pluginutil.NewContext([]string{
				"foo",
				"--az", "az1",
				"--network-name-1", "private",
				"--network-az-1", "az1",
				"--network-cidr-1", "10.180.132.0/22",
				"--network-gateway-1", "10.180.132.254",
				"--network-reserved-1", "10.180.134.0-10.180.135.250",
				"--network-reserved-1", "10.180.134.0-10.180.135.250",
				"--network-static-1", "MyStaticNetwork",
				"--network-dns-1", "10.148.20.6",
			}, cloudconfigs.CreateNetworkFlags([]cli.Flag{
				cli.StringFlag{Name: "az"},
			}, nop))

			validateCP := func(i, j int) error {
				return nil
			}
			cp := func(i, j int) interface{} {
				return nil
			}

			_, err := cloudconfigs.CreateNetworks(c, validateCP, cp)
			Ω(err).ShouldNot(HaveOccurred())
		})

		It("can create a network with more than one static address", func() {
			nop := func(f []cli.Flag, i int) []cli.Flag { return f }
			c = pluginutil.NewContext([]string{
				"foo",
				"--az", "az1",
				"--network-name-1", "private",
				"--network-az-1", "az1",
				"--network-cidr-1", "10.180.132.0/22",
				"--network-gateway-1", "10.180.132.254",
				"--network-reserved-1", "10.180.134.0-10.180.135.250",
				"--network-static-1", "MyStaticNetwork",
				"--network-static-1", "MyStaticNetwork",
				"--network-dns-1", "10.148.20.6",
			}, cloudconfigs.CreateNetworkFlags([]cli.Flag{
				cli.StringFlag{Name: "az"},
			}, nop))

			validateCP := func(i, j int) error {
				return nil
			}
			cp := func(i, j int) interface{} {
				return nil
			}

			_, err := cloudconfigs.CreateNetworks(c, validateCP, cp)
			Ω(err).ShouldNot(HaveOccurred())
		})
	})
})
