package main

import (
	"github.com/enaml-ops/omg-cli/pluginlib/cloudconfig"
	"github.com/enaml-ops/omg-cli/plugins/cloudconfigs/aws/plugin"
)

func main() {
	cloudconfig.Run(new(awsccplugin.AWSCloudConfig))
}
