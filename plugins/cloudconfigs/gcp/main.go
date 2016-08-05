package main

import (
	"github.com/enaml-ops/omg-cli/plugins/cloudconfigs/gcp/plugin"
	"github.com/enaml-ops/pluginlib/cloudconfig"
)

func main() {
	cloudconfig.Run(new(gcp.Plugin))
}
