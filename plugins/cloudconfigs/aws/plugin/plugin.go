package plugin

import (
	"fmt"
	"strconv"

	"gopkg.in/urfave/cli.v2"

	"github.com/enaml-ops/omg-cli/plugins/cloudconfigs"
	aws "github.com/enaml-ops/omg-cli/plugins/cloudconfigs/aws/cloud-config"
	v1 "github.com/enaml-ops/pluginlib/cloudconfigv1"
	"github.com/enaml-ops/pluginlib/pcli"
	"github.com/enaml-ops/pluginlib/pluginutil"
	"github.com/xchapter7x/lo"
)

// Plugin --
type Plugin struct {
	PluginVersion string
}

//CreateFlagnameWithSuffix - creates a CLI flag with flagname-index pattern
func CreateFlagnameWithSuffix(name string, suffix int) (flagname string) {
	return name + "-" + strconv.Itoa(suffix)
}

//GetFlags - Get flags associated with plugin
func (s *Plugin) GetFlags() []pcli.Flag {

	flags := cloudconfigs.CreateAZFlags()
	flags = cloudconfigs.CreateNetworkFlags(flags, networkFlags)
	flags = append(flags, pcli.CreateStringSliceFlag("aws-availablity-zone", "aws availablility zone"))
	return flags
}

func networkFlags(flags []pcli.Flag, i int) []pcli.Flag {
	flags = append(flags, pcli.CreateStringSliceFlag(cloudconfigs.CreateFlagnameWithSuffix("aws-subnet-name", i), fmt.Sprintf("aws subnet name %d", i)))
	flags = append(flags, pcli.CreateStringSliceFlag(cloudconfigs.CreateFlagnameWithSuffix("aws-security-group", i), fmt.Sprintf("list of security groups %d", i)))
	return flags
}

//GetMeta - Get metadata of the plugin
func (s *Plugin) GetMeta() v1.Meta {
	return v1.Meta{
		Name: "aws",
		Properties: map[string]interface{}{
			"version": s.PluginVersion,
		},
	}
}

//GetCloudConfig - get a serialized form of AWS cloud configuration
func (s *Plugin) GetCloudConfig(args []string) ([]byte, error) {
	c := pluginutil.NewContext(args, pluginutil.ToCliFlagArray(s.GetFlags()))
	cloudConfig := aws.NewAWSCloudConfig(c)
	b, err := cloudconfigs.GetDeploymentManifestBytes(cloudConfig)
	if err != nil {
		lo.G.Debug("cloud bytes call yielded error: ", err)
	}
	return b, err
}

//GetContext -
func (s *Plugin) GetContext(args []string) (c *cli.Context) {
	c = pluginutil.NewContext(args, pluginutil.ToCliFlagArray(s.GetFlags()))
	return
}
