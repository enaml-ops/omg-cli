package cloudconfig_test

import (
	"io/ioutil"

	"github.com/enaml-ops/enaml"
	"github.com/enaml-ops/omg-cli/plugins/cloudconfigs"
	. "github.com/enaml-ops/omg-cli/plugins/cloudconfigs/aws/cloud-config"
	. "github.com/enaml-ops/omg-cli/plugins/cloudconfigs/aws/plugin"
	"github.com/enaml-ops/omg-cli/plugins/cloudconfigs/utils"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gopkg.in/yaml.v2"
)

var _ = Describe("given CloudConfig Deployment for AWS", func() {
	Context("when calling CreateManifest", func() {
		var provider cloudconfigs.CloudConfigProvider
		var manifest *enaml.CloudConfigManifest
		var err error
		BeforeEach(func() {
			p := new(Plugin)
			c := p.GetContext([]string{"aws-cloud-config",
				"--az", "z1",
				"--az", "z2",
				"--az", "z3",
				"--aws-availablity-zone", "us-east-1a",
				"--aws-availablity-zone", "us-east-1b",
				"--aws-availablity-zone", "us-east-1c",
				"--network-name-1", "deployment",
				"--network-az-1", "z1",
				"--network-cidr-1", "10.0.16.0/20",
				"--network-gateway-1", "10.0.16.1",
				"--network-dns-1", "10.0.0.2",
				"--network-reserved-1", "10.0.16.2-10.0.16.10",
				"--network-static-1", "10.0.16.11",
				"--aws-subnet-name-1", "subnet-1",
				"--aws-security-group-1", "sg-deployment",
				"--network-name-2", "deployment",
				"--network-az-2", "z2",
				"--network-cidr-2", "10.0.32.0/20",
				"--network-gateway-2", "10.0.32.1",
				"--network-dns-2", "10.0.0.2",
				"--network-reserved-2", "10.0.32.2-10.0.32.10",
				"--network-static-2", "10.0.32.11",
				"--aws-subnet-name-2", "subnet-2",
				"--aws-security-group-2", "sg-deployment",
				"--network-name-3", "deployment",
				"--network-az-3", "z3",
				"--network-cidr-3", "10.0.48.0/20",
				"--network-gateway-3", "10.0.48.1",
				"--network-dns-3", "10.0.0.2",
				"--network-reserved-3", "10.0.48.2-10.0.48.10",
				"--network-static-3", "10.0.48.11",
				"--aws-subnet-name-3", "subnet-3",
				"--aws-security-group-3", "sg-deployment",
				"--network-name-4", "services",
				"--network-az-4", "z1",
				"--network-cidr-4", "10.0.64.0/20",
				"--network-gateway-4", "10.0.64.1",
				"--network-dns-4", "10.0.0.2",
				"--network-reserved-4", "10.0.64.2-10.0.64.10",
				"--network-static-4", "10.0.64.11",
				"--aws-subnet-name-4", "subnet-4",
				"--aws-security-group-4", "sg-services",
				"--network-name-5", "services",
				"--network-az-5", "z2",
				"--network-cidr-5", "10.0.80.0/20",
				"--network-gateway-5", "10.0.80.1",
				"--network-dns-5", "10.0.0.2",
				"--network-reserved-5", "10.0.80.2-10.0.80.10",
				"--network-static-5", "10.0.80.11",
				"--aws-subnet-name-5", "subnet-5",
				"--aws-security-group-5", "sg-services",
				"--network-name-6", "services",
				"--network-az-6", "z3",
				"--network-cidr-6", "10.0.96.0/20",
				"--network-gateway-6", "10.0.96.1",
				"--network-dns-6", "10.0.0.2",
				"--network-reserved-6", "10.0.96.2-10.0.96.10",
				"--network-static-6", "10.0.96.11",
				"--aws-subnet-name-6", "subnet-6",
				"--aws-security-group-6", "sg-services",
			})
			provider = NewAWSCloudConfig(c)
			manifest, err = cloudconfigs.CreateCloudConfigManifest(provider)
		})
		It("then it have a manifest with 3 azs", func() {
			Ω(err).ShouldNot(HaveOccurred())
			Ω(manifest.ContainsAZName("z1")).Should(BeTrue())
			Ω(manifest.ContainsAZName("z2")).Should(BeTrue())
			Ω(manifest.ContainsAZName("z3")).Should(BeTrue())

			bytes, err := ioutil.ReadFile("fixtures/aws-az-cloudconfig.yml")
			Ω(err).ShouldNot(HaveOccurred())
			azYml, err := yaml.Marshal(manifest.AZs)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(azYml).Should(MatchYAML(bytes))
		})
		It("then it 2 networks", func() {
			Ω(err).ShouldNot(HaveOccurred())
			bytes, err := ioutil.ReadFile("fixtures/aws-network-cloudconfig.yml")
			Ω(err).ShouldNot(HaveOccurred())
			networkYml, err := yaml.Marshal(manifest.Networks)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(networkYml).Should(MatchYAML(bytes))
		})
		It("then it should return vmtypes", func() {
			vmTypes, err := provider.CreateVMTypes()
			Ω(err).ShouldNot(HaveOccurred())
			Ω(vmTypes).Should(HaveLen(21))
			Ω(utils.GetVMTypeNames(vmTypes)).Should(ConsistOf("nano", "micro", "micro.ram", "small", "small.disk", "medium", "medium.mem", "medium.disk", "medium.cpu", "large", "large.mem", "large.disk", "large.cpu", "xlarge", "xlarge.mem", "xlarge.disk", "xlarge.cpu", "2xlarge", "2xlarge.mem", "2xlarge.disk", "2xlarge.cpu"))
		})
		It("then it return disk types", func() {
			diskTypes, err := provider.CreateDiskTypes()
			Ω(err).ShouldNot(HaveOccurred())
			Ω(diskTypes).Should(HaveLen(19))
			Ω(utils.GetDiskTypeNames(diskTypes)).Should(ConsistOf("1024", "2048", "5120", "10240", "20480", "30720", "51200", "76800", "102400", "153600", "204800", "307200", "512000", "768000", "1048576", "2097152", "5242880", "10485760", "16777216"))
		})
		It("then it have a manifest with a compilation", func() {
			Ω(err).ShouldNot(HaveOccurred())
			bytes, err := ioutil.ReadFile("fixtures/aws-compilation.yml")
			Ω(err).ShouldNot(HaveOccurred())
			compilationYml, err := yaml.Marshal(manifest.Compilation)
			Ω(err).ShouldNot(HaveOccurred())

			Ω(compilationYml).Should(MatchYAML(bytes))
		})
	})
})
