package cloudconfigs

import (
	"fmt"
	"strconv"

	"gopkg.in/urfave/cli.v2"
	"gopkg.in/yaml.v2"

	"github.com/enaml-ops/enaml"
)

const SupportedNetworkCount = 10

type BaseCloudConfig struct {
	Manifest *enaml.CloudConfigManifest
}

func GetDeploymentManifestBytes(provider CloudConfigProvider) ([]byte, error) {
	var manifest *enaml.CloudConfigManifest
	var err error
	var cloudConfigYml []byte
	if manifest, err = CreateCloudConfigManifest(provider); err != nil {
		return nil, err
	}
	if cloudConfigYml, err = yaml.Marshal(manifest); err != nil {
		return nil, err
	}
	return cloudConfigYml, nil
}

func CreateCloudConfigManifest(provider CloudConfigProvider) (*enaml.CloudConfigManifest, error) {
	var err error
	var azs []enaml.AZ
	var networks []enaml.DeploymentNetwork
	var vmTypes []enaml.VMType
	var diskTypes []enaml.DiskType
	var compilation *enaml.Compilation

	base := BaseCloudConfig{
		Manifest: &enaml.CloudConfigManifest{},
	}

	if azs, err = provider.CreateAZs(); err != nil {
		return nil, err
	}
	base.Manifest.AZs = azs

	if networks, err = provider.CreateNetworks(); err != nil {
		return nil, err
	}
	base.Manifest.Networks = networks

	if vmTypes, err = provider.CreateVMTypes(); err != nil {
		return nil, err
	}
	base.Manifest.VMTypes = vmTypes

	if diskTypes, err = provider.CreateDiskTypes(); err != nil {
		return nil, err
	}
	base.Manifest.DiskTypes = diskTypes

	if compilation, err = provider.CreateCompilation(); err != nil {
		return nil, err
	}
	base.Manifest.Compilation = compilation

	return base.Manifest, nil
}

func CreateFlagnameWithSuffix(name string, suffix int) (flagname string) {
	return name + "-" + strconv.Itoa(suffix)
}

func CheckRequiredLength(targetLength, index int, c *cli.Context, names ...string) error {
	var invalidNames []string
	for _, name := range names {
		formattedName := fmt.Sprintf(name, index)
		if len(c.StringSlice(formattedName)) != targetLength {
			invalidNames = append(invalidNames, formattedName)
		}
	}
	if len(invalidNames) > 0 {
		err := fmt.Errorf("Sorry you need to provide %s flags with %d element(s) to continue", invalidNames, targetLength)
		return err
	}
	return nil
}

func CreateNetworks(context *cli.Context, validateCloudPropertiesFunction func(int, int) error, cloudPropertiesFunction func(int, int) interface{}) ([]enaml.DeploymentNetwork, error) {

	networks := []enaml.DeploymentNetwork{}
	for i := 1; i <= SupportedNetworkCount; i++ {
		networkFlag := fmt.Sprintf("network-name-%d", i)
		if context.IsSet(networkFlag) {
			network := enaml.ManualNetwork{
				Name: context.String(networkFlag),
				Type: "manual",
			}
			azs := context.StringSlice(fmt.Sprintf("network-az-%d", i))
			if err := CheckRequiredLength(len(azs), i, context, "network-cidr-%d", "network-gateway-%d"); err != nil {
				return nil, err
			}
			ranges := context.StringSlice(fmt.Sprintf("network-cidr-%d", i))
			gateways := context.StringSlice(fmt.Sprintf("network-gateway-%d", i))
			dnsServers := context.StringSlice(fmt.Sprintf("network-dns-%d", i))
			reservedRanges := context.StringSlice(fmt.Sprintf("network-reserved-%d", i))
			staticIPs := context.StringSlice(fmt.Sprintf("network-static-%d", i))
			if err := validateCloudPropertiesFunction(len(azs), i); err != nil {
				return nil, err
			}
			multiAssignAZ := context.Bool("multi-assign-az")
			if multiAssignAZ {
				subnet := enaml.Subnet{
					AZs:             azs,
					Range:           ranges[0],
					Gateway:         gateways[0],
					DNS:             dnsServers,
					Reserved:        reservedRanges,
					Static:          staticIPs,
					CloudProperties: cloudPropertiesFunction(i, 0),
				}
				network.AddSubnet(subnet)
			} else {
				for index, az := range azs {
					subnet := enaml.Subnet{
						AZ:              az,
						Range:           ranges[index],
						Gateway:         gateways[index],
						DNS:             dnsServers,
						Reserved:        reservedRanges,
						Static:          staticIPs,
						CloudProperties: cloudPropertiesFunction(i, index),
					}
					network.AddSubnet(subnet)
				}
			}
			networks = append(networks, network)
		}
	}
	return networks, nil
}

func CreateCompilation(c *cli.Context) (*enaml.Compilation, error) {
	compilation := &enaml.Compilation{
		Workers:             8,
		ReuseCompilationVMs: true,
		AZ:                  c.StringSlice("network-az-1")[0],
		VMType:              "medium",
		Network:             c.String("network-name-1"),
	}
	return compilation, nil
}
