package vault

import (
	"fmt"

	"github.com/codegangsta/cli"
	"github.com/enaml-ops/enaml"
	"github.com/enaml-ops/omg-cli/pluginlib/product"
	"github.com/enaml-ops/omg-cli/pluginlib/util"
	"github.com/enaml-ops/omg-cli/plugins/products/vault/enaml-gen/consul"
	vaultlib "github.com/enaml-ops/omg-cli/plugins/products/vault/enaml-gen/vault"
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
type Plugin struct {
	NetworkName     string
	IPs             []string
	VMTypeName      string
	DiskTypeName    string
	AZs             []string
	StemcellName    string
	StemcellURL     string
	StemcellVersion string
	StemcellSHA     string
}

func (s *Plugin) GetFlags() (flags []cli.Flag) {
	return []cli.Flag{
		cli.StringSliceFlag{Name: "ip", Usage: "multiple static ips for each redis leader vm"},
		cli.StringSliceFlag{Name: "az", Usage: "list of AZ names to use"},
		cli.StringFlag{Name: "network", Usage: "the name of the network to use"},
		cli.StringFlag{Name: "vm-type", Usage: "name of your desired vm type"},
		cli.StringFlag{Name: "disk-type", Usage: "name of your desired disk type"},
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
	var err error
	c := pluginutil.NewContext(args, s.GetFlags())

	s.IPs = c.StringSlice("ip")
	s.AZs = c.StringSlice("az")
	s.NetworkName = c.String("network")
	s.StemcellName = c.String("stemcell-name")
	s.StemcellVersion = c.String("stemcell-ver")
	s.StemcellSHA = c.String("stemcell-sha")
	s.StemcellURL = c.String("stemcell-url")
	s.VMTypeName = c.String("vm-type")
	s.DiskTypeName = c.String("disk-type")

	if err = s.flagValidation(); err != nil {
		lo.G.Error("invalid arguments: ", err)
		lo.G.Panic("exiting due to invalid args")
	}

	if err = s.cloudconfigValidation(enaml.NewCloudConfigManifest(cloudConfig)); err != nil {
		lo.G.Error("invalid settings for cloud config on target bosh: ", err)
		lo.G.Panic("your deployment is not compatible with your cloud config, exiting")
	}
	lo.G.Debug("context", c)
	var dm = new(enaml.DeploymentManifest)
	dm.SetName("vault")
	dm.AddRemoteRelease("vault", BoshVaultReleaseVer, BoshVaultReleaseURL, BoshVaultReleaseSHA)
	dm.AddRemoteRelease("consul", BoshConsulReleaseVer, BoshConsulReleaseURL, BoshConsulReleaseSHA)
	dm.AddRemoteStemcell("bosh-aws-xen-ubuntu-trusty-go_agent", s.StemcellName, s.StemcellVersion, s.StemcellURL, s.StemcellSHA)

	dm.AddInstanceGroup(s.NewVaultInstanceGroup())
	dm.Update = enaml.Update{
		MaxInFlight:     1,
		UpdateWatchTime: "30000-300000",
		CanaryWatchTime: "30000-300000",
		Serial:          false,
		Canaries:        1,
	}
	return dm.Bytes()
}

func (s *Plugin) NewVaultInstanceGroup() (ig *enaml.InstanceGroup) {
	ig = &enaml.InstanceGroup{
		Name:               "vault",
		Instances:          len(s.IPs),
		VMType:             s.VMTypeName,
		AZs:                s.AZs,
		Stemcell:           s.StemcellName,
		PersistentDiskType: s.DiskTypeName,
		Jobs: []enaml.InstanceJob{
			s.createVaultJob(),
			s.createConsulJob(),
		},
		Networks: []enaml.Network{
			enaml.Network{Name: s.NetworkName, StaticIPs: s.IPs},
		},
		Update: enaml.Update{
			MaxInFlight: 1,
		},
	}
	return
}

func (s *Plugin) createVaultJob() enaml.InstanceJob {
	return enaml.InstanceJob{
		Name:    "vault",
		Release: "vault",
		Properties: &vaultlib.VaultJob{
			Vault: &vaultlib.Vault{
				Backend: &vaultlib.Backend{
					UseConsul: true,
				},
			},
		},
	}
}
func (s *Plugin) createConsulJob() enaml.InstanceJob {
	return enaml.InstanceJob{
		Name:    "consul",
		Release: "consul",
		Properties: &consul.ConsulJob{
			Consul: &consul.Consul{
				JoinHosts: s.IPs,
			},
		},
	}
}

func (s *Plugin) cloudconfigValidation(cloudConfig *enaml.CloudConfigManifest) (err error) {
	lo.G.Debug("running cloud config validation")

	for _, vmtype := range cloudConfig.VMTypes {
		err = fmt.Errorf("vm size %s does not exist in cloud config. options are: %v", s.VMTypeName, cloudConfig.VMTypes)
		if vmtype.Name == s.VMTypeName {
			err = nil
			break
		}
	}

	for _, disktype := range cloudConfig.DiskTypes {
		err = fmt.Errorf("disk size %s does not exist in cloud config. options are: %v", s.DiskTypeName, cloudConfig.DiskTypes)
		if disktype.Name == s.DiskTypeName {
			err = nil
			break
		}
	}

	for _, net := range cloudConfig.Networks {
		err = fmt.Errorf("network %s does not exist in cloud config. options are: %v", s.NetworkName, cloudConfig.Networks)
		if net.(map[interface{}]interface{})["name"] == s.NetworkName {
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

func (s *Plugin) flagValidation() (err error) {
	lo.G.Debug("validating given flags")

	if len(s.IPs) <= 0 {
		err = fmt.Errorf("no `ip` given")
	}
	if len(s.AZs) <= 0 {
		err = fmt.Errorf("no `az` given")
	}

	if s.NetworkName == "" {
		err = fmt.Errorf("no `network-name` given")
	}

	if s.VMTypeName == "" {
		err = fmt.Errorf("no `vm-type` given")
	}
	if s.DiskTypeName == "" {
		err = fmt.Errorf("no `disk-type` given")
	}

	if s.StemcellURL == "" {
		err = fmt.Errorf("no `stemcell-url` given")
	}

	if s.StemcellVersion == "" {
		err = fmt.Errorf("no `stemcell-ver` given")
	}

	if s.StemcellSHA == "" {
		err = fmt.Errorf("no `stemcell-sha` given")
	}
	return
}
