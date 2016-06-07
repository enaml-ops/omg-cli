package cloudfoundry

import (
	"github.com/enaml-ops/omg-cli/pluginlib/product"
	"github.com/enaml-ops/omg-cli/plugins/products/cloudfoundry/plugin"
)

func main() {
	product.Run(new(cloudfoundry.Plugin))
}
