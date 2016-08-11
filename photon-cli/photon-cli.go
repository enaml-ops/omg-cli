package photoncli

import (
	"fmt"

	"github.com/codegangsta/cli"
	"github.com/enaml-ops/enaml"
	"github.com/enaml-ops/omg-cli/plugins/products/bosh-init"
	"github.com/enaml-ops/omg-cli/plugins/products/bosh-init/enaml-gen/photoncpi"
	"github.com/enaml-ops/omg-cli/utils"
	"github.com/xchapter7x/lo"
)

func GetFlags() []cli.Flag {
	boshdefaults := boshinit.NewPhotonBoshBase()

	boshFlags := boshinit.BoshFlags(boshdefaults)
	photonFlags := []cli.Flag{
		cli.StringFlag{Name: "photon-target", Usage: "photon api endpoint http://PHOTON_CTRL_IP:9000"},
		cli.StringFlag{Name: "photon-user", Usage: "api admin user"},
		cli.StringFlag{Name: "photon-password", Usage: "api admin pass"},
		cli.BoolTFlag{Name: "photon-ignore-cert", Usage: "setting ignore cert or not"},
		cli.StringFlag{Name: "photon-project-id", Usage: "the photon project id"},
		cli.StringFlag{Name: "photon-machine-type", Value: "core-200", Usage: "photon instance type name"},
		cli.StringFlag{Name: "photon-network-id", Usage: "the network-id to deploy your bosh onto (THIS IS NOT THE NETWORK NAME)"},
	}
	boshFlags = append(boshFlags, photonFlags...)
	return boshFlags
}

func GetAction(boshInitDeploy func(string)) func(c *cli.Context) error {
	return func(c *cli.Context) (e error) {
		var boshBase *boshinit.BoshBase
		if boshBase, e = boshinit.NewBoshBase(c); e != nil {
			return
		}
		lo.G.Debug("Got boshbase", boshBase)
		utils.CheckRequired(c, "photon-target", "photon-project-id", "photon-user", "photon-password", "photon-network-id")

		provider := boshinit.NewPhotonIaaSProvider(&boshinit.PhotonBoshInitConfig{
			Photon: photoncpi.Photon{
				Target:     c.String("photon-target"),
				User:       c.String("photon-user"),
				Password:   c.String("photon-password"),
				IgnoreCert: c.Bool("photon-ignore-cert"),
				Project:    c.String("photon-project-id"),
			},
			NetworkName: c.String("photon-network-id"),
			MachineType: c.String("photon-machine-type"),
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
