package boshinit_test

import (
	. "github.com/bosh-ops/bosh-install/deployments/bosh-init"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/xchapter7x/enaml"
)

var _ = Describe("NewAWSBosh", func() {
	Describe("given the function", func() {
		Context("when called w/ valid parameters", func() {
			var boshConfig = BoshInitConfig{
				Name:                  "bosh",
				BoshReleaseVersion:    "256.2",
				BoshPrivateIP:         "10.0.0.6",
				BoshCPIReleaseVersion: "52",
				GoAgentVersion:        "3012",
				BoshReleaseSHA:        "ff2f4e16e02f66b31c595196052a809100cfd5a8",
				BoshCPIReleaseSHA:     "dc4a0cca3b33dce291e4fbeb9e9948b6a7be3324",
				GoAgentSHA:            "3380b55948abe4c437dee97f67d2d8df4eec3fc1",
				BoshInstanceSize:      "m3.xlarge",
				BoshAvailabilityZone:  "us-east-1c",
				AWSSubnet:             "subnet-xxxxxx",
				AWSElasticIP:          "1.0.2.3",
				BoshDirectorName:      "my-bosh",
				AWSPEMFilePath:        "./some.pem",
				AWSAccessKeyID:        "xxxxxxx",
				AWSSecretKey:          "xxxxxxxxxxxxxxxxxxxx",
				AWSRegion:             "us-east-1",
			}
			var manifest *enaml.DeploymentManifest

			BeforeEach(func() {
				manifest = NewAWSBosh(boshConfig)
			})

			It("then it should be using the aws stemcell", func() {
				Ω(manifest.ResourcePools[0].Stemcell.URL).ShouldNot(ContainSubstring("azure"))
				Ω(manifest.ResourcePools[0].Stemcell.URL).Should(ContainSubstring("aws"))
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
				}()).Should(ConsistOf("director", "nats", "hm", "postgres", "blobstore", "registry", "ntp", "agent", "aws"))
			})

			It("then it should properly define job templates", func() {
				Ω(len(manifest.Jobs[0].Templates)).Should(Equal(7))
				Ω(func() (r []string) {
					for _, v := range manifest.Jobs[0].Templates {
						r = append(r, v.Name)
					}
					return
				}()).Should(ConsistOf("nats", "postgres", "blobstore", "director", "health_monitor", "registry", "aws_cpi"))
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
