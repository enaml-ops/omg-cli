package photoncli

import (
	"errors"

	"gopkg.in/urfave/cli.v2"

	"github.com/enaml-ops/omg-cli/plugins/products/bosh-init"
	"github.com/enaml-ops/omg-cli/plugins/products/bosh-init/enaml-gen/photoncpi"
	"github.com/enaml-ops/omg-cli/utils"
	"github.com/enaml-ops/pluginlib/pcli"
	"github.com/xchapter7x/lo"
)

func GetFlags() []pcli.Flag {
	boshdefaults := boshinit.NewPhotonBoshBase()

	boshFlags := boshinit.BoshFlags(boshdefaults)
	photonFlags := []pcli.Flag{
		pcli.CreateStringFlag("photon-target", "photon api endpoint http://PHOTON_CTRL_IP:9000"),
		pcli.CreateStringFlag("photon-user", "api admin user"),
		pcli.CreateStringFlag("photon-password", "api admin pass"),
		pcli.CreateBoolTFlag("photon-ignore-cert", "setting ignore cert or not"),
		pcli.CreateStringFlag("photon-project-id", "the photon project id"),
		pcli.CreateStringFlag("photon-machine-type", "photon instance type name", "core-200"),
		pcli.CreateStringFlag("photon-network-id", "the network-id to deploy your bosh onto (THIS IS NOT THE NETWORK NAME)"),
	}
	boshFlags = append(boshFlags, photonFlags...)
	return boshFlags
}

func GetAction(boshInitDeploy func(string)) func(c *cli.Context) error {
	return func(c *cli.Context) (e error) {
		var boshBase *boshinit.BoshBase
		var err error

		if boshBase, err = boshinit.NewBoshBase(c); err != nil {
			lo.G.Error(err.Error())
			return err
		}

		lo.G.Debug("Got boshbase", boshBase)
		if err = utils.CheckRequired(c, "photon-target", "photon-project-id", "photon-network-id"); err != nil {
			lo.G.Error(err.Error())
			return err
		}

		user := c.IsSet("photon-user")
		pass := c.IsSet("photon-password")

		if user != pass {
			lo.G.Error("--photon-user and --photon-password must be specified together")
			return errors.New("--photon-user and --photon-password must be specified together")
		}

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

		if err = boshBase.HandleDeployment(provider, boshInitDeploy); err != nil {
			return err
		}

		return nil
	}
}
