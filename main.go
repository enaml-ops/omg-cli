package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"sort"
	"strings"

	"github.com/enaml-ops/omg-cli/aws-cli"
	"github.com/enaml-ops/omg-cli/azure-cli"
	"github.com/enaml-ops/omg-cli/bosh"
	gcpcli "github.com/enaml-ops/omg-cli/gcp-cli"
	"github.com/enaml-ops/omg-cli/photon-cli"
	"github.com/enaml-ops/omg-cli/utils"
	"github.com/enaml-ops/omg-cli/vsphere-cli"
	"github.com/enaml-ops/pluginlib/pcli"
	"github.com/enaml-ops/pluginlib/registry"
	"github.com/enaml-ops/pluginlib/pluginutil"
	"github.com/olekukonko/tablewriter"
	"github.com/pivotalservices/gtils/osutils"
	"github.com/xchapter7x/lo"
	"gopkg.in/urfave/cli.v2"
)

// Version is the version of omg-cli.
var Version string

// CloudConfigPluginsDir is the directory where registered cloud config plugins are stored.
var CloudConfigPluginsDir = "./.plugins/cloudconfig"

// ProductPluginsDir is the directory where registered product plugins are stored.
var ProductPluginsDir = "./.plugins/product"

func main() {
	app := (&cli.App{})
	app.Version = Version
	app.Commands = []*cli.Command{
		{
			Name:   "azure",
			Usage:  "azure [--flags] - deploy a bosh to azure",
			Action: azurecli.GetAction(BoshInitDeploy),
			Flags:  pluginutil.ToCliFlagArray(azurecli.GetFlags()),
		},
		{
			Name:   "aws",
			Usage:  "aws [--flags] - deploy a bosh to aws",
			Action: awscli.GetAction(BoshInitDeploy),
			Flags:  pluginutil.ToCliFlagArray(awscli.GetFlags()),
		},
		{
			Name:   "gcp",
			Usage:  "gcp [--flags] - deploy a bosh to GCP",
			Action: gcpcli.GetAction(BoshInitDeploy),
			Flags:  pluginutil.ToCliFlagArray(gcpcli.GetFlags()),
		},
		{
			Name:   "photon",
			Usage:  "photon [--flags] - deploy a bosh to photon",
			Action: photoncli.GetAction(BoshInitDeploy),
			Flags:  pluginutil.ToCliFlagArray(photoncli.GetFlags()),
		},
		{
			Name:   "vsphere",
			Usage:  "vsphere [--flags] - deploy a bosh to vsphere",
			Action: vspherecli.GetAction(BoshInitDeploy),
			Flags:  pluginutil.ToCliFlagArray(vspherecli.GetFlags()),
		},
		{
			Name: "list-cloudconfigs",
			Action: func(c *cli.Context) error {
				fmt.Println("Cloud Configs:")
				table := tablewriter.NewWriter(os.Stdin)
				table.SetHeader([]string{"Name", "Command", "Properties"})
				var data [][]string
				formatProperties := func(p map[string]interface{}) string {
					var res string
					for n, v := range p {
						res += fmt.Sprintf("%v: %v\n", n, v)
					}
					return res
				}

				for _, plgn := range registry.ListCloudConfigs() {
					row := []string{
						plgn.Name,
						path.Base(plgn.Path),
						formatProperties(plgn.Properties),
					}
					data = append(data, row)
				}
				table.AppendBulk(data)
				table.Render()
				return nil
			},
		},
		{
			Name: "list-products",
			Action: func(c *cli.Context) error {
				ListProducts(os.Stdout, registry.ListProducts())
				return nil
			},
		},
		{
			Name:  "product-meta",
			Usage: "product-meta <prod-name> - show product metadata",
			Action: func(c *cli.Context) error {
				if c.Args().Len() > 0 {
					return productMeta(c.Args().First())
				}
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
			Flags: pluginutil.ToCliFlagArray([]pcli.Flag{
				pcli.CreateStringFlag("type", "define if the plugin to be registered is a cloudconfig or a product", "product"),
				pcli.CreateStringFlag("pluginpath", "the path to the plugin you wish to register"),
			}),
		},
		{
			Name:        "deploy-cloudconfig",
			Usage:       "deploy-cloudconfig <cloudconfig-name> [--flags] - deploy a cloudconfig to bosh",
			Flags:       pluginutil.ToCliFlagArray(bosh.GetAuthFlags()),
			Subcommands: utils.GetCloudConfigCommands(CloudConfigPluginsDir),
		},
		{
			Name:        "deploy-product",
			Usage:       "deploy-product <prod-name> [--flags] - deploy a product via bosh",
			Flags:       pluginutil.ToCliFlagArray(bosh.GetAuthFlags()),
			Subcommands: utils.GetProductCommands(ProductPluginsDir),
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		fmt.Println("ERROR:", err)
		os.Exit(1)
	}
}

func init() {
	if strings.ToLower(os.Getenv("LOG_LEVEL")) != "debug" {
		log.SetOutput(ioutil.Discard)
	}
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

func formatProps(props map[string]interface{}) string {
	var keys []string
	for k := range props {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	buf := &bytes.Buffer{}

	// output version first, then the remainder of the properties (sorted)
	ver, ok := props["version"]
	if ok {
		fmt.Fprintln(buf, "version:", ver)
	}

	for _, key := range keys {
		// don't output version twice
		if ok && key == "version" {
			continue
		}
		fmt.Fprintf(buf, "%v: %v\n", key, props[key])
	}
	return buf.String()
}

// ListProducts writes a formatted list of products to w.
func ListProducts(w io.Writer, products map[string]registry.Record) {
	table := tablewriter.NewWriter(w)
	table.SetHeader([]string{"Name", "Command", "Properties"})

	var keys []string
	for k := range products {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, key := range keys {
		row := []string{
			products[key].Name,
			path.Base(products[key].Path),
			formatProps(products[key].Properties),
		}
		table.Append(row)
	}
	table.Render()
}

func productMeta(product string) error {
	record, ok := registry.ListProducts()[product]
	if !ok {
		return fmt.Errorf("product '%s' not found", product)
	}
	fmt.Println(formatProps(record.Properties))
	return nil
}
