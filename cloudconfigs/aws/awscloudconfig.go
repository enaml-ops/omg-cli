package main

import (
	"github.com/bosh-ops/bosh-install/plugin/cloudconfig"
	"github.com/codegangsta/cli"
	"github.com/xchapter7x/enaml"
)

func main() {
	cloudconfig.Run(new(AWSCloudConfig))
}

type AWSCloudConfig struct{}

func (s *AWSCloudConfig) GetFlags() (flags []cli.Flag) {
	return
}

func (s *AWSCloudConfig) GetMeta() cloudconfig.Meta {
	return cloudconfig.Meta{
		Name: "aws",
	}
}

func (s *AWSCloudConfig) GetCloudConfig(args []string) enaml.CloudConfigManifest {
	return enaml.CloudConfigManifest{}
}
