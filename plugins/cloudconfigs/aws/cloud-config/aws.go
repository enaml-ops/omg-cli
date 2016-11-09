package cloudconfig

import (
	"fmt"

	"gopkg.in/urfave/cli.v2"
	"gopkg.in/yaml.v2"

	"github.com/enaml-ops/enaml"
	"github.com/enaml-ops/enaml/cloudproperties/aws"
	"github.com/enaml-ops/omg-cli/plugins/cloudconfigs"
	"github.com/enaml-ops/omg-cli/plugins/cloudconfigs/aws/cloud-config/generated"
)

type AWSCloudConfig struct {
	Context *cli.Context
}

type AZCloudProperties struct {
	AvailabiltyZone string `yaml:"availability_zone"`
}

func NewAWSCloudConfig(c *cli.Context) cloudconfigs.CloudConfigProvider {

	provider := &AWSCloudConfig{
		Context: c,
	}
	return provider
}

func (c *AWSCloudConfig) CreateAZs() ([]enaml.AZ, error) {
	azNames := c.Context.StringSlice("az")
	availablityZones := c.Context.StringSlice("aws-availablity-zone")
	azs := []enaml.AZ{}

	for i, azName := range azNames {
		az := enaml.AZ{
			Name: azName,
			CloudProperties: AZCloudProperties{
				AvailabiltyZone: availablityZones[i],
			},
		}
		azs = append(azs, az)
	}
	return azs, nil
}

/*func NewAWSCloudConfig(region string, securityGroupList []string, subnets []SubnetBucket) (awsCloudConfig *enaml.CloudConfigManifest) {
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
}*/

func (c *AWSCloudConfig) CreateCompilation() (*enaml.Compilation, error) {
	return cloudconfigs.CreateCompilation(c.Context)
}

func (c *AWSCloudConfig) CreateNetworks() ([]enaml.DeploymentNetwork, error) {
	networks, err := cloudconfigs.CreateNetworks(c.Context, c.validateCloudProperties, c.networkCloudProperties)
	return networks, err
}

func (c *AWSCloudConfig) networkCloudProperties(i, index int) interface{} {
	networkNames := c.Context.StringSlice(fmt.Sprintf("aws-subnet-name-%d", i))
	securityGroupNames := c.Context.StringSlice(fmt.Sprintf("aws-security-group-%d", i))
	return awscloudproperties.Network{
		Subnet:         networkNames[index],
		SecurityGroups: securityGroupNames,
	}
}

func (c *AWSCloudConfig) validateCloudProperties(length, i int) error {
	return cloudconfigs.CheckRequiredLength(length, i, c.Context, "aws-subnet-name-%d")
}

func (c *AWSCloudConfig) CreateVMTypes() ([]enaml.VMType, error) {

	var vmTypes []enaml.VMType
	if fileBytes, err := generated.Asset("files/vm_types.yml"); err == nil {
		if err = yaml.Unmarshal(fileBytes, &vmTypes); err == nil {
			return vmTypes, nil
		} else {
			return nil, err
		}
	} else {
		return nil, err
	}
}

func (c *AWSCloudConfig) CreateDiskTypes() ([]enaml.DiskType, error) {
	var diskTypes []enaml.DiskType
	if fileBytes, err := generated.Asset("files/disk_types.yml"); err == nil {
		if err = yaml.Unmarshal(fileBytes, &diskTypes); err == nil {
			return diskTypes, nil
		} else {
			return nil, err
		}
	} else {
		return nil, err
	}
}
