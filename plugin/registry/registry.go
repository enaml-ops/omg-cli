package registry

import (
	"log"
	"os/exec"

	"github.com/enaml-ops/omg-cli/plugin/cloudconfig"
	"github.com/codegangsta/cli"
	"github.com/hashicorp/go-plugin"
)

var (
	cloudconfigs map[string]registryRecord
	products     map[string]registryRecord
)

func init() {
	cloudconfigs = make(map[string]registryRecord)
	products = make(map[string]registryRecord)
}

type registryRecord struct {
	Name       string
	Path       string
	Properties map[string]interface{}
}

func RegisterProduct(pluginpath string) error {
	return nil
}

func ListCloudConfigs() map[string]registryRecord {
	return cloudconfigs
}

func RegisterCloudConfig(pluginpath string) ([]cli.Flag, error) {
	client, ccPlugin := GetCloudConfigReference(pluginpath)
	defer client.Kill()
	meta := ccPlugin.GetMeta()
	cloudconfigs[meta.Name] = registryRecord{
		Name:       meta.Name,
		Path:       pluginpath,
		Properties: meta.Properties,
	}
	return ccPlugin.GetFlags(), nil
}

func GetCloudConfigReference(pluginpath string) (*plugin.Client, cloudconfig.CloudConfigDeployer) {
	client := plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig: cloudconfig.HandshakeConfig,
		Plugins: map[string]plugin.Plugin{
			cloudconfig.PluginsMapHash: new(cloudconfig.CloudConfigPlugin),
		},
		Cmd: exec.Command(pluginpath, "plugin"),
	})

	rpcClient, err := client.Client()

	if err != nil {
		log.Fatal(err)
	}
	raw, err := rpcClient.Dispense(cloudconfig.PluginsMapHash)

	if err != nil {
		log.Fatal(err)
	}
	return client, raw.(cloudconfig.CloudConfigDeployer)
}
