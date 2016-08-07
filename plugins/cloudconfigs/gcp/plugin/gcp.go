package plugin

import (
	"fmt"
	"strings"

	"github.com/codegangsta/cli"
	"github.com/enaml-ops/enaml"
	"github.com/enaml-ops/omg-cli/plugins/cloudconfigs"
)

type VMProperties struct {
	MachineType    string `yaml:"machine_type,omitempty"`
	RootDiskSizeGB int    `yaml:"root_disk_size_gb"`
	RootDiskType   string `yaml:"root_disk_type"`
}
type NetworkCloudProperties struct {
	NetworkName    string   `yaml:"network_name"`
	SubNetworkName string   `yaml:"subnetwork_name"`
	Tags           []string `yaml:"tags"`
}
type GCPCloudConfig struct {
	Context *cli.Context
}

func NewGCPCloudConfig(c *cli.Context) cloudconfigs.CloudConfigProvider {

	provider := &GCPCloudConfig{
		Context: c,
	}
	return provider
}

func (c *GCPCloudConfig) networkCloudProperties(i, index int) interface{} {
	gcpNetworkNames := c.Context.StringSlice(fmt.Sprintf("gcp-network-name-%d", i))
	gcpSubNetworkNames := c.Context.StringSlice(fmt.Sprintf("gcp-subnetwork-name-%d", i))
	gcpNetworkTags := c.Context.StringSlice(fmt.Sprintf("gcp-network-tag-%d", i))
	return NetworkCloudProperties{
		NetworkName:    gcpNetworkNames[index],
		SubNetworkName: gcpSubNetworkNames[index],
		Tags:           strings.Split(gcpNetworkTags[index], ","),
	}
}

func (c *GCPCloudConfig) validateCloudProperties(length, i int) error {
	return cloudconfigs.CheckRequiredLength(length, i, c.Context, "gcp-network-name-%d", "gcp-subnetwork-name-%d", "gcp-network-tag-%d")
}

func (c *GCPCloudConfig) CreateNetworks() ([]enaml.DeploymentNetwork, error) {
	networks, err := cloudconfigs.CreateNetworks(c.Context, c.validateCloudProperties, c.networkCloudProperties)
	return networks, err
}
func (c *GCPCloudConfig) CreateAZs() ([]enaml.AZ, error) {
	azNames := c.Context.StringSlice("az")
	gcpAZNames := c.Context.StringSlice("gcp-availability-zone")

	if len(azNames) != len(gcpAZNames) {
		err := fmt.Errorf("Sorry you need to provide the same number of az and gcp-availability-zone flags")
		return nil, err
	}
	azs := []enaml.AZ{}

	for i, azName := range azNames {
		az := enaml.AZ{
			Name: azName,
			CloudProperties: map[string]string{
				"availability_zone": gcpAZNames[i],
			},
		}
		azs = append(azs, az)
	}
	return azs, nil
}

func (c *GCPCloudConfig) CreateVMTypes() ([]enaml.VMType, error) {
	vmTypes := []enaml.VMType{
		enaml.VMType{
			Name: "small",
			CloudProperties: VMProperties{
				MachineType:    "n1-standard-1",
				RootDiskSizeGB: 30,
				RootDiskType:   "pd-standard",
			},
		},
		enaml.VMType{
			Name: "medium",
			CloudProperties: VMProperties{
				MachineType:    "n1-standard-2",
				RootDiskSizeGB: 50,
				RootDiskType:   "pd-standard",
			},
		},
		enaml.VMType{
			Name: "large.memory",
			CloudProperties: VMProperties{
				MachineType:    "n1-highmem-4",
				RootDiskSizeGB: 50,
				RootDiskType:   "pd-standard",
			},
		},
		enaml.VMType{
			Name: "large.cpu",
			CloudProperties: VMProperties{
				MachineType:    "n1-highcpu-4",
				RootDiskSizeGB: 30,
				RootDiskType:   "pd-standard",
			},
		},
	}
	return vmTypes, nil
}

func (c *GCPCloudConfig) CreateDiskTypes() ([]enaml.DiskType, error) {
	diskTypes := []enaml.DiskType{
		enaml.DiskType{
			Name:     "small",
			DiskSize: 3000,
			CloudProperties: VMProperties{
				RootDiskSizeGB: 3,
				RootDiskType:   "pd-standard",
			},
		},
		enaml.DiskType{
			Name:     "medium",
			DiskSize: 30000,
			CloudProperties: VMProperties{
				RootDiskSizeGB: 50,
				RootDiskType:   "pd-standard",
			},
		},
		enaml.DiskType{
			Name:     "large",
			DiskSize: 50000,
			CloudProperties: VMProperties{
				RootDiskSizeGB: 50,
				RootDiskType:   "pd-standard",
			},
		},
	}
	return diskTypes, nil
}

func (c *GCPCloudConfig) CreateCompilation() (*enaml.Compilation, error) {
	return cloudconfigs.CreateCompilation(c.Context)
}
