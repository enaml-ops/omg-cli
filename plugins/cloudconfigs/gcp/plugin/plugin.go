package plugin

import (
	"fmt"

	"github.com/codegangsta/cli"
	"github.com/enaml-ops/omg-cli/plugins/cloudconfigs"
	"github.com/enaml-ops/pluginlib/cloudconfig"
	"github.com/enaml-ops/pluginlib/util"
)

type Plugin struct{}

func (s *Plugin) GetFlags() []cli.Flag {
	flags := cloudconfigs.CreateAZFlags()
	flags = cloudconfigs.CreateNetworkFlags(flags, networkFlags)
	flags = append(flags, cli.StringSliceFlag{Name: "gcp-availability-zone", Usage: "gcp availability_zone name for az"})
	return flags
}

func networkFlags(flags []cli.Flag, i int) []cli.Flag {
	flags = append(flags, cli.StringSliceFlag{Name: cloudconfigs.CreateFlagnameWithSuffix("gcp-network-name", i), Usage: fmt.Sprintf("gcp network name for network %d", i)})
	flags = append(flags, cli.StringSliceFlag{Name: cloudconfigs.CreateFlagnameWithSuffix("gcp-subnetwork-name", i), Usage: fmt.Sprintf("gcp subnetwork name for network %d", i)})
	flags = append(flags, cli.StringSliceFlag{Name: cloudconfigs.CreateFlagnameWithSuffix("gcp-network-tag", i), Usage: fmt.Sprintf("comma delimited list of gcp network tags for network %d", i)})
	return flags
}

//GetMeta - Get metadata of the plugin
func (s *Plugin) GetMeta() cloudconfig.Meta {
	return cloudconfig.Meta{
		Name: "gcp",
	}
}

//GetCloudConfig - get a serialized form of AWS cloud configuration
func (s *Plugin) GetCloudConfig(args []string) (b []byte) {
	var err error
	c := pluginutil.NewContext(args, s.GetFlags())
	cloudConfig := NewGCPCloudConfig(c)
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
