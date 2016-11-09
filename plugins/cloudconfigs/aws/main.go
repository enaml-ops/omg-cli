package main

import (
	"github.com/enaml-ops/omg-cli/plugins/cloudconfigs/aws/plugin"
	v1 "github.com/enaml-ops/pluginlib/cloudconfigv1"
)

var Version string = "v0.0.0"

func main() {
	v1.Run(&plugin.Plugin{
		PluginVersion: Version,
	})
}
