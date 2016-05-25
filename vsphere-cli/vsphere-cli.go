package vspherecli

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/codegangsta/cli"
	"github.com/enaml-ops/enaml"
	"github.com/enaml-ops/omg-cli/deployments/bosh-init"
	"github.com/enaml-ops/omg-cli/utils"
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

// GetFlags returns the available CLI flags
func GetFlags() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{Name: "bosh-release-ver", Value: "256.2", Usage: "the version of the bosh release you wish to use (found on bosh.io)"},
		cli.StringFlag{Name: "bosh-cpi-release-ver", Value: "52", Usage: "the bosh cpi version you wish to use (found on bosh.io)"},
		cli.StringFlag{Name: "go-agent-ver", Value: "3012", Usage: "the go agent version you wish to use (found on bosh.io)"},
		cli.StringFlag{Name: "bosh-release-sha", Value: "ff2f4e16e02f66b31c595196052a809100cfd5a8", Usage: "sha1 of the bosh release being used (found on bosh.io)"},
		cli.StringFlag{Name: "bosh-cpi-release-sha", Value: "dc4a0cca3b33dce291e4fbeb9e9948b6a7be3324", Usage: "sha1 of the cpi release being used (found on bosh.io)"},
		cli.StringFlag{Name: "go-agent-sha", Value: "3380b55948abe4c437dee97f67d2d8df4eec3fc1", Usage: "sha1 of the go agent being use (found on bosh.io)"},
		cli.StringFlag{Name: "director-name", Value: "my-bosh", Usage: "the name of your director"},
		cli.BoolFlag{Name: "print-manifest", Usage: "if you would simply like to output a manifest the set this flag as true."},
	}
}

// GetAction returns a function action that can be registered with the CLI
func GetAction(boshInitDeploy func(string)) func(c *cli.Context) error {
	return func(c *cli.Context) (e error) {
		checkRequired("aws-subnet", c)
		checkRequired("aws-elastic-ip", c)
		checkRequired("aws-pem-path", c)
		checkRequired("aws-access-key", c)
		checkRequired("aws-secret", c)
		checkRequired("aws-region", c)

		manifest := boshinit.NewAWSBosh(boshinit.BoshInitConfig{
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
			AWSSubnet:             c.String("aws-subnet"),
			AWSElasticIP:          c.String("aws-elastic-ip"),
			BoshDirectorName:      c.String("director-name"),
			AWSPEMFilePath:        c.String("aws-pem-path"),
			AWSAccessKeyID:        c.String("aws-access-key"),
			AWSSecretKey:          c.String("aws-secret"),
			AWSRegion:             c.String("aws-region"),
			AWSSecurityGroups:     utils.ClearDefaultStringSliceValue(c.StringSlice("aws-security-group")...),
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
