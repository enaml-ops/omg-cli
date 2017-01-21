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
	flags = append(flags, pcli.CreateStringSliceFlag("gcp-availability-zone", "gcp availability_zone name for az"))
	return flags
}

func networkFlags(flags []pcli.Flag, i int) []pcli.Flag {
	flags = append(flags,
		pcli.CreateStringSliceFlag(cloudconfigs.CreateFlagnameWithSuffix("gcp-network-name", i), fmt.Sprintf("gcp network name for network %d", i)),
		pcli.CreateStringSliceFlag(cloudconfigs.CreateFlagnameWithSuffix("gcp-subnetwork-name", i), fmt.Sprintf("gcp subnetwork name for network %d", i)),
		pcli.CreateStringSliceFlag(cloudconfigs.CreateFlagnameWithSuffix("gcp-network-tag", i), fmt.Sprintf("gcp network tags for network %d", i)))
	return flags
}

//GetMeta - Get metadata of the plugin
func (s *Plugin) GetMeta() cloudconfig.Meta {
	return cloudconfig.Meta{
		Name: "gcp",
		Properties: map[string]interface{}{
			"version": s.PluginVersion,
		},
	}
}

//GetCloudConfig - get a serialized form of GCP cloud configuration
func (s *Plugin) GetCloudConfig(args []string) ([]byte, error) {
	c := pluginutil.NewContext(args, pluginutil.ToCliFlagArray(s.GetFlags()))
	cloudConfig := NewGCPCloudConfig(c)
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
