package cloudfoundry_test

import (
	"github.com/enaml-ops/omg-cli/plugins/products/cloudfoundry/enaml-gen/etcd"
	"github.com/enaml-ops/omg-cli/plugins/products/cloudfoundry/enaml-gen/etcd_metrics_server"
	. "github.com/enaml-ops/omg-cli/plugins/products/cloudfoundry/plugin"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Etcd Partition", func() {
	Context("when initialized WITHOUT a complete set of arguments", func() {
		It("HasValidValues should return false", func() {
			plugin := new(Plugin)
			c := plugin.GetContext([]string{
				"cloudfoundry",
				"--metron-secret", "metronsecret",
				"--metron-zone", "metronzoneguid",
				"--syslog-address", "syslog-server",
				"--syslog-port", "10601",
				"--syslog-transport", "tcp",
				"--nats-user", "nats",
				"--nats-pass", "pass",
				"--nats-machine-ip", "1.0.0.5",
				"--nats-machine-ip", "1.0.0.6",
			})
			Ω(NewEtcdPartition(c).HasValidValues()).Should(Equal(false))
		})
	})
	Context("when initialized WITH a complete set of arguments", func() {
		var etcdPartition InstanceGrouper
		BeforeEach(func() {
			plugin := new(Plugin)
			c := plugin.GetContext([]string{
				"cloudfoundry",
				"--stemcell-name", "cool-ubuntu-animal",
				"--az", "eastprod-1",
				"--etcd-machine-ip", "1.0.0.7",
				"--etcd-machine-ip", "1.0.0.8",
				"--etcd-network", "foundry-net",
				"--etcd-vm-type", "blah",
				"--etcd-disk-type", "blah-disk",
				"--metron-secret", "metronsecret",
				"--metron-zone", "metronzoneguid",
				"--syslog-address", "syslog-server",
				"--syslog-port", "10601",
				"--syslog-transport", "tcp",
				"--nats-user", "nats",
				"--nats-pass", "pass",
				"--nats-machine-ip", "1.0.0.5",
				"--nats-machine-ip", "1.0.0.6",
			})
			etcdPartition = NewEtcdPartition(c)
		})
		It("HasValidValues should return true", func() {
			Ω(etcdPartition.HasValidValues()).Should(Equal(true))
		})
		It("then it should allow the user to configure the etcd IPs", func() {
			ig := etcdPartition.ToInstanceGroup()
			network := ig.Networks[0]
			Ω(len(network.StaticIPs)).Should(Equal(2))
			Ω(network.StaticIPs).Should(ConsistOf("1.0.0.7", "1.0.0.8"))
		})
		It("then it should have 2 instances", func() {
			ig := etcdPartition.ToInstanceGroup()
			Ω(ig.Instances).Should(Equal(2))
		})
		It("then it should configure the correct number of instances automatically from the count of given IPs", func() {
			ig := etcdPartition.ToInstanceGroup()
			network := ig.Networks[0]
			Ω(len(network.StaticIPs)).Should(Equal(ig.Instances))
		})

		It("then it should allow the user to configure the AZs", func() {
			ig := etcdPartition.ToInstanceGroup()
			Ω(len(ig.AZs)).Should(Equal(1))
			Ω(ig.AZs[0]).Should(Equal("eastprod-1"))
		})

		It("then it should allow the user to configure vm-type", func() {
			ig := etcdPartition.ToInstanceGroup()
			Ω(ig.VMType).ShouldNot(BeEmpty())
			Ω(ig.VMType).Should(Equal("blah"))
		})

		It("then it should allow the user to configure network to use", func() {
			ig := etcdPartition.ToInstanceGroup()
			network := ig.GetNetworkByName("foundry-net")
			Ω(network).ShouldNot(BeNil())
		})

		It("then it should allow the user to configure the used stemcell", func() {
			ig := etcdPartition.ToInstanceGroup()
			Ω(ig.Stemcell).ShouldNot(BeEmpty())
			Ω(ig.Stemcell).Should(Equal("cool-ubuntu-animal"))
		})

		It("then it should allow the user to configure disk type to use", func() {
			ig := etcdPartition.ToInstanceGroup()
			Ω(ig.PersistentDiskType).Should(Equal("blah-disk"))
		})
		It("then it should have update max in-flight 1 and serial", func() {
			ig := etcdPartition.ToInstanceGroup()
			Ω(ig.Update.MaxInFlight).Should(Equal(1))
			Ω(ig.Update.Serial).Should(Equal(false))
		})
		It("then it should then have 4 jobs", func() {
			ig := etcdPartition.ToInstanceGroup()
			Ω(len(ig.Jobs)).Should(Equal(4))
		})
		It("then it should then have etcd job", func() {
			ig := etcdPartition.ToInstanceGroup()
			job := ig.GetJobByName("etcd")
			Ω(job).ShouldNot(BeNil())
			props, _ := job.Properties.(*etcd.Etcd)
			Ω(props.RequireSsl).Should(Equal(false))
			Ω(props.PeerRequireSsl).Should(Equal(false))
			Ω(props.Machines).Should(ConsistOf("1.0.0.7", "1.0.0.8"))
		})
		It("then it should then have etcd_metrics_server job", func() {
			ig := etcdPartition.ToInstanceGroup()
			job := ig.GetJobByName("etcd_metrics_server")
			Ω(job).ShouldNot(BeNil())
			props, _ := job.Properties.(*etcd_metrics_server.EtcdMetricsServer)
			Ω(props.Nats.Machines).Should(ConsistOf("1.0.0.5", "1.0.0.6"))
			Ω(props.Nats.Username).Should(Equal("nats"))
			Ω(props.Nats.Password).Should(Equal("pass"))
		})
	})
})
