package cloudfoundry_test

import (
	"github.com/enaml-ops/omg-cli/plugins/products/cloudfoundry/enaml-gen/uaa"
	. "github.com/enaml-ops/omg-cli/plugins/products/cloudfoundry/plugin"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("UAA Partition", func() {
	Context("when initialized WITHOUT a complete set of arguments", func() {
		It("then HasValidValues returns false", func() {
			plugin := new(Plugin)
			c := plugin.GetContext([]string{
				"cloudfoundry",
			})
			uaaPartition := NewUAAPartition(c)
			Ω(uaaPartition.HasValidValues()).Should(BeFalse())
		})
	})
	Context("when initialized WITH a complete set of arguments", func() {
		var uaaPartition InstanceGrouper
		BeforeEach(func() {
			plugin := new(Plugin)
			c := plugin.GetContext([]string{
				"cloudfoundry",
				"--network", "foundry-net",
				"--stemcell-name", "cool-ubuntu-animal",
				"--az", "eastprod-1",
				"--network", "foundry-net",
				"--uaa-vm-type", "blah",
				"--uaa-instances", "1",
				"--system-domain", "sys.test.com",
				"--consul-ip", "1.0.0.1",
				"--consul-ip", "1.0.0.2",
				"--consul-encryption-key", "consulencryptionkey",
				"--consul-ca-cert", "consul-ca-cert",
				"--consul-agent-cert", "consul-agent-cert",
				"--consul-agent-key", "consul-agent-key",
				"--consul-server-cert", "consulservercert",
				"--consul-server-key", "consulserverkey",
				"--metron-secret", "metronsecret",
				"--metron-zone", "metronzoneguid",
				"--syslog-address", "syslog-server",
				"--syslog-port", "10601",
				"--syslog-transport", "tcp",
				"--etcd-machine-ip", "1.0.0.7",
				"--etcd-machine-ip", "1.0.0.8",
			})
			uaaPartition = NewUAAPartition(c)
		})
		It("then HasValidValues should return true", func() {
			Ω(uaaPartition.HasValidValues()).Should(Equal(true))
		})
		It("then it should not configure static ips for uaaPartition", func() {
			ig := uaaPartition.ToInstanceGroup()
			network := ig.Networks[0]
			Ω(len(network.StaticIPs)).Should(Equal(0))
		})
		It("then it should have 1 instances", func() {
			ig := uaaPartition.ToInstanceGroup()
			Ω(ig.Instances).Should(Equal(1))
		})
		It("then it should allow the user to configure the AZs", func() {
			ig := uaaPartition.ToInstanceGroup()
			Ω(len(ig.AZs)).Should(Equal(1))
			Ω(ig.AZs[0]).Should(Equal("eastprod-1"))
		})

		It("then it should allow the user to configure vm-type", func() {
			ig := uaaPartition.ToInstanceGroup()
			Ω(ig.VMType).ShouldNot(BeEmpty())
			Ω(ig.VMType).Should(Equal("blah"))
		})

		It("then it should allow the user to configure network to use", func() {
			ig := uaaPartition.ToInstanceGroup()
			network := ig.GetNetworkByName("foundry-net")
			Ω(network).ShouldNot(BeNil())
		})

		It("then it should allow the user to configure the used stemcell", func() {
			ig := uaaPartition.ToInstanceGroup()
			Ω(ig.Stemcell).ShouldNot(BeEmpty())
			Ω(ig.Stemcell).Should(Equal("cool-ubuntu-animal"))
		})
		It("then it should have update max in-flight 1 and serial false", func() {
			ig := uaaPartition.ToInstanceGroup()
			Ω(ig.Update.MaxInFlight).Should(Equal(1))
			Ω(ig.Update.Serial).Should(Equal(false))
		})

		It("then it should then have 5 jobs", func() {
			ig := uaaPartition.ToInstanceGroup()
			Ω(len(ig.Jobs)).Should(Equal(5))
		})
		It("then it should then have uaa job", func() {
			ig := uaaPartition.ToInstanceGroup()
			job := ig.GetJobByName("uaa")
			Ω(job).ShouldNot(BeNil())
			props, _ := job.Properties.(*uaa.Uaa)
			Ω(props.Login).ShouldNot(BeNil())

		})
	})
})
