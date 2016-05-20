package main

import (
	"github.com/bosh-ops/bosh-install/cloudconfigs/aws/plugin"
	"github.com/bosh-ops/bosh-install/plugin/cloudconfig"
)

func main() {
	cloudconfig.Run(new(awsccplugin.AWSCloudConfig))
}
