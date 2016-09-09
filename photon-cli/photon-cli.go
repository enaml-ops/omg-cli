package photoncli

import (
	"errors"

	"gopkg.in/urfave/cli.v2"
	"github.com/enaml-ops/omg-cli/plugins/products/bosh-init"
	"github.com/enaml-ops/omg-cli/plugins/products/bosh-init/enaml-gen/photoncpi"
	"github.com/enaml-ops/omg-cli/utils"
	"github.com/xchapter7x/lo"
)

func GetFlags() []cli.Flag {
	boshdefaults := boshinit.NewPhotonBoshBase(new(boshinit.BoshBase))

	boshFlags := boshinit.BoshFlags(boshdefaults)
	photonFlags := []cli.Flag{
		&cli.StringFlag{Name: "photon-target", Usage: "photon api endpoint http://PHOTON_CTRL_IP:9000"},
		&cli.StringFlag{Name: "photon-user", Usage: "api admin user"},
		&cli.StringFlag{Name: "photon-password", Usage: "api admin pass"},
		&cli.BoolFlag{Value: true,Name: "photon-ignore-cert", Usage: "setting ignore cert or not"},
		&cli.StringFlag{Name: "photon-project-id", Usage: "the photon project id"},
		&cli.StringFlag{Name: "photon-machine-type", Value: "core-200", Usage: "photon instance type name"},
		&cli.StringFlag{Name: "photon-network-id", Usage: "the network-id to deploy your bosh onto (THIS IS NOT THE NETWORK NAME)"},
	}
	boshFlags = append(boshFlags, photonFlags...)
	return boshFlags
}

func GetAction(boshInitDeploy func(string)) func(c *cli.Context) error {
	return func(c *cli.Context) (e error) {
		var b *boshinit.BoshBase
		var err error

		if b, err = boshinit.NewBoshBase(c); err != nil {
			lo.G.Error(err.Error())
			return err
		}
		boshBase := boshinit.NewPhotonBoshBase(b)

		if boshBase.CPIJobName == "" {
			lo.G.Error("sorry we could not proceed bc you did not set a cpijobname in your code.")
			return err
		}

		lo.G.Debug("Got boshbase", boshBase)
		if err := utils.CheckRequired(c, "photon-target", "photon-project-id", "photon-network-id"); err != nil {
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

		if err := boshBase.HandleDeployment(provider, boshInitDeploy); err != nil {
			return err
		}

		return nil
	}
}
