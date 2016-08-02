package gcp

import (
	"fmt"

	"github.com/codegangsta/cli"
	"github.com/enaml-ops/enaml"
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
			return
		}
		lo.G.Debug("Got boshbase", boshBase)
		utils.CheckRequired(c, "gcp-network-name", "gcp-subnetwork-name", "gcp-default-zone", "gcp-project", "gcp-machine-type", "gcp-disk-type")

		provider := boshinit.NewGCPIaaSProvider(&boshinit.GCPBoshInitConfig{
			NetworkName:    c.String("gcp-network-name"),
			SubnetworkName: c.String("gcp-subnetwork-name"),
			DefaultZone:    c.String("gcp-default-zone"),
			Project:        c.String("gcp-project"),
			MachineType:    c.String("gcp-machine-type"),
			DiskType:       c.String("gcp-disk-type"),
		}, boshBase)

		manifest := provider.CreateDeploymentManifest()

		lo.G.Debug("Got manifest", manifest)
		if yamlString, err := enaml.Paint(manifest); err == nil {

			if c.Bool("print-manifest") {
				fmt.Println(yamlString)

			} else {
				utils.DeployYaml(yamlString, boshInitDeploy)
			}
		} else {
			e = err
		}
		return
	}
}
