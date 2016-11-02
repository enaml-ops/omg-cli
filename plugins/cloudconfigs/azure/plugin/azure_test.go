package plugin_test

import (
	"io/ioutil"

	"github.com/enaml-ops/enaml"
	"github.com/enaml-ops/omg-cli/plugins/cloudconfigs"
	. "github.com/enaml-ops/omg-cli/plugins/cloudconfigs/azure/plugin"
	"github.com/enaml-ops/omg-cli/plugins/cloudconfigs/utils"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gopkg.in/yaml.v2"
)

var _ = Describe("given AzureCloud Config", func() {
	Context("when calling CreateManifest", func() {
		var provider cloudconfigs.CloudConfigProvider
		var manifest *enaml.CloudConfigManifest
		var err error
		BeforeEach(func() {
			p := new(Plugin)
			c := p.GetContext([]string{"azure-cloud-config",
				"--az", "z1",
				"--az", "z2",
				"--multi-assign-az", //This is required if you want to bypass validations about multiple cidr/gateway required error message
				"--network-name-1", "default",
				"--network-az-1", "z1",
				"--network-az-1", "z2",
				"--network-cidr-1", "10.0.0.0/26",
				"--network-gateway-1", "10.0.0.1",
				"--network-dns-1", "169.254.169.254",
				"--network-dns-1", "8.8.8.8",
				"--azure-virtual-network-name-1", "boshnet",
				"--azure-subnet-name-1", "boshsub",
				"--network-name-2", "concourse",
				"--network-az-2", "z1",
				"--network-az-2", "z2",
				"--network-cidr-2", "10.0.0.64/26",
				"--network-gateway-2", "10.0.0.65",
				"--network-dns-2", "169.254.169.254",
				"--network-dns-2", "8.8.8.8",
				"--azure-virtual-network-name-2", "boshnet",
				"--azure-subnet-name-2", "concoursesub"})
			provider = NewAzureCloudConfig(c)
			manifest, err = cloudconfigs.CreateCloudConfigManifest(provider)
		})
		It("then it should have a manifest with 2 azs", func() {
			Ω(err).ShouldNot(HaveOccurred())
			Ω(manifest.ContainsAZName("z1")).Should(BeTrue())
			Ω(manifest.ContainsAZName("z2")).Should(BeTrue())
			bytes, err := ioutil.ReadFile("fixtures/azure-azs.yml")
			Ω(err).ShouldNot(HaveOccurred())
			azYml, err := yaml.Marshal(manifest.AZs)
			Ω(err).ShouldNot(HaveOccurred())

			Ω(azYml).Should(MatchYAML(bytes))
		})
		It("then it should return vmtypes", func() {
			vmTypes, err := provider.CreateVMTypes()
			Ω(err).ShouldNot(HaveOccurred())
			Ω(vmTypes).Should(HaveLen(21))
			Ω(utils.GetVMTypeNames(vmTypes)).Should(ConsistOf("nano", "micro", "micro.ram", "small", "small.disk", "medium", "medium.mem", "medium.disk", "medium.cpu", "large", "large.mem", "large.disk", "large.cpu", "xlarge", "xlarge.mem", "xlarge.disk", "xlarge.cpu", "2xlarge", "2xlarge.mem", "2xlarge.disk", "2xlarge.cpu"))
		})
		It("then it have a manifest with 2 network", func() {
			Ω(err).ShouldNot(HaveOccurred())
			bytes, err := ioutil.ReadFile("fixtures/azure-networks.yml")
			Ω(err).ShouldNot(HaveOccurred())
			networkYml, err := yaml.Marshal(manifest.Networks)
			Ω(err).ShouldNot(HaveOccurred())

			Ω(networkYml).Should(MatchYAML(bytes))
		})
		It("then it return disk types", func() {
			diskTypes, err := provider.CreateDiskTypes()
			Ω(err).ShouldNot(HaveOccurred())
			Ω(diskTypes).Should(HaveLen(19))
			Ω(utils.GetDiskTypeNames(diskTypes)).Should(ConsistOf("1024", "2048", "5120", "10240", "20480", "30720", "51200", "76800", "102400", "153600", "204800", "307200", "512000", "768000", "1048576", "2097152", "5242880", "10485760", "16777216"))
		})
		It("then it have a manifest with a compilation", func() {
			Ω(err).ShouldNot(HaveOccurred())
			bytes, err := ioutil.ReadFile("fixtures/azure-compilation.yml")
			Ω(err).ShouldNot(HaveOccurred())
			compilationYml, err := yaml.Marshal(manifest.Compilation)
			Ω(err).ShouldNot(HaveOccurred())

			Ω(compilationYml).Should(MatchYAML(bytes))
		})
	})
})
