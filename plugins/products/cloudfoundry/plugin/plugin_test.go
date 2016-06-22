package cloudfoundry_test

import (
	"github.com/enaml-ops/enaml"
	. "github.com/enaml-ops/omg-cli/plugins/products/cloudfoundry/plugin"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/xchapter7x/lo"
)

var _ = Describe("Cloud Foundry Plugin", func() {
	Describe("given a cf with Go-Routers", func() {

		Context("when the plugin is called by a operator WITHOUT a complete set of arguments", func() {
			It("then it should return the error and exit", func() {
				Ω(func() {
					cf := new(Plugin)
					cf.GetProduct([]string{
						"cloudfounddry",
						"--router-ip", "1.0.0.1",
						"--router-ip", "1.0.0.2",
						"--router-network", "foundry-net",
					}, []byte(``))
				}).Should(Panic())
			})
		})
		XContext("when the plugin is called by a operator with a complete set of arguments", func() {
			var deploymentManifest *enaml.DeploymentManifest
			BeforeEach(func() {
				cf := new(Plugin)
				dm := cf.GetProduct([]string{
					"cloudfounddry",
					"--router-ip", "1.0.0.1",
					"--router-ip", "1.0.0.2",
					"--router-network", "foundry-net",
				}, []byte(``))
				deploymentManifest = enaml.NewDeploymentManifest(dm)
			})
			It("then it should allow the user to configure the IPs", func() {
				lo.G.Debug("do something here: ", deploymentManifest)
				job := deploymentManifest.GetJobByName("router-partition")
				network := job.Networks[0]
				Ω(len(network.StaticIPs)).Should(Equal(2))
				Ω(network.StaticIPs).Should(ConsistOf("1.0.0.1", "1.0.0.2"))
			})

			It("then it should allow the user to configure the AZs", func() {
				Ω(Plugin{}).Should(BeNil())
			})

			It("then it should allow the user to configure vm-type", func() {
				Ω(Plugin{}).Should(BeNil())
			})

			It("then it should allow the user to configure network to use", func() {
				Ω(Plugin{}).Should(BeNil())
			})

			It("then it should allow the user to configure the used stemcell", func() {
				Ω(Plugin{}).Should(BeNil())
			})

			It("then it should allow the user to configure the cert & key used", func() {
				Ω(Plugin{}).Should(BeNil())
			})

			It("then it should allow the user to configure if we enable ssl", func() {
				Ω(Plugin{}).Should(BeNil())
			})

			It("then it should allow the user to configure the nats pool to use", func() {
				Ω(Plugin{}).Should(BeNil())
			})

			It("then it should allow the user to configure the loggregator pool to use", func() {
				Ω(Plugin{}).Should(BeNil())
			})
		})
	})
})
