package awscli

import (
	"gopkg.in/urfave/cli.v2"
	"github.com/enaml-ops/omg-cli/plugins/products/bosh-init"
	"github.com/enaml-ops/omg-cli/utils"
	"github.com/xchapter7x/lo"
)

func GetFlags() []cli.Flag {
	boshdefaults := boshinit.GetAWSBoshBase()

	boshFlags := boshinit.BoshFlags(boshdefaults)
	awsFlags := []cli.Flag{
		&cli.StringFlag{Name: "aws-instance-size", Value: "m3.xlarge", Usage: "the size of aws instance you wish to create"},
		&cli.StringFlag{Name: "aws-availability-zone", Value: "us-east-1c", Usage: "the ec2 az you wish to deploy to"},
		&cli.StringFlag{Name: "aws-subnet", Value: "", Usage: "your target vpc subnet"},
		&cli.StringFlag{Name: "aws-pem-path", Value: "", Usage: "your aws pem file path"},
		&cli.StringFlag{Name: "aws-access-key", Value: "", Usage: "aws account access key"},
		&cli.StringFlag{Name: "aws-keyname", Value: "bosh", Usage: "aws keyname"},
		&cli.StringFlag{Name: "aws-secret", Value: "", Usage: "aws account secret key"},
		&cli.StringFlag{Name: "aws-region", Value: "us-east-1", Usage: "ec2 region to deploy on"},
		&cli.StringSliceFlag{Name: "aws-security-group", Value: cli.NewStringSlice("bosh"), Usage: "this is for security groups to apply to your VM. you can add as many security group flags as you like"},
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
			AWSSecurityGroups:   utils.ClearDefaultStringSliceValue(c.StringSlice("aws-security-group")...),
		}, boshBase)

		if err := boshBase.HandleDeployment(provider, boshInitDeploy); err != nil {
			return err
		}

		return nil
	}
}
