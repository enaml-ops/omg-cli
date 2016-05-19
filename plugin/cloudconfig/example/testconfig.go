package main

import (
	"github.com/bosh-ops/bosh-install/plugin/cloudconfig"
	"github.com/codegangsta/cli"
	"github.com/xchapter7x/enaml"
)

func main() {
	cloudconfig.Run(new(MyCloudConfig))
}

type MyCloudConfig struct{}

func (s *MyCloudConfig) GetFlags() (flags []cli.Flag) {
	return
}

func (s *MyCloudConfig) GetMeta() cloudconfig.Meta {
	return cloudconfig.Meta{
		Name: "myfakecloudconfig",
	}
}

func (s *MyCloudConfig) GetCloudConfig(args []string) enaml.CloudConfigManifest {
	return enaml.CloudConfigManifest{}
}
