package pluginutil_test

import (
	. "github.com/bosh-ops/bosh-install/plugin/util"
	"github.com/codegangsta/cli"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("NewContext function", func() {
	Context("when called with valid args and flags", func() {
		It("then it should return a properly init'd cli.context", func() {
			ctx := NewContext([]string{"test", "--this", "that"}, []cli.Flag{
				cli.StringFlag{Name: "this"},
			})
			Ω(ctx.String("this")).Should(Equal("that"))
		})
	})
})
