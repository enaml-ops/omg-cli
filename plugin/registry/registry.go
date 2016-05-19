package registry

import (
	"log"
	"os/exec"

	"github.com/bosh-ops/bosh-install/plugin/cloudconfig"
	"github.com/hashicorp/go-plugin"
	"github.com/xchapter7x/lo"
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

func RegisterCloudConfig(pluginpath string) error {
	client := plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig: cloudconfig.HandshakeConfig,
		Plugins: map[string]plugin.Plugin{
			cloudconfig.PluginsMapHash: new(cloudconfig.CloudConfigPlugin),
		},
		Cmd: exec.Command(pluginpath, "plugin"),
	})
	defer client.Kill()
	rpcClient, err := client.Client()

	if err != nil {
		log.Fatal(err)
	}
	raw, err := rpcClient.Dispense(cloudconfig.PluginsMapHash)
	lo.G.Debug("we have a raw: ", raw)

	if err != nil {
		log.Fatal(err)
	}
	cloudconfigPlugin := raw.(cloudconfig.CloudConfigDeployer)
	meta := cloudconfigPlugin.GetMeta()
	cloudconfigs[meta.Name] = registryRecord{
		Name:       meta.Name,
		Path:       pluginpath,
		Properties: meta.Properties,
	}
	return nil
}
