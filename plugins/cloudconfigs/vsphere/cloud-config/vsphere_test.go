package cloudconfig_test

import (
	"fmt"

	"github.com/enaml-ops/enaml"
	"github.com/enaml-ops/enaml/cloudproperties/aws"
	. "github.com/enaml-ops/omg-cli/plugins/cloudconfigs/vsphere/cloud-config"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("given CloudConfig Deployment for vSphere", func() {
	var vsphereConfig *enaml.CloudConfigManifest
	BeforeEach(func() {
		vsphereConfig = NewVSphereCloudConfig()
	})

	Context("when AZs are defined", func() {
		It("then each AZ definition should map to a unique vsphere AZ", func() {
			err := checkUniqueAZs(vsphereConfig.AZs)
			Ω(err).ShouldNot(HaveOccurred())
		})
	})
})

func checkUniqueAZs(azs []enaml.AZ) error {
	exists := make(map[string]int)
	for _, v := range azs {
		awsAZ := v.CloudProperties.(awscloudproperties.AZ).AvailabilityZoneName
		if _, alreadyExists := exists[awsAZ]; alreadyExists {
			return fmt.Errorf("duplicate az assignment to: %s", awsAZ)
		}
		exists[awsAZ] = 1
	}
	return nil
}
