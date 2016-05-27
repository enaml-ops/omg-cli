package cloudconfig_test

import (
	"errors"
	"fmt"

	"github.com/enaml-ops/enaml"
	"github.com/enaml-ops/enaml/cloudproperties/aws"
	. "github.com/enaml-ops/omg-cli/plugins/cloudconfigs/aws/cloud-config"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("given CloudConfig Deployment for AWS", func() {
	var awsConfig *enaml.CloudConfigManifest
	BeforeEach(func() {
		awsConfig = NewAWSCloudConfig(awscloudproperties.USWest, map[string]string{"us-east-1b": "subnet-asdfasdfa", "us-east-2c": "subnet-4444444"}, []string{"security-group-for-something-secure"})
	})

	Context("when AZs are defined", func() {
		It("then each AZ definition should map to a unique aws AZ", func() {
			err := checkUniqueAZs(awsConfig.AZs)
			Ω(err).ShouldNot(HaveOccurred())
		})
	})

	Context("when a user of the iaas would like to define vm types", func() {
		It("then there should be 2 vm type options available", func() {
			Ω(len(awsConfig.VMTypes)).Should(Equal(2))
		})

		It("then they should have the option of a small VM configuration", func() {
			_, err := getVmTypeByName(SmallVMName, awsConfig.VMTypes)
			Ω(err).ShouldNot(HaveOccurred())
		})

		Context("when the vmtype is small", func() {
			var vm enaml.VMType
			BeforeEach(func() {
				vm, _ = getVmTypeByName(SmallVMName, awsConfig.VMTypes)
			})
			It("then it should use a t2.micro size aws instance", func() {
				Ω(vm.CloudProperties.(awscloudproperties.VMType).InstanceType).Should(Equal(SmallVMSize))
			})

			It("then it should use a properly configured ephemeral disk", func() {
				properSmallDiskSize := SmallEphemeralDiskSize
				properDiskType := SmallDiskType
				Ω(vm.CloudProperties.(awscloudproperties.VMType).EphemeralDisk.Size).Should(Equal(properSmallDiskSize))
				Ω(vm.CloudProperties.(awscloudproperties.VMType).EphemeralDisk.DiskType).Should(Equal(properDiskType))
			})
		})

		It("then they should have the option of a large VM configuration", func() {
			_, err := getVmTypeByName(MediumVMName, awsConfig.VMTypes)
			Ω(err).ShouldNot(HaveOccurred())
		})

		Context("when the vmtype is large", func() {
			var vm enaml.VMType
			BeforeEach(func() {
				vm, _ = getVmTypeByName(MediumVMName, awsConfig.VMTypes)
			})
			It("then it should use a m3.medium size aws instance", func() {
				Ω(vm.CloudProperties.(awscloudproperties.VMType).InstanceType).Should(Equal(MediumVMSize))
			})

			It("then it should use a properly configured ephemeral disk", func() {
				properMediumDiskSize := MediumEphemeralDiskSize
				properDiskType := MediumDiskType
				Ω(vm.CloudProperties.(awscloudproperties.VMType).EphemeralDisk.Size).Should(Equal(properMediumDiskSize))
				Ω(vm.CloudProperties.(awscloudproperties.VMType).EphemeralDisk.DiskType).Should(Equal(properDiskType))
			})
		})
	})

	Context("when a user of the iaas would like to assign disk", func() {
		It("then all disk types should be properly configured", func() {
			for _, v := range awsConfig.DiskTypes {
				Ω(v.Name).ShouldNot(BeEmpty())
				Ω(v.DiskSize).Should(BeNumerically(">", 0))
				Ω(v.CloudProperties).ShouldNot(BeNil())
			}
		})
		It("then they should have the option of a small capacity configuration", func() {
			err := checkDiskTypeExists(awsConfig.DiskTypes, DiskSmallName)
			Ω(err).ShouldNot(HaveOccurred())
		})
		It("then they should have the option of a medium capacity configuration", func() {
			err := checkDiskTypeExists(awsConfig.DiskTypes, DiskMediumName)
			Ω(err).ShouldNot(HaveOccurred())
		})
		It("then they should have the option of a large capacity configuration", func() {
			err := checkDiskTypeExists(awsConfig.DiskTypes, DiskLargeName)
			Ω(err).ShouldNot(HaveOccurred())
		})
	})

	Context("when a user of the iaas would like to assign a network", func() {
		It("then they should have a private and a vip network", func() {
			var networkList []string
			for _, v := range awsConfig.Networks {
				switch v.(type) {
				case enaml.ManualNetwork:
					networkList = append(networkList, v.(enaml.ManualNetwork).Name)
				case enaml.VIPNetwork:
					networkList = append(networkList, v.(enaml.VIPNetwork).Name)
				case enaml.DynamicNetwork:
					networkList = append(networkList, v.(enaml.DynamicNetwork).Name)
				}
			}
			Ω(networkList).Should(ContainElement(PrivateNetworkName))
			Ω(networkList).Should(ContainElement(VIPNetworkName))
		})

		Context("when a user of the iaas is assigning the private network", func() {
			var privateNetwork enaml.ManualNetwork
			BeforeEach(func() {
				for _, v := range awsConfig.Networks {
					var name string
					switch v.(type) {
					case enaml.ManualNetwork:
						name = v.(enaml.ManualNetwork).Name
					case enaml.VIPNetwork:
						name = v.(enaml.VIPNetwork).Name
					case enaml.DynamicNetwork:
						name = v.(enaml.DynamicNetwork).Name
					}
					if name == PrivateNetworkName {
						privateNetwork = v.(enaml.ManualNetwork)
					}
				}
			})

			It("then they should have one subnet for each configured AZ", func() {
				Ω(len(privateNetwork.Subnets)).Should(Equal(len(awsConfig.AZs)))
			})

			It("then each subnet should be configured w/ required fields", func() {
				for _, v := range privateNetwork.Subnets {
					Ω(v.Range).ShouldNot(BeEmpty())
					Ω(v.AZ).ShouldNot(BeEmpty())
					Ω(v.Gateway).ShouldNot(BeEmpty())
					Ω(v.DNS).ShouldNot(BeEmpty())
				}
			})
		})
	})
})

func getVmTypeByName(name string, vmTypes []enaml.VMType) (res enaml.VMType, err error) {
	err = errors.New("no type found")
	for _, k := range vmTypes {
		if k.Name == name {
			err = nil
			res = k
		}
	}
	return
}

func checkUniqueAZs(azs []enaml.AZ) error {
	exists := make(map[string]int)
	for _, v := range azs {
		awsAZ := v.CloudProperties.(awscloudproperties.AZ).AvailabilityZoneName
		if _, alreadyExists := exists[awsAZ]; alreadyExists {
			return errors.New(fmt.Sprintf("duplicate az assignment to: %s", awsAZ))
		}
		exists[awsAZ] = 1
	}
	return nil
}

func checkDiskTypeExists(dsk []enaml.DiskType, name string) (err error) {
	err = errors.New(name + " capacity configuration not found")
	for _, v := range dsk {
		if v.Name == name {
			err = nil
			break
		}
	}
	return
}
