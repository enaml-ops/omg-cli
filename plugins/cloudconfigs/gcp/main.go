package main

import (
	"github.com/enaml-ops/omg-cli/plugins/cloudconfigs/gcp/plugin"
	"github.com/enaml-ops/pluginlib/cloudconfigv1"
)

var Version string = "v0.0.0"

func main() {
	cloudconfig.Run(&plugin.Plugin{
		PluginVersion: Version,
	})
}
