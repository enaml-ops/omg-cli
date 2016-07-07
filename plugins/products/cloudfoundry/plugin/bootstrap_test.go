package cloudfoundry_test

import (
	"github.com/enaml-ops/enaml"
	"github.com/enaml-ops/omg-cli/plugins/products/cf-mysql/enaml-gen/bootstrap"
	. "github.com/enaml-ops/omg-cli/plugins/products/cloudfoundry/plugin"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("given a bootstrap partition", func() {
	Context("when initialized WITHOUT a complete set of arguments", func() {
		var ig InstanceGrouper
		BeforeEach(func() {
			p := new(Plugin)
			c := p.GetContext([]string{"cloudfoundry"})
			ig = NewBootstrapPartition(c)
		})

		It("should not be nil", func() {
			Ω(ig).ShouldNot(BeNil())
		})

		It("should contain a bootstrap job", func() {
			group := ig.ToInstanceGroup()
			Ω(group.GetJobByName("bootstrap")).ShouldNot(BeNil())
		})

		It("should not have valid values", func() {
			Ω(ig.HasValidValues()).Should(BeFalse())
		})
	})

	Context("when initialized with a complete set of arguments", func() {
		var ig InstanceGrouper
		var dm *enaml.DeploymentManifest
		BeforeEach(func() {
			p := new(Plugin)
			c := p.GetContext([]string{
				"cloudfoundry",
				"--az", "z1",
				"--stemcell-name", "cool-ubuntu-animal",
				"--network", "foundry-net",
				"--mysql-ip", "10.0.0.26",
				"--mysql-ip", "10.0.0.27",
				"--mysql-ip", "10.0.0.28",
				"--mysql-bootstrap-username", "user",
				"--mysql-bootstrap-password", "pass",
			})
			ig = NewBootstrapPartition(c)

			dm = new(enaml.DeploymentManifest)
			dm.AddInstanceGroup(ig.ToInstanceGroup())
		})

		It("should have valid values", func() {
			Ω(ig.HasValidValues()).Should(BeTrue())
		})

		It("should have the correct VM type and lifecycle", func() {
			group := dm.GetInstanceGroupByName("bootstrap")
			Ω(group.Lifecycle).Should(Equal("errand"))
			Ω(group.VMType).Should(Equal("errand"))
		})

		It("should have a single instance", func() {
			group := dm.GetInstanceGroupByName("bootstrap")
			Ω(group.Instances).Should(Equal(1))
		})

		It("should have update max in flight 1", func() {
			group := dm.GetInstanceGroupByName("bootstrap")
			Ω(group.Update.MaxInFlight).Should(Equal(1))
		})

		It("should allow the user to configure the AZs", func() {
			group := dm.GetInstanceGroupByName("bootstrap")
			Ω(len(group.AZs)).Should(Equal(1))
			Ω(group.AZs[0]).Should(Equal("z1"))
		})

		It("should allow the user to configure the used stemcell", func() {
			group := dm.GetInstanceGroupByName("bootstrap")
			Ω(group.Stemcell).Should(Equal("cool-ubuntu-animal"))
		})

		It("should allow the user to configure the network to use", func() {
			group := dm.GetInstanceGroupByName("bootstrap")
			Ω(len(group.Networks)).Should(Equal(1))
			Ω(group.Networks[0].Name).Should(Equal("foundry-net"))
		})

		It("should have a valid bootstrap job", func() {
			group := dm.GetInstanceGroupByName("bootstrap")
			job := group.GetJobByName("bootstrap")
			Ω(job.Release).Should(Equal(CFMysqlReleaseName))

			props := job.Properties.(*bootstrap.Bootstrap)
			Ω(props.ClusterIps).Should(ConsistOf("10.0.0.26", "10.0.0.27", "10.0.0.28"))
			Ω(props.DatabaseStartupTimeout).Should(Equal(1200))
			Ω(props.BootstrapEndpoint.Username).Should(Equal("user"))
			Ω(props.BootstrapEndpoint.Password).Should(Equal("pass"))
		})

	})
})
