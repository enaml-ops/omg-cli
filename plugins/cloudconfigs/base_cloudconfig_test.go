package cloudconfigs_test

import (
	"github.com/enaml-ops/enaml"
	"github.com/enaml-ops/omg-cli/plugins/cloudconfigs"
	"github.com/enaml-ops/pluginlib/pcli"
	"github.com/enaml-ops/pluginlib/pluginutil"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gopkg.in/urfave/cli.v2"
)

var _ = Describe("cloud config", func() {
	Context("when creating a cloud config with a single AZ", func() {

		validateCP := func(i, j int) error {
			return nil
		}
		cp := func(i, j int) interface{} {
			return nil
		}

		Context("when specifying multiple DNS addresses", func() {
			var c *cli.Context
			var err error
			var nets []enaml.DeploymentNetwork
			BeforeEach(func() {
				nop := func(f []pcli.Flag, i int) []pcli.Flag { return f }
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
				}, pluginutil.ToCliFlagArray(cloudconfigs.CreateNetworkFlags([]pcli.Flag{pcli.CreateStringFlag("az", "az flag")}, nop)))
				nets, err = cloudconfigs.CreateNetworks(c, validateCP, cp)
			})

			It("can create a network with more than one DNS address", func() {
				Ω(err).ShouldNot(HaveOccurred())
			})

			It("should create multiple DNS records in the network list", func() {
				Ω(len(nets)).Should(Equal(1))
				mn := nets[0].(*enaml.ManualNetwork)
				Ω(len(mn.Subnets)).Should(Equal(1))
				Ω(len(mn.Subnets[0].DNS)).Should(Equal(2))
			})
		})

		Context("when specifying multiple reserved addresses", func() {
			var c *cli.Context
			var err error
			var nets []enaml.DeploymentNetwork
			BeforeEach(func() {
				nop := func(f []pcli.Flag, i int) []pcli.Flag { return f }
				c = pluginutil.NewContext([]string{
					"foo",
					"--az", "az1",
					"--network-name-1", "private",
					"--network-az-1", "az1",
					"--network-cidr-1", "10.180.132.0/22",
					"--network-gateway-1", "10.180.132.254",
					"--network-reserved-1", "10.180.134.0-10.180.135.250",
					"--network-reserved-1", "10.180.134.0-10.180.135.251",
					"--network-static-1", "MyStaticNetwork",
					"--network-dns-1", "10.148.20.6",
				}, pluginutil.ToCliFlagArray(cloudconfigs.CreateNetworkFlags([]pcli.Flag{pcli.CreateStringFlag("az", "az flag")}, nop)))
				nets, err = cloudconfigs.CreateNetworks(c, validateCP, cp)
			})

			It("can create a network with more than one reserved address", func() {
				Ω(err).ShouldNot(HaveOccurred())
			})

			It("should create multiple reserved records in the network list", func() {
				Ω(len(nets)).Should(Equal(1))
				mn := nets[0].(*enaml.ManualNetwork)
				Ω(len(mn.Subnets)).Should(Equal(1))
				Ω(len(mn.Subnets[0].Reserved)).Should(Equal(2))
			})
		})

		Context("when specifying multiple static addresses", func() {
			var c *cli.Context
			var err error
			var nets []enaml.DeploymentNetwork
			BeforeEach(func() {
				nop := func(f []pcli.Flag, i int) []pcli.Flag { return f }
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
				}, pluginutil.ToCliFlagArray(cloudconfigs.CreateNetworkFlags([]pcli.Flag{pcli.CreateStringFlag("az", "az flag")}, nop)))
				nets, err = cloudconfigs.CreateNetworks(c, validateCP, cp)
			})

			It("can create a network with more than one static address", func() {
				Ω(err).ShouldNot(HaveOccurred())
			})

			It("should create multiple static records in the network list", func() {
				Ω(len(nets)).Should(Equal(1))
				mn := nets[0].(*enaml.ManualNetwork)
				Ω(len(mn.Subnets)).Should(Equal(1))
				Ω(len(mn.Subnets[0].Static)).Should(Equal(2))
			})
		})
	})
})
