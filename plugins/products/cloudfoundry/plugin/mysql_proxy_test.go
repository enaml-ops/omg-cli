package cloudfoundry_test

import (
	"github.com/enaml-ops/omg-cli/plugins/products/cf-mysql/enaml-gen/proxy"
	. "github.com/enaml-ops/omg-cli/plugins/products/cloudfoundry/plugin"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("MySQL Proxy Partition", func() {
	Context("when initialized WITHOUT a complete set of arguments", func() {
		It("then it should return the error and exit", func() {
			plugin := new(Plugin)
			c := plugin.GetContext([]string{
				"cloudfoundry",
			})
			mySQLProxy := NewMySQLProxyPartition(c)
			Ω(mySQLProxy).ShouldNot(BeNil())
			Ω(mySQLProxy.HasValidValues()).Should(BeFalse())
		})
	})
	Context("when initialized WITH a complete set of arguments", func() {
		var mysqlProxyPartition InstanceGrouper
		BeforeEach(func() {
			plugin := new(Plugin)
			c := plugin.GetContext([]string{
				"cloudfoundry",
				"--stemcell-name", "cool-ubuntu-animal",
				"--az", "eastprod-1",
				"--mysql-ip", "1.0.10.1",
				"--mysql-ip", "1.0.10.2",
				"--mysql-proxy-ip", "1.0.10.3",
				"--mysql-proxy-ip", "1.0.10.4",
				"--network", "foundry-net",
				"--mysql-proxy-vm-type", "blah",
				"--mysql-proxy-external-host", "mysqlhostname",
				"--mysql-proxy-api-username", "apiuser",
				"--mysql-proxy-api-password", "apipassword",
				"--syslog-address", "syslog-server",
				"--syslog-port", "10601",
				"--syslog-transport", "tcp",
				"--nats-user", "nats",
				"--nats-pass", "pass",
				"--nats-machine-ip", "1.0.0.5",
				"--nats-machine-ip", "1.0.0.6",
			})
			mysqlProxyPartition = NewMySQLProxyPartition(c)
		})
		It("then it should have all valid values", func() {
			Ω(mysqlProxyPartition.HasValidValues()).Should(BeTrue())
		})
		It("then it should allow the user to configure the mysql proxy IPs", func() {
			ig := mysqlProxyPartition.ToInstanceGroup()
			Ω(len(ig.Networks)).Should(Equal(1))
			network := ig.Networks[0]
			Ω(len(network.StaticIPs)).Should(Equal(2))
			Ω(network.StaticIPs).Should(ConsistOf("1.0.10.3", "1.0.10.4"))
		})
		It("then it should have 2 instances", func() {
			ig := mysqlProxyPartition.ToInstanceGroup()
			Ω(ig.Instances).Should(Equal(2))
		})
		It("then it should configure the correct number of instances automatically from the count of given IPs", func() {
			ig := mysqlProxyPartition.ToInstanceGroup()
			network := ig.Networks[0]
			Ω(len(network.StaticIPs)).Should(Equal(ig.Instances))
		})

		It("then it should allow the user to configure the AZs", func() {
			ig := mysqlProxyPartition.ToInstanceGroup()
			Ω(len(ig.AZs)).Should(Equal(1))
			Ω(ig.AZs[0]).Should(Equal("eastprod-1"))
		})

		It("then it should allow the user to configure vm-type", func() {
			ig := mysqlProxyPartition.ToInstanceGroup()
			Ω(ig.VMType).ShouldNot(BeEmpty())
			Ω(ig.VMType).Should(Equal("blah"))
		})

		It("then it should allow the user to configure network to use", func() {
			ig := mysqlProxyPartition.ToInstanceGroup()
			network := ig.GetNetworkByName("foundry-net")
			Ω(network).ShouldNot(BeNil())
		})

		It("then it should allow the user to configure the used stemcell", func() {
			ig := mysqlProxyPartition.ToInstanceGroup()
			Ω(ig.Stemcell).ShouldNot(BeEmpty())
			Ω(ig.Stemcell).Should(Equal("cool-ubuntu-animal"))
		})
		It("then it should have update max in-flight 1 and serial", func() {
			ig := mysqlProxyPartition.ToInstanceGroup()
			Ω(ig.Update.MaxInFlight).Should(Equal(1))
			Ω(ig.Update.Serial).Should(Equal(false))
		})

		It("then it should then have 1 job", func() {
			ig := mysqlProxyPartition.ToInstanceGroup()
			Ω(len(ig.Jobs)).Should(Equal(1))
		})
		It("then it should then have proxy job", func() {
			ig := mysqlProxyPartition.ToInstanceGroup()
			job := ig.GetJobByName("proxy")
			Ω(job).ShouldNot(BeNil())
			Ω(job.Release).Should(Equal("cf-mysql"))
			props, _ := job.Properties.(*proxy.Proxy)
			Ω(props.ApiUsername).Should(Equal("apiuser"))
			Ω(props.ApiPassword).Should(Equal("apipassword"))
			Ω(props.ExternalHost).Should(Equal("mysqlhostname"))
			Ω(props.ClusterIps).Should(ConsistOf("1.0.10.1", "1.0.10.2"))
			Ω(props.SyslogAggregator.Address).Should(Equal("syslog-server"))
			Ω(props.SyslogAggregator.Port).Should(Equal(10601))
			Ω(props.SyslogAggregator.Transport).Should(Equal("tcp"))
			Ω(props.Nats).ShouldNot(BeNil())
			Ω(props.Nats.Port).Should(Equal(4222))
			Ω(props.Nats.User).Should(Equal("nats"))
			Ω(props.Nats.Password).Should(Equal("pass"))
			Ω(props.Nats.Machines).Should(ConsistOf("1.0.0.5", "1.0.0.6"))
		})

	})
})
