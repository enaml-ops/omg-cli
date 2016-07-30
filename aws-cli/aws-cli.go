package awscli

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/codegangsta/cli"
	"github.com/enaml-ops/enaml"
	"github.com/enaml-ops/omg-cli/plugins/products/bosh-init"
	"github.com/enaml-ops/omg-cli/utils"
	"github.com/xchapter7x/lo"
)

func deployYaml(myYaml string, boshInitDeploy func(string)) {
	fmt.Println("deploying your bosh")
	content := []byte(myYaml)
	tmpfile, err := ioutil.TempFile("", "bosh-init-deployment")
	defer os.Remove(tmpfile.Name())

	if err != nil {
		log.Fatal(err)
	}

	if _, err := tmpfile.Write(content); err != nil {
		log.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		log.Fatal(err)
	}
	boshInitDeploy(tmpfile.Name())
}

func checkRequired(name string, c *cli.Context) {
	if c.String(name) == "" {
		fmt.Println("Sorry you need to provide " + name)
		os.Exit(1)
	}
}

func GetFlags() []cli.Flag {
	boshdefaults := boshinit.GetAWSBoshBase()

	boshFlags := boshinit.BoshFlags(boshdefaults)
	awsFlags := []cli.Flag{
		cli.StringFlag{Name: "aws-instance-size", Value: "m3.xlarge", Usage: "the size of aws instance you wish to create"},
		cli.StringFlag{Name: "aws-availability-zone", Value: "us-east-1c", Usage: "the ec2 az you wish to deploy to"},
		cli.StringFlag{Name: "aws-subnet", Value: "", Usage: "your target vpc subnet"},
		cli.StringFlag{Name: "aws-pem-path", Value: "", Usage: "your aws pem file path"},
		cli.StringFlag{Name: "aws-access-key", Value: "", Usage: "aws account access key"},
		cli.StringFlag{Name: "aws-keyname", Value: "bosh", Usage: "aws keyname"},
		cli.StringFlag{Name: "aws-secret", Value: "", Usage: "aws account secret key"},
		cli.StringFlag{Name: "aws-region", Value: "us-east-1", Usage: "ec2 region to deploy on"},
		cli.StringSliceFlag{Name: "aws-security-group", Value: &cli.StringSlice{"bosh"}, Usage: "this is for security groups to apply to your VM. you can add as many security group flags as you like"},
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
			return
		}
		lo.G.Debug("Got boshbase", boshBase)
		checkRequired("aws-subnet", c)
		checkRequired("aws-pem-path", c)
		checkRequired("aws-access-key", c)
		checkRequired("aws-secret", c)
		checkRequired("aws-region", c)

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

		manifest := provider.CreateDeploymentManifest()

		lo.G.Debug("Got manifest", manifest)
		if yamlString, err := enaml.Paint(manifest); err == nil {

			if c.Bool("print-manifest") {
				fmt.Println(yamlString)

			} else {
				deployYaml(yamlString, boshInitDeploy)
			}
		} else {
			e = err
		}
		return
	}
}
