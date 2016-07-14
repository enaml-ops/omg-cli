package pluginutil_test

import (
	"net/http"
	"os"

	"github.com/codegangsta/cli"
	. "github.com/enaml-ops/omg-cli/pluginlib/util"
	"github.com/enaml-ops/omg-cli/pluginlib/util/utilfakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("given: a VaultOverlay", func() {
	Describe("given a properly initialized vaultoverlay targeting a vault", func() {
		var vault VaultUnmarshaler

		BeforeEach(func() {
			doer := new(utilfakes.FakeDoer)
			b, _ := os.Open("fixtures/vault.json")
			doer.DoReturns(&http.Response{
				Body: b,
			}, nil)
			vault = NewVaultUnmarshal("domain.com", "my-really-long-token", doer)
		})

		Context("when calling unmarshalflags on a context that was not given the flag value from the cli", func() {
			var ctx *cli.Context

			BeforeEach(func() {
				flgs := []cli.Flag{
					cli.StringFlag{Name: "knock"},
				}
				vault.UnmarshalFlags("secret/move-along-nothing-to-see-here", flgs)
				ctx = NewContext([]string{"mycoolapp"}, flgs)
			})

			It("should set the value in the flag using the given vault hash", func() {
				Ω(ctx.String("knock")).Should(Equal("knocks"))
			})
		})

		Context("when calling unmarshalflags on a context that was given the flag value from the cli", func() {
			var ctx *cli.Context

			BeforeEach(func() {
				flgs := []cli.Flag{
					cli.StringFlag{Name: "knock"},
				}
				vault.UnmarshalFlags("secret/move-along-nothing-to-see-here", flgs)
				ctx = NewContext([]string{"mycoolapp", "--knock", "something-different"}, flgs)
			})

			It("should overwrite the default vault value with the cli flag value given", func() {
				Ω(ctx.String("knock")).ShouldNot(Equal("knocks"))
				Ω(ctx.String("knock")).Should(Equal("something-different"))
			})
		})

		Context("when calling unmarshalflags on a context that was given a stringslice value from the cli", func() {
			var ctx *cli.Context

			BeforeEach(func() {
				doer := new(utilfakes.FakeDoer)
				b, _ := os.Open("fixtures/vaultslice.json")

				doer.DoReturns(&http.Response{
					Body: b,
				}, nil)
				vault = NewVaultUnmarshal("domain.com", "my-really-long-token", doer)
				flgs := []cli.Flag{
					cli.StringSliceFlag{Name: "knock-slice"},
					cli.StringFlag{Name: "stuff"},
				}
				vault.UnmarshalFlags("secret/move-along-nothing-to-see-here", flgs)
				ctx = NewContext([]string{"mycoolapp", "--stuff", "with-val"}, flgs)
			})

			It("should overwrite the value in the flag using the given vault hash", func() {
				Ω(ctx.StringSlice("knock-slice")).Should(ConsistOf("knocks-1", "knocks-2", "knocks-3"))
			})
		})

		Context("when calling unmarshalflags on a context which was not defined with the flag contained in vault", func() {
			var ctx *cli.Context

			BeforeEach(func() {
				flgs := []cli.Flag{
					cli.StringFlag{Name: "badda"},
				}
				vault.UnmarshalFlags("secret/move-along-nothing-to-see-here", flgs)
				ctx = NewContext([]string{"mycoolapp"}, flgs)
			})

			It("then it should not set or create the flag in the context", func() {
				Ω(ctx.String("knock")).Should(BeEmpty())
			})
		})
	})
})
