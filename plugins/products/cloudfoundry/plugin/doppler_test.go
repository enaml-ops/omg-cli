package cloudfoundry_test

import (
	"github.com/enaml-ops/omg-cli/plugins/products/cloudfoundry/enaml-gen/consul_agent"
	. "github.com/enaml-ops/omg-cli/plugins/products/cloudfoundry/plugin"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Doppler Partition", func() {
	Context("when initialized WITHOUT a complete set of arguments", func() {
		It("then HasValidValues should be false", func() {
			plugin := new(Plugin)
			c := plugin.GetContext([]string{
				"cloudfoundry",
				"--network", "foundry-net",
				"--metron-secret", "metronsecret",
				"--metron-zone", "metronzoneguid",
				"--syslog-address", "syslog-server",
				"--syslog-port", "10601",
				"--syslog-transport", "tcp",
				"--etcd-machine-ip", "1.0.0.7",
				"--etcd-machine-ip", "1.0.0.8",
			})
			doppler := NewDopplerPartition(c)
			Ω(doppler.HasValidValues()).Should(BeFalse())
		})
	})
	Context("when initialized WITH a complete set of arguments", func() {
		var doppler InstanceGrouper
		BeforeEach(func() {
			plugin := new(Plugin)
			c := plugin.GetContext([]string{
				"cloudfoundry",
				"--stemcell-name", "cool-ubuntu-animal",
				"--az", "eastprod-1",
				"--doppler-ip", "1.0.11.1",
				"--doppler-ip", "1.0.11.2",
				"--network", "foundry-net",
				"--doppler-vm-type", "blah",
				"--metron-secret", "metronsecret",
				"--metron-zone", "metronzoneguid",
				"--syslog-address", "syslog-server",
				"--syslog-port", "10601",
				"--syslog-transport", "tcp",
				"--etcd-machine-ip", "1.0.0.7",
				"--etcd-machine-ip", "1.0.0.8",
				"--doppler-zone", "dopplerzone",
				"--doppler-status-user", "doppleruser",
				"--doppler-status-password", "dopplerpassword",
				"--doppler-status-port", "5768",
				"--doppler-drain-buffer-size", "100",
				"--doppler-shared-secret", "secret",
				"--system-domain", "sys.test.com",
				"--cc-bulk-api-password", "bulk-pwd",
			})
			doppler = NewDopplerPartition(c)
		})
		It("then HasValidValues should be true", func() {
			Ω(doppler.HasValidValues()).Should(Equal(true))
		})
		It("then it should allow the user to configure the doppler IPs", func() {
			ig := doppler.ToInstanceGroup()
			network := ig.Networks[0]
			Ω(len(network.StaticIPs)).Should(Equal(2))
			Ω(network.StaticIPs).Should(ConsistOf("1.0.11.1", "1.0.11.2"))
		})
		It("then it should have 2 instances", func() {
			ig := doppler.ToInstanceGroup()
			Ω(ig.Instances).Should(Equal(2))
		})
		It("then it should configure the correct number of instances automatically from the count of given IPs", func() {
			ig := doppler.ToInstanceGroup()
			network := ig.Networks[0]
			Ω(len(network.StaticIPs)).Should(Equal(ig.Instances))
		})

		It("then it should allow the user to configure the AZs", func() {
			ig := doppler.ToInstanceGroup()
			Ω(len(ig.AZs)).Should(Equal(1))
			Ω(ig.AZs[0]).Should(Equal("eastprod-1"))
		})

		It("then it should allow the user to configure vm-type", func() {
			ig := doppler.ToInstanceGroup()
			Ω(ig.VMType).ShouldNot(BeEmpty())
			Ω(ig.VMType).Should(Equal("blah"))
		})

		It("then it should allow the user to configure network to use", func() {
			ig := doppler.ToInstanceGroup()
			network := ig.GetNetworkByName("foundry-net")
			Ω(network).ShouldNot(BeNil())
		})

		It("then it should allow the user to configure the used stemcell", func() {
			ig := doppler.ToInstanceGroup()
			Ω(ig.Stemcell).ShouldNot(BeEmpty())
			Ω(ig.Stemcell).Should(Equal("cool-ubuntu-animal"))
		})
		It("then it should have update max in-flight 1 and serial false", func() {
			ig := doppler.ToInstanceGroup()
			Ω(ig.Update.MaxInFlight).Should(Equal(1))
			Ω(ig.Update.Serial).Should(Equal(false))
		})

		XIt("then it should then have 4 jobs", func() {
			ig := doppler.ToInstanceGroup()
			Ω(len(ig.Jobs)).Should(Equal(4))
		})
		XIt("then it should then have x job", func() {
			ig := doppler.ToInstanceGroup()
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
