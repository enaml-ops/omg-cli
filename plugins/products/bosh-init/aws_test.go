package boshinit_test

import (
	"github.com/enaml-ops/enaml"
	. "github.com/enaml-ops/omg-cli/plugins/products/bosh-init"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("NewAWSBosh", func() {
	Describe("given the function", func() {
		Context("when called w/ valid parameters", func() {
			var boshConfig = BoshInitConfig{
				BoshInstanceSize:     "m3.xlarge",
				BoshAvailabilityZone: "us-east-1c",
				AWSSubnet:            "subnet-xxxxxx",
				AWSPEMFilePath:       "./some.pem",
				AWSAccessKeyID:       "xxxxxxx",
				AWSSecretKey:         "xxxxxxxxxxxxxxxxxxxx",
				AWSRegion:            "us-east-1",
			}
			var boshBase = &BoshBase{
				CPIName:            "aws_cpi",
				BoshReleaseVersion: "256.2",
				PrivateIP:          "10.0.0.6",
				PublicIP:           "1.0.2.3",
				CPIReleaseVersion:  "52",
				GOAgentVersion:     "3012",
				BoshReleaseSHA:     "ff2f4e16e02f66b31c595196052a809100cfd5a8",
				CPIReleaseSHA:      "dc4a0cca3b33dce291e4fbeb9e9948b6a7be3324",
				GOAgentSHA:         "3380b55948abe4c437dee97f67d2d8df4eec3fc1",
				NetworkCIDR:        "10.0.0.0/24",
				NetworkGateway:     "10.0.0.1",
				NetworkDNS:         []string{"10.0.0.2"},
				DirectorName:       "my-bosh",
			}
			var manifest *enaml.DeploymentManifest

			BeforeEach(func() {
				manifest = NewAWSBosh(boshConfig, boshBase)
			})

			It("then it should be using the aws stemcell", func() {
				Ω(manifest.ResourcePools[0].Stemcell.URL).ShouldNot(ContainSubstring("azure"))
				Ω(manifest.ResourcePools[0].Stemcell.URL).Should(ContainSubstring("aws"))
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
				}()).Should(ConsistOf("director", "nats", "hm", "postgres", "blobstore", "registry", "ntp", "agent", "aws"))
			})

			It("then it should properly define job templates", func() {
				Ω(len(manifest.Jobs[0].Templates)).Should(Equal(8))
				Ω(func() (r []string) {
					for _, v := range manifest.Jobs[0].Templates {
						r = append(r, v.Name)
					}
					return
				}()).Should(ConsistOf("nats", "postgres", "blobstore", "director", "health_monitor", "uaa", "registry", "aws_cpi"))
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
