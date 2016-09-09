package gcpcli

import (
	"github.com/enaml-ops/omg-cli/plugins/products/bosh-init"
	"github.com/enaml-ops/omg-cli/utils"
	"github.com/enaml-ops/pluginlib/pcli"
	"github.com/xchapter7x/lo"
	"gopkg.in/urfave/cli.v2"
)

func GetFlags() []pcli.Flag {
	boshdefaults := boshinit.NewGCPBoshBase()

	boshFlags := boshinit.BoshFlags(boshdefaults)
	gcpFlags := []pcli.Flag{
		pcli.CreateStringFlag("gcp-network-name", "the GCP network name"),
		pcli.CreateStringFlag("gcp-subnetwork-name", "the GCP subnetwork"),
		pcli.CreateStringFlag("gcp-default-zone", "the default GCP zone"),
		pcli.CreateStringFlag("gcp-project", "the GCP project"),
		pcli.CreateStringFlag("gcp-machine-type", "GCP machine type", "n1-standard-4"),
		pcli.CreateStringFlag("gcp-disk-type", "root disk type property", "pd-standard"),
	}
	boshFlags = append(boshFlags, gcpFlags...)
	return boshFlags
}

func GetAction(boshInitDeploy func(string)) func(c *cli.Context) error {
	return func(c *cli.Context) (e error) {
		var boshBase *boshinit.BoshBase
		if boshBase, e = boshinit.NewBoshBase(c); e != nil {
			lo.G.Error(e.Error())
			return e
		}
		lo.G.Debug("Got boshbase", boshBase)
		if err := utils.CheckRequired(c, "gcp-network-name", "gcp-subnetwork-name", "gcp-default-zone", "gcp-project", "gcp-machine-type", "gcp-disk-type"); err != nil {
			lo.G.Error(err.Error())
			return err
		}

		provider := boshinit.NewGCPIaaSProvider(&boshinit.GCPBoshInitConfig{
			NetworkName:    c.String("gcp-network-name"),
			SubnetworkName: c.String("gcp-subnetwork-name"),
			DefaultZone:    c.String("gcp-default-zone"),
			Project:        c.String("gcp-project"),
			MachineType:    c.String("gcp-machine-type"),
			DiskType:       c.String("gcp-disk-type"),
		}, boshBase)

		if err := boshBase.HandleDeployment(provider, boshInitDeploy); err != nil {
			return err
		}

		return nil
	}
}
