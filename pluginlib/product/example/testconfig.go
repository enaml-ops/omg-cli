package main

import (
	"github.com/enaml-ops/pluginlib/cred"
	"github.com/enaml-ops/pluginlib/pcli"
	"github.com/enaml-ops/pluginlib/productv1"
)

func main() {
	product.Run(new(MyProduct))
}

type MyProduct struct{}

func (s *MyProduct) GetFlags() (flags []pcli.Flag) {
	return
}

func (s *MyProduct) GetMeta() product.Meta {
	return product.Meta{
		Name: "myfakeproduct",
	}
}

func (s *MyProduct) GetProduct(args []string, cloudconfig []byte, cs cred.Store) ([]byte, error) {
	return []byte(""), nil
}
