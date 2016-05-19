package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"

	"github.com/bosh-ops/bosh-install/aws-cli"
	"github.com/bosh-ops/bosh-install/azure-cli"
	"github.com/bosh-ops/bosh-install/plugin/registry"
	"github.com/codegangsta/cli"
	"github.com/xchapter7x/enaml"
	"github.com/xchapter7x/lo"
)

var Version string
var CloudConfigPluginsDir = "./.plugins/cloudconfig"
var ProductPluginsDir = "./.plugins/product"
var cloudConfigCommands []cli.Command
var productCommands []cli.Command
var productList []string
var cloudconfigList []string

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
			Flags:       getBoshAuthFlags(),
			Subcommands: cloudConfigCommands,
		},
		{
			Name:        "deploy-product",
			Usage:       "deploy-product <prod-name> [--flags] - deploy a product via bosh",
			Flags:       getBoshAuthFlags(),
			Subcommands: productCommands,
		},
	}
	app.Run(os.Args)
}

func init() {

	if strings.ToLower(os.Getenv("LOG_LEVEL")) != "debug" {
		log.SetOutput(ioutil.Discard)
	}
	registerCloudConfig()
}

func registerCloudConfig() {
	files, _ := ioutil.ReadDir(CloudConfigPluginsDir)
	for _, f := range files {
		lo.G.Debug("registering: ", f.Name())
		pluginPath := path.Join(CloudConfigPluginsDir, f.Name())
		flags, _ := registry.RegisterCloudConfig(pluginPath)

		cloudConfigCommands = append(cloudConfigCommands, cli.Command{
			Name:  f.Name(),
			Usage: "deploy the " + f.Name() + " cloud config",
			Flags: flags,
			Action: func(c *cli.Context) error {
				client, cc := registry.GetCloudConfigReference(pluginPath)
				defer client.Kill()
				manifest := cc.GetCloudConfig(c)
				processManifest(c, manifest)
				return nil
			},
		})
	}
	lo.G.Debug("registered cloud configs: ", registry.ListCloudConfigs())
}

func getBoshAuthFlags() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{Name: "bosh-url", Value: "https://mybosh.com", Usage: "this is the url or ip of your bosh director"},
		cli.StringFlag{Name: "bosh-user", Value: "bosh", Usage: "this is the username for your bosh director"},
		cli.StringFlag{Name: "bosh-pass", Value: "", Usage: "this is the pasword for your bosh director"},
		cli.BoolFlag{Name: "print-manifest", Usage: "if you would simply like to output a manifest the set this flag as true."},
	}
}

func processManifest(c *cli.Context, manifest enaml.CloudConfigManifest) (e error) {
	if yamlString, err := enaml.Cloud(&manifest); err == nil {

		if c.Bool("print-manifest") {
			fmt.Println(yamlString)

		} else {
			fmt.Println("TODO: do something with my manifest here", manifest)
		}
	} else {
		e = err
	}
	return
}
