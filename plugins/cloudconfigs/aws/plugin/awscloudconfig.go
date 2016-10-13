package awsccplugin

import (
	"strconv"

	aws "github.com/enaml-ops/omg-cli/plugins/cloudconfigs/aws/cloud-config"
	"github.com/enaml-ops/pluginlib/cloudconfig"
	"github.com/enaml-ops/pluginlib/pcli"
	"github.com/enaml-ops/pluginlib/pluginutil"
	"github.com/xchapter7x/lo"
	"gopkg.in/urfave/cli.v2"
)

//AZCountSupported - number of supported AZ's, setting to 10 to start with.
const AZCountSupported = 10

// AWSCloudConfig --
type AWSCloudConfig struct {
	PluginVersion string
}

//CreateFlagnameWithSuffix - creates a CLI flag with flagname-index pattern
func CreateFlagnameWithSuffix(name string, suffix int) (flagname string) {
	return name + "-" + strconv.Itoa(suffix)
}

//GetFlags - Get flags associated with plugin
func (s *AWSCloudConfig) GetFlags() (flags []pcli.Flag) {
	flags = []pcli.Flag{
		pcli.CreateStringFlag("aws-region", "aws region"),
		pcli.CreateStringSliceFlag("aws-security-group", "list of security groups"),
	}
	for i := 1; i <= AZCountSupported; i++ {
		flags = append(flags, pcli.CreateStringFlag(CreateFlagnameWithSuffix("bosh-az-name", i), "name for bosh availablility zone in cloud config"))
		flags = append(flags, pcli.CreateStringFlag(CreateFlagnameWithSuffix("cidr", i), "cidr range for the given network"))
		flags = append(flags, pcli.CreateStringFlag(CreateFlagnameWithSuffix("gateway", i), "gateway for given network"))
		flags = append(flags, pcli.CreateStringSliceFlag(CreateFlagnameWithSuffix("dns", i), "dns for given network"))
		flags = append(flags, pcli.CreateStringFlag(CreateFlagnameWithSuffix("aws-az-name", i), "aws az name for given network"))
		flags = append(flags, pcli.CreateStringFlag(CreateFlagnameWithSuffix("aws-subnet-name", i), "aws subnet name for given network"))
		flags = append(flags, pcli.CreateStringSliceFlag(CreateFlagnameWithSuffix("bosh-reserve-range", i), "bosh reserve range for given network"))
		flags = append(flags, pcli.CreateStringSliceFlag(CreateFlagnameWithSuffix("bosh-static-range", i), "bosh static range for given network"))
	}
	return
}

//GetMeta - Get metadata of the plugin
func (s *AWSCloudConfig) GetMeta() cloudconfig.Meta {
	return cloudconfig.Meta{
		Name: "aws",
		Properties: map[string]interface{}{
			"version": s.PluginVersion,
		},
	}
}

//GetCloudConfig - get a serialized form of AWS cloud configuration
func (s *AWSCloudConfig) GetCloudConfig(args []string) (b []byte) {
	var err error
	c := pluginutil.NewContext(args, pluginutil.ToCliFlagArray(s.GetFlags()))
	cloud := aws.NewAWSCloudConfig(
		c.String("aws-region"),
		c.StringSlice("aws-security-group"),
		getSubnetBucketList(c),
	)
	if b, err = cloud.Bytes(); err != nil {
		lo.G.Error("cloud bytes call yielded error: ", err)
	}
	return b
}

func getSubnetBucketList(c *cli.Context) (bucket []aws.SubnetBucket) {

	for i := 1; i <= AZCountSupported; i++ {
		tmpBucket := aws.SubnetBucket{
			BoshAZName:       c.String(CreateFlagnameWithSuffix("bosh-az-name", i)),
			Cidr:             c.String(CreateFlagnameWithSuffix("cidr", i)),
			Gateway:          c.String(CreateFlagnameWithSuffix("gateway", i)),
			DNS:              c.StringSlice(CreateFlagnameWithSuffix("dns", i)),
			AWSAZName:        c.String(CreateFlagnameWithSuffix("aws-az-name", i)),
			AWSSubnetName:    c.String(CreateFlagnameWithSuffix("aws-subnet-name", i)),
			BoshReserveRange: c.StringSlice(CreateFlagnameWithSuffix("bosh-reserve-range", i)),
			BoshStaticRange:  c.StringSlice(CreateFlagnameWithSuffix("bosh-static-range", i)),
		}
		if isValidSubnetBucket(tmpBucket) {
			bucket = append(bucket, tmpBucket)
		}
	}
	return
}

func isValidSubnetBucket(bucket aws.SubnetBucket) bool {
	return (bucket.BoshAZName != "" &&
		bucket.Cidr != "" &&
		bucket.Gateway != "" &&
		len(bucket.DNS) > 0 &&
		bucket.AWSAZName != "" &&
		bucket.AWSSubnetName != "" &&
		len(bucket.BoshReserveRange) > 0)
}
