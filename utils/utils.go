package utils

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"time"

	"github.com/codegangsta/cli"
	"github.com/enaml-ops/enaml"
	"github.com/enaml-ops/enaml/enamlbosh"
	"github.com/enaml-ops/omg-cli/pluginlib/registry"
	"github.com/enaml-ops/omg-cli/pluginlib/util"
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
				return processCloudConfig(c, manifest)
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
			Flags: pluginutil.ToCliFlagArray(flags),
			Action: func(c *cli.Context) (err error) {
				var cloudConfig *enaml.CloudConfigManifest
				client, productDeployment := registry.GetProductReference(pluginPath)
				defer client.Kill()
				boshclient := enamlbosh.NewClient(c.Parent().String("bosh-user"), c.Parent().String("bosh-pass"), c.Parent().String("bosh-url"), c.Parent().Int("bosh-port"))
				httpClient := defaultHTTPClient(c.Parent().Bool("ssl-ignore"), c.Parent().String("bosh-user"), c.Parent().String("bosh-pass"))

				if cloudConfig, err = boshclient.GetCloudConfig(httpClient); err == nil {
					var cloudConfigBytes []byte
					var task enamlbosh.BoshTask
					cloudConfigBytes, err = cloudConfig.Bytes()
					deploymentManifest := productDeployment.GetProduct(c.Parent().Args(), cloudConfigBytes)
					task, err = processProductDeployment(c, deploymentManifest, true)
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
func ProcessProductBytes(manifest []byte, printManifest bool, user, pass, url string, port int, httpClient HttpClientDoer, poll bool) (task enamlbosh.BoshTask, err error) {
	if printManifest {
		yamlString := string(manifest)
		UIPrint(yamlString)

	} else {
		dm := enaml.NewDeploymentManifest(manifest)
		boshclient := enamlbosh.NewClient(user, pass, url, port)
		ProcessRemoteBoshAssets(dm, boshclient, httpClient, true)
		UIPrint("Uploading product deployment...")

		if task, err = boshclient.PostDeployment(*dm, httpClient); err == nil {
			UIPrint("upload complete.")
			lo.G.Debug("res: ", task, err)

			switch task.State {
			case enamlbosh.StatusCancelled, enamlbosh.StatusError:
				err = fmt.Errorf("task is in failed state: ", task)

			default:
				if poll {
					err = PollTaskAndWait(task, boshclient, httpClient, -1)
				}
			}

		} else {
			lo.G.Error("error: ", err)
		}
	}
	return
}

func ProcessRemoteBoshAssets(dm *enaml.DeploymentManifest, boshClient *enamlbosh.Client, httpClient HttpClientDoer, poll bool) (err error) {
	var errStemcells error
	var errReleases error
	var remoteStemcells []enaml.Stemcell
	defer UIPrint("remote asset check complete.")
	UIPrint("Checking product deployment for remote assets...")

	if remoteStemcells, err = ProcessStemcellsToBeUploaded(dm.Stemcells, boshClient, httpClient); err == nil {
		if errStemcells = ProcessRemoteStemcells(remoteStemcells, boshClient, httpClient, poll); errStemcells != nil {
			lo.G.Info("issues processing stemcell: ", errStemcells)
		}
	}

	if errReleases = ProcessRemoteReleases(dm.Releases, boshClient, httpClient, poll); errReleases != nil {
		lo.G.Info("issues processing release: ", errReleases)
	}

	if errReleases != nil || errStemcells != nil {
		err = fmt.Errorf("stemcell err: %v   release err: %v", errStemcells, errReleases)
	}
	return
}

func isRemoteStemcell(stemcell enaml.Stemcell) bool {
	return stemcell.URL != "" && stemcell.SHA1 != ""
}

//ProcessRemoteStemcells - upload any remote stemcells given
func ProcessRemoteStemcells(scl []enaml.Stemcell, boshClient *enamlbosh.Client, httpClient HttpClientDoer, poll bool) (err error) {
	defer UIPrint("remote stemcells complete")
	UIPrint("Checking for remote stemcells...")

	for _, stemcell := range scl {
		var task enamlbosh.BoshTask

		if isRemoteStemcell(stemcell) {
			if task, err = boshClient.PostRemoteStemcell(stemcell, httpClient); err == nil {
				lo.G.Debug("task: ", task, err)

				switch task.State {
				case enamlbosh.StatusCancelled, enamlbosh.StatusError:
					err = fmt.Errorf("task is in failed state: ", task)

				default:
					if poll {
						err = PollTaskAndWait(task, boshClient, httpClient, -1)
					}
				}
			}
		} else {
			UIPrint(fmt.Sprintf("Remote stemcells...%s [%s] already exists", stemcell.Name, stemcell.Version))
		}
	}
	return
}

//ProcessStemcellsToBeUploaded - finds list of remote stemcells that have not been uploaded
func ProcessStemcellsToBeUploaded(stemcells []enaml.Stemcell, boshClient *enamlbosh.Client, httpClient HttpClientDoer) (remoteStemcells []enaml.Stemcell, err error) {
	var exists bool
	for _, stemcell := range stemcells {
		if isRemoteStemcell(stemcell) {
			if exists, err = boshClient.CheckRemoteStemcell(stemcell, httpClient); err == nil && !exists {
				remoteStemcells = append(remoteStemcells, stemcell)
			}
		}
	}
	return
}

//ProcessRemoteReleases - upload any remote Releases given
func ProcessRemoteReleases(rl []enaml.Release, boshClient *enamlbosh.Client, httpClient HttpClientDoer, poll bool) (err error) {
	defer UIPrint("remote releases complete")
	UIPrint("Checking for remote releases...")
	var isRemoteRelease = func(release enaml.Release) bool {
		return release.URL != "" && release.SHA1 != ""
	}

	for _, release := range rl {
		var task enamlbosh.BoshTask

		if isRemoteRelease(release) {

			if task, err = boshClient.PostRemoteRelease(release, httpClient); err == nil {
				lo.G.Debug("task: ", task, err)

				switch task.State {
				case enamlbosh.StatusCancelled, enamlbosh.StatusError:
					err = fmt.Errorf("task is in failed state: ", task)

				default:
					if poll {
						err = PollTaskAndWait(task, boshClient, httpClient, -1)
					}
				}
			}
		}
	}
	return
}

//PollTaskAndWait - will poll the give task until its status is cancelled, done
//or error. a -1 tries value indicates infinite
func PollTaskAndWait(task enamlbosh.BoshTask, boshClient *enamlbosh.Client, httpClient HttpClientDoer, tries int) (err error) {
	defer UIPrint(fmt.Sprintf("Task %s is %s", task.Description, task.State))
	UIPrint("Polling task...")
	var cnt = 0

Loop:
	for {

		if task, err = boshClient.GetTask(task.ID, httpClient); err != nil {
			break

		} else {

			switch task.State {
			case enamlbosh.StatusDone:
				UIPrint("task state %s", task.State)
				break Loop

			case enamlbosh.StatusCancelled, enamlbosh.StatusError:
				lo.G.Error("task error: ", task.State, task.Description)
				err = fmt.Errorf("%s - %s", task.State, task.Description)
				break Loop

			default:
				UIPrint(fmt.Sprintf("task is %s - %s", task.State, task.Description))
				time.Sleep(1 * time.Second)
			}
		}
		cnt += 1

		if tries != -1 && cnt >= tries {
			UIPrint("hit poll limit, exiting task poller without error")
			break Loop
		}
	}
	return
}

func processProductDeployment(c *cli.Context, manifest []byte, poll bool) (enamlbosh.BoshTask, error) {
	httpClient := defaultHTTPClient(c.Parent().Bool("ssl-ignore"), c.Parent().String("bosh-user"), c.Parent().String("bosh-pass"))
	return ProcessProductBytes(
		manifest,
		c.Parent().Bool("print-manifest"),
		c.Parent().String("bosh-user"),
		c.Parent().String("bosh-pass"),
		c.Parent().String("bosh-url"),
		c.Parent().Int("bosh-port"),
		httpClient,
		poll,
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
			httpClient := defaultHTTPClient(c.Parent().Bool("ssl-ignore"), c.Parent().String("bosh-user"), c.Parent().String("bosh-pass"))
			var res *http.Response
			var err error
			res, err = httpClient.Do(req)
			if err != nil {
				lo.G.Error("res: ", res)
				lo.G.Error("error: ", err)
				e = err
			} else {
				defer res.Body.Close()
				body, err := ioutil.ReadAll(res.Body)
				if err != nil {
					e = err
				} else {
					sbody := string(body)
					lo.G.Debugf("%d HTTP response:\n%s", res.StatusCode, sbody)
					if res.StatusCode >= 400 {
						e = fmt.Errorf("%s error pushing cloud config to BOSH: %s", res.Status, sbody)
					}
				}
			}
		} else {
			e = err
		}
	}
	return
}

//defaultHTTPClient - generates a client which can ignore ssl warnings as well
//as correctly assign the host:port on bosh rest call redirects.
func defaultHTTPClient(sslIngore bool, user, pass string) *http.Client {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: sslIngore},
	}
	client := &http.Client{Transport: tr}
	client.CheckRedirect = func(req *http.Request, via []*http.Request) (err error) {
		req.URL, err = url.Parse(req.URL.Scheme + "://" + via[0].URL.Host + req.URL.Path)
		req.SetBasicAuth(user, pass)
		lo.G.Debug("new req: ", req)
		return nil
	}
	return client
}
