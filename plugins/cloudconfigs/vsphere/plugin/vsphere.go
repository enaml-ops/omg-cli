package plugin

import (
	"fmt"
	"strings"

	"github.com/codegangsta/cli"
	"github.com/enaml-ops/enaml"
	"github.com/enaml-ops/omg-cli/plugins/cloudconfigs"
)

type AZCloudProperties struct {
	DataCenters []DataCenter `yaml:"datacenters"`
	/*datacenters:
	  - name: my-dc
	    clusters:
	    - {my-vsphere-cluster: {resource_pool: my-vsphere-res-pool}}*/
}
type DataCenter struct {
	Name     string                    `yaml:"name"`
	Clusters []map[string]ResourcePool `yaml:"clusters"`
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

func (c *VSphereCloudConfig) CreateNetworks() ([]enaml.DeploymentNetwork, error) {
	context := c.Context
	networks := []enaml.DeploymentNetwork{}
	for i := 1; i <= SupportedNetworkCount; i++ {
		networkFlag := fmt.Sprintf("network-name-%d", i)
		if context.IsSet(networkFlag) {
			network := enaml.ManualNetwork{
				Name: context.String(networkFlag),
				Type: "manual",
			}
			azs := context.StringSlice(fmt.Sprintf("network-az-%d", i))
			if err := cloudconfigs.CheckRequiredLength(len(azs), i, context, "network-cidr-%d", "network-gateway-%d", "network-dns-%d", "network-reserved-%d", "network-static-%d"); err != nil {
				return nil, err
			}
			ranges := context.StringSlice(fmt.Sprintf("network-cidr-%d", i))
			gateways := context.StringSlice(fmt.Sprintf("network-gateway-%d", i))
			dnsServers := context.StringSlice(fmt.Sprintf("network-dns-%d", i))
			reservedRanges := context.StringSlice(fmt.Sprintf("network-reserved-%d", i))
			staticIPs := context.StringSlice(fmt.Sprintf("network-static-%d", i))
			if err := cloudconfigs.CheckRequiredLength(len(azs), i, context, "vsphere-network-name-%d"); err != nil {
				return nil, err
			}
			photonNetworkNames := context.StringSlice(fmt.Sprintf("vsphere-network-name-%d", i))
			for index, az := range azs {
				subnet := enaml.Subnet{
					AZ:       az,
					Range:    ranges[index],
					Gateway:  gateways[index],
					DNS:      strings.Split(dnsServers[index], ","),
					Reserved: strings.Split(reservedRanges[index], ","),
					Static:   strings.Split(staticIPs[index], ","),
					CloudProperties: NetworkCloudProperties{
						NetworkName: photonNetworkNames[index],
					},
				}
				network.AddSubnet(subnet)
			}
			networks = append(networks, network)
		}
	}
	return networks, nil
}

func clusterConfig(clusterName, resourcePoolName string) (clusters []map[string]ResourcePool) {
	cluster := make(map[string]ResourcePool)
	cluster[clusterName] = ResourcePool{
		ResourcePool: resourcePoolName,
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
	if len(azNames) != len(resourcePools) {
		err := fmt.Errorf("Sorry you need to provide the same number of az and vsphere-resource-pool flags")
		return nil, err
	}

	azs := []enaml.AZ{}

	for i, azName := range azNames {
		azCloudProperties := AZCloudProperties{}
		dataCenter := DataCenter{
			Name:     datacenters[i],
			Clusters: clusterConfig(clusters[i], resourcePools[i]),
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
	vmTypes := []enaml.VMType{
		enaml.VMType{
			Name: "small",
			CloudProperties: VMProperties{
				CPU:  1,
				RAM:  2048,
				Disk: 30000,
			},
		},
		enaml.VMType{
			Name: "medium",
			CloudProperties: VMProperties{
				CPU:  2,
				RAM:  4096,
				Disk: 50000,
			},
		},
		enaml.VMType{
			Name: "large.memory",
			CloudProperties: VMProperties{
				CPU:  4,
				RAM:  65536,
				Disk: 50000,
			},
		},
		enaml.VMType{
			Name: "large.cpu",
			CloudProperties: VMProperties{
				CPU:  4,
				RAM:  4096,
				Disk: 30000,
			},
		},
	}
	return vmTypes, nil
}

func (c *VSphereCloudConfig) CreateDiskTypes() ([]enaml.DiskType, error) {
	diskTypes := []enaml.DiskType{
		enaml.DiskType{
			Name:            "small",
			DiskSize:        3000,
			CloudProperties: VMProperties{},
		},
		enaml.DiskType{
			Name:            "medium",
			DiskSize:        30000,
			CloudProperties: VMProperties{},
		},
		enaml.DiskType{
			Name:            "large",
			DiskSize:        50000,
			CloudProperties: VMProperties{},
		},
	}
	return diskTypes, nil
}

func (c *VSphereCloudConfig) CreateCompilation() (*enaml.Compilation, error) {
	compilation := &enaml.Compilation{
		Workers:             8,
		ReuseCompilationVMs: true,
		AZ:                  c.Context.StringSlice("network-az-1")[0],
		VMType:              "medium",
		Network:             c.Context.String("network-name-1"),
	}
	return compilation, nil
}
