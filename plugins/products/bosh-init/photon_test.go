package boshinit_test

import (
	"fmt"
	"io/ioutil"

	"github.com/enaml-ops/enaml"
	. "github.com/enaml-ops/omg-cli/plugins/products/bosh-init"
	"github.com/enaml-ops/omg-cli/plugins/products/bosh-init/enaml-gen/blobstore"
	"github.com/enaml-ops/omg-cli/plugins/products/bosh-init/enaml-gen/director"
	"github.com/enaml-ops/omg-cli/plugins/products/bosh-init/enaml-gen/photoncpi"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	yaml "gopkg.in/yaml.v2"
)

var _ = Describe("NewPhotonBosh", func() {
	Describe("given a NewPhotonBoshBase function", func() {
		Context("When creating the bosh job", func() {
			cfg := &PhotonBoshInitConfig{}
			base := &BoshBase{}
			provider := NewPhotonIaaSProvider(cfg, base)

			It("sets the CPI job name", func() {
				manifest, err := provider.CreateDeploymentManifest()
				Ω(err).ShouldNot(HaveOccurred())

				job := manifest.GetJobByName("bosh")
				Ω(job).ShouldNot(BeNil())

				director := job.Properties["director"].(*director.Director)
				Ω(director.CpiJob).Should(Equal(PhotonCPIJobName))
			})
		})

		Context("when called with a valid boshbase for uaa", func() {
			var boshbase *BoshBase
			const controlDirectorName = "fake-director-name"
			BeforeEach(func() {
				boshbase = NewPhotonBoshBase()
				boshbase.Mode = "uaa"
				boshbase.DirectorName = controlDirectorName
				boshbase.UAAReleaseSHA = "uaa-release.com"
				boshbase.UAAReleaseURL = "uaa-release-lkashdjlkgahsdg"
				boshbase.NetworkCIDR = "10.0.0.1/24"
				boshbase.NetworkGateway = "10.0.0.254"
				boshbase.PrivateIP = "10.0.0.4"
				boshbase.PersistentDiskSize = 32768
			})
			It("then it should give us a boshbase with properly set director name", func() {
				Ω(boshbase.DirectorName).Should(Equal(controlDirectorName))
			})

			It("then it should give us a boshbase with properly set uaa mode", func() {
				Ω(boshbase.Mode).Should(Equal("uaa"))
			})
		})
	})

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

		Context("when using a bosh config for cpi version 0.9.0", func() {
			cfg := &PhotonBoshInitConfig{}
			const (
				controlURL               = "file://example-cpi"
				controlSHA               = "slkjdaslkdjlakjdsk"
				controlPrivateIP         = "1.2.3.4"
				controlDirectorPassword  = "blah"
				controlNatsAgentPassword = "bleh"
			)
			var boshBase = NewPhotonBoshBase()
			boshBase.PrivateIP = controlPrivateIP
			boshBase.DirectorPassword = controlDirectorPassword
			boshBase.NatsPassword = controlNatsAgentPassword
			boshBase.CPIReleaseURL = controlURL
			boshBase.CPIReleaseSHA = controlSHA

			var provider IAASManifestProvider
			var blobstoreOptions map[string]interface{}
			var bs map[string]interface{}
			var props *photoncpi.PhotoncpiJob
			var dm *enaml.DeploymentManifest
			var err error
			BeforeEach(func() {
				provider = NewPhotonIaaSProvider(cfg, boshBase)
				cp := provider.CreateCloudProvider()
				Ω(cp.Template.Name).Should(Equal("cpi"))
				props = cp.Properties.(*photoncpi.PhotoncpiJob)
				blobstoreOptions = props.Blobstore.Options.(map[string]interface{})
				dm, err = provider.(*PhotonBosh).CreateDeploymentManifest()
				bs = dm.GetJobByName("bosh").Properties["blobstore"].(map[string]interface{})
			})

			It("then it should not error", func() {
				Ω(err).ShouldNot(HaveOccurred())
			})

			It("then it should set a valid blobstore port", func() {
				Ω(bs).Should(HaveKeyWithValue("provider", "dav"))
			})

			It("then it should set valid options", func() {
				Ω(bs["options"]).Should(HaveKeyWithValue("endpoint", fmt.Sprintf("http://%v:25250", controlPrivateIP)))
				Ω(bs["options"]).Should(HaveKeyWithValue("user", "agent"))
				Ω(bs["options"]).Should(HaveKeyWithValue("password", controlNatsAgentPassword))
			})

			It("then it should set a valid blobstore ip address", func() {
				Ω(bs).Should(HaveKeyWithValue("address", controlPrivateIP))
			})

			It("then it should set a valid blobstore port", func() {
				Ω(bs).Should(HaveKeyWithValue("port", 25250))
			})

			It("then it should set valid blobstore agent credentials", func() {
				Ω(bs["agent"].(blobstore.Agent).User).Should(Equal("agent"))
				Ω(bs["agent"].(blobstore.Agent).Password).Should(Equal(controlNatsAgentPassword))
			})

			It("then it should set valid director credentials", func() {
				Ω(bs["director"].(blobstore.Director).User).Should(Equal("director"))
				Ω(bs["director"].(blobstore.Director).Password).Should(Equal(controlDirectorPassword))
			})

			It("then it should set a valid blobstore path", func() {
				Ω(blobstoreOptions).Should(HaveKeyWithValue("blobstore_path", "/var/vcap/micro_bosh/data/cache"))
			})

			It("then it should properly set the cloud provider", func() {
				Ω(props.Blobstore).ShouldNot(BeNil())
				Ω(props.Blobstore.Provider).Should(Equal("local"))
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
			boshBase.NetworkDNS = []string{"10.0.0.2"}
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
			var boshBase = NewPhotonBoshBase()
			boshBase.Mode = "uaa"
			boshBase.PrivateIP = controlPrivateIP
			boshBase.PublicIP = "1.0.2.3"
			boshBase.NetworkCIDR = "10.0.0.0/24"
			boshBase.NetworkGateway = "10.0.0.1"
			boshBase.NetworkDNS = []string{"10.0.0.2"}
			boshBase.NtpServers = []string{controlNTP}
			boshBase.MBusPassword = controlMbusPass
			boshBase.NatsPassword = controlNatsPass
			boshBase.PersistentDiskSize = 32768

			var provider IAASManifestProvider

			BeforeEach(func() {
				provider = NewPhotonIaaSProvider(cfg, boshBase)
			})

			It("includes a CPI job", func() {
				t := provider.CreateCPITemplate()
				Ω(t.Name).Should(Equal(PhotonCPIJobName))
				Ω(t.Release).Should(Equal(PhotonCPIReleaseName))
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
				r, err := provider.CreateResourcePool()
				Ω(err).ShouldNot(HaveOccurred())
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
				manifest, err := provider.CreateDeploymentManifest()
				Ω(manifest).ShouldNot(BeNil())
				Ω(err).ShouldNot(HaveOccurred())
			})

			It("builds a manifest with a bosh job containing a single private network", func() {
				manifest, err := provider.CreateDeploymentManifest()
				Ω(manifest).ShouldNot(BeNil())
				Ω(err).ShouldNot(HaveOccurred())

				var privateNets int
				for _, n := range manifest.Jobs[0].Networks {
					if n.Name == "private" {
						privateNets++
					}
				}
				Ω(privateNets).Should(Equal(1))
			})
		})
	})
})
