package main

import (
	"github.com/enaml-ops/pluginlib/cloudconfigv1"
	"github.com/enaml-ops/pluginlib/pcli"
)

func main() {
	cloudconfig.Run(new(MyCloudConfig))
}

type MyCloudConfig struct{}

func (s *MyCloudConfig) GetFlags() (flags []pcli.Flag) {
	return
}

func (s *MyCloudConfig) GetMeta() cloudconfig.Meta {
	return cloudconfig.Meta{
		Name: "myfakecloudconfig",
	}
}

func (s *MyCloudConfig) GetCloudConfig(args []string) ([]byte, error) {
	return []byte(""), nil
}
