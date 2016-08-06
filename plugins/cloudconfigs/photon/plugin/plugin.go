package plugin

import (
	"fmt"

	"github.com/codegangsta/cli"
	"github.com/enaml-ops/omg-cli/plugins/cloudconfigs"
	"github.com/enaml-ops/pluginlib/cloudconfig"
	"github.com/enaml-ops/pluginlib/util"
)

const SupportedNetworkCount = 10
const SupportedAZCount = 3

type Plugin struct{}

func (s *Plugin) GetFlags() []cli.Flag {
	flags := []cli.Flag{
		cli.StringSliceFlag{Name: "az", Usage: "az name"},
	}
	for i := 1; i <= SupportedNetworkCount; i++ {
		flags = append(flags, cli.StringFlag{Name: cloudconfigs.CreateFlagnameWithSuffix("network-name", i), Usage: "network name"})
		flags = append(flags, cli.StringSliceFlag{Name: cloudconfigs.CreateFlagnameWithSuffix("network-az", i), Usage: fmt.Sprintf("az of network %d", i)})
		flags = append(flags, cli.StringSliceFlag{Name: cloudconfigs.CreateFlagnameWithSuffix("network-cidr", i), Usage: fmt.Sprintf("range of network %d", i)})
		flags = append(flags, cli.StringSliceFlag{Name: cloudconfigs.CreateFlagnameWithSuffix("network-gateway", i), Usage: fmt.Sprintf("gateway of network %d", i)})
		flags = append(flags, cli.StringSliceFlag{Name: cloudconfigs.CreateFlagnameWithSuffix("network-dns", i), Usage: fmt.Sprintf("comma delimited list of DNS servers for network %d", i)})
		flags = append(flags, cli.StringSliceFlag{Name: cloudconfigs.CreateFlagnameWithSuffix("network-reserved", i), Usage: fmt.Sprintf("comma delimited list of reserved network ranges for network %d", i)})
		flags = append(flags, cli.StringSliceFlag{Name: cloudconfigs.CreateFlagnameWithSuffix("network-static", i), Usage: fmt.Sprintf("comma delimited list of static IP addresses for network %d", i)})
		//Photon specifics
		flags = append(flags, cli.StringSliceFlag{Name: cloudconfigs.CreateFlagnameWithSuffix("photon-network-name", i), Usage: fmt.Sprintf("photon network name for network %d", i)})
	}
	return flags
}

//GetMeta - Get metadata of the plugin
func (s *Plugin) GetMeta() cloudconfig.Meta {
	return cloudconfig.Meta{
		Name: "photon",
	}
}

//GetCloudConfig - get a serialized form of AWS cloud configuration
func (s *Plugin) GetCloudConfig(args []string) (b []byte) {
	var err error
	c := pluginutil.NewContext(args, s.GetFlags())
	cloudConfig := NewPhotonCloudConfig(c)
	if b, err = cloudconfigs.GetDeploymentManifestBytes(cloudConfig); err != nil {
		panic(err)
	}
	return b
}

//GetContext -
func (s *Plugin) GetContext(args []string) (c *cli.Context) {
	c = pluginutil.NewContext(args, s.GetFlags())
	return
}
