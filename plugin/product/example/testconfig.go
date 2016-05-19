package main

import (
	"github.com/bosh-ops/bosh-install/plugin/product"
	"github.com/codegangsta/cli"
	"github.com/xchapter7x/enaml"
)

func main() {
	product.Run(new(MyProduct))
}

type MyProduct struct{}

func (s *MyProduct) GetFlags() (flags []cli.Flag) {
	return
}

func (s *MyProduct) GetMeta() product.Meta {
	return product.Meta{
		Name: "myfakeproduct",
	}
}

func (s *MyProduct) GetProduct(args []string) enaml.DeploymentManifest {
	return enaml.DeploymentManifest{}
}
