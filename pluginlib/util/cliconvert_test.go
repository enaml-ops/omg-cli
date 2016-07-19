package pluginutil_test

import (
	"github.com/codegangsta/cli"
	"github.com/enaml-ops/omg-cli/pluginlib/pcli"
	. "github.com/enaml-ops/omg-cli/pluginlib/util"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("given ToCliFlagArray", func() {
	Context("when called with a []pcli.Flag", func() {
		controlFlags := []pcli.Flag{
			pcli.StringFlag{Name: "blahstring"},
			pcli.StringSliceFlag{Name: "blahslice"},
			pcli.IntFlag{Name: "blahint"},
			pcli.BoolFlag{Name: "blahbool"},
			pcli.BoolTFlag{Name: "blahboolt"},
		}
		It("then it should convert to a []cli.Flag", func() {
			cliFlags := ToCliFlagArray(controlFlags)
			Ω(cliFlags).ShouldNot(Equal(make([]cli.Flag, 0)))
			Ω(len(cliFlags)).Should(Equal(len(controlFlags)))
		})
	})
})
