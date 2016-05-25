package main

import (
	"github.com/codegangsta/cli"
	"github.com/enaml-ops/omg-cli/pluginlib/product"
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

func (s *MyProduct) GetProduct(args []string, cloudconfig []byte) []byte {
	return []byte("")
}
