package main

import (
	"fmt"
	"os"

	"github.com/bosh-ops/bosh-install/aws-cli"
	"github.com/codegangsta/cli"
)

var Version string

func main() {
	app := cli.NewApp()
	app.Version = Version
	app.Commands = []cli.Command{
		{
			Name:  "azure",
			Usage: "azure [--flags] - deploy a bosh to azure",
			Action: func(c *cli.Context) {
				fmt.Println("not yet implemented")
				os.Exit(1)
			},
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
