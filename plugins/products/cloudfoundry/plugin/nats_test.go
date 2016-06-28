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
		// var natsPartition InstanceGroupFactory

		BeforeEach(func() {
			plugin := new(Plugin)
			c := plugin.GetContext([]string{
				"cloudfoundry",
				"--stemcell-name", "trusty",
				"--az", "eastprod-1",
				"--nats-machine-ip", "1.0.0.2",
				"--nats-network", "foundry-net",
				"--nats-vm-type", "blah",
				"--metron-secret", "metronsecret",
				"--metron-zone", "metronzoneguid",
				"--etcd-machine-ip", "1.0.0.7",
				"--etcd-machine-ip", "1.0.0.8",
			})
			_, err = NewNatsPartition(c)
		})
		It("then it should not return an error", func() {
			// igf := natsPartition.ToInstanceGroup()
			//b, _ := yaml.Marshal(igf)
			//fmt.Print(string(b))
			Ω(err).Should(BeNil())
		})
	})
})
