package awsccplugin

import (
	"strings"

	aws "github.com/bosh-ops/bosh-install/cloudconfigs/aws/cloud-config"
	"github.com/bosh-ops/bosh-install/plugin/cloudconfig"
	"github.com/codegangsta/cli"
	"github.com/xchapter7x/lo"
)

type AWSCloudConfig struct{}

func (s *AWSCloudConfig) GetFlags() (flags []cli.Flag) {
	return []cli.Flag{
		cli.StringSliceFlag{Name: "az-subnet-map", Usage: "comma separated list of az-subnet maps (ex: `us-east-1c:subnet-123456`)"},
		cli.StringFlag{Name: "region", Usage: "aws region"},
		cli.StringFlag{Name: "security-groups", Usage: "comma separated list of security groups"},
	}
}

func (s *AWSCloudConfig) GetMeta() cloudconfig.Meta {
	return cloudconfig.Meta{
		Name: "aws",
	}
}

func parseAZSubnetSlice(azSubnetSlice []string) (azSubnetMap map[string]string) {
	azSubnetMap = make(map[string]string)
	for _, v := range azSubnetSlice {
		ss := strings.SplitN(v, ":", 2)
		azSubnetMap[ss[0]] = ss[1]
	}
	return
}

func (s *AWSCloudConfig) GetCloudConfig(args []string) (b []byte) {
	var err error
	cloud := aws.NewAWSCloudConfig(
		"bosh",
		parseAZSubnetSlice([]string{"us-east-1c:subnet-12345us-east-1c:subnet-1234566"}),
		[]string{"bosh"},
	)
	if b, err = cloud.Bytes(); err != nil {
		lo.G.Error("cloud bytes call yielded error: ", err)
	}
	return b
}
