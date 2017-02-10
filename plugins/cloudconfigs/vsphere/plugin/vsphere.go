package plugin

import (
	"fmt"

	"gopkg.in/urfave/cli.v2"
	"gopkg.in/yaml.v2"

	"github.com/enaml-ops/enaml"
	"github.com/enaml-ops/omg-cli/plugins/cloudconfigs"
	"github.com/enaml-ops/omg-cli/plugins/cloudconfigs/vsphere/plugin/generated"
)

type AZCloudProperties struct {
	DataCenters []DataCenter `yaml:"datacenters"`
}
type DataCenter struct {
	Name     string                   `yaml:"name"`
	Clusters []map[string]interface{} `yaml:"clusters"`
}
type ResourcePool struct {
	ResourcePool string `yaml:"resource_pool"`
}
type VMProperties struct {
	CPU        int      `yaml:"cpu,omitempty"`
	RAM        int      `yaml:"ram,omitempty"`
	Disk       int      `yaml:"disk,omitempty"`
	DataStores []string `yaml:"datastores,omitempty"`
}
type NetworkCloudProperties struct {
	NetworkName string `yaml:"name"`
}
type VSphereCloudConfig struct {
	Context *cli.Context
}

func NewVSphereCloudConfig(c *cli.Context) cloudconfigs.CloudConfigProvider {

	provider := &VSphereCloudConfig{
		Context: c,
	}
	return provider
}

func (c *VSphereCloudConfig) networkCloudProperties(i, index int) interface{} {
	networkNames := c.Context.StringSlice(fmt.Sprintf("vsphere-network-name-%d", i))
	return NetworkCloudProperties{
		NetworkName: networkNames[index],
	}
}

func (c *VSphereCloudConfig) validateCloudProperties(length, i int) error {
	multiAssignAZ := c.Context.Bool("multi-assign-az")
	if multiAssignAZ {
		return cloudconfigs.CheckRequiredLength(1, i, c.Context, "vsphere-network-name-%d")
	} else {
		return cloudconfigs.CheckRequiredLength(length, i, c.Context, "vsphere-network-name-%d")
	}

}

func (c *VSphereCloudConfig) CreateNetworks() ([]enaml.DeploymentNetwork, error) {
	networks, err := cloudconfigs.CreateNetworks(c.Context, c.validateCloudProperties, c.networkCloudProperties)
	return networks, err
}

func clusterConfig(clusterName, resourcePoolName string) (clusters []map[string]interface{}) {
	cluster := make(map[string]interface{})
	if resourcePoolName != "" {
		cluster[clusterName] = ResourcePool{
			ResourcePool: resourcePoolName,
		}
	} else {
		cluster[clusterName] = make(map[string]string, 0)
	}
	clusters = append(clusters, cluster)
	return
}
func (c *VSphereCloudConfig) CreateAZs() ([]enaml.AZ, error) {
	azNames := c.Context.StringSlice("az")
	datacenters := c.Context.StringSlice("vsphere-datacenter")
	clusters := c.Context.StringSlice("vsphere-cluster")
	resourcePools := c.Context.StringSlice("vsphere-resource-pool")
	if len(azNames) != len(datacenters) {
		err := fmt.Errorf("Sorry you need to provide the same number of az and vsphere-datacenter flags")
		return nil, err
	}
	if len(azNames) != len(clusters) {
		err := fmt.Errorf("Sorry you need to provide the same number of az and vsphere-cluster flags")
		return nil, err
	}
	if len(resourcePools) > 0 {
		if len(azNames) != len(resourcePools) {
			err := fmt.Errorf("Sorry you need to provide the same number of az and vsphere-resource-pool flags")
			return nil, err
		}
	}

	azs := []enaml.AZ{}

	for i, azName := range azNames {
		azCloudProperties := AZCloudProperties{}
		var resourcePoolName = ""
		if len(resourcePools) > 0 {
			resourcePoolName = resourcePools[i]
		}
		dataCenter := DataCenter{
			Name:     datacenters[i],
			Clusters: clusterConfig(clusters[i], resourcePoolName),
		}
		azCloudProperties.DataCenters = append(azCloudProperties.DataCenters, dataCenter)
		az := enaml.AZ{
			Name:            azName,
			CloudProperties: azCloudProperties,
		}
		azs = append(azs, az)
	}
	return azs, nil
}

func (c *VSphereCloudConfig) CreateVMTypes() ([]enaml.VMType, error) {

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

func (c *VSphereCloudConfig) CreateDiskTypes() ([]enaml.DiskType, error) {
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

func (c *VSphereCloudConfig) CreateCompilation() (*enaml.Compilation, error) {
	return cloudconfigs.CreateCompilation(c.Context)
}
