package awsccplugin

import (
	aws "github.com/bosh-ops/bosh-install/cloudconfigs/aws/cloud-config"
	"github.com/bosh-ops/bosh-install/plugin/cloudconfig"
	"github.com/codegangsta/cli"
	"github.com/xchapter7x/lo"
)

type AWSCloudConfig struct{}

func (s *AWSCloudConfig) GetFlags() (flags []cli.Flag) {
	return []cli.Flag{}
}

func (s *AWSCloudConfig) GetMeta() cloudconfig.Meta {
	return cloudconfig.Meta{
		Name: "aws",
	}
}

func (s *AWSCloudConfig) GetCloudConfig(args []string) (b []byte) {
	var err error
	cloud := aws.NewAWSCloudConfig("bosh", []string{"bosh", "bosh"}, []string{"bosh", "bosh"}, []string{"bosh", "bosh"})
	if b, err = cloud.Bytes(); err != nil {
		lo.G.Error("cloud bytes call yielded error: ", err)
	}
	return b
}
