package plugin

import (
	"fmt"

	"github.com/enaml-ops/enaml"
	"github.com/enaml-ops/omg-cli/plugins/cloudconfigs"
	"github.com/enaml-ops/omg-cli/plugins/cloudconfigs/azure/plugin/generated"
	"gopkg.in/urfave/cli.v2"
	"gopkg.in/yaml.v2"
)

type networkCloudProperties struct {
	NetworkName string `yaml:"virtual_network_name"`
	SubnetName  string `yaml:"subnet_name"`
}
type azureCloudConfig struct {
	Context *cli.Context
}

// NewAzureCloudConfig creates an Azure cloud config from the specified context
func NewAzureCloudConfig(c *cli.Context) cloudconfigs.CloudConfigProvider {
	provider := &azureCloudConfig{
		Context: c,
	}
	return provider
}

func (c *azureCloudConfig) networkCloudProperties(i, index int) interface{} {
	azureNetworkNames := c.Context.StringSlice(fmt.Sprintf("azure-virtual-network-name-%d", i))
	azureSubnetNames := c.Context.StringSlice(fmt.Sprintf("azure-subnet-name-%d", i))
	return networkCloudProperties{
		NetworkName: azureNetworkNames[index],
		SubnetName:  azureSubnetNames[index],
	}
}
func (c *azureCloudConfig) validateCloudProperties(length, i int) error {
	multiAssignAZ := c.Context.Bool("multi-assign-az")
	if multiAssignAZ {
		return cloudconfigs.CheckRequiredLength(1, i, c.Context, "azure-virtual-network-name-%d", "azure-subnet-name-%d")
	}
	return cloudconfigs.CheckRequiredLength(length, i, c.Context, "azure-virtual-network-name-%d", "azure-subnet-name-%d")
}

//CreateNetworks : Create Azure specific network configuration
func (c *azureCloudConfig) CreateNetworks() ([]enaml.DeploymentNetwork, error) {
	return cloudconfigs.CreateNetworks(c.Context, c.validateCloudProperties, c.networkCloudProperties)
}

//CreateAZs : Create Azure specific availability zone configuration
func (c *azureCloudConfig) CreateAZs() ([]enaml.AZ, error) {
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
func (c *azureCloudConfig) CreateDiskTypes() ([]enaml.DiskType, error) {
	fileBytes, err := generated.Asset("files/disk_types.yml")
	if err != nil {
		return nil, err
	}
	var diskTypes []enaml.DiskType
	err = yaml.Unmarshal(fileBytes, &diskTypes)
	return diskTypes, err
}

//CreateVMTypes : Create Azure specific VM types
func (c *azureCloudConfig) CreateVMTypes() ([]enaml.VMType, error) {
	fileBytes, err := generated.Asset("files/vm_types.yml")
	if err != nil {
		return nil, err
	}
	var vmTypes []enaml.VMType
	err = yaml.Unmarshal(fileBytes, &vmTypes)
	return vmTypes, err
}

//CreateCompilation : Create Azure specific VM types
func (c *azureCloudConfig) CreateCompilation() (*enaml.Compilation, error) {
	return cloudconfigs.CreateCompilation(c.Context)
}
