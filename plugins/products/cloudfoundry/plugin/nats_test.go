package cloudfoundry_test

import (
	. "github.com/enaml-ops/omg-cli/plugins/products/cloudfoundry/plugin"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	//"gopkg.in/yaml.v2"
	//"fmt"
)

var _ = Describe("Nats Partition", func() {

	Context("when initialized WITHOUT a complete set of arguments", func() {
		It("should return the error and exit", func() {
			plugin := new(Plugin)
			c := plugin.GetContext([]string{
				"cloudfoundry",
				"--az", "eastprod-1",
			})
			_, err := NewNatsPartition(c)
			Ω(err).ShouldNot(BeNil())
		})
	})

	Context("when initialized WITH a complete set of arguments", func() {
		var err error
		 var natsPartition InstanceGrouper

		BeforeEach(func() {
			plugin := new(Plugin)
			c := plugin.GetContext([]string{
				"cloudfoundry",
				"--stemcell-name", "trusty",
				"--az", "eastprod-1",
				"--nats-machine-ip", "10.0.0.2",
				"--nats-machine-ip", "10.0.0.3",
				"--nats-network", "foundry-net",
				"--nats-vm-type", "blah",
				"--metron-secret", "metronsecret",
				"--metron-zone", "metronzoneguid",
				"--etcd-machine-ip", "10.0.0.7",
				"--etcd-machine-ip", "10.0.0.8",
			})
			natsPartition, err = NewNatsPartition(c)
		})
		It("then it should not return an error", func() {
			Ω(err).Should(BeNil())
		})
		It("should have 2 instances ", func() {
			igf := natsPartition.ToInstanceGroup()
			Ω(igf.Instances).Should(Equal(2))
		})
		It("should have the IP ranges set correctly", func() {
			igf := natsPartition.ToInstanceGroup()
			//b, _ := yaml.Marshal(igf)
			//fmt.Print(string(b))
			networks := igf.Networks
			Ω(len(networks)).Should(Equal(1))
			Ω(len(networks[0].StaticIPs)).Should(Equal(2))
			Ω(networks[0].StaticIPs).Should(ConsistOf("10.0.0.2", "10.0.0.3"))
		})
		It("then it should allow the user to configure the AZs", func() {
			ig := natsPartition.ToInstanceGroup()
			Ω(len(ig.AZs)).Should(Equal(1))
			Ω(ig.AZs[0]).Should(Equal("eastprod-1"))
		})
		It("then it should allow the user to configure network to use", func() {
			ig := natsPartition.ToInstanceGroup()
			network := ig.GetNetworkByName("foundry-net")
			Ω(network).ShouldNot(BeNil())
		})

		It("then it should allow the user to configure the used stemcell", func() {
			ig := natsPartition.ToInstanceGroup()
			Ω(ig.Stemcell).ShouldNot(BeEmpty())
			Ω(ig.Stemcell).Should(Equal("trusty"))
		})
		It("then it should have update max in-flight 1 and serial", func() {
			ig := natsPartition.ToInstanceGroup()
			Ω(ig.Update.MaxInFlight).Should(Equal(1))
			Ω(ig.Update.Serial).Should(Equal(true))
		})
		It("then it should then have 3 jobs", func() {
			ig := natsPartition.ToInstanceGroup()
			Ω(len(ig.Jobs)).Should(Equal(3))
		})
	})
})
