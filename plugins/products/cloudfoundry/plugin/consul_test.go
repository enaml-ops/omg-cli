package cloudfoundry_test

import (
	"github.com/enaml-ops/omg-cli/plugins/products/cloudfoundry/enaml-gen/consul_agent"
	. "github.com/enaml-ops/omg-cli/plugins/products/cloudfoundry/plugin"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Consul Partition", func() {
	Context("when initialized WITHOUT a complete set of arguments", func() {
		It("then it should return the error and exit", func() {
			plugin := new(Plugin)
			c := plugin.GetContext([]string{
				"cloudfoundry",
				"--consul-ip", "1.0.0.1",
				"--consul-ip", "1.0.0.2",
				"--consul-network", "foundry-net",
			})
			_, err := NewConsulPartition(c)
			Ω(err).ShouldNot(BeNil())
		})
	})
	Context("when initialized WITH a complete set of arguments", func() {
		var err error
		var consul InstanceGroupFactory
		BeforeEach(func() {
			plugin := new(Plugin)
			c := plugin.GetContext([]string{
				"cloudfoundry",
				"--stemcell-name", "cool-ubuntu-animal",
				"--az", "eastprod-1",
				"--consul-ip", "1.0.0.1",
				"--consul-ip", "1.0.0.2",
				"--consul-network", "foundry-net",
				"--consul-vm-type", "blah",
				"--consul-encryption-key", "encyption-key",
				"--consul-ca-cert", "ca-cert",
				"--consul-agent-cert", "agent-cert",
				"--consul-agent-key", "agent-key",
				"--consul-server-cert", "server-cert",
				"--consul-server-key", "server-key",
			})
			consul, err = NewConsulPartition(c)
		})
		It("then it should not return an error", func() {
			Ω(err).Should(BeNil())
		})
		It("then it should allow the user to configure the consul IPs", func() {
			ig := consul.ToInstanceGroup()
			network := ig.Networks[0]
			Ω(len(network.StaticIPs)).Should(Equal(2))
			Ω(network.StaticIPs).Should(ConsistOf("1.0.0.1", "1.0.0.2"))
		})
		It("then it should have 2 instances", func() {
			ig := consul.ToInstanceGroup()
			Ω(ig.Instances).Should(Equal(2))
		})
		It("then it should configure the correct number of instances automatically from the count of given IPs", func() {
			ig := consul.ToInstanceGroup()
			network := ig.Networks[0]
			Ω(len(network.StaticIPs)).Should(Equal(ig.Instances))
		})

		It("then it should allow the user to configure the AZs", func() {
			ig := consul.ToInstanceGroup()
			Ω(len(ig.AZs)).Should(Equal(1))
			Ω(ig.AZs[0]).Should(Equal("eastprod-1"))
		})

		It("then it should allow the user to configure vm-type", func() {
			ig := consul.ToInstanceGroup()
			Ω(ig.VMType).ShouldNot(BeEmpty())
			Ω(ig.VMType).Should(Equal("blah"))
		})

		It("then it should allow the user to configure network to use", func() {
			ig := consul.ToInstanceGroup()
			network := ig.GetNetworkByName("foundry-net")
			Ω(network).ShouldNot(BeNil())
		})

		It("then it should allow the user to configure the used stemcell", func() {
			ig := consul.ToInstanceGroup()
			Ω(ig.Stemcell).ShouldNot(BeEmpty())
			Ω(ig.Stemcell).Should(Equal("cool-ubuntu-animal"))
		})

		XIt("then it should then have 3 jobs", func() {
			ig := consul.ToInstanceGroup()
			Ω(len(ig.Jobs)).Should(Equal(3))
		})
		It("then it should then have consul agent job", func() {
			ig := consul.ToInstanceGroup()
			job := ig.GetJobByName("consul_agent")
			Ω(job).ShouldNot(BeNil())
			props, _ := job.Properties.(*consul_agent.Consul)
			Ω(props.ServerKey).Should(Equal("server-key"))
			Ω(props.ServerCert).Should(Equal("server-cert"))
			Ω(props.AgentCert).Should(Equal("agent-cert"))
			Ω(props.AgentKey).Should(Equal("agent-key"))
			Ω(props.CaCert).Should(Equal("ca-cert"))
			Ω(props.EncryptKeys).Should(Equal([]string{"encyption-key"}))
			agent := props.Agent
			Ω(agent.Servers.Lan).Should(Equal([]string{"1.0.0.1", "1.0.0.2"}))
		})

	})
})
