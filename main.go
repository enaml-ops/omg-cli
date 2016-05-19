package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/bosh-ops/bosh-install/aws-cli"
	"github.com/bosh-ops/bosh-install/azure-cli"
	"github.com/bosh-ops/bosh-install/plugin/registry"
	"github.com/codegangsta/cli"
	"github.com/xchapter7x/lo"
)

var Version string
var CloudConfigPluginsDir = "./.plugins/cloudconfig"
var ProductPluginsDir = "./.plugins/product"
var cloudConfigCommands []cli.Command
var productCommands []cli.Command
var productList []string
var cloudconfigList []string

func init() {
	files, _ := ioutil.ReadDir(CloudConfigPluginsDir)
	for _, f := range files {
		lo.G.Debug("registering: ", f.Name())
		registry.RegisterCloudConfig(path.Join(CloudConfigPluginsDir, f.Name()))
	}
	lo.G.Debug("registered cloud configs: ", registry.ListCloudConfigs())

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
				for _, plgn := range registry.ListCloudConfigs() {
					fmt.Println(plgn.Name, " - ", plgn.Path, " - ", plgn.Properties)
				}
				return nil
			},
		},
		{
			Name: "list-products",
			Action: func(c *cli.Context) error {
				fmt.Println("Products:")
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
