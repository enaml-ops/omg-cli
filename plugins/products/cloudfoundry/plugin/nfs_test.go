package cloudfoundry_test

import (
	"github.com/enaml-ops/omg-cli/plugins/products/cloudfoundry/enaml-gen/debian_nfs_server"
	. "github.com/enaml-ops/omg-cli/plugins/products/cloudfoundry/plugin"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("NFS Partition", func() {
	Context("when initialized WITHOUT a complete set of arguments", func() {
		It("then it should return the error and exit", func() {
			plugin := new(Plugin)
			c := plugin.GetContext([]string{
				"cloudfoundry",
				"--metron-secret", "metronsecret",
				"--metron-zone", "metronzoneguid",
				"--syslog-address", "syslog-server",
				"--syslog-port", "10601",
				"--syslog-transport", "tcp",
				"--etcd-machine-ip", "1.0.0.7",
				"--etcd-machine-ip", "1.0.0.8",
			})
			Ω(NewNFSPartition(c).HasValidValues()).Should(Equal(false))
		})
	})
	Context("when initialized WITH a complete set of arguments", func() {
		var err error
		var nfsPartition InstanceGrouper
		BeforeEach(func() {
			plugin := new(Plugin)
			c := plugin.GetContext([]string{
				"cloudfoundry",
				"--stemcell-name", "cool-ubuntu-animal",
				"--az", "eastprod-1",
				"--nfs-ip", "1.0.0.1",
				"--nfs-network", "foundry-net",
				"--nfs-vm-type", "blah",
				"--nfs-disk-type", "blah-disk",
				"--nfs-allow-from-network-cidr", "1.0.0.0/22",
				"--metron-secret", "metronsecret",
				"--metron-zone", "metronzoneguid",
				"--syslog-address", "syslog-server",
				"--syslog-port", "10601",
				"--syslog-transport", "tcp",
				"--etcd-machine-ip", "1.0.0.7",
				"--etcd-machine-ip", "1.0.0.8",
			})
			nfsPartition = NewNFSPartition(c)
		})
		It("then it should not return an error", func() {
			Ω(err).Should(BeNil())
		})
		It("then it should not return an error", func() {
			Ω(nfsPartition.ToInstanceGroup().Name).Should(Equal("nfs_server-partition"))
		})

		It("then it should allow the user to configure the nfs IP", func() {
			ig := nfsPartition.ToInstanceGroup()
			network := ig.Networks[0]
			Ω(len(network.StaticIPs)).Should(Equal(1))
			Ω(network.StaticIPs).Should(ConsistOf("1.0.0.1"))
		})
		It("then it should have 1 instances", func() {
			ig := nfsPartition.ToInstanceGroup()
			Ω(ig.Instances).Should(Equal(1))
		})
		It("then it should configure the correct number of instances automatically from the count of given IPs", func() {
			ig := nfsPartition.ToInstanceGroup()
			network := ig.Networks[0]
			Ω(len(network.StaticIPs)).Should(Equal(ig.Instances))
		})

		It("then it should allow the user to configure the AZs", func() {
			ig := nfsPartition.ToInstanceGroup()
			Ω(len(ig.AZs)).Should(Equal(1))
			Ω(ig.AZs[0]).Should(Equal("eastprod-1"))
		})

		It("then it should allow the user to configure disk-type", func() {
			ig := nfsPartition.ToInstanceGroup()
			Ω(ig.PersistentDiskType).ShouldNot(BeEmpty())
			Ω(ig.PersistentDiskType).Should(Equal("blah-disk"))
		})
		It("then it should allow the user to configure vm-type", func() {
			ig := nfsPartition.ToInstanceGroup()
			Ω(ig.VMType).ShouldNot(BeEmpty())
			Ω(ig.VMType).Should(Equal("blah"))
		})

		It("then it should allow the user to configure network to use", func() {
			ig := nfsPartition.ToInstanceGroup()
			network := ig.GetNetworkByName("foundry-net")
			Ω(network).ShouldNot(BeNil())
		})

		It("then it should allow the user to configure the used stemcell", func() {
			ig := nfsPartition.ToInstanceGroup()
			Ω(ig.Stemcell).ShouldNot(BeEmpty())
			Ω(ig.Stemcell).Should(Equal("cool-ubuntu-animal"))
		})
		It("then it should have update max in-flight 1 and serial", func() {
			ig := nfsPartition.ToInstanceGroup()
			Ω(ig.Update.MaxInFlight).Should(Equal(1))
			Ω(ig.Update.Serial).Should(Equal(true))
		})

		It("then it should then have 3 jobs", func() {
			ig := nfsPartition.ToInstanceGroup()
			Ω(len(ig.Jobs)).Should(Equal(3))
		})
		It("then it should then have debian_nfs_server job", func() {
			ig := nfsPartition.ToInstanceGroup()
			job := ig.GetJobByName("debian_nfs_server")
			Ω(job).ShouldNot(BeNil())
			props, _ := job.Properties.(*debian_nfs_server.NfsServer)
			Ω(props.AllowFromEntries).Should(ConsistOf("1.0.0.0/22"))
		})

	})
})
