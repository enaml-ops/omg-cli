package main

import (
	"gopkg.in/urfave/cli.v2"
	"github.com/enaml-ops/pluginlib/cloudconfig"
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

func (s *MyCloudConfig) GetCloudConfig(args []string) []byte {
	return []byte("")
}
