package cloudfoundry_test

import (
	"github.com/enaml-ops/omg-cli/plugins/products/cloudfoundry/enaml-gen/smoke-tests"
	. "github.com/enaml-ops/omg-cli/plugins/products/cloudfoundry/plugin"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Consul Partition", func() {
	Context("when initialized WITHOUT a complete set of arguments", func() {
		It("then HasValidValues should return false", func() {
			plugin := new(Plugin)
			c := plugin.GetContext([]string{
				"cloudfoundry",
				"--network", "foundry-net",
				"--app-domain", "apps.test.com",
			})
			smoke := NewSmokeErrand(c)
			Ω(smoke.HasValidValues()).Should(BeFalse())
		})
	})
	Context("when initialized WITH a complete set of arguments", func() {
		var smokeErrand InstanceGrouper
		BeforeEach(func() {
			plugin := new(Plugin)
			c := plugin.GetContext([]string{
				"cloudfoundry",
				"--stemcell-name", "cool-ubuntu-animal",
				"--az", "eastprod-1",
				"--network", "foundry-net",
				"--errand-vm-type", "blah",
				"--uaa-login-protocol", "https",
				"--smoke-tests-password", "password",
				"--system-domain", "sys.test.com",
				"--app-domain", "apps.test.com",
			})
			smokeErrand = NewSmokeErrand(c)
		})
		It("then HasValidValues should be true", func() {
			Ω(smokeErrand.HasValidValues()).Should(BeTrue())
		})
		It("then it should have 1 instances", func() {
			ig := smokeErrand.ToInstanceGroup()
			Ω(ig.Instances).Should(Equal(1))
		})
		It("then it should allow the user to configure the AZs", func() {
			ig := smokeErrand.ToInstanceGroup()
			Ω(len(ig.AZs)).Should(Equal(1))
			Ω(ig.AZs[0]).Should(Equal("eastprod-1"))
		})

		It("then it should allow the user to configure vm-type", func() {
			ig := smokeErrand.ToInstanceGroup()
			Ω(ig.VMType).ShouldNot(BeEmpty())
			Ω(ig.VMType).Should(Equal("blah"))
		})

		It("then it should allow the user to configure network to use", func() {
			ig := smokeErrand.ToInstanceGroup()
			network := ig.GetNetworkByName("foundry-net")
			Ω(network).ShouldNot(BeNil())
		})

		It("then it should allow the user to configure the used stemcell", func() {
			ig := smokeErrand.ToInstanceGroup()
			Ω(ig.Stemcell).ShouldNot(BeEmpty())
			Ω(ig.Stemcell).Should(Equal("cool-ubuntu-animal"))
		})
		It("then it should have update max in-flight 1", func() {
			ig := smokeErrand.ToInstanceGroup()
			Ω(ig.Update.MaxInFlight).Should(Equal(1))
			Ω(ig.Update.Serial).Should(Equal(false))
		})
		It("then it should have lifecycle of errand", func() {
			ig := smokeErrand.ToInstanceGroup()
			Ω(ig.Lifecycle).Should(Equal("errand"))
		})

		It("then it should then have 1 jobs", func() {
			ig := smokeErrand.ToInstanceGroup()
			Ω(len(ig.Jobs)).Should(Equal(1))
		})
		It("then it should then have smoke-tests job", func() {
			ig := smokeErrand.ToInstanceGroup()
			job := ig.GetJobByName("smoke-tests")
			Ω(job).ShouldNot(BeNil())
			props, _ := job.Properties.(*smoke_tests.SmokeTests)
			Ω(props.AppsDomain).Should(Equal("apps.test.com"))
			Ω(props.Api).Should(Equal("https://api.sys.test.com"))
			Ω(props.Org).Should(Equal("CF_SMOKE_TEST_ORG"))
			Ω(props.Space).Should(Equal("CF_SMOKE_TEST_SPACE"))
			Ω(props.User).Should(Equal("smoke_tests"))
			Ω(props.Password).Should(Equal("password"))
			Ω(props.UseExistingOrg).Should(BeFalse())
			Ω(props.UseExistingSpace).Should(BeFalse())
		})
	})
})
