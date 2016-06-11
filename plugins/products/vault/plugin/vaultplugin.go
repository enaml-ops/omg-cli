package vault

import (
	"fmt"

	"github.com/codegangsta/cli"
	"github.com/enaml-ops/enaml"
	"github.com/enaml-ops/omg-cli/pluginlib/product"
	"github.com/enaml-ops/omg-cli/pluginlib/util"
	"github.com/xchapter7x/lo"
)

const (
	BoshVaultReleaseURL  = "https://bosh.io/d/github.com/cloudfoundry-community/vault-boshrelease"
	BoshVaultReleaseVer  = "0.3.0"
	BoshVaultReleaseSHA  = "bd1ae82104dcf36abf64875fc5a46e1661bf2eac"
	BoshConsulReleaseURL = "https://bosh.io/d/github.com/cloudfoundry-community/consul-boshrelease"
	BoshConsulReleaseVer = "20"
	BoshConsulReleaseSHA = "9a0591c6b4d88d7d756ea933e14ddf112d05f334"
)

type jobBucket struct {
	JobName   string
	JobType   int
	Instances int
}
type Plugin struct{}

func (s *Plugin) GetFlags() (flags []cli.Flag) {
	return []cli.Flag{
		cli.StringSliceFlag{Name: "ip", Usage: "multiple static ips for each redis leader vm"},
		cli.StringFlag{Name: "disk-size", Value: "4096", Usage: "size of disk on VMs"},
		cli.StringFlag{Name: "network-name", Usage: "name of your target network"},
		cli.StringFlag{Name: "vm-size", Usage: "name of your desired vm size"},
		cli.StringFlag{Name: "stemcell-url", Usage: "the url of the stemcell you wish to use"},
		cli.StringFlag{Name: "stemcell-ver", Usage: "the version number of the stemcell you wish to use"},
		cli.StringFlag{Name: "stemcell-sha", Usage: "the sha of the stemcell you will use"},
		cli.StringFlag{Name: "stemcell-name", Value: "trusty", Usage: "the name of the stemcell you will use"},
	}
}

func (s *Plugin) GetMeta() product.Meta {
	return product.Meta{
		Name: "vault",
	}
}

func (s *Plugin) GetProduct(args []string, cloudConfig []byte) (b []byte) {
	c := pluginutil.NewContext(args, s.GetFlags())

	if err := s.flagValidation(c); err != nil {
		lo.G.Error("invalid arguments: ", err)
		lo.G.Panic("exiting due to invalid args")
	}

	if err := s.cloudconfigValidation(c, enaml.NewCloudConfigManifest(cloudConfig)); err != nil {
		lo.G.Error("invalid settings for cloud config on target bosh: ", err)
		lo.G.Panic("your deployment is not compatible with your cloud config, exiting")
	}
	lo.G.Debug("context", c)
	var dm = new(enaml.DeploymentManifest)
	dm.SetName("enaml-vault")
	dm.AddRemoteRelease("vault", BoshVaultReleaseVer, BoshVaultReleaseURL, BoshVaultReleaseSHA)
	dm.AddRemoteRelease("consul", BoshConsulReleaseVer, BoshConsulReleaseURL, BoshConsulReleaseSHA)
	dm.AddRemoteStemcell(c.String("stemcell-name"), c.String("stemcell-name"), c.String("stemcell-ver"), c.String("stemcell-url"), c.String("stemcell-sha"))
	dm.AddJob(NewVaultJob("vault", c.String("network-name"), c.String("disk-size"), c.String("vm-size"), c.StringSlice("ip")))
	return dm.Bytes()
}

func (s *Plugin) cloudconfigValidation(c *cli.Context, cloudConfig *enaml.CloudConfigManifest) (err error) {
	lo.G.Debug("running cloud config validation")
	var vmsize = c.String("vm-size")
	var disksize = c.String("disk-size")
	var netname = c.String("network-name")

	for _, vmtype := range cloudConfig.VMTypes {
		err = fmt.Errorf("vm size %s does not exist in cloud config. options are: %v", vmsize, cloudConfig.VMTypes)
		if vmtype.Name == vmsize {
			err = nil
			break
		}
	}

	for _, disktype := range cloudConfig.DiskTypes {
		err = fmt.Errorf("disk size %s does not exist in cloud config. options are: %v", disksize, cloudConfig.DiskTypes)
		if disktype.Name == disksize {
			err = nil
			break
		}
	}

	for _, net := range cloudConfig.Networks {
		err = fmt.Errorf("network %s does not exist in cloud config. options are: %v", netname, cloudConfig.Networks)
		if net.(map[interface{}]interface{})["name"] == netname {
			err = nil
			break
		}
	}

	if len(cloudConfig.VMTypes) == 0 {
		err = fmt.Errorf("no vm sizes found in cloud config")
	}

	if len(cloudConfig.DiskTypes) == 0 {
		err = fmt.Errorf("no disk sizes found in cloud config")
	}

	if len(cloudConfig.Networks) == 0 {
		err = fmt.Errorf("no networks found in cloud config")
	}
	return
}

func (s *Plugin) flagValidation(c *cli.Context) (err error) {
	lo.G.Debug("validating given flags")

	if len(c.StringSlice("ip")) <= 0 {
		err = fmt.Errorf("no `ip` given")
	}

	if len(c.String("network-name")) <= 0 {
		err = fmt.Errorf("no `network-name` given")
	}

	if len(c.String("vm-size")) <= 0 {
		err = fmt.Errorf("no `vm-size` given")
	}

	if len(c.String("stemcell-url")) <= 0 {
		err = fmt.Errorf("no `stemcell-url` given")
	}

	if len(c.String("stemcell-ver")) <= 0 {
		err = fmt.Errorf("no `stemcell-ver` given")
	}

	if len(c.String("stemcell-sha")) <= 0 {
		err = fmt.Errorf("no `stemcell-sha` given")
	}
	return
}

func NewVaultJob(name, networkName, disk, vmSize string, ips []string) (job enaml.Job) {
	network := enaml.Network{
		Name:      networkName,
		StaticIPs: ips,
	}
	properties := enaml.Properties{
		"consul": map[string]interface{}{
			"join_hosts": ips,
		},
		"vault": map[string]interface{}{
			"backend": map[string]interface{}{
				"use_consul": true,
			},
		},
	}

	job = enaml.Job{
		Name:       name,
		Properties: properties,
		Instances:  len(ips),
		Networks: []enaml.Network{
			network,
		},
		Templates: []enaml.Template{
			enaml.Template{Name: "vault", Release: "vault"},
			enaml.Template{Name: "consul", Release: "consul"},
		},
		PersistentDisk: disk,
		ResourcePool:   vmSize,
	}
	lo.G.Debug("job", job)
	return
}
