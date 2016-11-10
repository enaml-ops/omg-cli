package awscli

import (
	"github.com/enaml-ops/omg-cli/plugins/products/bosh-init"
	"github.com/enaml-ops/omg-cli/utils"
	"github.com/enaml-ops/pluginlib/pcli"
	"github.com/xchapter7x/lo"
	"gopkg.in/urfave/cli.v2"
)

func GetFlags() []pcli.Flag {
	boshdefaults := boshinit.GetAWSBoshBase()

	boshFlags := boshinit.BoshFlags(boshdefaults)
	awsFlags := []pcli.Flag{
		pcli.CreateStringFlag("aws-instance-size", "the size of aws instance you wish to create", "m3.xlarge"),
		pcli.CreateStringFlag("aws-availability-zone", "the ec2 az you wish to deploy to", "us-east-1a"),
		pcli.CreateStringFlag("aws-subnet", "your target vpc subnet"),
		pcli.CreateStringFlag("aws-pem-path", "your aws pem file path"),
		pcli.CreateStringFlag("aws-access-key", "aws account access key"),
		pcli.CreateStringFlag("aws-keyname", "aws keyname", "bosh"),
		pcli.CreateStringFlag("aws-secret", "aws account secret key"),
		pcli.CreateStringFlag("aws-region", "ec2 region to deploy on", "us-east-1"),
		pcli.CreateStringSliceFlag("aws-security-group", "this is for security groups to apply to your VM. you can add as many security group flags as you like", "bosh"),
	}
	for _, flag := range awsFlags {
		boshFlags = append(boshFlags, flag)
	}
	return boshFlags
}

func GetAction(boshInitDeploy func(string)) func(c *cli.Context) error {
	return func(c *cli.Context) (e error) {
		var boshBase *boshinit.BoshBase
		if boshBase, e = boshinit.NewBoshBase(c); e != nil {
			lo.G.Error(e.Error())
			return e
		}
		lo.G.Debug("Got boshbase", boshBase)
		if err := utils.CheckRequired(c, "aws-subnet", "aws-pem-path", "aws-access-key", "aws-secret", "aws-region"); err != nil {
			lo.G.Error(err.Error())
			return err
		}

		provider := boshinit.NewAWSIaaSProvider(boshinit.AWSInitConfig{
			AWSInstanceSize:     c.String("aws-instance-size"),
			AWSAvailabilityZone: c.String("aws-availability-zone"),
			AWSSubnet:           c.String("aws-subnet"),
			AWSPEMFilePath:      c.String("aws-pem-path"),
			AWSAccessKeyID:      c.String("aws-access-key"),
			AWSSecretKey:        c.String("aws-secret"),
			AWSRegion:           c.String("aws-region"),
			AWSKeyName:          c.String("aws-keyname"),
			AWSSecurityGroups:   c.StringSlice("aws-security-group"),
		}, boshBase)

		if err := boshBase.HandleDeployment(provider, boshInitDeploy); err != nil {
			return err
		}
		return nil
	}
}
