package redis

import (
	"github.com/enaml-ops/omg-cli/pluginlib/product"
	"github.com/enaml-ops/omg-cli/plugins/products/redis/plugin"
)

func main() {
	product.Run(new(redis.Plugin))
}
