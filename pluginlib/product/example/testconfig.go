package main

import (
	"github.com/enaml-ops/pluginlib/pcli"
	"github.com/enaml-ops/pluginlib/product"
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

func (s *MyProduct) GetProduct(args []string, cloudconfig []byte) ([]byte, error) {
	return []byte(""), nil
}
