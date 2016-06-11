package vault

import (
	"github.com/enaml-ops/omg-cli/pluginlib/product"
	"github.com/enaml-ops/omg-cli/plugins/products/vault/plugin"
)

func main() {
	product.Run(new(vault.Plugin))
}
