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

type SubnetBucket struct {
	BoshAZName string
	Cidr string
	Gateway string
	DNS[] string
	AWSAZName string
	AWSSubnetName string
	BoshReserveRange []string
}

func NewAWSCloudConfig(region string, securityGroupList []string, subnets []SubnetBucket) (awsCloudConfig *enaml.CloudConfigManifest) {
	if err := validateFlags(region, subnets, securityGroupList); err != nil {
		lo.G.Error(err)
		return
	}
	DefaultSecurityGroups = securityGroupList
	Region = region

	awsCloudConfig = &enaml.CloudConfigManifest{}
	AddAZs(awsCloudConfig, subnets)
	AddDisk(awsCloudConfig)
	AddNetwork(awsCloudConfig, subnets)
	AddVMTypes(awsCloudConfig)

	for _, azname := range subnets {
		AddCompilation(awsCloudConfig, azname.BoshAZName, MediumVMName, PrivateNetworkName)
		break
	}
	return
}

func validateFlags(region string, subnets []SubnetBucket, securityGroupList []string) error {

	if len(securityGroupList) == 0 {
		return errors.New("!!!!!!!!!!\n\nyou should give at least one security group\n\n!!!!!!!!!!!")
	}

	if len(subnets) < 1 {
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

func AddAZs(cfg *enaml.CloudConfigManifest, subnets []SubnetBucket) {
	for _, az := range subnets {
		cfg.AddAZ(enaml.AZ{
			Name: az.BoshAZName,
			CloudProperties: awscloudproperties.AZ{
				AvailabilityZoneName: az.AWSAZName,
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

func AddNetwork(cfg *enaml.CloudConfigManifest, subnets []SubnetBucket) {
	privateNetwork := enaml.NewManualNetwork(PrivateNetworkName)

	for _, subnet := range subnets {
		privateNetwork.AddSubnet(createSubnet(subnet))
	}
	cfg.AddNetwork(privateNetwork)
	cfg.AddNetwork(enaml.NewVIPNetwork(VIPNetworkName))
}

func createSubnet(subnetBucket SubnetBucket) enaml.Subnet {
	subnet := enaml.NewSubnet(subnetBucket.Cidr, subnetBucket.Gateway, subnetBucket.BoshAZName)

	for _, dns := range subnetBucket.DNS {
		subnet.AddDNS(dns)
	}

	for _, r := range subnetBucket.BoshReserveRange {
		subnet.AddReserved(r)
	}
	subnet.CloudProperties = awscloudproperties.Network{
		Subnet: subnetBucket.AWSSubnetName,
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
