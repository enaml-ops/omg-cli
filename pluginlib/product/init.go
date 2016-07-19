package product

import (
	"encoding/gob"

	"github.com/enaml-ops/omg-cli/pluginlib/pcli"
)

func init() {
	gob.Register(pcli.StringSliceFlag{})
	gob.Register(pcli.StringFlag{})
	gob.Register(pcli.BoolFlag{})
	gob.Register(pcli.BoolTFlag{})
	gob.Register(pcli.IntFlag{})
}
