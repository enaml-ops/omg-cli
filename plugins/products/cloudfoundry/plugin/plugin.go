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
		cli.StringSliceFlag{Name: "router-ip", Usage: "a list of the router ips you wish to use"},
		cli.StringFlag{Name: "router-network", Usage: "the name of the network you wish to place your routers in"},
		cli.StringFlag{},
		cli.StringFlag{},
	}
}

func (s *Plugin) GetMeta() product.Meta {
	return product.Meta{
		Name: "cf-shortcut",
	}
}

func (s *Plugin) GetProduct(args []string, cloudConfig []byte) (b []byte) {
	c := pluginutil.NewContext(args, s.GetFlags())
	dm := enaml.NewDeploymentManifest([]byte(``))

	if goRouterPartition, err := NewGoRouterPartition(c); err == nil {
		ig := goRouterPartition.ToInstanceGroup()
		lo.G.Debug("grrrrr", ig)
		dm.AddInstanceGroup(ig)

	} else {
		lo.G.Error("invalid go router group response:", err)
		lo.G.Panic("!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!! invalid go router group response:", err)
	}
	return dm.Bytes()
}
