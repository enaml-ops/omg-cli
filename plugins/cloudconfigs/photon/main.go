package main

import (
	"github.com/enaml-ops/omg-cli/plugins/cloudconfigs/photon/plugin"
	"github.com/enaml-ops/pluginlib/cloudconfig"
)

var Version string = "v0.0.0"

func main() {
	cloudconfig.Run(&plugin.Plugin{
		PluginVersion: Version,
	})
}
