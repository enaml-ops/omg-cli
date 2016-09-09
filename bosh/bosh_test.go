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
	"io/ioutil"
	"net"
	"net/http"
	"net/url"

	"github.com/enaml-ops/enaml"
	"github.com/enaml-ops/enaml/enamlbosh"
	"github.com/enaml-ops/pluginlib/cloudconfig"
	"github.com/enaml-ops/pluginlib/pcli"
	"github.com/enaml-ops/pluginlib/product"
	"github.com/enaml-ops/pluginlib/util"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
	"gopkg.in/urfave/cli.v2"
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

const controlUUID = "31631ff9-ac41-4eba-a944-04c820633e7f"
const basicAuthBoshInfo = `{"name":"enaml-bosh","uuid":"` + controlUUID + `","version":"1.3232.2.0 (00000000)","user":null,"cpi":"aws_cpi","user_authentication":{"type":"basic","options":{}},"features":{"dns":{"status":false,"extras":{"domain_name":null}},"compiled_package_cache":{"status":false,"extras":{"provider":null}},"snapshots":{"status":false}}}`

type FakeCloudConfigDeployer struct{}

func (f FakeCloudConfigDeployer) GetMeta() cloudconfig.Meta {
	return cloudconfig.Meta{Name: "test"}
}

func (f FakeCloudConfigDeployer) GetFlags() []pcli.Flag {
	return []pcli.Flag{}
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
			server.AppendHandlers(
				// have our test server respond to /info acting as a basic-auth bosh
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/info"),
					ghttp.RespondWith(http.StatusOK, basicAuthBoshInfo),
				),
			)
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
				}, pluginutil.ToCliFlagArray(GetAuthFlags()))
				err = CloudConfigAction(c, ccd)
			})

			It("returns without error", func() {
				Ω(err).ShouldNot(HaveOccurred())
			})

			It("doesn't attempt to upload a cloud config", func() {
				Ω(server.ReceivedRequests()).Should(BeEmpty())
			})
		})

		Context("when called without the --print-manifest option", func() {
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
				}, pluginutil.ToCliFlagArray(GetAuthFlags()))
			})

			It("makes a request and returns without error", func() {
				err = CloudConfigAction(c, ccd)
				Ω(err).ShouldNot(HaveOccurred())
				// one request to the /info endpoint (when creating the client),
				// and one request to /cloud_configs
				Ω(len(server.ReceivedRequests())).Should(Equal(2))
			})
		})
	})

	Describe("ProductAction", func() {
		var pd product.ProductDeployer
		var server *ghttp.Server

		BeforeEach(func() {
			pd = FakeProductDeployer{}
			server = ghttp.NewTLSServer()
			server.AppendHandlers(
				// have our test server respond to /info acting as a basic-auth bosh
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/info"),
					ghttp.RespondWith(http.StatusOK, basicAuthBoshInfo),
				),
			)
		})

		AfterEach(func() {
			server.Close()
		})

		Context("when called with the --print-manifest option", func() {
			var c *cli.Context
			var err error
			var deploymentPostBody []byte
			var oldPrint = UIPrint

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
				}, pluginutil.ToCliFlagArray(GetAuthFlags()))

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
						ghttp.VerifyRequest("GET", "/info"),
						ghttp.RespondWith(http.StatusOK, basicAuthBoshInfo),
					),
					ghttp.CombineHandlers(
						ghttp.VerifyRequest("POST", "/deployments"),
						ghttp.RespondWithJSONEncoded(http.StatusOK, task),
						// capture the POST to /deployments so we can verify it later
						func(w http.ResponseWriter, req *http.Request) {
							body, _ := ioutil.ReadAll(req.Body)
							req.Body.Close()
							deploymentPostBody = body
						},
					),
					ghttp.CombineHandlers(
						ghttp.VerifyRequest("GET", "/tasks/42"),
						ghttp.RespondWithJSONEncoded(http.StatusOK, task),
					),
				)
				UIPrint = func(stuff ...interface{}) (int, error) {
					deploymentPostBody = []byte(stuff[0].(string))
					return 0, nil
				}
			})

			AfterEach(func() {
				UIPrint = oldPrint
			})

			It("decorates the deployment with the bosh UUID", func() {
				ProductAction(c, pd)

				Ω(deploymentPostBody).ShouldNot(BeNil())
				dm := enaml.NewDeploymentManifest(deploymentPostBody)
				Ω(dm.DirectorUUID).Should(Equal(controlUUID))
			})

			It("returns without error", func() {
				err = ProductAction(c, pd)
				Ω(err).ShouldNot(HaveOccurred())

				// it should GET the /info and /cloud_config, but not POST the product
				Ω(len(server.ReceivedRequests())).Should(Equal(3))
			})
		})

		Context("when called without the --print-manifest option", func() {
			var c *cli.Context
			var err error

			var deploymentPostBody []byte

			BeforeEach(func() {
				port, url := portAndURL(server)
				c = pluginutil.NewContext([]string{
					"foo",
					"--ssl-ignore",
					"--bosh-url", url,
					"--bosh-port", port,
					"--bosh-user", basicAuthUser,
					"--bosh-pass", basicAuthPass,
				}, pluginutil.ToCliFlagArray(GetAuthFlags()))

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
						ghttp.VerifyRequest("GET", "/info"),
						ghttp.RespondWith(http.StatusOK, basicAuthBoshInfo),
					),
					ghttp.CombineHandlers(
						ghttp.VerifyRequest("POST", "/deployments"),
						ghttp.RespondWithJSONEncoded(http.StatusOK, task),
						// capture the POST to /deployments so we can verify it later
						func(w http.ResponseWriter, req *http.Request) {
							body, _ := ioutil.ReadAll(req.Body)
							req.Body.Close()
							deploymentPostBody = body
						},
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

			It("decorates the deployment with the bosh UUID", func() {
				ProductAction(c, pd)

				Ω(deploymentPostBody).ShouldNot(BeNil())
				dm := enaml.NewDeploymentManifest(deploymentPostBody)
				Ω(dm.DirectorUUID).Should(Equal(controlUUID))
			})
		})
	})
})
