package plugin

import (
	"fmt"

	"github.com/codegangsta/cli"
	"github.com/enaml-ops/enaml"
	"github.com/enaml-ops/omg-cli/plugins/cloudconfigs"
)

type VMProperties struct {
	MachineType    string `yaml:"vm_flavor,omitempty"`
	RootDiskSizeGB int    `yaml:"vm_attached_disk_size_gb,omitempty"`
	RootDiskType   string `yaml:"disk_flavor"`
}

type NetworkCloudProperties struct {
	NetworkName string `yaml:"network_id"`
}

type PhotonCloudConfig struct {
	Context *cli.Context
}

func NewPhotonCloudConfig(c *cli.Context) cloudconfigs.CloudConfigProvider {
	return &PhotonCloudConfig{
		Context: c,
	}
}

func (c *PhotonCloudConfig) networkCloudProperties(i, index int) interface{} {
	photonNetworkNames := c.Context.StringSlice(fmt.Sprintf("photon-network-name-%d", i))
	return NetworkCloudProperties{
		NetworkName: photonNetworkNames[index],
	}
}

func (c *PhotonCloudConfig) validateCloudProperties(length, i int) error {
	return cloudconfigs.CheckRequiredLength(length, i, c.Context, "photon-network-name-%d")
}

func (c *PhotonCloudConfig) CreateNetworks() ([]enaml.DeploymentNetwork, error) {
	networks, err := cloudconfigs.CreateNetworks(c.Context, c.validateCloudProperties, c.networkCloudProperties)
	return networks, err
}

func (c *PhotonCloudConfig) CreateAZs() ([]enaml.AZ, error) {
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

func (c *PhotonCloudConfig) CreateVMTypes() ([]enaml.VMType, error) {
	vmTypes := []enaml.VMType{
		enaml.VMType{
			Name: "small",
			CloudProperties: VMProperties{
				MachineType:    "core-200",
				RootDiskSizeGB: 30,
				RootDiskType:   "core-200",
			},
		},
		enaml.VMType{
			Name: "medium",
			CloudProperties: VMProperties{
				MachineType:    "core-200",
				RootDiskSizeGB: 50,
				RootDiskType:   "core-200",
			},
		},
		enaml.VMType{
			Name: "large.memory",
			CloudProperties: VMProperties{
				MachineType:    "core-200",
				RootDiskSizeGB: 50,
				RootDiskType:   "core-200",
			},
		},
		enaml.VMType{
			Name: "large.cpu",
			CloudProperties: VMProperties{
				MachineType:    "core-200",
				RootDiskSizeGB: 30,
				RootDiskType:   "core-200",
			},
		},
	}
	return vmTypes, nil
}

func (c *PhotonCloudConfig) CreateDiskTypes() ([]enaml.DiskType, error) {
	diskTypes := []enaml.DiskType{
		enaml.DiskType{
			Name:     "small",
			DiskSize: 3000,
			CloudProperties: VMProperties{
				RootDiskType: "core-200",
			},
		},
		enaml.DiskType{
			Name:     "medium",
			DiskSize: 30000,
			CloudProperties: VMProperties{
				RootDiskType: "core-200",
			},
		},
		enaml.DiskType{
			Name:     "large",
			DiskSize: 50000,
			CloudProperties: VMProperties{
				RootDiskType: "core-200",
			},
		},
	}
	return diskTypes, nil
}

func (c *PhotonCloudConfig) CreateCompilation() (*enaml.Compilation, error) {
	return cloudconfigs.CreateCompilation(c.Context)
}
