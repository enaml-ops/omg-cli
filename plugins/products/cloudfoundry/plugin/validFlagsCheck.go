package cloudfoundry

import (
	"github.com/codegangsta/cli"
	"github.com/xchapter7x/lo"
)

func hasValidStringFlags(c *cli.Context, flaglist []string) bool {
	for _, v := range flaglist {

		if c.String(v) == "" {
			lo.G.Error("empty flag value for required field: ", v)
			return false
		}
	}
	return true
}

func hasValidStringSliceFlags(c *cli.Context, flaglist []string) bool {

	for _, v := range flaglist {

		if len(c.StringSlice(v)) > 0 {
			lo.G.Error("empty flag value for required field: az")
			return false
		}
	}
	return true
}
