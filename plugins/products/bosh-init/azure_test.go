package boshinit_test

import (
	"github.com/enaml-ops/enaml"
	. "github.com/enaml-ops/omg-cli/plugins/products/bosh-init"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("NewAzureBosh", func() {
	Describe("given the function", func() {
		Context("when called w/ valid parameters", func() {
			var boshConfig = AzureInitConfig{
				AzureInstanceSize:         "Standard_D1",
				AzureVnet:                 "something",
				AzureSubnet:               "sub-somthing",
				AzureSubscriptionID:       "azure-subscription-id",
				AzureTenantID:             "azure-tenant-id",
				AzureClientID:             "azure-client-id",
				AzureClientSecret:         "azure-client-secret",
				AzureResourceGroup:        "azure-resource-group",
				AzureStorageAccount:       "azure-storage-account",
				AzureDefaultSecurityGroup: "azure-security-group",
				AzureSSHPubKey:            "azure-ssh-pub-key",
				AzureSSHUser:              "azure-ssh-user",
				AzureEnvironment:          "AzureCloud",
				AzurePrivateKeyPath:       "./bosh",
			}
			var boshBase = GetAzureDefaults()
			boshBase.PublicIP = "x.x.x.x"
			boshBase.Mode = "uaa"
			boshBase.DirectorName = "my-bosh"

			var manifest *enaml.DeploymentManifest
			var err error
			BeforeEach(func() {
				manifest, err = NewAzureIaaSProvider(boshConfig, boshBase).CreateDeploymentManifest()
			})

			It("then it should not error", func() {
				Ω(err).ShouldNot(HaveOccurred())
			})
			It("then it should be using the azure stemcell", func() {
				Ω(manifest.ResourcePools[0].Stemcell.URL).Should(ContainSubstring("azure"))
			})

			It("then it should have the correct job config to deploy a bosh", func() {
				Ω(len(manifest.Jobs)).Should(Equal(1))
			})

			XIt("then it should properly define job properties", func() {
				Ω(len(manifest.Jobs[0].Properties)).Should(Equal(9))
				Ω(func() (r []interface{}) {
					for n, _ := range manifest.Jobs[0].Properties {
						r = append(r, n)
					}
					return
				}()).Should(ConsistOf("director", "nats", "hm", "postgres", "blobstore", "registry", "ntp", "agent", "azure"))
			})

			It("then it should properly define job templates", func() {
				Ω(len(manifest.Jobs[0].Templates)).Should(Equal(8))
				Ω(func() (r []string) {
					for _, v := range manifest.Jobs[0].Templates {
						r = append(r, v.Name)
					}
					return
				}()).Should(ConsistOf("nats", "postgres", "blobstore", "director", "health_monitor", "registry", "uaa", "cpi"))
			})

			It("then it should properly define job networks", func() {
				Ω(len(manifest.Jobs[0].Networks)).Should(Equal(2))
				Ω(func() (r []string) {
					for _, v := range manifest.Jobs[0].Networks {
						r = append(r, v.Name)
					}
					return
				}()).Should(ConsistOf("private", "public"))
			})
		})
	})
})
