package plugin_test

import (
	"io/ioutil"

	"github.com/enaml-ops/enaml"
	"github.com/enaml-ops/omg-cli/plugins/cloudconfigs"
	"github.com/enaml-ops/omg-cli/plugins/cloudconfigs/utils"
	. "github.com/enaml-ops/omg-cli/plugins/cloudconfigs/vsphere/plugin"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gopkg.in/yaml.v2"
)

var _ = Describe("given vSphere Cloud Config", func() {
	Context("when calling CreateManifest", func() {
		var provider cloudconfigs.CloudConfigProvider
		var manifest *enaml.CloudConfigManifest
		var err error
		BeforeEach(func() {
			p := new(Plugin)
			c := p.GetContext([]string{"photon-cloud-config",
				"--az", "z1",
				"--az", "z2",
				"--az", "z3",
				"--network-name-1", "bosh",
				"--network-az-1", "z1",
				"--network-cidr-1", "10.0.0.0/26",
				"--network-gateway-1", "10.0.0.1",
				"--network-dns-1", "169.254.169.254",
				"--network-dns-1", "8.8.8.8",
				"--network-reserved-1", "10.0.0.1-10.0.0.2",
				"--network-reserved-1", "10.0.0.60-10.0.0.63",
				"--network-reserved-2", "10.0.0.65-10.0.0.70",
				"--network-reserved-2", "10.0.0.122-10.0.0.127",
				"--vsphere-network-name-1", "vnet",
				"--network-name-2", "concourse",
				"--network-az-2", "z1",
				"--network-cidr-2", "10.0.0.64/26",
				"--network-gateway-2", "10.0.0.65",
				"--network-dns-2", "169.254.169.254",
				"--network-dns-2", "8.8.8.8",
				"--network-static-1", "10.0.0.4",
				"--network-static-1", "10.0.0.10",
				"--network-static-2", "10.0.0.72",
				"--network-static-2", "10.0.0.73",
				"--network-static-2", "10.0.0.74",
				"--network-static-2", "10.0.0.75",
				"--vsphere-network-name-2", "vnet",
				"--vsphere-datacenter", "dc",
				"--vsphere-datacenter", "dc",
				"--vsphere-datacenter", "dc",
				"--vsphere-cluster", "vsphere-cluster1",
				"--vsphere-cluster", "vsphere-cluster2",
				"--vsphere-cluster", "vsphere-cluster3",
				"--vsphere-resource-pool", "vsphere-res-pool1",
				"--vsphere-resource-pool", "vsphere-res-pool1",
				"--vsphere-resource-pool", "vsphere-res-pool1",
			})
			provider = NewVSphereCloudConfig(c)
			manifest, err = cloudconfigs.CreateCloudConfigManifest(provider)
		})
		It("then it have a manifest with 3 azs", func() {
			Ω(err).ShouldNot(HaveOccurred())
			Ω(manifest.ContainsAZName("z1")).Should(BeTrue())
			Ω(manifest.ContainsAZName("z2")).Should(BeTrue())
			Ω(manifest.ContainsAZName("z3")).Should(BeTrue())

			bytes, err := ioutil.ReadFile("fixtures/vcenter-azs.yml")
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
			bytes, err := ioutil.ReadFile("fixtures/vcenter-networks.yml")
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
			bytes, err := ioutil.ReadFile("fixtures/vcenter-compilation.yml")
			Ω(err).ShouldNot(HaveOccurred())
			compilationYml, err := yaml.Marshal(manifest.Compilation)
			Ω(err).ShouldNot(HaveOccurred())

			Ω(compilationYml).Should(MatchYAML(bytes))
		})
	})
})
