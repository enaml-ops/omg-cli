package plugin

import (
	"fmt"

	"github.com/enaml-ops/omg-cli/plugins/cloudconfigs"
	"github.com/enaml-ops/pluginlib/cloudconfigv1"
	"github.com/enaml-ops/pluginlib/pcli"
	"github.com/enaml-ops/pluginlib/pluginutil"
	"gopkg.in/urfave/cli.v2"
)

type Plugin struct {
	PluginVersion string
}

func (s *Plugin) GetFlags() []pcli.Flag {
	flags := cloudconfigs.CreateAZFlags()
	flags = cloudconfigs.CreateNetworkFlags(flags, networkFlags)
	flags = append(flags, pcli.CreateStringSliceFlag("vsphere-datacenter", "vsphere datacenter name"))
	flags = append(flags, pcli.CreateStringSliceFlag("vsphere-cluster", "vsphere cluster name"))
	flags = append(flags, pcli.CreateStringSliceFlag("vsphere-resource-pool", "vsphere resource pool name"))
	return flags
}

func networkFlags(flags []pcli.Flag, i int) []pcli.Flag {
	flags = append(flags, pcli.CreateStringSliceFlag(cloudconfigs.CreateFlagnameWithSuffix("vsphere-network-name", i), fmt.Sprintf("vsphere network name for network %d", i)))
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
func (s *Plugin) GetCloudConfig(args []string) ([]byte, error) {
	c := pluginutil.NewContext(args, pluginutil.ToCliFlagArray(s.GetFlags()))
	cloudConfig := NewVSphereCloudConfig(c)
	b, err := cloudconfigs.GetDeploymentManifestBytes(cloudConfig)
	if err != nil {
		return nil, err
	}
	return b, nil
}

//GetContext -
func (s *Plugin) GetContext(args []string) (c *cli.Context) {
	c = pluginutil.NewContext(args, pluginutil.ToCliFlagArray(s.GetFlags()))
	return
}
