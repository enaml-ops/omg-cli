package pluginutil

import (
	"github.com/codegangsta/cli"
	"github.com/enaml-ops/omg-cli/pluginlib/pcli"
)

func ToCliFlagArray(fs []pcli.Flag) (cliFlags []cli.Flag) {
	for _, f := range fs {
		cliFlags = append(cliFlags, f.ToCli().(cli.Flag))
	}
	return
}
