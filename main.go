package main

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/codegangsta/cli"
	"github.com/enaml-ops/enaml"
	"github.com/enaml-ops/omg-cli/aws-cli"
	"github.com/enaml-ops/omg-cli/azure-cli"
	"github.com/enaml-ops/omg-cli/plugin/registry"
	"github.com/enaml-ops/omg-cli/utils/boshapi"
	"github.com/pivotalservices/gtils/osutils"
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
			Name:  "register-plugin",
			Usage: "register-plugin -type [cloudconfig, product] -pluginpath <plugin-binary>",
			Action: func(c *cli.Context) (err error) {

				if c.String("type") != "" && c.String("pluginpath") != "" {
					err = registerPlugin(c.String("type"), c.String("pluginpath"))
				}
				return
			},
			Flags: []cli.Flag{
				cli.StringFlag{Name: "type", Value: "product", Usage: "define if the plugin to be registered is a cloud-config or a product"},
				cli.StringFlag{Name: "pluginpath", Value: "", Usage: "the path to the plugin you wish to register"},
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

func registerPlugin(typename, pluginpath string) (err error) {
	var srcPlugin *os.File

	if srcPlugin, err = os.Open(pluginpath); err == nil {
		defer srcPlugin.Close()

		switch typename {
		case "cloudconfig":
			dstFilepath := path.Join(CloudConfigPluginsDir, path.Base(pluginpath))
			err = copyPlugin(srcPlugin, dstFilepath)

		case "product":
			dstFilepath := path.Join(ProductPluginsDir, path.Base(pluginpath))
			err = copyPlugin(srcPlugin, dstFilepath)

		default:
			err = errors.New("invalid type selected")
			lo.G.Error("error: ", err)
		}
	}
	return
}

func copyPlugin(src io.Reader, dst string) (err error) {
	var dstPlugin *os.File
	if dstPlugin, err = osutils.SafeCreate(dst); err == nil {
		defer dstPlugin.Close()
		_, err = io.Copy(dstPlugin, src)
		os.Chmod(dst, 755)
	}
	return
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
				lo.G.Debug("running the cloud config plugin")
				client, cc := registry.GetCloudConfigReference(pluginPath)
				defer client.Kill()
				lo.G.Debug("we found client and cloud config: ", client, cc)
				lo.G.Debug("meta", cc.GetMeta())
				lo.G.Debug("args: ", c.Parent().Args())
				manifest := cc.GetCloudConfig(c.Parent().Args())
				lo.G.Debug("we found a manifest and context: ", manifest, c)
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
		cli.IntFlag{Name: "bosh-port", Value: 25555, Usage: "this is the port of your bosh director"},
		cli.StringFlag{Name: "bosh-user", Value: "bosh", Usage: "this is the username for your bosh director"},
		cli.StringFlag{Name: "bosh-pass", Value: "", Usage: "this is the pasword for your bosh director"},
		cli.BoolFlag{Name: "print-manifest", Usage: "if you would simply like to output a manifest the set this flag as true."},
	}
}

func processManifest(c *cli.Context, manifest []byte) (e error) {

	if c.Parent().Bool("print-manifest") {
		yamlString := string(manifest)
		fmt.Println(yamlString)

	} else {
		ccm := enaml.NewCloudConfigManifest(manifest)
		boshclient := boshapi.NewClient(c.Parent().String("bosh-user"), c.Parent().String("bosh-pass"), c.Parent().String("bosh-url"), c.Parent().Int("bosh-port"))
		if req, err := boshclient.NewCloudConfigRequest(*ccm); err == nil {
			httpClient := new(http.Client)

			if res, err := httpClient.Do(req); err != nil {
				lo.G.Error("res: ", res)
				lo.G.Error("error: ", err)
				e = err
			}
		} else {
			e = err
		}
	}
	return
}
