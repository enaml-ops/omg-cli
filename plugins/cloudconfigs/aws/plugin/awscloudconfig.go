package awsccplugin

import (
	"strings"

	"github.com/codegangsta/cli"
	"github.com/enaml-ops/omg-cli/pluginlib/cloudconfig"
	"github.com/enaml-ops/omg-cli/pluginlib/util"
	aws "github.com/enaml-ops/omg-cli/plugins/cloudconfigs/aws/cloud-config"
	"github.com/xchapter7x/lo"
)

type AWSCloudConfig struct{}

func (s *AWSCloudConfig) GetFlags() (flags []cli.Flag) {
	return []cli.Flag{
		cli.StringSliceFlag{Name: "az-subnet-map", Usage: "comma separated list of az-subnet maps (ex: `us-east-1c:subnet-123456`)"},
		cli.StringFlag{Name: "region", Usage: "aws region"},
		cli.StringSliceFlag{Name: "security-group", Usage: "list of security groups"},
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
	c := pluginutil.NewContext(args, s.GetFlags())
	cloud := aws.NewAWSCloudConfig(
		c.String("region"),
		parseAZSubnetSlice(c.StringSlice("az-subnet-map")),
		c.StringSlice("security-group"),
	)
	if b, err = cloud.Bytes(); err != nil {
		lo.G.Error("cloud bytes call yielded error: ", err)
	}
	return b
}
