package main

import (
	"github.com/enaml-ops/omg-cli/plugins/cloudconfigs/aws/plugin"
	"github.com/enaml-ops/pluginlib/cloudconfig"
)

func main() {
	cloudconfig.Run(new(awsccplugin.AWSCloudConfig))
}
