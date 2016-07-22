package boshinit_test

import (
	"github.com/enaml-ops/enaml"
	. "github.com/enaml-ops/omg-cli/plugins/products/bosh-init"
	"github.com/enaml-ops/omg-cli/plugins/products/bosh-init/enaml-gen/vsphere_cpi"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("NewVSphereBosh", func() {
	Describe("given the function", func() {
		Context("when called w/ valid parameters", func() {
			var boshConfig = BoshInitConfig{
				VSphereAddress:                    "172.16.1.2",
				VSphereUser:                       "vsadmin",
				VSpherePassword:                   "secret",
				VSphereDatacenterName:             "PCF_DC1",
				VSphereVMFolder:                   "pcf_vms",
				VSphereTemplateFolder:             "pcf_templates",
				VSphereDatastorePattern:           "DS1",
				VSpherePersistentDatastorePattern: "DS1_Persistent",
				VSphereDiskPath:                   "pcf_disks",
				VSphereClusters:                   []string{"PCF1"},
				VSphereNetworks: []Network{Network{
					Name:    "PCF_Net1",
					Range:   "172.16.0.0/23",
					Gateway: "172.16.1.1",
					DNS:     []string{"172.16.1.2"},
				}},
			}
			var boshBase = &BoshBase{
				Mode:               "uaa",
				BoshReleaseVersion: "256.2",
				PrivateIP:          "172.16.1.6",
				CPIReleaseVersion:  "22",
				GOAgentVersion:     "3232.4",
				BoshReleaseSHA:     "ff2f4e16e02f66b31c595196052a809100cfd5a8",
				CPIReleaseSHA:      "dd1827e5f4dfc37656017c9f6e48441f51a7ab73",
				GOAgentSHA:         "27ec32ddbdea13e3025700206388ae5882a23c67",
				NetworkCIDR:        "10.0.0.0/24",
				NetworkGateway:     "10.0.0.1",
				NetworkDNS:         []string{"10.0.0.2"},
				DirectorName:       "my-bosh",
			}
			var manifest *enaml.DeploymentManifest

			BeforeEach(func() {
				manifest = NewVSphereBosh(boshConfig, boshBase)
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

			XIt("then it should properly define vcenter properties", func() {
				Ω(manifest.Jobs[0].Properties).Should(HaveKey("vcenter"))
				var vcenter vsphere_cpi.Vcenter

				Ω(vcenter.Address).Should(Equal("172.16.1.2"))
				Ω(vcenter.User).Should(Equal("vsadmin"))
				Ω(vcenter.Password).Should(Equal("secret"))
				Ω(vcenter.Datacenters).Should(HaveLen(1))
				dc := vcenter.Datacenters.(VSphereDatacenters)[0]
				Ω(dc.Name).Should(Equal("PCF_DC1"))
				Ω(dc.DatastorePattern).Should(Equal("DS1"))
				Ω(dc.PersistentDatastorePattern).Should(Equal("DS1_Persistent"))
				Ω(dc.DiskPath).Should(Equal("pcf_disks"))
				Ω(dc.TemplateFolder).Should(Equal("pcf_templates"))
				Ω(dc.VMFolder).Should(Equal("pcf_vms"))
				Ω(dc.Clusters).Should(HaveLen(1))
				Ω(dc.Clusters[0]).Should(Equal("PCF1"))
			})

			Context("When PersistentDatastorePattern isn't specified", func() {
				BeforeEach(func() {
					boshConfig.VSpherePersistentDatastorePattern = ""
					manifest = NewVSphereBosh(boshConfig, boshBase)
				})
				XIt("then it should fallback to DatastorePattern", func() {
					var vcenter vsphere_cpi.Vcenter
					dc := vcenter.Datacenters.(VSphereDatacenters)[0]
					Ω(dc.DatastorePattern).Should(Equal("DS1"))
					Ω(dc.PersistentDatastorePattern).Should(Equal("DS1"))
				})
			})
		})
	})
})
