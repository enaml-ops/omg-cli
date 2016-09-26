package plugin

import (
	"fmt"

	"github.com/enaml-ops/omg-cli/plugins/cloudconfigs"
	"github.com/enaml-ops/pluginlib/cloudconfig"
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
	return flags
}

func networkFlags(flags []pcli.Flag, i int) []pcli.Flag {
	flags = append(flags, pcli.CreateStringSliceFlag(cloudconfigs.CreateFlagnameWithSuffix("photon-network-name", i), fmt.Sprintf("photon network name for network %d", i)))
	return flags
}

//GetMeta - Get metadata of the plugin
func (s *Plugin) GetMeta() cloudconfig.Meta {
	return cloudconfig.Meta{
		Name: "photon",
		Properties: map[string]interface{}{
			"version": s.PluginVersion,
		},
	}
}

//GetCloudConfig - get a serialized form of AWS cloud configuration
func (s *Plugin) GetCloudConfig(args []string) (b []byte) {
	var err error
	c := pluginutil.NewContext(args, pluginutil.ToCliFlagArray(s.GetFlags()))
	cloudConfig := NewPhotonCloudConfig(c)
	if b, err = cloudconfigs.GetDeploymentManifestBytes(cloudConfig); err != nil {
		panic(err)
	}
	return b
}

//GetContext -
func (s *Plugin) GetContext(args []string) (c *cli.Context) {
	c = pluginutil.NewContext(args, pluginutil.ToCliFlagArray(s.GetFlags()))
	return
}
