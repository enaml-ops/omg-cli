package plugin

import (
	"fmt"

	"gopkg.in/urfave/cli.v2"
	"github.com/enaml-ops/omg-cli/plugins/cloudconfigs"
	"github.com/enaml-ops/pluginlib/cloudconfig"
	"github.com/enaml-ops/pluginlib/util"
)

type Plugin struct {
	PluginVersion string
}

func (s *Plugin) GetFlags() []cli.Flag {
	flags := cloudconfigs.CreateAZFlags()
	flags = cloudconfigs.CreateNetworkFlags(flags, networkFlags)
	flags = append(flags, &cli.StringSliceFlag{Name: "vsphere-datacenter", Usage: "vsphere datacenter name"})
	flags = append(flags, &cli.StringSliceFlag{Name: "vsphere-cluster", Usage: "vsphere cluster name"})
	flags = append(flags, &cli.StringSliceFlag{Name: "vsphere-resource-pool", Usage: "vsphere resource pool name"})

	return flags
}

func networkFlags(flags []cli.Flag, i int) []cli.Flag {
	flags = append(flags, &cli.StringSliceFlag{Name: cloudconfigs.CreateFlagnameWithSuffix("vsphere-network-name", i), Usage: fmt.Sprintf("vsphere network name for network %d", i)})
	return flags
}

//GetMeta - Get metadata of the plugin
func (s *Plugin) GetMeta() cloudconfig.Meta {
	return cloudconfig.Meta{
		Name: "vsphere",
		Properties: map[string]interface{}{
			"version": s.PluginVersion,
		},
	}
}

//GetCloudConfig - get a serialized form of vCenter cloud configuration
func (s *Plugin) GetCloudConfig(args []string) (b []byte) {
	var err error
	c := pluginutil.NewContext(args, s.GetFlags())
	cloudConfig := NewVSphereCloudConfig(c)
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
