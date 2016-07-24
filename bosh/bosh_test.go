// Note: this test is in `package bosh` and not package `bosh_test`
// because we need to be able to override some internal behavior.
//
// The package is coded up to reuse the same enamlbosh client for all
// calls (so that we don't need to repeatedly acquire UAA tokens).
//
// This works well for the CLI tool, but we want to be able to test
// both UAA and basic auth, so for tests we want to be able to
// create new clients instead of reusing them.

package bosh

import (
	"net"
	"net/http"
	"net/url"

	"github.com/codegangsta/cli"
	"github.com/enaml-ops/enaml/enamlbosh"
	"github.com/enaml-ops/pluginlib/cloudconfig"
	"github.com/enaml-ops/pluginlib/pcli"
	"github.com/enaml-ops/pluginlib/product"
	"github.com/enaml-ops/pluginlib/util"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
)

const (
	controlCloudConfig = "MyCloudConfig"
	controlProduct     = "MyProduct"
	basicAuthUser      = "user"
	basicAuthPass      = "pass"
	clientID           = "clientid"
	clientSecret       = "clientsecret"
)

const tokenResponse = `{
  "access_token":"abcdef01234567890",
  "token_type":"bearer",
  "refresh_token":"0987654321fedcba",
  "expires_in":3599,
  "scope":"opsman.user uaa.admin scim.read opsman.admin scim.write",
  "jti":"foo"
}`

type FakeCloudConfigDeployer struct{}

func (f FakeCloudConfigDeployer) GetMeta() cloudconfig.Meta {
	return cloudconfig.Meta{Name: "test"}
}

func (f FakeCloudConfigDeployer) GetFlags() []cli.Flag {
	return []cli.Flag{}
}

func (f FakeCloudConfigDeployer) GetCloudConfig(args []string) []byte {
	return []byte(controlCloudConfig)
}

type FakeProductDeployer struct{}

func (f FakeProductDeployer) GetMeta() product.Meta {
	return product.Meta{Name: "test"}
}

func (f FakeProductDeployer) GetFlags() []pcli.Flag {
	return []pcli.Flag{}
}

func (f FakeProductDeployer) GetProduct(args []string, cloudConfig []byte) []byte {
	return []byte(controlProduct)
}

func portAndURL(s *ghttp.Server) (port, host string) {
	u, _ := url.Parse(s.URL())
	host, port, _ = net.SplitHostPort(u.Host)
	return port, u.Scheme + "://" + host
}

var _ = Describe("bosh", func() {
	BeforeEach(func() {
		// install a no-op print function to keep test output clean
		UIPrint = func(a ...interface{}) (int, error) { return 0, nil }

		// make sure we generate a new client using the CLI context
		// intead of reusing an old one
		boshclient = nil
	})

	Describe("CloudConfigAction", func() {
		var ccd cloudconfig.CloudConfigDeployer
		var server *ghttp.Server

		BeforeEach(func() {
			ccd = FakeCloudConfigDeployer{}
			server = ghttp.NewTLSServer()
		})

		AfterEach(func() {
			server.Close()
		})

		Context("when called with the --print-manifest option", func() {
			var c *cli.Context
			var err error

			BeforeEach(func() {
				port, url := portAndURL(server)
				c = pluginutil.NewContext([]string{
					"foo",
					"--print-manifest",
					"--ssl-ignore",
					"--bosh-url", url,
					"--bosh-port", port,
					"--bosh-user", basicAuthUser,
					"--bosh-pass", basicAuthPass,
				}, GetAuthFlags())
				err = CloudConfigAction(c, ccd)
			})

			It("returns without error", func() {
				Ω(err).ShouldNot(HaveOccurred())
			})

			It("doesn't attempt to upload a cloud config", func() {
				Ω(server.ReceivedRequests()).Should(BeEmpty())
			})
		})

		Context("when called without UAA flags", func() {
			var c *cli.Context
			var err error

			BeforeEach(func() {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyBasicAuth(basicAuthUser, basicAuthPass),
						ghttp.RespondWithJSONEncoded(http.StatusOK, struct{}{}),
					),
				)
				port, url := portAndURL(server)
				c = pluginutil.NewContext([]string{
					"foo",
					"--ssl-ignore",
					"--bosh-url", url,
					"--bosh-port", port,
					"--bosh-user", basicAuthUser,
					"--bosh-pass", basicAuthPass,
				}, GetAuthFlags())

			})

			It("makes a request and returns without error", func() {
				err = CloudConfigAction(c, ccd)
				Ω(err).ShouldNot(HaveOccurred())
				Ω(len(server.ReceivedRequests())).Should(Equal(1))
			})
		})

		Context("when called with UAA flags", func() {
			var c *cli.Context
			var err error

			BeforeEach(func() {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest("POST", "/oauth/token"),
						ghttp.RespondWith(http.StatusOK, tokenResponse, http.Header{
							"Content-Type": []string{"application/json"}}),
					),
					ghttp.CombineHandlers(
						ghttp.RespondWith(http.StatusOK, ""),
					),
				)

				port, url := portAndURL(server)
				c = pluginutil.NewContext([]string{
					"foo",
					"--ssl-ignore",
					"--bosh-url", url,
					"--bosh-port", port,
					"--bosh-user", basicAuthUser,
					"--bosh-pass", basicAuthPass,
					"--bosh-client-id", clientID,
					"--bosh-client-secret", clientSecret,
					"--uaa-url", server.URL(),
				}, GetAuthFlags())
			})

			It("makes a request using a UAA token", func() {
				err = CloudConfigAction(c, ccd)
				Ω(err).ShouldNot(HaveOccurred())

				Ω(server.ReceivedRequests()).ShouldNot(BeEmpty())
				lastReq := server.ReceivedRequests()[len(server.ReceivedRequests())-1]
				Ω(lastReq.Header.Get("Authorization")).Should(ContainSubstring("Bearer"))
			})
		})
	})

	Describe("ProductAction", func() {
		var pd product.ProductDeployer
		var server *ghttp.Server

		BeforeEach(func() {
			pd = FakeProductDeployer{}
			server = ghttp.NewTLSServer()
		})

		AfterEach(func() {
			server.Close()
		})

		Context("when called with the --print-manifest option", func() {
			var c *cli.Context
			var err error

			BeforeEach(func() {
				port, url := portAndURL(server)
				c = pluginutil.NewContext([]string{
					"foo",
					"--print-manifest",
					"--ssl-ignore",
					"--bosh-url", url,
					"--bosh-port", port,
					"--bosh-user", basicAuthUser,
					"--bosh-pass", basicAuthPass,
				}, GetAuthFlags())

				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyBasicAuth(basicAuthUser, basicAuthPass),
						ghttp.VerifyRequest("GET", "/cloud_configs"),
						ghttp.RespondWithJSONEncoded(http.StatusOK,
							[]enamlbosh.CloudConfigResponseBody{
								{Properties: "response"},
							}),
					),
				)
			})

			It("returns without error", func() {
				err = ProductAction(c, pd)
				Ω(err).ShouldNot(HaveOccurred())

				// it should GET the cloud config, but not POST the product
				Ω(len(server.ReceivedRequests())).Should(Equal(1))
			})
		})

		Context("when called without UAA flags", func() {
			var c *cli.Context
			var err error

			BeforeEach(func() {
				port, url := portAndURL(server)
				c = pluginutil.NewContext([]string{
					"foo",
					"--ssl-ignore",
					"--bosh-url", url,
					"--bosh-port", port,
					"--bosh-user", basicAuthUser,
					"--bosh-pass", basicAuthPass,
				}, GetAuthFlags())

				task := enamlbosh.BoshTask{
					State: enamlbosh.StatusDone,
					ID:    42,
				}

				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyBasicAuth(basicAuthUser, basicAuthPass),
						ghttp.VerifyRequest("GET", "/cloud_configs"),
						ghttp.RespondWithJSONEncoded(http.StatusOK,
							[]enamlbosh.CloudConfigResponseBody{
								{Properties: "response"},
							}),
					),
					ghttp.CombineHandlers(
						ghttp.VerifyRequest("POST", "/deployments"),
						ghttp.RespondWithJSONEncoded(http.StatusOK, task),
					),
					ghttp.CombineHandlers(
						ghttp.VerifyRequest("GET", "/tasks/42"),
						ghttp.RespondWithJSONEncoded(http.StatusOK, task),
					),
				)
			})

			It("returns without error", func() {
				err = ProductAction(c, pd)
				Ω(err).ShouldNot(HaveOccurred())
			})
		})
	})
})
