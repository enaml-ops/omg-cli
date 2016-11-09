package plugin

import (
	"fmt"

	"github.com/enaml-ops/omg-cli/plugins/cloudconfigs"
	"github.com/enaml-ops/pluginlib/cloudconfigv1"
	"github.com/enaml-ops/pluginlib/pcli"
	"github.com/enaml-ops/pluginlib/pluginutil"
	"github.com/xchapter7x/lo"
	"gopkg.in/urfave/cli.v2"
)

//Plugin : version
type Plugin struct {
	PluginVersion string
}

func (s *Plugin) GetFlags() []pcli.Flag {
	flags := cloudconfigs.CreateAZFlags()
	flags = cloudconfigs.CreateNetworkFlags(flags, networkFlags)
	return flags
}

func networkFlags(flags []pcli.Flag, i int) []pcli.Flag {
	flags = append(flags, pcli.CreateStringSliceFlag(cloudconfigs.CreateFlagnameWithSuffix("azure-virtual-network-name", i), fmt.Sprintf("Azure virtual network name : %d", i)))
	flags = append(flags, pcli.CreateStringSliceFlag(cloudconfigs.CreateFlagnameWithSuffix("azure-subnet-name", i), fmt.Sprintf("Azure network subnet name : %d", i)))
	return flags
}

//GetMeta - Get metadata of the plugin
func (s *Plugin) GetMeta() cloudconfig.Meta {
	return cloudconfig.Meta{
		Name: "azure",
		Properties: map[string]interface{}{
			"version": s.PluginVersion,
		},
	}
}

//GetCloudConfig - get a serialized form of azure cloud configuration
func (s *Plugin) GetCloudConfig(args []string) ([]byte, error) {
	c := pluginutil.NewContext(args, pluginutil.ToCliFlagArray(s.GetFlags()))
	cloudConfig := NewAzureCloudConfig(c)
	b, err := cloudconfigs.GetDeploymentManifestBytes(cloudConfig)
	if err != nil {
		lo.G.Debug("error getting Azure cloud config:", err)
	}
	return b, err
}

//GetContext -
func (s *Plugin) GetContext(args []string) (c *cli.Context) {
	c = pluginutil.NewContext(args, pluginutil.ToCliFlagArray(s.GetFlags()))
	return
}
