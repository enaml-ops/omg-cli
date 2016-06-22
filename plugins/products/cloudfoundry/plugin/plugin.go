package cloudfoundry

import (
	"github.com/codegangsta/cli"
	"github.com/enaml-ops/enaml"
	"github.com/enaml-ops/omg-cli/pluginlib/product"
	"github.com/enaml-ops/omg-cli/pluginlib/util"
)

func (s *Plugin) GetFlags() (flags []cli.Flag) {
	return []cli.Flag{}
}

func (s *Plugin) GetMeta() product.Meta {
	return product.Meta{
		Name: "cf-shortcut",
	}
}

func (s *Plugin) GetProduct(args []string, cloudConfig []byte) (b []byte) {
	_ = pluginutil.NewContext(args, s.GetFlags())
	dm := enaml.NewDeploymentManifest([]byte(legacyYamlManifest))
	return dm.Bytes()
}
