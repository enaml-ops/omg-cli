package cloudfoundry_test

import (
	"github.com/codegangsta/cli"
	"github.com/enaml-ops/omg-cli/plugins/products/cloudfoundry/enaml-gen/consul_agent"
	. "github.com/enaml-ops/omg-cli/plugins/products/cloudfoundry/plugin"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Consul Agent", func() {
	Context("when initialized WITHOUT a complete set of arguments", func() {
		It("then hasValidValues should return false", func() {
			plugin := new(Plugin)
			c := plugin.GetContext([]string{
				"cloudfoundry",
			})
			consulAgent := NewConsulAgent(c, []string{})
			Ω(consulAgent.HasValidValues()).Should(BeFalse())
		})
	})
	Context("when initialized WITH a complete set of arguments", func() {

		var c *cli.Context
		BeforeEach(func() {
			plugin := new(Plugin)
			c = plugin.GetContext([]string{
				"cloudfoundry",
				"--consul-ip", "1.0.0.1",
				"--consul-ip", "1.0.0.2",
				"--consul-encryption-key", "encyption-key",
				"--consul-ca-cert", "ca-cert",
				"--consul-agent-cert", "agent-cert",
				"--consul-agent-key", "agent-key",
				"--consul-server-cert", "server-cert",
				"--consul-server-key", "server-key",
			})

		})
		It("then hasValidValues should return true for consul with server false", func() {
			consulAgent := NewConsulAgent(c, []string{})
			Ω(consulAgent.HasValidValues()).Should(BeTrue())
		})
		It("then hasValidValues should return true for consul with server true", func() {
			consulAgent := NewConsulAgentServer(c)
			Ω(consulAgent.HasValidValues()).Should(BeTrue())
		})
		It("then job properties are set properly for server false", func() {
			consulAgent := NewConsulAgent(c, []string{})
			job := consulAgent.CreateJob()
			Ω(job).ShouldNot(BeNil())
			props := job.Properties.(*consul_agent.ConsulAgentJob)
			Ω(props.Consul.Agent.Servers.Lan).Should(ConsistOf("1.0.0.1", "1.0.0.2"))
			Ω(props.Consul.AgentCert).Should(Equal("agent-cert"))
			Ω(props.Consul.AgentKey).Should(Equal("agent-key"))
			Ω(props.Consul.ServerCert).Should(Equal("server-cert"))
			Ω(props.Consul.ServerKey).Should(Equal("server-key"))
			Ω(props.Consul.EncryptKeys).Should(ConsistOf("encyption-key"))
			Ω(props.Consul.Agent.Domain).Should(Equal("cf.internal"))
			Ω(props.Consul.Agent.Mode).Should(BeNil())
		})
		It("then job properties are set properly etcd service", func() {
			consulAgent := NewConsulAgent(c, []string{"etcd"})
			job := consulAgent.CreateJob()
			Ω(job).ShouldNot(BeNil())
			props := job.Properties.(*consul_agent.ConsulAgentJob)
			etcdMap := make(map[string]map[string]string)
			etcdMap["etcd"] = make(map[string]string)
			Ω(props.Consul.Agent.Services).Should(Equal(etcdMap))
		})
		It("then job properties are set properly etcd and uaa service", func() {
			consulAgent := NewConsulAgent(c, []string{"etcd", "uaa"})
			job := consulAgent.CreateJob()
			Ω(job).ShouldNot(BeNil())
			props := job.Properties.(*consul_agent.ConsulAgentJob)
			servicesMap := make(map[string]map[string]string)
			servicesMap["etcd"] = make(map[string]string)
			servicesMap["uaa"] = make(map[string]string)
			Ω(props.Consul.Agent.Services).Should(Equal(servicesMap))
		})
		It("then job properties are set properly for server true", func() {
			consulAgent := NewConsulAgentServer(c)
			job := consulAgent.CreateJob()
			Ω(job).ShouldNot(BeNil())
			props := job.Properties.(*consul_agent.ConsulAgentJob)
			Ω(props.Consul.Agent.Servers.Lan).Should(ConsistOf("1.0.0.1", "1.0.0.2"))
			Ω(props.Consul.AgentCert).Should(Equal("agent-cert"))
			Ω(props.Consul.AgentKey).Should(Equal("agent-key"))
			Ω(props.Consul.ServerCert).Should(Equal("server-cert"))
			Ω(props.Consul.ServerKey).Should(Equal("server-key"))
			Ω(props.Consul.EncryptKeys).Should(ConsistOf("encyption-key"))
			Ω(props.Consul.Agent.Domain).Should(Equal("cf.internal"))
			Ω(props.Consul.Agent.Mode).Should(Equal("server"))
		})

	})
})
