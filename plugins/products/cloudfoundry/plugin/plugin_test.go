package cloudfoundry_test

import (
	"github.com/enaml-ops/enaml"
	. "github.com/enaml-ops/omg-cli/plugins/products/cloudfoundry/plugin"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Cloud Foundry Plugin", func() {
	Describe("given a cf with Go-Routers", func() {

		Context("when the plugin is called by a operator WITHOUT a complete set of arguments", func() {
			It("then it should return the error and exit", func() {
				Ω(func() {
					cf := new(Plugin)
					cf.GetProduct([]string{
						"cloudfoundry",
						"--router-ip", "1.0.0.1",
						"--router-ip", "1.0.0.2",
						"--router-network", "foundry-net",
					}, []byte(``))
				}).Should(Panic())
			})
		})
		Context("when the plugin is called by a operator with a complete set of arguments", func() {
			var deploymentManifest *enaml.DeploymentManifest
			BeforeEach(func() {
				cf := new(Plugin)
				dm := cf.GetProduct([]string{
					"cloudfoundry",
					"--stemcell-name", "cool-ubuntu-animal",
					"--az", "eastprod-1",
					"--router-ip", "1.0.0.1",
					"--router-ip", "1.0.0.2",
					"--router-network", "foundry-net",
					"--router-vm-type", "blah",
					"--router-ssl-cert", "blah",
					"--router-ssl-key", "blah",
					"--nats-user", "nats",
					"--nats-pass", "pass",
					"--nats-machine-ip", "1.0.0.5",
					"--nats-machine-ip", "1.0.0.6",
					"--etcd-machine-ip", "1.0.0.7",
					"--etcd-machine-ip", "1.0.0.8",
				}, []byte(``))
				deploymentManifest = enaml.NewDeploymentManifest(dm)
			})
			It("then it should allow the user to configure the router IPs", func() {
				ig := deploymentManifest.GetInstanceGroupByName("router-partition")
				network := ig.Networks[0]
				Ω(len(network.StaticIPs)).Should(Equal(2))
				Ω(network.StaticIPs).Should(ConsistOf("1.0.0.1", "1.0.0.2"))
			})

			It("then it should configure the correct number of instances automatically from the count of given IPs", func() {
				ig := deploymentManifest.GetInstanceGroupByName("router-partition")
				network := ig.Networks[0]
				Ω(len(network.StaticIPs)).Should(Equal(ig.Instances))
			})

			It("then it should allow the user to configure the AZs", func() {
				ig := deploymentManifest.GetInstanceGroupByName("router-partition")
				Ω(len(ig.AZs)).Should(Equal(1))
				Ω(ig.AZs[0]).Should(Equal("eastprod-1"))
			})

			It("then it should allow the user to configure vm-type", func() {
				ig := deploymentManifest.GetInstanceGroupByName("router-partition")
				Ω(ig.VMType).ShouldNot(BeEmpty())
				Ω(ig.VMType).Should(Equal("blah"))
			})

			XIt("then it should allow the user to configure network to use", func() {
				Ω(Plugin{}).Should(BeNil())
			})

			XIt("then it should allow the user to configure the used stemcell", func() {
				Ω(Plugin{}).Should(BeNil())
			})

			XIt("then it should allow the user to configure the cert & key used", func() {
				Ω(Plugin{}).Should(BeNil())
			})

			XIt("then it should allow the user to configure if we enable ssl", func() {
				Ω(Plugin{}).Should(BeNil())
			})

			XIt("then it should allow the user to configure the nats pool to use", func() {
				Ω(Plugin{}).Should(BeNil())
			})

			XIt("then it should allow the user to configure the loggregator pool to use", func() {
				Ω(Plugin{}).Should(BeNil())
			})
		})
	})
})
