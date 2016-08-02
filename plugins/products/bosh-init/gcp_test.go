package boshinit_test

import (
	"fmt"
	"io/ioutil"

	. "github.com/enaml-ops/omg-cli/plugins/products/bosh-init"
	"github.com/enaml-ops/omg-cli/plugins/products/bosh-init/enaml-gen/google_cpi"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	yaml "gopkg.in/yaml.v2"
)

var _ = Describe("NewGCPBosh", func() {
	Describe("given NewGCPBosh", func() {
		Context("when using a bosh config with a CPI URL", func() {
			cfg := &GCPBoshInitConfig{}
			const (
				controlURL = "file://example-cpi"
				controlSHA = "slkjdaslkdjlakjdsk"
			)
			var boshBase = NewGCPBoshBase()
			boshBase.CPIReleaseURL = controlURL
			boshBase.CPIReleaseSHA = controlSHA

			var provider IAASManifestProvider

			BeforeEach(func() {
				provider = NewGCPIaaSProvider(cfg, boshBase)
			})

			It("creates a CPI release with the specified URL and SHA", func() {
				release := provider.CreateCPIRelease()
				Ω(release.Name).Should(Equal(GCPCPIReleaseName))
				Ω(release.URL).Should(Equal(controlURL))
				Ω(release.SHA1).Should(Equal(controlSHA))
			})
		})

		Context("when using a bosh config without a CPI URL", func() {
			cfg := &GCPBoshInitConfig{}
			var boshBase = NewGCPBoshBase()
			var provider IAASManifestProvider

			BeforeEach(func() {
				provider = NewGCPIaaSProvider(cfg, boshBase)
			})

			It("creates a CPI release the default URL and SHA", func() {
				release := provider.CreateCPIRelease()
				Ω(release.Name).Should(Equal(GCPCPIReleaseName))
				Ω(release.URL).Should(Equal(GCPCPIURL))
				Ω(release.SHA1).Should(Equal(GCPCPISHA))
			})
		})

		Context("when using a bosh config with a CPI URL", func() {
			cfg := &GCPBoshInitConfig{}
			const (
				controlURL = "file://example-cpi"
				controlSHA = "slkjdaslkdjlakjdsk"
			)
			var boshBase = NewGCPBoshBase()
			boshBase.CPIReleaseURL = controlURL
			boshBase.CPIReleaseSHA = controlSHA

			var provider IAASManifestProvider

			BeforeEach(func() {
				provider = NewGCPIaaSProvider(cfg, boshBase)
			})

			It("creates a CPI release with the specified URL and SHA", func() {
				release := provider.CreateCPIRelease()
				Ω(release.Name).Should(Equal(GCPCPIReleaseName))
				Ω(release.URL).Should(Equal(controlURL))
				Ω(release.SHA1).Should(Equal(controlSHA))
			})
		})

		Context("when using a bosh init config with networking", func() {
			cfg := &GCPBoshInitConfig{
				NetworkName:    "dwallraff-vnet",
				SubnetworkName: "dwallraff-subnet-bosh-us-east1",
			}
			var boshBase = NewGCPBoshBase()
			boshBase.NetworkCIDR = "10.0.0.0/24"
			boshBase.NetworkGateway = "10.0.0.1"
			boshBase.PrivateIP = "10.0.0.4"

			var provider IAASManifestProvider

			BeforeEach(func() {
				provider = NewGCPIaaSProvider(cfg, boshBase)
			})

			It("creates the VIP network", func() {
				net := provider.CreateVIPNetwork()
				Ω(net.Name).Should(Equal("vip"))
				Ω(net.Type).Should(Equal("vip"))
			})

			It("creates the private network", func() {
				net := provider.CreateManualNetwork()

				bytes, err := ioutil.ReadFile("fixtures/gcp-manual-network.yml")
				Ω(err).ShouldNot(HaveOccurred())
				netYml, err := yaml.Marshal(net.Subnets)
				Ω(err).ShouldNot(HaveOccurred())

				Ω(netYml).Should(MatchYAML(bytes))
			})
		})

		Context("when called w/ valid parameters", func() {
			const (
				controlPrivateIP = "10.0.0.6"
				controlProject   = "myproject"
				controlZone      = "zone"
				controlNTP       = "169.254.0.1"
				controlMbusPass  = "mbuspassword"
				controlNatsPass  = "natspassword"
			)
			cfg := &GCPBoshInitConfig{
				Project:     controlProject,
				DefaultZone: controlZone,
			}
			var boshBase = &BoshBase{
				Mode:           "uaa",
				CPIJobName:     "bosh-google-cpi",
				PrivateIP:      controlPrivateIP,
				PublicIP:       "1.0.2.3",
				CPIReleaseSHA:  "dc4a0cca3b33dce291e4fbeb9e9948b6a7be3324",
				NetworkCIDR:    "10.0.0.0/24",
				NetworkGateway: "10.0.0.1",
				NetworkDNS:     []string{"10.0.0.2"},
				DirectorName:   "my-bosh",
				NtpServers:     []string{controlNTP},
				MBusPassword:   controlMbusPass,
				NatsPassword:   controlNatsPass,
			}

			var provider IAASManifestProvider

			BeforeEach(func() {
				provider = NewGCPIaaSProvider(cfg, boshBase)
			})

			It("includes a CPI job", func() {
				t := provider.CreateCPITemplate()
				Ω(t.Name).Should(Equal(GCPCPIJobName))
				Ω(t.Release).Should(Equal(GCPCPIReleaseName))
			})

			It("includes the job network", func() {
				n := provider.CreateJobNetwork()
				Ω(n.Name).Should(Equal("private"))
				Ω(n.StaticIPs).Should(ConsistOf(controlPrivateIP))
				Ω(n.Default).Should(ConsistOf("dns", "gateway"))
			})

			It("includes the disk pool", func() {
				d := provider.CreateDiskPool()
				Ω(d.Name).Should(Equal("disks"))
				Ω(d.DiskSize).Should(Equal(32768))

				const controlYML = `type: pd-standard`
				yml, _ := yaml.Marshal(d.CloudProperties)
				Ω(yml).Should(MatchYAML(controlYML))
			})

			It("includes the resource pool", func() {
				r := provider.CreateResourcePool()
				Ω(r.Name).Should(Equal("vms"))
				Ω(r.Network).Should(Equal("private"))
				Ω(r.Stemcell.URL).Should(Equal(GCPStemcellURL))
				Ω(r.Stemcell.SHA1).Should(Equal(GCPStemcellSHA))

				controlYML, err := ioutil.ReadFile("fixtures/gcp-resource-pool-properties.yml")
				Ω(err).ShouldNot(HaveOccurred())

				bytes, err := yaml.Marshal(r.CloudProperties)
				Ω(err).ShouldNot(HaveOccurred())
				Ω(bytes).Should(MatchYAML(controlYML))
			})

			It("creates the cloud provider", func() {
				cp := provider.CreateCloudProvider()
				Ω(cp.Template.Name).Should(Equal(GCPCPIReleaseName))
				Ω(cp.Template.Release).Should(Equal(GCPCPIReleaseName))

				Ω(cp.SSHTunnel.Host).Should(Equal(controlPrivateIP))
				Ω(cp.SSHTunnel.Port).Should(Equal(22))
				Ω(cp.SSHTunnel.PrivateKeyPath).Should(Equal("~/.ssh/bosh"))
				Ω(cp.SSHTunnel.User).Should(Equal("bosh"))

				Ω(cp.MBus).Should(Equal(fmt.Sprintf("https://mbus:%s@%s:6868", controlMbusPass, controlPrivateIP)))

				props := cp.Properties.(*google_cpi.GoogleCpiJob)
				Ω(props.Google.Project).Should(Equal(controlProject))
				Ω(props.Google.DefaultZone).Should(Equal(controlZone))
				Ω(props.Agent.Mbus).Should(Equal(fmt.Sprintf("https://mbus:%s@0.0.0.0:6868", controlMbusPass)))
				Ω(props.Blobstore.Provider).Should(Equal("local"))
				Ω(props.Blobstore.Path).Should(Equal("/var/vcap/micro_bosh/data/cache"))
				Ω(props.Ntp).Should(ConsistOf(controlNTP))
			})

			It("creates the job properties", func() {
				props := provider.CreateCPIJobProperties()
				controlYML, err := ioutil.ReadFile("fixtures/gcp-properties.yml")
				Ω(err).ShouldNot(HaveOccurred())

				bytes, err := yaml.Marshal(props)
				Ω(err).ShouldNot(HaveOccurred())
				Ω(bytes).Should(MatchYAML(controlYML))
			})

			It("builds a valid manifest", func() {
				manifest := provider.CreateDeploymentManifest()
				Ω(manifest).ShouldNot(BeNil())
			})
		})
	})
})
