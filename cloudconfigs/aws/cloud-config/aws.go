package cloudconfig

import (
	"errors"

	"github.com/xchapter7x/enaml"
	"github.com/xchapter7x/enaml/cloudproperties/aws"
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
	AZ1Name               string
	AZ2Name               string
	SubnetPropertyName1   string
	SubnetPropertyName2   string
	DefaultSecurityGroups []string
)

func NewAWSCloudConfig(region string, azList, subnetList, securityGroupList []string) (awsCloudConfig *enaml.CloudConfigManifest) {
	if err := validateFlags(region, azList, subnetList, securityGroupList); err != nil {
		lo.G.Error(err)
		return
	}
	DefaultSecurityGroups = securityGroupList
	SubnetPropertyName1 = subnetList[0]
	SubnetPropertyName2 = subnetList[1]
	AZ1Name = azList[0]
	AZ2Name = azList[1]
	Region = region

	awsCloudConfig = &enaml.CloudConfigManifest{}
	AddAZs(awsCloudConfig)
	AddDisk(awsCloudConfig)
	AddNetwork(awsCloudConfig)
	AddVMTypes(awsCloudConfig)
	AddCompilation(awsCloudConfig, AZ1Name, MediumVMName, PrivateNetworkName)
	return
}

func validateFlags(region string, azList, subnetList, securityGroupList []string) error {

	if len(securityGroupList) == 0 {
		return errors.New("!!!!!!!!!!\n\nyou should give at least one security group\n\n!!!!!!!!!!!")
	}

	if len(subnetList) < 1 {
		return errors.New("!!!!!!!!!!!\n\nyou have not given any subnets\n\n!!!!!!!!!!!")
	}

	if len(azList) < 1 {
		return errors.New("!!!!!!!!!!!\n\nyou have not given any AZs to use\n\n!!!!!!!!!!!")
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

func AddAZs(cfg *enaml.CloudConfigManifest) {
	cfg.AddAZ(enaml.AZ{
		Name: AZ1Name,
		CloudProperties: awscloudproperties.AZ{
			AvailabilityZoneName: Region + "a",
			SecurityGroups:       DefaultSecurityGroups,
		},
	})
	cfg.AddAZ(enaml.AZ{
		Name: AZ2Name,
		CloudProperties: awscloudproperties.AZ{
			AvailabilityZoneName: Region + "b",
			SecurityGroups:       DefaultSecurityGroups,
		},
	})
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

func AddNetwork(cfg *enaml.CloudConfigManifest) {
	octet1 := "10.0.0"
	octet2 := "10.10.0"
	dns := octet1 + ".2"
	privateNetwork := enaml.NewManualNetwork(PrivateNetworkName)
	privateNetwork.AddSubnet(createSubnet(octet1, dns, AZ1Name, SubnetPropertyName1))
	privateNetwork.AddSubnet(createSubnet(octet2, dns, AZ2Name, SubnetPropertyName2))
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
