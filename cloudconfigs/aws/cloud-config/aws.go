package cloudconfig

import (
	"errors"

	"github.com/enaml-ops/enaml"
	"github.com/enaml-ops/enaml/cloudproperties/aws"
	"github.com/xchapter7x/lo"
)

const (
	DefaultDiskType         = "gp2"
	DiskSmallName           = "small"
	DiskMediumName          = "medium"
	DiskLargeName           = "large"
	SmallVMName             = "small"
	SmallVMSize             = "t2.micro"
	MediumVMName            = "medium"
	MediumVMSize            = "m3.medium"
	MediumDiskType          = DefaultDiskType
	MediumEphemeralDiskSize = 30000
	SmallDiskType           = DefaultDiskType
	SmallEphemeralDiskSize  = 3000
	PrivateNetworkName      = "private"
	VIPNetworkName          = "vip"
)

var (
	Region                string
	DefaultSecurityGroups []string
)

func NewAWSCloudConfig(region string, azSubnetMap map[string]string, securityGroupList []string) (awsCloudConfig *enaml.CloudConfigManifest) {
	if err := validateFlags(region, azSubnetMap, securityGroupList); err != nil {
		lo.G.Error(err)
		return
	}
	DefaultSecurityGroups = securityGroupList
	Region = region

	awsCloudConfig = &enaml.CloudConfigManifest{}
	AddAZs(awsCloudConfig, azSubnetMap)
	AddDisk(awsCloudConfig)
	AddNetwork(awsCloudConfig, azSubnetMap)
	AddVMTypes(awsCloudConfig)

	for azname, _ := range azSubnetMap {
		AddCompilation(awsCloudConfig, azname, MediumVMName, PrivateNetworkName)
		break
	}
	return
}

func validateFlags(region string, azSubnetMap map[string]string, securityGroupList []string) error {

	if len(securityGroupList) == 0 {
		return errors.New("!!!!!!!!!!\n\nyou should give at least one security group\n\n!!!!!!!!!!!")
	}

	if len(azSubnetMap) < 1 {
		return errors.New("!!!!!!!!!!!\n\nyou have not given any subnets\n\n!!!!!!!!!!!")
	}

	if region == "" {
		return errors.New("!!!!!!!!!!!\n\nyou have not given a region to use\n\n!!!!!!!!!!!")
	}
	return nil
}

func AddCompilation(cfg *enaml.CloudConfigManifest, az string, vmtype string, network string) {
	cfg.SetCompilation(&enaml.Compilation{
		Workers:             5,
		ReuseCompilationVMs: true,
		AZ:                  az,
		VMType:              vmtype,
		Network:             network,
	})
}

func AddAZs(cfg *enaml.CloudConfigManifest, azSubnetMap map[string]string) {
	for azname, _ := range azSubnetMap {
		cfg.AddAZ(enaml.AZ{
			Name: azname,
			CloudProperties: awscloudproperties.AZ{
				AvailabilityZoneName: azname,
				SecurityGroups:       DefaultSecurityGroups,
			},
		})
	}
}

func AddDisk(cfg *enaml.CloudConfigManifest) {
	cfg.AddDiskType(createDiskType(DiskSmallName, 3000, DefaultDiskType))
	cfg.AddDiskType(createDiskType(DiskMediumName, 20000, DefaultDiskType))
	cfg.AddDiskType(createDiskType(DiskLargeName, 50000, DefaultDiskType))
}

func createDiskType(name string, size int, typename string) enaml.DiskType {
	return enaml.DiskType{
		Name:     name,
		DiskSize: size,
		CloudProperties: awscloudproperties.EphemeralDisk{
			DiskType: typename,
		}}
}

func AddNetwork(cfg *enaml.CloudConfigManifest, azSubnetMap map[string]string) {
	octet1 := "10.0.0"
	dns := octet1 + ".2"
	privateNetwork := enaml.NewManualNetwork(PrivateNetworkName)

	for azname, subnetname := range azSubnetMap {
		privateNetwork.AddSubnet(createSubnet(octet1, dns, azname, subnetname))
	}
	cfg.AddNetwork(privateNetwork)
	cfg.AddNetwork(enaml.NewVIPNetwork(VIPNetworkName))
}

func createSubnet(octet, dns, azname, subnetPropertyName string) enaml.Subnet {
	subnet := enaml.NewSubnet(octet, azname)
	subnet.AddDNS(dns)
	subnet.AddReserved(octet + ".1-" + octet + ".10")
	subnet.CloudProperties = awscloudproperties.Network{
		Subnet: subnetPropertyName,
	}
	return subnet
}

func AddVMTypes(cfg *enaml.CloudConfigManifest) {
	cfg.AddVMType(enaml.VMType{
		Name:            SmallVMName,
		CloudProperties: NewVMCloudProperty(SmallVMSize, SmallDiskType, SmallEphemeralDiskSize),
	})
	cfg.AddVMType(enaml.VMType{
		Name:            MediumVMName,
		CloudProperties: NewVMCloudProperty(MediumVMSize, MediumDiskType, MediumEphemeralDiskSize),
	})
}

func NewVMCloudProperty(instanceType, diskType string, diskSize int) awscloudproperties.VMType {
	return awscloudproperties.VMType{
		InstanceType: instanceType,
		EphemeralDisk: awscloudproperties.EphemeralDisk{
			Size:     diskSize,
			DiskType: diskType,
		},
	}
}
