package cloudfoundry_test

import (
	"github.com/enaml-ops/enaml"

	grtrlib "github.com/enaml-ops/omg-cli/plugins/products/cloudfoundry/enaml-gen/gorouter"
	"github.com/enaml-ops/omg-cli/plugins/products/cloudfoundry/enaml-gen/metron_agent"
	. "github.com/enaml-ops/omg-cli/plugins/products/cloudfoundry/plugin"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Go-Router Partition", func() {
	Context("when initialized WITHOUT a complete set of arguments", func() {
		It("then it should return the error and exit", func() {
			cf := new(Plugin)
			c := cf.GetContext([]string{
				"cloudfoundry",
				"--router-ip", "1.0.0.1",
				"--router-ip", "1.0.0.2",
				"--router-network", "foundry-net",
			})
			gr := NewGoRouterPartition(c)
			Ω(gr.HasValidValues()).Should(BeFalse())
		})
	})
	Context("when the plugin is called by a operator with a complete set of arguments", func() {
		var deploymentManifest *enaml.DeploymentManifest
		BeforeEach(func() {
			cf := new(Plugin)
			c := cf.GetContext([]string{
				"cloudfoundry",
				"--stemcell-name", "cool-ubuntu-animal",
				"--az", "eastprod-1",
				"--router-ip", "1.0.0.1",
				"--router-ip", "1.0.0.2",
				"--router-network", "foundry-net",
				"--router-vm-type", "blah",
				"--router-ssl-cert", "blah",
				"--router-ssl-key", "blah",
				"--router-pass", "blabadebleblahblah",
				"--metron-secret", "metronsecret",
				"--metron-zone", "metronzoneguid",
				"--nats-user", "nats",
				"--nats-pass", "pass",
				"--nats-machine-ip", "1.0.0.5",
				"--nats-machine-ip", "1.0.0.6",
				"--etcd-machine-ip", "1.0.0.7",
				"--etcd-machine-ip", "1.0.0.8",
				"--router-enable-ssl",
			})
			gr := NewGoRouterPartition(c)
			deploymentManifest = new(enaml.DeploymentManifest)
			deploymentManifest.AddInstanceGroup(gr.ToInstanceGroup())
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

		It("then it should allow the user to configure network to use", func() {
			ig := deploymentManifest.GetInstanceGroupByName("router-partition")
			network := ig.GetNetworkByName("foundry-net")
			Ω(network).ShouldNot(BeNil())
		})

		It("then it should allow the user to configure the used stemcell", func() {
			ig := deploymentManifest.GetInstanceGroupByName("router-partition")
			Ω(ig.Stemcell).ShouldNot(BeEmpty())
			Ω(ig.Stemcell).Should(Equal("cool-ubuntu-animal"))
		})

		It("then it should allow the user to configure if we enable ssl", func() {
			ig := deploymentManifest.GetInstanceGroupByName("router-partition")
			job := ig.GetJobByName("gorouter")
			properties := job.Properties.(*grtrlib.Gorouter)
			Ω(properties.Router.EnableSsl).Should(BeTrue())
		})

		It("then it should allow the user to configure the nats pool to use", func() {
			ig := deploymentManifest.GetInstanceGroupByName("router-partition")
			job := ig.GetJobByName("gorouter")
			properties := job.Properties.(*grtrlib.Gorouter)
			Ω(properties.Nats.Machines).Should(ConsistOf("1.0.0.5", "1.0.0.6"))
			Ω(properties.Nats.User).Should(Equal("nats"))
			Ω(properties.Nats.Password).Should(Equal("pass"))
			Ω(properties.Nats.Port).Should(Equal(4222))
		})

		It("then it should allow the user to configure the loggregator pool to use", func() {
			ig := deploymentManifest.GetInstanceGroupByName("router-partition")
			job := ig.GetJobByName("metron_agent")
			properties := job.Properties.(*metron_agent.MetronAgent)
			Ω(properties.Loggregator.Etcd.Machines).Should(ConsistOf("1.0.0.7", "1.0.0.8"))
		})

		It("then it should allow the user to configure the metron agent", func() {
			ig := deploymentManifest.GetInstanceGroupByName("router-partition")
			job := ig.GetJobByName("metron_agent")
			properties := job.Properties.(*metron_agent.MetronAgent)
			Ω(properties.MetronAgent.Zone).Should(Equal("metronzoneguid"))
			Ω(properties.MetronAgent.Deployment).Should(Equal(DeploymentName))
			Ω(properties.MetronEndpoint.SharedSecret).Should(Equal("metronsecret"))
		})

		It("then it should allow the user to configure the router user/pass", func() {
			ig := deploymentManifest.GetInstanceGroupByName("router-partition")
			job := ig.GetJobByName("gorouter")
			properties := job.Properties.(*grtrlib.Gorouter)
			Ω(properties.Router.Status.User).Should(Equal("router_status"))
			Ω(properties.Router.Status.Password).Should(Equal("blabadebleblahblah"))
		})

		XIt("then it should allow the user to configure the cert & key used from a file", func() {
			ig := deploymentManifest.GetInstanceGroupByName("router-partition")
			job := ig.GetJobByName("gorouter")
			Ω(job.Properties).Should(
				HaveKeyWithValue("router",
					HaveKeyWithValue("ssl_cert", Equal("router_status"))))

			Ω(job.Properties).Should(
				HaveKeyWithValue("router",
					HaveKeyWithValue("ssl_key", Equal("blabadebleblahblah"))))
		})
	})
})
