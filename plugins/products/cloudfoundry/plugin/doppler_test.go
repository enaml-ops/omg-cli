package cloudfoundry_test

import (
	"github.com/enaml-ops/omg-cli/plugins/products/cloudfoundry/enaml-gen/doppler"
	"github.com/enaml-ops/omg-cli/plugins/products/cloudfoundry/enaml-gen/syslog_drain_binder"
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
			dopplerPartition := NewDopplerPartition(c)
			Ω(dopplerPartition.HasValidValues()).Should(BeFalse())
		})
	})
	Context("when initialized WITH a complete set of arguments", func() {
		var dopplerPartition InstanceGrouper
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
				"--doppler-drain-buffer-size", "100",
				"--doppler-shared-secret", "secret",
				"--system-domain", "sys.test.com",
				"--cc-bulk-api-password", "bulk-pwd",
			})
			dopplerPartition = NewDopplerPartition(c)
		})
		It("then HasValidValues should be true", func() {
			Ω(dopplerPartition.HasValidValues()).Should(Equal(true))
		})
		It("then it should allow the user to configure the doppler IPs", func() {
			ig := dopplerPartition.ToInstanceGroup()
			network := ig.Networks[0]
			Ω(len(network.StaticIPs)).Should(Equal(2))
			Ω(network.StaticIPs).Should(ConsistOf("1.0.11.1", "1.0.11.2"))
		})
		It("then it should have 2 instances", func() {
			ig := dopplerPartition.ToInstanceGroup()
			Ω(ig.Instances).Should(Equal(2))
		})
		It("then it should configure the correct number of instances automatically from the count of given IPs", func() {
			ig := dopplerPartition.ToInstanceGroup()
			network := ig.Networks[0]
			Ω(len(network.StaticIPs)).Should(Equal(ig.Instances))
		})

		It("then it should allow the user to configure the AZs", func() {
			ig := dopplerPartition.ToInstanceGroup()
			Ω(len(ig.AZs)).Should(Equal(1))
			Ω(ig.AZs[0]).Should(Equal("eastprod-1"))
		})

		It("then it should allow the user to configure vm-type", func() {
			ig := dopplerPartition.ToInstanceGroup()
			Ω(ig.VMType).ShouldNot(BeEmpty())
			Ω(ig.VMType).Should(Equal("blah"))
		})

		It("then it should allow the user to configure network to use", func() {
			ig := dopplerPartition.ToInstanceGroup()
			network := ig.GetNetworkByName("foundry-net")
			Ω(network).ShouldNot(BeNil())
		})

		It("then it should allow the user to configure the used stemcell", func() {
			ig := dopplerPartition.ToInstanceGroup()
			Ω(ig.Stemcell).ShouldNot(BeEmpty())
			Ω(ig.Stemcell).Should(Equal("cool-ubuntu-animal"))
		})
		It("then it should have update max in-flight 1 and serial false", func() {
			ig := dopplerPartition.ToInstanceGroup()
			Ω(ig.Update.MaxInFlight).Should(Equal(1))
			Ω(ig.Update.Serial).Should(Equal(false))
		})

		It("then it should then have 4 jobs", func() {
			ig := dopplerPartition.ToInstanceGroup()
			Ω(len(ig.Jobs)).Should(Equal(4))
		})
		It("then it should then have doppler job", func() {
			ig := dopplerPartition.ToInstanceGroup()
			job := ig.GetJobByName("doppler")
			Ω(job).ShouldNot(BeNil())
			props, _ := job.Properties.(*doppler.DopplerJob)
			Ω(props.Doppler.Zone).Should(Equal("dopplerzone"))
			Ω(props.Doppler.MessageDrainBufferSize).Should(Equal(100))
			Ω(props.Loggregator.Etcd.Machines).Should(ConsistOf("1.0.0.7", "1.0.0.8"))
			Ω(props.DopplerEndpoint.SharedSecret).Should(Equal("secret"))
		})
		It("then it should then have syslog_drain_binder job", func() {
			ig := dopplerPartition.ToInstanceGroup()
			job := ig.GetJobByName("syslog_drain_binder")
			Ω(job).ShouldNot(BeNil())
			props, _ := job.Properties.(*syslog_drain_binder.SyslogDrainBinderJob)
			Ω(props.Ssl.SkipCertVerify).Should(Equal(true))
			Ω(props.SystemDomain).Should(Equal("sys.test.com"))
			Ω(props.Cc.BulkApiPassword).Should(Equal("bulk-pwd"))
			Ω(props.Cc.SrvApiUri).Should(Equal("https://api.sys.test.com"))
			Ω(props.Loggregator.Etcd.Machines).Should(ConsistOf("1.0.0.7", "1.0.0.8"))

		})
	})
})
