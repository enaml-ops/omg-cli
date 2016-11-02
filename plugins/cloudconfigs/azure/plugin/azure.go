package plugin

import (
	"fmt"

	"github.com/enaml-ops/enaml"
	"github.com/enaml-ops/omg-cli/plugins/cloudconfigs"
	"github.com/enaml-ops/omg-cli/plugins/cloudconfigs/azure/plugin/generated"
	"gopkg.in/urfave/cli.v2"
	"gopkg.in/yaml.v2"
)

type NetworkCloudProperties struct {
	NetworkName string `yaml:"virtual_network_name"`
	SubnetName  string `yaml:"subnet_name"`
}
type AzureCloudConfig struct {
	Context *cli.Context
}

//NewAzureCloudConfig : Create Azure cloud config
func NewAzureCloudConfig(c *cli.Context) cloudconfigs.CloudConfigProvider {
	provider := &AzureCloudConfig{
		Context: c,
	}
	return provider
}

func (c *AzureCloudConfig) networkCloudProperties(i, index int) interface{} {
	azureNetworkNames := c.Context.StringSlice(fmt.Sprintf("azure-virtual-network-name-%d", i))
	azureSubnetNames := c.Context.StringSlice(fmt.Sprintf("azure-subnet-name-%d", i))
	return NetworkCloudProperties{
		NetworkName: azureNetworkNames[index],
		SubnetName:  azureSubnetNames[index],
	}
}
func (c *AzureCloudConfig) validateCloudProperties(length, i int) error {
	multiAssignAZ := c.Context.Bool("multi-assign-az")
	if multiAssignAZ {
		return cloudconfigs.CheckRequiredLength(1, i, c.Context, "azure-virtual-network-name-%d", "azure-subnet-name-%d")
	}
	return cloudconfigs.CheckRequiredLength(length, i, c.Context, "azure-virtual-network-name-%d", "azure-subnet-name-%d")
}

//CreateNetworks : Create Azure specific network configuration
func (c *AzureCloudConfig) CreateNetworks() ([]enaml.DeploymentNetwork, error) {
	networks, err := cloudconfigs.CreateNetworks(c.Context, c.validateCloudProperties, c.networkCloudProperties)
	return networks, err
}

//CreateAZs : Create Azure specific availability zone configuration
func (c *AzureCloudConfig) CreateAZs() ([]enaml.AZ, error) {
	azNames := c.Context.StringSlice("az")
	azs := []enaml.AZ{}
	for _, azName := range azNames {
		az := enaml.AZ{
			Name: azName,
		}
		azs = append(azs, az)
	}
	return azs, nil
}

//CreateDiskTypes : Returns Azure disk types
func (c *AzureCloudConfig) CreateDiskTypes() ([]enaml.DiskType, error) {
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

//CreateVMTypes : Create Azure specific VM types
func (c *AzureCloudConfig) CreateVMTypes() ([]enaml.VMType, error) {
	var vmTypes []enaml.VMType
	if fileBytes, err := generated.Asset("files/vm_types.yml"); err == nil {
		if err = yaml.Unmarshal(fileBytes, &vmTypes); err == nil {
			return vmTypes, nil
		}
		return nil, err
	} else {
		return nil, err
	}
}

//CreateCompilation : Create Azure specific VM types
func (c *AzureCloudConfig) CreateCompilation() (*enaml.Compilation, error) {
	return cloudconfigs.CreateCompilation(c.Context)
}
