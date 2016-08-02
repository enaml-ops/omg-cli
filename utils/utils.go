package utils

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"

	"github.com/codegangsta/cli"
	"github.com/enaml-ops/omg-cli/bosh"
	"github.com/enaml-ops/pluginlib/registry"
	"github.com/enaml-ops/pluginlib/util"
	"github.com/xchapter7x/lo"
)

// ClearDefaultStringSliceValue - this is simply to work around a defect in the
// cli package, where the default is appended to rather than replaced by user
// defined flags for StringSliceFlag values.
func ClearDefaultStringSliceValue(stringSliceArgs ...string) (res []string) {
	if isJustDefault(stringSliceArgs) {
		res = stringSliceArgs

	} else {
		res = stringSliceArgs[1:]
	}
	return
}

func isJustDefault(stringSliceArgs []string) bool {
	return len(stringSliceArgs) == 1
}

// GetCloudConfigCommands builds a list of CLI commands depending on
// which cloud config plugins are installed.
func GetCloudConfigCommands(target string) (commands []cli.Command) {
	files, _ := ioutil.ReadDir(target)
	for _, f := range files {
		lo.G.Debug("registering: ", f.Name())
		pluginPath := path.Join(target, f.Name())
		flags, _ := registry.RegisterCloudConfig(pluginPath)

		commands = append(commands, cli.Command{
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

				return bosh.CloudConfigAction(c.Parent(), cc)
			},
		})
	}
	lo.G.Debug("registered cloud configs: ", registry.ListCloudConfigs())
	return
}

// GetProductCommands builds a list of CLI commands depending on which
// product plugins are installed.
func GetProductCommands(target string) (commands []cli.Command) {
	files, _ := ioutil.ReadDir(target)
	for _, f := range files {
		lo.G.Debug("registering: ", f.Name())
		pluginPath := path.Join(target, f.Name())
		flags, _ := registry.RegisterProduct(pluginPath)

		commands = append(commands, cli.Command{
			Name:  f.Name(),
			Usage: "deploy the " + f.Name() + " product",
			Flags: pluginutil.ToCliFlagArray(flags),
			Action: func(c *cli.Context) (err error) {
				client, productDeployment := registry.GetProductReference(pluginPath)
				defer client.Kill()
				return bosh.ProductAction(c.Parent(), productDeployment)
			},
		})
	}
	lo.G.Debug("registered product plugins: ", registry.ListProducts())
	return
}

func ConvertToCLIStringSliceFlag(values []string) *cli.StringSlice {
	cliSlice := &cli.StringSlice{}
	for _, value := range values {
		cliSlice.Set(value)
	}
	return cliSlice
}

func CheckRequired(names []string, c *cli.Context) {
	var invalidNames []string
	for _, name := range names {
		if c.String(name) == "" {
			invalidNames = append(invalidNames, name)
		}
	}
	if len(invalidNames) > 0 {
		fmt.Println("Sorry you need to provide", invalidNames, "flags to continue")
		os.Exit(1)
	}
}

func DeployYaml(myYaml string, boshInitDeploy func(string)) {
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
