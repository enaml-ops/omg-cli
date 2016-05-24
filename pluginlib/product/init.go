package product

import (
	"encoding/gob"

	"github.com/codegangsta/cli"
)

func init() {
	gob.Register(cli.StringSliceFlag{})
	gob.Register(cli.StringFlag{})
	gob.Register(cli.BoolFlag{})
	gob.Register(cli.BoolTFlag{})
	gob.Register(cli.DurationFlag{})
	gob.Register(cli.GenericFlag{})
	gob.Register(cli.IntFlag{})
	gob.Register(cli.IntSliceFlag{})
}
