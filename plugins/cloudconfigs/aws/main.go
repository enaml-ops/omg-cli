package main

import (
	"github.com/enaml-ops/omg-cli/plugins/cloudconfigs/aws/plugin"
	"github.com/enaml-ops/pluginlib/cloudconfig"
)

var Version string = "v0.0.0"

func main() {
	cloudconfig.Run(&awsccplugin.AWSCloudConfig{
		PluginVersion: Version,
	})
}
