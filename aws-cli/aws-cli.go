package awscli

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/bosh-ops/bosh-install/deployments/bosh-init-aws"
	"github.com/codegangsta/cli"
	"github.com/xchapter7x/enaml"
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
	return []cli.Flag{
		cli.StringFlag{Name: "name", Value: "bosh", Usage: "the vm name to be created in your ec2 account"},
		cli.StringFlag{Name: "bosh-release-ver", Value: "256.2", Usage: "the version of the bosh release you wish to use (found on bosh.io)"},
		cli.StringFlag{Name: "bosh-private-ip", Value: "10.0.0.6", Usage: "the private ip for the bosh vm to be created in ec2"},
		cli.StringFlag{Name: "bosh-cpi-release-ver", Value: "52", Usage: "the bosh cpi version you wish to use (found on bosh.io)"},
		cli.StringFlag{Name: "go-agent-ver", Value: "3012", Usage: "the go agent version you wish to use (found on bosh.io)"},
		cli.StringFlag{Name: "bosh-release-sha", Value: "ff2f4e16e02f66b31c595196052a809100cfd5a8", Usage: "sha1 of the bosh release being used (found on bosh.io)"},
		cli.StringFlag{Name: "bosh-cpi-release-sha", Value: "dc4a0cca3b33dce291e4fbeb9e9948b6a7be3324", Usage: "sha1 of the cpi release being used (found on bosh.io)"},
		cli.StringFlag{Name: "go-agent-sha", Value: "3380b55948abe4c437dee97f67d2d8df4eec3fc1", Usage: "sha1 of the go agent being use (found on bosh.io)"},
		cli.StringFlag{Name: "aws-instance-size", Value: "m3.xlarge", Usage: "the size of aws instance you wish to create"},
		cli.StringFlag{Name: "aws-availability-zone", Value: "us-east-1c", Usage: "the ec2 az you wish to deploy to"},
		cli.StringFlag{Name: "director-name", Value: "my-bosh", Usage: "the name of your director"},
		cli.StringFlag{Name: "aws-subnet", Value: "", Usage: "your target vpc subnet"},
		cli.StringFlag{Name: "aws-elastic-ip", Value: "", Usage: "your elastic ip to assign to the bosh vm"},
		cli.StringFlag{Name: "aws-pem-path", Value: "", Usage: "your aws pem file path"},
		cli.StringFlag{Name: "aws-access-key", Value: "", Usage: "aws account access key"},
		cli.StringFlag{Name: "aws-secret", Value: "", Usage: "aws account secret key"},
		cli.StringFlag{Name: "aws-region", Value: "us-east-1", Usage: "ec2 region to deploy on"},
		cli.BoolFlag{Name: "print-manifest", Usage: "if you would simply like to output a manifest the set this flag as true."},
	}
}

func GetAction(boshInitDeploy func(string)) func(c *cli.Context) error {
	return func(c *cli.Context) (e error) {
		checkRequired("aws-subnet", c)
		checkRequired("aws-elastic-ip", c)
		checkRequired("aws-pem-path", c)
		checkRequired("aws-access-key", c)
		checkRequired("aws-secret", c)
		checkRequired("aws-region", c)

		manifest := boshinitaws.NewBoshInit(boshinitaws.BoshInitConfig{
			Name:                  c.String("name"),
			BoshReleaseVersion:    c.String("bosh-release-ver"),
			BoshPrivateIP:         c.String("bosh-private-ip"),
			BoshCPIReleaseVersion: c.String("bosh-cpi-release-ver"),
			GoAgentVersion:        c.String("go-agent-ver"),
			BoshReleaseSHA:        c.String("bosh-release-sha"),
			BoshCPIReleaseSHA:     c.String("bosh-cpi-release-sha"),
			GoAgentSHA:            c.String("go-agent-sha"),
			BoshInstanceSize:      c.String("aws-instance-size"),
			BoshAvailabilityZone:  c.String("aws-availability-zone"),
			BoshAWSSubnet:         c.String("aws-subnet"),
			AWSElasticIP:          c.String("aws-elastic-ip"),
			BoshDirectorName:      c.String("director-name"),
			AWSPEMFilePath:        c.String("aws-pem-path"),
			AWSAccessKeyID:        c.String("aws-access-key"),
			AWSSecretKey:          c.String("aws-secret"),
			AWSRegion:             c.String("aws-region"),
		})

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
