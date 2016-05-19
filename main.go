package main

import (
	"fmt"
	"os"

	"github.com/bosh-ops/bosh-install/aws-cli"
	"github.com/bosh-ops/bosh-install/azure-cli"
	"github.com/codegangsta/cli"
)

var Version string
var cloudConfigCommands []cli.Command
var productCommands []cli.Command
var productList []string
var cloudconfigList []string

func init() {
	cloudConfigCommands = append(cloudConfigCommands, cli.Command{
		Name:  "test",
		Usage: "add a new template",
		Action: func(c *cli.Context) error {
			fmt.Println("no cloud config plugins supported yet: ", c.Args().First())
			return nil
		},
	})

	productCommands = append(productCommands, cli.Command{
		Name:  "test",
		Usage: "add a new template",
		Action: func(c *cli.Context) error {
			fmt.Println("no product plugins supported yet: ", c.Args().First())
			return nil
		},
	})
}

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
		{
			Name: "list-cloudconfigs",
			Action: func(c *cli.Context) error {
				fmt.Println("Cloud Configs:")
				for _, v := range cloudconfigList {
					fmt.Println(v)
				}
				return nil
			},
		},
		{
			Name: "list-products",
			Action: func(c *cli.Context) error {
				fmt.Println("Products:")
				for _, v := range productList {
					fmt.Println(v)
				}
				return nil
			},
		},
		{
			Name:        "deploy-cloudconfig",
			Usage:       "deploy-cloudconfig <cloudconfig-name> [--flags] - deploy a cloudconfig to bosh",
			Subcommands: cloudConfigCommands,
		},
		{
			Name:        "deploy-product",
			Usage:       "deploy-product <prod-name> [--flags] - deploy a product via bosh",
			Subcommands: productCommands,
		},
	}
	app.Run(os.Args)
}
