package utils

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"path"

	"github.com/codegangsta/cli"
	"github.com/enaml-ops/enaml"
	"github.com/enaml-ops/enaml/enamlbosh"
	"github.com/enaml-ops/omg-cli/pluginlib/registry"
	"github.com/xchapter7x/lo"
)

//ClearDefaultStringSliceValue - this is simply to work around a defect in the
//cli package, where the default is appended to rather than replaced by user
//defined flags for StringSliceFlag values.
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
				manifest := cc.GetCloudConfig(c.Parent().Args())
				lo.G.Debug("we found a manifest and context: ", manifest, c)
				processManifest(c, manifest)
				return nil
			},
		})
	}
	lo.G.Debug("registered cloud configs: ", registry.ListCloudConfigs())
	return
}

func processManifest(c *cli.Context, manifest []byte) (e error) {

	if c.Parent().Bool("print-manifest") {
		yamlString := string(manifest)
		fmt.Println(yamlString)

	} else {
		ccm := enaml.NewCloudConfigManifest(manifest)
		boshclient := enamlbosh.NewClient(c.Parent().String("bosh-user"), c.Parent().String("bosh-pass"), c.Parent().String("bosh-url"), c.Parent().Int("bosh-port"))
		if req, err := boshclient.NewCloudConfigRequest(*ccm); err == nil {
			tr := &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: c.Parent().Bool("ssl-ignore")},
			}
			httpClient := &http.Client{Transport: tr}

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
