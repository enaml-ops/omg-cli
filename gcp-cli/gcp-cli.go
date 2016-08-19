package gcpcli

import (
	"github.com/codegangsta/cli"
	"github.com/enaml-ops/omg-cli/plugins/products/bosh-init"
	"github.com/enaml-ops/omg-cli/utils"
	"github.com/xchapter7x/lo"
)

func GetFlags() []cli.Flag {
	boshdefaults := boshinit.NewGCPBoshBase()

	boshFlags := boshinit.BoshFlags(boshdefaults)
	gcpFlags := []cli.Flag{
		cli.StringFlag{Name: "gcp-network-name", Usage: "the GCP network name"},
		cli.StringFlag{Name: "gcp-subnetwork-name", Usage: "the GCP subnetwork"},
		cli.StringFlag{Name: "gcp-default-zone", Usage: "the default GCP zone"},
		cli.StringFlag{Name: "gcp-project", Usage: "the GCP project"},
		cli.StringFlag{Name: "gcp-machine-type", Value: "n1-standard-4", Usage: "GCP machine type"},
		cli.StringFlag{Name: "gcp-disk-type", Value: "pd-standard", Usage: "root disk type property"},
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
