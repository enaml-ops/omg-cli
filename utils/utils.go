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

var UIPrint = fmt.Println

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
				processCloudConfig(c, manifest)
				return nil
			},
		})
	}
	lo.G.Debug("registered cloud configs: ", registry.ListCloudConfigs())
	return
}

func GetProductCommands(target string) (commands []cli.Command) {
	files, _ := ioutil.ReadDir(target)
	for _, f := range files {
		lo.G.Debug("registering: ", f.Name())
		pluginPath := path.Join(target, f.Name())
		flags, _ := registry.RegisterProduct(pluginPath)

		commands = append(commands, cli.Command{
			Name:  f.Name(),
			Usage: "deploy the " + f.Name() + " product",
			Flags: flags,
			Action: func(c *cli.Context) (err error) {
				var cloudConfig *enaml.CloudConfigManifest
				client, productDeployment := registry.GetProductReference(pluginPath)
				defer client.Kill()
				boshclient := enamlbosh.NewClient(c.Parent().String("bosh-user"), c.Parent().String("bosh-pass"), c.Parent().String("bosh-url"), c.Parent().Int("bosh-port"))
				httpClient := defaultHTTPClient(c.Parent().Bool("ssl-ignore"))

				if cloudConfig, err = boshclient.GetCloudConfig(httpClient); err == nil {
					var cloudConfigBytes []byte
					var task []enamlbosh.BoshTask
					cloudConfigBytes, err = cloudConfig.Bytes()
					deploymentManifest := productDeployment.GetProduct(c.Parent().Args(), cloudConfigBytes)
					task, err = processProductDeployment(c, deploymentManifest)
					lo.G.Debug("bosh task: ", task)
				}
				return
			},
		})
	}
	lo.G.Debug("registered product plugins: ", registry.ListProducts())
	return
}

//ProcessProductBytes - upload a product deployments bytes to bosh
func ProcessProductBytes(manifest []byte, printManifest bool, user, pass, url string, port int, httpClient HttpClientDoer) (task []enamlbosh.BoshTask, err error) {
	if printManifest {
		yamlString := string(manifest)
		UIPrint(yamlString)

	} else {
		dm := enaml.NewDeploymentManifest(manifest)
		boshclient := enamlbosh.NewClient(user, pass, url, port)

		if err = processRemoteBoshAssets(dm, boshclient, httpClient); err == nil {
			UIPrint("Uploading product deployment...")

			if task, err = boshclient.PostDeployment(*dm, httpClient); err == nil {
				UIPrint("upload complete.")
				lo.G.Debug("res: ", task, err)

			} else {
				lo.G.Error("error: ", err)
			}

		} else {
			lo.G.Error("error: ", err)
		}
	}
	return
}

func processRemoteBoshAssets(dm *enaml.DeploymentManifest, boshClient *enamlbosh.Client, httpClient HttpClientDoer) (err error) {
	defer UIPrint("remote asset check complete.")
	UIPrint("Checking product deployment for remote assets...")

	if err = ProcessRemoteStemcells(dm.Stemcells, boshClient, httpClient); err == nil {
		err = ProcessRemoteReleases(dm.Releases, boshClient, httpClient)
	}
	return
}

//ProcessRemoteStemcells - upload any remote stemcells given
func ProcessRemoteStemcells(scl []enaml.Stemcell, boshClient *enamlbosh.Client, httpClient HttpClientDoer) (err error) {
	defer UIPrint("remote stemcells complete")
	UIPrint("Checking for remote stemcells...")
	var isRemoteStemcell = func(stemcell enaml.Stemcell) bool {
		return stemcell.URL != "" && stemcell.SHA1 != ""
	}

	for _, stemcell := range scl {
		var task []enamlbosh.BoshTask

		if isRemoteStemcell(stemcell) {
			task, err = boshClient.PostRemoteStemcell(stemcell, httpClient)
			lo.G.Debug("task: ", task, err)
		}
	}
	return
}

//ProcessRemoteReleases - upload any remote Releases given
func ProcessRemoteReleases(rl []enaml.Release, boshClient *enamlbosh.Client, httpClient HttpClientDoer) (err error) {
	defer UIPrint("remote releases complete")
	UIPrint("Checking for remote releases...")
	var isRemoteRelease = func(release enaml.Release) bool {
		return release.URL != "" && release.SHA1 != ""
	}

	for _, release := range rl {
		var task []enamlbosh.BoshTask

		if isRemoteRelease(release) {
			task, err = boshClient.PostRemoteRelease(release, httpClient)
			lo.G.Debug("task: ", task, err)
		}
	}
	return
}

func processProductDeployment(c *cli.Context, manifest []byte) ([]enamlbosh.BoshTask, error) {
	httpClient := defaultHTTPClient(c.Parent().Bool("ssl-ignore"))
	return ProcessProductBytes(
		manifest,
		c.Parent().Bool("print-manifest"),
		c.Parent().String("bosh-user"),
		c.Parent().String("bosh-pass"),
		c.Parent().String("bosh-url"),
		c.Parent().Int("bosh-port"),
		httpClient,
	)
}

func processCloudConfig(c *cli.Context, manifest []byte) (e error) {

	if c.Parent().Bool("print-manifest") {
		yamlString := string(manifest)
		fmt.Println(yamlString)

	} else {
		ccm := enaml.NewCloudConfigManifest(manifest)
		boshclient := enamlbosh.NewClient(c.Parent().String("bosh-user"), c.Parent().String("bosh-pass"), c.Parent().String("bosh-url"), c.Parent().Int("bosh-port"))
		if req, err := boshclient.NewCloudConfigRequest(*ccm); err == nil {
			httpClient := defaultHTTPClient(c.Parent().Bool("ssl-ignore"))

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

func defaultHTTPClient(sslIngore bool) *http.Client {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: sslIngore},
	}
	return &http.Client{Transport: tr}
}
