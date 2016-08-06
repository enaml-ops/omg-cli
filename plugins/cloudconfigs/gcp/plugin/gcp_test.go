package plugin_test

import (
	"io/ioutil"

	"github.com/enaml-ops/enaml"
	"github.com/enaml-ops/omg-cli/plugins/cloudconfigs"
	. "github.com/enaml-ops/omg-cli/plugins/cloudconfigs/gcp/plugin"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gopkg.in/yaml.v2"
)

var _ = Describe("given GCPCloud Config", func() {
	Context("when calling CreateManifest", func() {
		var provider cloudconfigs.CloudConfigProvider
		var manifest *enaml.CloudConfigManifest
		var err error
		BeforeEach(func() {
			p := new(Plugin)
			c := p.GetContext([]string{"gcp-cloud-config",
				"--az", "z1",
				"--az", "z2",
				"--az", "z3",
				"--gcp-availability-zone", "test1",
				"--gcp-availability-zone", "test2",
				"--gcp-availability-zone", "test3",
				"--network-name-1", "bosh",
				"--network-az-1", "test1",
				"--network-cidr-1", "10.0.0.0/26",
				"--network-gateway-1", "10.0.0.1",
				"--network-dns-1", "169.254.169.254,8.8.8.8",
				"--network-reserved-1", "10.0.0.1-10.0.0.2,10.0.0.60-10.0.0.63",
				"--network-static-1", "10.0.0.4,10.0.0.10",
				"--gcp-network-name-1", "dwallraff-vnet",
				"--gcp-subnetwork-name-1", "dwallraff-subnet-bosh-us-east1",
				"--gcp-network-tag-1", "nat-traverse,no-ip",
				"--network-name-2", "concourse",
				"--network-az-2", "test1",
				"--network-cidr-2", "10.0.0.64/26",
				"--network-gateway-2", "10.0.0.65",
				"--network-dns-2", "169.254.169.254,8.8.8.8",
				"--network-reserved-2", "10.0.0.65-10.0.0.70,10.0.0.122-10.0.0.127",
				"--network-static-2", "10.0.0.72,10.0.0.73,10.0.0.74,10.0.0.75",
				"--gcp-network-name-2", "dwallraff-vnet",
				"--gcp-subnetwork-name-2", "dwallraff-subnet-concourse-us-east1-c",
				"--gcp-network-tag-2", "nat-traverse,no-ip",
			})
			provider = NewGCPCloudConfig(c)
			manifest, err = cloudconfigs.CreateCloudConfigManifest(provider)
		})
		It("then it have a manifest with 3 azs", func() {
			Ω(err).ShouldNot(HaveOccurred())
			Ω(manifest.ContainsAZName("z1")).Should(BeTrue())
			Ω(manifest.ContainsAZName("z2")).Should(BeTrue())
			Ω(manifest.ContainsAZName("z3")).Should(BeTrue())

			bytes, err := ioutil.ReadFile("fixtures/gcp-azs.yml")
			Ω(err).ShouldNot(HaveOccurred())
			azYml, err := yaml.Marshal(manifest.AZs)
			Ω(err).ShouldNot(HaveOccurred())

			Ω(azYml).Should(MatchYAML(bytes))
		})
		It("then it should return vmtypes", func() {
			vmTypes, err := provider.CreateVMTypes()
			Ω(err).ShouldNot(HaveOccurred())

			bytes, err := ioutil.ReadFile("fixtures/gcp-vmtypes.yml")
			Ω(err).ShouldNot(HaveOccurred())
			vmTypesYml, err := yaml.Marshal(vmTypes)
			Ω(err).ShouldNot(HaveOccurred())

			Ω(vmTypesYml).Should(MatchYAML(bytes))
		})
		It("then it have a manifest with 2 network", func() {
			Ω(err).ShouldNot(HaveOccurred())
			bytes, err := ioutil.ReadFile("fixtures/gcp-networks.yml")
			Ω(err).ShouldNot(HaveOccurred())
			networkYml, err := yaml.Marshal(manifest.Networks)
			Ω(err).ShouldNot(HaveOccurred())

			Ω(networkYml).Should(MatchYAML(bytes))
		})
		It("then it have a manifest disk types", func() {
			Ω(err).ShouldNot(HaveOccurred())
			bytes, err := ioutil.ReadFile("fixtures/gcp-disktypes.yml")
			Ω(err).ShouldNot(HaveOccurred())
			diskTypeYml, err := yaml.Marshal(manifest.DiskTypes)
			Ω(err).ShouldNot(HaveOccurred())

			Ω(diskTypeYml).Should(MatchYAML(bytes))
		})
		It("then it have a manifest with a compilation", func() {
			Ω(err).ShouldNot(HaveOccurred())
			bytes, err := ioutil.ReadFile("fixtures/gcp-compilation.yml")
			Ω(err).ShouldNot(HaveOccurred())
			compilationYml, err := yaml.Marshal(manifest.Compilation)
			Ω(err).ShouldNot(HaveOccurred())

			Ω(compilationYml).Should(MatchYAML(bytes))
		})
	})
})
