package awsccplugin

import (
	"strconv"

	"gopkg.in/urfave/cli.v2"
	aws "github.com/enaml-ops/omg-cli/plugins/cloudconfigs/aws/cloud-config"
	"github.com/enaml-ops/pluginlib/cloudconfig"
	"github.com/enaml-ops/pluginlib/util"
	"github.com/xchapter7x/lo"
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
func (s *AWSCloudConfig) GetFlags() (flags []cli.Flag) {
	flags = []cli.Flag{
		&cli.StringFlag{Name: "aws-region", Usage: "aws region"},
		&cli.StringSliceFlag{Name: "aws-security-group", Usage: "list of security groups"},
	}
	for i := 1; i <= AZCountSupported; i++ {
		flags = append(flags, &cli.StringFlag{Name: CreateFlagnameWithSuffix("bosh-az-name", i), Usage: "name for bosh availablility zone in cloud config"})
		flags = append(flags, &cli.StringFlag{Name: CreateFlagnameWithSuffix("cidr", i), Usage: "cidr range for the given network"})
		flags = append(flags, &cli.StringFlag{Name: CreateFlagnameWithSuffix("gateway", i), Usage: "gateway for given network"})
		flags = append(flags, &cli.StringSliceFlag{Name: CreateFlagnameWithSuffix("dns", i), Usage: "dns for given network"})
		flags = append(flags, &cli.StringFlag{Name: CreateFlagnameWithSuffix("aws-az-name", i), Usage: "aws az name for given network"})
		flags = append(flags, &cli.StringFlag{Name: CreateFlagnameWithSuffix("aws-subnet-name", i), Usage: "aws subnet name for given network"})
		flags = append(flags, &cli.StringSliceFlag{Name: CreateFlagnameWithSuffix("bosh-reserve-range", i), Usage: "bosh reserve range for given network"})
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
	c := pluginutil.NewContext(args, s.GetFlags())
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
