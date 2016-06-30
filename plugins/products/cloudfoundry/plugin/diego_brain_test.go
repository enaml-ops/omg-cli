package cloudfoundry_test

import (
	"github.com/enaml-ops/enaml"
	"github.com/enaml-ops/omg-cli/plugins/products/cloudfoundry/enaml-gen/auctioneer"
	. "github.com/enaml-ops/omg-cli/plugins/products/cloudfoundry/plugin"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("given a Diego Brain Partition", func() {
	Context("when initialized WITHOUT a complete set of arguments", func() {
		var ig InstanceGrouper
		BeforeEach(func() {
			cf := new(Plugin)
			c := cf.GetContext([]string{""})
			ig = NewDiegoBrainPartition(c)
		})

		It("then it should contain the appropriate jobs", func() {
			group := ig.ToInstanceGroup()
			Ω(group.GetJobByName("auctioneer")).ShouldNot(BeNil())
			Ω(group.GetJobByName("cc_uploader")).ShouldNot(BeNil())
			Ω(group.GetJobByName("converger")).ShouldNot(BeNil())
			Ω(group.GetJobByName("file_server")).ShouldNot(BeNil())
			Ω(group.GetJobByName("nsync")).ShouldNot(BeNil())
			Ω(group.GetJobByName("route_emitter")).ShouldNot(BeNil())
			Ω(group.GetJobByName("ssh_proxy")).ShouldNot(BeNil())
			Ω(group.GetJobByName("stager")).ShouldNot(BeNil())
			Ω(group.GetJobByName("tps")).ShouldNot(BeNil())
			Ω(group.GetJobByName("consul_agent")).ShouldNot(BeNil())
			Ω(group.GetJobByName("metron_agent")).ShouldNot(BeNil())
			Ω(group.GetJobByName("statsd-injector")).ShouldNot(BeNil())
		})

		It("then it should not validate", func() {
			Ω(ig.HasValidValues()).Should(BeFalse())
		})
	})

	Context("when initialized with a complete set of arguments", func() {
		var deploymentManifest *enaml.DeploymentManifest
		var grouper InstanceGrouper

		BeforeEach(func() {
			cf := new(Plugin)
			c := cf.GetContext([]string{
				"cloudfoundry",
				"--az", "eastprod-1",
				"--stemcell-name", "cool-ubuntu-animal",
				"--network", "foundry-net",
				"--diego-brain-ip", "10.0.0.39",
				"--diego-brain-ip", "10.0.0.40",
				"--diego-brain-vm-type", "brainvmtype",
				"--diego-brain-disk-type", "braindisktype",
				"--auctioneer-ca-cert", "cacert",
				"--auctioneer-client-cert", "clientcert",
				"--auctioneer-client-key", "clientkey",
			})
			grouper = NewDiegoBrainPartition(c)
			deploymentManifest = new(enaml.DeploymentManifest)
			deploymentManifest.AddInstanceGroup(grouper.ToInstanceGroup())
		})

		It("then it should validate", func() {
			Ω(grouper.HasValidValues()).Should(BeTrue())
		})

		It("then it should allow the user to configure the AZs", func() {
			ig := deploymentManifest.GetInstanceGroupByName("diego_brain-partition")
			Ω(len(ig.AZs)).Should(Equal(1))
			Ω(ig.AZs[0]).Should(Equal("eastprod-1"))
		})

		It("then it should allow the user to configure the used stemcell", func() {
			ig := deploymentManifest.GetInstanceGroupByName("diego_brain-partition")
			Ω(ig.Stemcell).Should(Equal("cool-ubuntu-animal"))
		})

		It("then it should allow the user to configure the network to use", func() {
			ig := deploymentManifest.GetInstanceGroupByName("diego_brain-partition")
			network := ig.GetNetworkByName("foundry-net")
			Ω(network).ShouldNot(BeNil())
			Ω(len(network.StaticIPs)).Should(Equal(2))
			Ω(network.StaticIPs[0]).Should(Equal("10.0.0.39"))
			Ω(network.StaticIPs[1]).Should(Equal("10.0.0.40"))
		})

		It("then it should allow the user to configure the VM type", func() {
			ig := deploymentManifest.GetInstanceGroupByName("diego_brain-partition")
			Ω(ig.VMType).Should(Equal("brainvmtype"))
		})

		It("then it should allow the user to configure the disk type", func() {
			ig := deploymentManifest.GetInstanceGroupByName("diego_brain-partition")
			Ω(ig.PersistentDiskType).Should(Equal("braindisktype"))
		})

		It("then it should configure the correct number of instances automatically from the count of IPs", func() {
			ig := deploymentManifest.GetInstanceGroupByName("diego_brain-partition")
			Ω(len(ig.Networks)).Should(Equal(1))
			Ω(len(ig.Networks[0].StaticIPs)).Should(Equal(ig.Instances))
		})

		It("then it should have update max-in-flight 1", func() {
			ig := deploymentManifest.GetInstanceGroupByName("diego_brain-partition")
			Ω(ig.Update.MaxInFlight).Should(Equal(1))
		})

		It("then it should allow the user to configure the auctioneer SSL", func() {
			ig := deploymentManifest.GetInstanceGroupByName("diego_brain-partition")
			job := ig.GetJobByName("auctioneer")
			a := job.Properties.(*auctioneer.Auctioneer)
			Ω(a.Bbs.CaCert).Should(Equal("cacert"))
			Ω(a.Bbs.ClientCert).Should(Equal("clientcert"))
			Ω(a.Bbs.ClientKey).Should(Equal("clientkey"))
		})
	})
})
