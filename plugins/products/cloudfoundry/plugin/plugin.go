package cloudfoundry

import (
	"github.com/codegangsta/cli"
	"github.com/enaml-ops/enaml"
	"github.com/enaml-ops/omg-cli/pluginlib/product"
	"github.com/enaml-ops/omg-cli/pluginlib/util"
	"github.com/xchapter7x/lo"
)

func (s *Plugin) GetFlags() (flags []cli.Flag) {
	return []cli.Flag{
		cli.StringFlag{Name: "stemcell-name", Usage: "the name of your desired stemcell"},
		cli.StringSliceFlag{Name: "az", Usage: "list of AZ names to use"},
		cli.StringSliceFlag{Name: "router-ip", Usage: "a list of the router ips you wish to use"},
		cli.StringFlag{Name: "router-network", Usage: "the name of the network you wish to place your routers in"},
		cli.StringFlag{Name: "router-vm-type", Usage: "the name of your desired vm size"},
		cli.StringFlag{Name: "router-ssl-cert-file", Usage: "the file location of your go router ssl cert"},
		cli.StringFlag{Name: "router-ssl-cert", Usage: "the go router ssl cert"},
		cli.StringFlag{Name: "router-ssl-key-file", Usage: "the file location of your go router ssl key"},
		cli.StringFlag{Name: "router-ssl-key", Usage: "the go router ssl key"},
		cli.StringFlag{Name: "router-user", Value: "router_status", Usage: "the username of the go-routers"},
		cli.StringFlag{Name: "router-pass", Usage: "the password of the go-routers"},
		cli.BoolFlag{Name: "router-enable-ssl", Usage: "enable or disable ssl on your routers"},
		cli.StringFlag{Name: "nats-user", Value: "nats", Usage: "username for your nats pool"},
		cli.StringFlag{Name: "nats-pass", Value: "nats-password", Usage: "password for your nats pool"},
		cli.StringSliceFlag{Name: "nats-machine-ip", Usage: "ip of a nats node vm"},
		cli.StringSliceFlag{Name: "etcd-machine-ip", Usage: "ip of a etcd node vm"},
		cli.StringFlag{Name: "metron-zone", Usage: "zone guid for the metron agent"},
		cli.StringFlag{Name: "metron-secret", Usage: "shared secret for the metron agent endpoint"},
	}
}

func (s *Plugin) GetMeta() product.Meta {
	return product.Meta{
		Name: "cloudfoundry",
	}
}

func (s *Plugin) GetProduct(args []string, cloudConfig []byte) (b []byte) {
	c := pluginutil.NewContext(args, s.GetFlags())
	dm := enaml.NewDeploymentManifest([]byte(``))
	dm.SetName(DeploymentName)
	dm.AddRelease(enaml.Release{Name: CFReleaseName, Version: CFReleaseVersion})
	dm.AddStemcell(enaml.Stemcell{OS: StemcellName, Version: StemcellVersion, Alias: StemcellAlias})

	if goRouterPartition, err := NewGoRouterPartition(c); err == nil {
		ig := goRouterPartition.ToInstanceGroup()
		lo.G.Debug("instance-group: ", ig)
		dm.AddInstanceGroup(ig)

	} else {
		lo.G.Error("invalid go router group response:", err)
		lo.G.Panic("!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!! invalid go router group response:", err)
	}
	return dm.Bytes()
}
