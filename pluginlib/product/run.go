package product

import (
	"os"

	"github.com/hashicorp/go-plugin"
)

func Run(cc ProductDeployer) {
	if len(os.Args) >= 2 && os.Args[1] != "" {
		plugin.Serve(&plugin.ServeConfig{
			HandshakeConfig: HandshakeConfig,
			Plugins: map[string]plugin.Plugin{
				PluginsMapHash: NewProductPlugin(cc),
			},
		})
		return
	}
}
