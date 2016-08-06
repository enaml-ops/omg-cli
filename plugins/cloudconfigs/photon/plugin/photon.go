package plugin

import (
	"fmt"
	"strings"

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

	provider := &PhotonCloudConfig{
		Context: c,
	}
	return provider
}

func (c *PhotonCloudConfig) CreateNetworks() ([]enaml.DeploymentNetwork, error) {
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
			cloudconfigs.CheckRequiredLength(len(azs), i, context, "network-cidr-%d", "network-gateway-%d", "network-dns-%d", "network-reserved-%d", "network-static-%d")
			ranges := context.StringSlice(fmt.Sprintf("network-cidr-%d", i))
			gateways := context.StringSlice(fmt.Sprintf("network-gateway-%d", i))
			dnsServers := context.StringSlice(fmt.Sprintf("network-dns-%d", i))
			reservedRanges := context.StringSlice(fmt.Sprintf("network-reserved-%d", i))
			staticIPs := context.StringSlice(fmt.Sprintf("network-static-%d", i))
			cloudconfigs.CheckRequiredLength(len(azs), i, context, "photon-network-name-%d")
			photonNetworkNames := context.StringSlice(fmt.Sprintf("photon-network-name-%d", i))
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
	compilation := &enaml.Compilation{
		Workers:             8,
		ReuseCompilationVMs: true,
		AZ:                  c.Context.StringSlice("network-az-1")[0],
		VMType:              "medium",
		Network:             c.Context.String("network-name-1"),
	}
	return compilation, nil
}
