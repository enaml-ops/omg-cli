package boshinit_test

import (
	"errors"

	"github.com/enaml-ops/enaml"
	. "github.com/enaml-ops/omg-cli/plugins/products/bosh-init"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("NewAWSBosh", func() {
	Describe("given the function", func() {
		Context("when called w/ valid parameters", func() {
			var boshConfig = BoshInitConfig{
				Name:                      "bosh",
				BoshReleaseVersion:        "256.2",
				BoshPrivateIP:             "10.0.0.4",
				BoshCPIReleaseVersion:     "11",
				GoAgentVersion:            "3012",
				BoshReleaseSHA:            "ff2f4e16e02f66b31c595196052a809100cfd5a8",
				BoshCPIReleaseSHA:         "dc4a0cca3b33dce291e4fbeb9e9948b6a7be3324",
				GoAgentSHA:                "3380b55948abe4c437dee97f67d2d8df4eec3fc1",
				BoshInstanceSize:          "Standard_D1",
				BoshDirectorName:          "my-bosh",
				BoshCIDR:                  "10.0.0.0/24",
				BoshGateway:               "10.0.0.1",
				BoshDNS:                   []string{"168.63.129.16"},
				AzurePublicIP:             "x.x.x.x",
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
			var manifest *enaml.DeploymentManifest

			BeforeEach(func() {
				manifest = NewAzureBosh(boshConfig)
			})

			It("then it should create valid azure job properties", func() {
				Ω(func() error {
					for n, v := range manifest.Jobs[0].Properties {
						if n == "azure" && v == nil {
							return errors.New("azure is nil, thats not right")
						}
					}
					return nil
				}()).ShouldNot(HaveOccurred())
			})

			It("then it should be using the azure stemcell", func() {
				Ω(manifest.ResourcePools[0].Stemcell.URL).Should(ContainSubstring("azure"))
			})

			It("then it should have the correct job config to deploy a bosh", func() {
				Ω(len(manifest.Jobs)).Should(Equal(1))
			})

			It("then it should properly define job properties", func() {
				Ω(len(manifest.Jobs[0].Properties)).Should(Equal(9))
				Ω(func() (r []string) {
					for n, _ := range manifest.Jobs[0].Properties {
						r = append(r, n)
					}
					return
				}()).Should(ConsistOf("director", "nats", "hm", "postgres", "blobstore", "registry", "ntp", "agent", "azure"))
			})

			It("then it should properly define job templates", func() {
				Ω(len(manifest.Jobs[0].Templates)).Should(Equal(7))
				Ω(func() (r []string) {
					for _, v := range manifest.Jobs[0].Templates {
						r = append(r, v.Name)
					}
					return
				}()).Should(ConsistOf("nats", "postgres", "blobstore", "director", "health_monitor", "registry", "cpi"))
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
