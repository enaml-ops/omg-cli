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
			pcli.NewFlag(pcli.StringFlag, "blahstring", "", ""),
			pcli.NewFlag(pcli.StringSliceFlag, "blahslice", "", ""),
			pcli.NewFlag(pcli.IntFlag, "blahint", "", ""),
			pcli.NewFlag(pcli.BoolFlag, "blahbool", "", ""),
			pcli.NewFlag(pcli.BoolTFlag, "blahboolt", "", ""),
		}

		It("then it should convert to a []cli.Flag", func() {
			cliFlags := ToCliFlagArray(controlFlags)
			Ω(len(cliFlags)).Should(Equal(len(controlFlags)))

			Ω(func() {
				_ = cliFlags[0].(cli.StringFlag)
				_ = cliFlags[1].(cli.StringSliceFlag)
				_ = cliFlags[2].(cli.IntFlag)
				_ = cliFlags[3].(cli.BoolFlag)
				_ = cliFlags[4].(cli.BoolTFlag)
			}).ShouldNot(Panic())
		})
	})

	Context("when converting a string slice flag with an empty value", func() {
		f := pcli.NewFlag(pcli.StringSliceFlag, "my-slice", "", "")
		It("should leave the value empty on the returned type", func() {
			cliFlags := ToCliFlagArray([]pcli.Flag{f})
			Ω(len(cliFlags)).Should(Equal(1))
			sf := cliFlags[0].(cli.StringSliceFlag)
			Ω(sf.Value.Value()).Should(BeEmpty())
		})
	})

	Context("when converting a string slice flag with a non-empty value", func() {
		const controlString = "blah"
		f := pcli.NewFlag(pcli.StringSliceFlag, "my-slice", "", controlString)
		cliFlags := ToCliFlagArray([]pcli.Flag{f})
		It("should set the default value on the returned type", func() {
			Ω(len(cliFlags)).Should(Equal(1))
			sf := cliFlags[0].(cli.StringSliceFlag)
			Ω(sf.Value.Value()).Should(ConsistOf(controlString))
		})
	})
})
