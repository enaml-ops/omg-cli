package boshinit_test

import (
	"io/ioutil"

	"github.com/enaml-ops/enaml"
	. "github.com/enaml-ops/omg-cli/plugins/products/bosh-init"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	yaml "gopkg.in/yaml.v2"
)

var _ = Describe("NewVSphereBosh", func() {
	Describe("given the function", func() {
		Context("when called w/ valid parameters", func() {
			var boshConfig = VSphereInitConfig{
				VSphereAddress:        "172.16.1.2",
				VSphereUser:           "vsadmin",
				VSpherePassword:       "secret",
				VSphereDatacenterName: "PCF_DC1",
				VSphereVMFolder:       "pcf_vms",
				VSphereTemplateFolder: "pcf_templates",
				VSphereDataStore:      "DS1",
				VSphereDiskPath:       "pcf_disks",
				VSphereClusters:       []string{"PCF1"},
				VSphereResourcePool:   []string{"RP01"},
				VSphereNetworks: []Network{Network{
					Name:    "PCF_Net1",
					Range:   "172.16.0.0/23",
					Gateway: "172.16.1.1",
					DNS:     []string{"172.16.1.2"},
				}},
			}
			var boshBase = GetVSphereDefaults()
			boshBase.Mode = "uaa"
			boshBase.DirectorName = "my-bosh"
			boshBase.NetworkCIDR = "172.16.0.0/23"
			boshBase.NetworkGateway = "172.16.1.1"
			boshBase.NetworkDNS = []string{"172.16.1.2"}
			boshBase.PrivateIP = "172.16.1.6"

			var manifest *enaml.DeploymentManifest
			var err error
			BeforeEach(func() {
				manifest, err = NewVSphereIaaSProvider(boshConfig, boshBase).CreateDeploymentManifest()
			})

			It("then error should not occur", func() {
				Ω(err).ShouldNot(HaveOccurred())
			})

			It("then it should be using the vsphere esx stemcell", func() {
				Ω(manifest.ResourcePools[0].Stemcell.URL).ShouldNot(ContainSubstring("aws"))
				Ω(manifest.ResourcePools[0].Stemcell.URL).ShouldNot(ContainSubstring("azure"))
				Ω(manifest.ResourcePools[0].Stemcell.URL).Should(ContainSubstring("esx"))
			})

			It("then it should have the correct job config to deploy a bosh", func() {
				Ω(len(manifest.Jobs)).Should(Equal(1))
			})

			XIt("then it should properly define job properties", func() {
				Ω(len(manifest.Jobs[0].Properties)).Should(Equal(9))
				Ω(func() (r []interface{}) {
					for n := range manifest.Jobs[0].Properties {
						r = append(r, n)
					}
					return
				}()).Should(ConsistOf("nats", "postgres", "blobstore", "director", "hm", "vcenter", "agent", "ntp", "registry"))
			})

			It("then it should properly define job templates", func() {
				Ω(len(manifest.Jobs[0].Templates)).Should(Equal(8))
				Ω(func() (r []string) {
					for _, v := range manifest.Jobs[0].Templates {
						r = append(r, v.Name)
					}
					return
				}()).Should(ConsistOf("nats", "postgres", "blobstore", "director", "health_monitor", "vsphere_cpi", "uaa", "registry"))
			})

			It("then it should properly define job networks", func() {
				Ω(manifest.Jobs[0].Networks).Should(HaveLen(1))
				net := manifest.Jobs[0].Networks[0]
				Ω(net.Name).Should(Equal("private"))
				Ω(net.StaticIPs).Should(ContainElement("172.16.1.6"))
			})

			It("then it should properly define networks", func() {
				Ω(manifest.Networks).Should(HaveLen(1))
				net := manifest.Networks[0].(enaml.ManualNetwork)
				Ω(net.Name).Should(Equal("private"))
				Ω(net.Type).Should(Equal("manual"))
				Ω(net.Subnets).Should(HaveLen(1))
				subnet := net.Subnets[0]
				Ω(subnet.DNS).Should(HaveLen(1))
				Ω(subnet.DNS[0]).Should(Equal("172.16.1.2"))
				Ω(subnet.Gateway).Should(Equal("172.16.1.1"))
				Ω(subnet.Range).Should(Equal("172.16.0.0/23"))
				cloudprops := subnet.CloudProperties.(VSpherecloudpropertiesNetwork)
				Ω(cloudprops.Name).Should(Equal("PCF_Net1"))
			})

			It("then it should properly define vcenter properties", func() {
				controlYaml, err := ioutil.ReadFile("./fixtures/vsphere-manifest.yml")
				Ω(err).ShouldNot(HaveOccurred())
				b, err := yaml.Marshal(manifest.ResourcePools[0].CloudProperties)
				Ω(err).ShouldNot(HaveOccurred())
				Ω(b).Should(MatchYAML(string(controlYaml)))
			})
			It("then it should properly define vcenter properties with no resource pools", func() {
				boshConfig.VSphereResourcePool = nil
				boshConfig.VSphereClusters = []string{"PCF1", "PCF2"}
				manifest, err = NewVSphereIaaSProvider(boshConfig, boshBase).CreateDeploymentManifest()

				controlYaml, err := ioutil.ReadFile("./fixtures/vsphere-manifest-no-resourcepools.yml")
				Ω(err).ShouldNot(HaveOccurred())
				b, err := yaml.Marshal(manifest.ResourcePools[0].CloudProperties)
				Ω(err).ShouldNot(HaveOccurred())
				Ω(b).Should(MatchYAML(string(controlYaml)))
			})
		})
	})
})
