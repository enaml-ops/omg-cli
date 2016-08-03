package boshinit_test

import (
	"io/ioutil"

	. "github.com/enaml-ops/omg-cli/plugins/products/bosh-init"
	"github.com/enaml-ops/omg-cli/plugins/products/bosh-init/enaml-gen/photoncpi"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	yaml "gopkg.in/yaml.v2"
)

var _ = Describe("NewPhotonBosh", func() {
	Describe("given NewPhotonBosh", func() {
		Context("when using a bosh config with a CPI URL", func() {
			cfg := &PhotonBoshInitConfig{}
			const (
				controlURL = "file://example-cpi"
				controlSHA = "slkjdaslkdjlakjdsk"
			)
			var boshBase = NewPhotonBoshBase()
			boshBase.CPIReleaseURL = controlURL
			boshBase.CPIReleaseSHA = controlSHA

			var provider IAASManifestProvider

			BeforeEach(func() {
				provider = NewPhotonIaaSProvider(cfg, boshBase)
			})

			It("creates a CPI release with the specified URL and SHA", func() {
				release := provider.CreateCPIRelease()
				Ω(release.Name).Should(Equal(PhotonCPIReleaseName))
				Ω(release.URL).Should(Equal(controlURL))
				Ω(release.SHA1).Should(Equal(controlSHA))
			})
		})

		Context("when using a bosh config without a CPI URL", func() {
			cfg := &PhotonBoshInitConfig{}
			var boshBase = NewPhotonBoshBase()
			var provider IAASManifestProvider

			BeforeEach(func() {
				provider = NewPhotonIaaSProvider(cfg, boshBase)
			})

			It("creates a CPI release the default URL and SHA", func() {
				release := provider.CreateCPIRelease()
				Ω(release.Name).Should(Equal(PhotonCPIReleaseName))
				Ω(release.URL).Should(Equal(PhotonCPIURL))
				Ω(release.SHA1).Should(Equal(PhotonCPISHA))
			})
		})

		Context("when using a bosh config with a CPI URL", func() {
			cfg := &PhotonBoshInitConfig{}
			const (
				controlURL = "file://example-cpi"
				controlSHA = "slkjdaslkdjlakjdsk"
			)
			var boshBase = NewPhotonBoshBase()
			boshBase.CPIReleaseURL = controlURL
			boshBase.CPIReleaseSHA = controlSHA

			var provider IAASManifestProvider

			BeforeEach(func() {
				provider = NewPhotonIaaSProvider(cfg, boshBase)
			})

			It("creates a CPI release with the specified URL and SHA", func() {
				release := provider.CreateCPIRelease()
				Ω(release.Name).Should(Equal(PhotonCPIReleaseName))
				Ω(release.URL).Should(Equal(controlURL))
				Ω(release.SHA1).Should(Equal(controlSHA))
			})
		})

		Context("when using a bosh init config with networking", func() {
			cfg := &PhotonBoshInitConfig{
				NetworkName: "dwallraff-vnet",
			}
			var boshBase = NewPhotonBoshBase()
			boshBase.NetworkCIDR = "10.0.0.0/24"
			boshBase.NetworkGateway = "10.0.0.1"
			boshBase.PrivateIP = "10.0.0.4"

			var provider IAASManifestProvider

			BeforeEach(func() {
				provider = NewPhotonIaaSProvider(cfg, boshBase)
			})

			It("creates the VIP network", func() {
				net := provider.CreateVIPNetwork()
				Ω(net.Name).Should(Equal("vip"))
				Ω(net.Type).Should(Equal("vip"))
			})

			It("creates the private network", func() {
				net := provider.CreateManualNetwork()

				bytes, err := ioutil.ReadFile("fixtures/photon-manual-network.yml")
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
			cfg := &PhotonBoshInitConfig{
				Photon: photoncpi.Photon{
					Target:     "http://PHOTON_CTRL_IP:9000",
					User:       "PHOTON_USER",
					Password:   "PHOTON_PASSWD",
					IgnoreCert: "PHOTON_IGNORE_CERT",
					Project:    "PHOTON_PROJ_ID",
				},
				MachineType: "n1-standard-4",
			}
			var boshBase = &BoshBase{
				Mode:           "uaa",
				CPIJobName:     "bosh-photon-cpi",
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
				provider = NewPhotonIaaSProvider(cfg, boshBase)
			})

			It("includes a CPI job", func() {
				t := provider.CreateCPITemplate()
				Ω(t.Name).Should(Equal(PhotonCPIJobName))
				Ω(t.Release).Should(Equal(PhotonCPIReleaseName))
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

				const controlYML = `disk_flavor: core-200`
				yml, _ := yaml.Marshal(d.CloudProperties)
				Ω(yml).Should(MatchYAML(controlYML))
			})

			It("includes the resource pool", func() {
				r := provider.CreateResourcePool()
				Ω(r.Name).Should(Equal("vms"))
				Ω(r.Network).Should(Equal("private"))
				Ω(r.Stemcell.URL).Should(Equal(PhotonStemcellURL))
				Ω(r.Stemcell.SHA1).Should(Equal(PhotonStemcellSHA))

				controlYML, err := ioutil.ReadFile("fixtures/photon-resource-pool-properties.yml")
				Ω(err).ShouldNot(HaveOccurred())

				bytes, err := yaml.Marshal(r.CloudProperties)
				Ω(err).ShouldNot(HaveOccurred())
				Ω(bytes).Should(MatchYAML(controlYML))
			})

			It("creates the job properties", func() {
				props := provider.CreateCPIJobProperties()
				controlYML, err := ioutil.ReadFile("fixtures/photon-properties.yml")
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
