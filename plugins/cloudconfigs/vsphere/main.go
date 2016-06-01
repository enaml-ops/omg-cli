package main

import (
	"github.com/enaml-ops/omg-cli/pluginlib/cloudconfig"
	"github.com/enaml-ops/omg-cli/plugins/cloudconfigs/vsphere/plugin"
)

func main() {
	cloudconfig.Run(new(vsphereccplugin.VSphereCloudConfig))
}
