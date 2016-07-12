package concourseplugin

import (
	"github.com/codegangsta/cli"
	"github.com/enaml-ops/omg-cli/pluginlib/product"
	"github.com/enaml-ops/omg-cli/pluginlib/util"
	"github.com/xchapter7x/lo"
)

type ConcoursePlugin struct{}

func (s *ConcoursePlugin) GetFlags() (flags []cli.Flag) {
	return generateFlags()
}

func (s *ConcoursePlugin) GetMeta() product.Meta {
	return product.Meta{
		Name: "concourse",
	}
}

func (s *ConcoursePlugin) GetProduct(args []string, cloudConfig []byte) (b []byte) {
	if len(cloudConfig) == 0 {
		lo.G.Debug("concourseplugin: empty cloud config")
		panic("cloud config cannot be empty")
	}
	c := pluginutil.NewContext(args, s.GetFlags())
	dm := NewDeploymentManifest(c, cloudConfig)
	return dm.Bytes()
}
