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

func (s *MyCloudConfig) GetAction() func(c *cli.Context) error {
	return func(*cli.Context) error {
		return nil
	}
}

func (s *MyCloudConfig) GetMeta() cloudconfig.Meta {
	return cloudconfig.Meta{
		Name: "myfakecloudconfig",
	}
}

func (s *MyCloudConfig) GetCloudConfig() enaml.CloudConfigManifest {
	return enaml.CloudConfigManifest{}
}
