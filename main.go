package main

import (
	"os"

	"github.com/bosh-ops/bosh-install/aws-cli"
	"github.com/bosh-ops/bosh-install/azure-cli"
	"github.com/codegangsta/cli"
)

var Version string

func main() {
	app := cli.NewApp()
	app.Version = Version
	app.Commands = []cli.Command{
		{
			Name:   "azure",
			Usage:  "azure [--flags] - deploy a bosh to azure",
			Action: azurecli.GetAction(BoshInitDeploy),
			Flags:  azurecli.GetFlags(),
		},
		{
			Name:   "aws",
			Usage:  "aws [--flags] - deploy a bosh to aws",
			Action: awscli.GetAction(BoshInitDeploy),
			Flags:  awscli.GetFlags(),
		},
	}
	app.Run(os.Args)
}
