package utils_test

import (
	"fmt"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/codegangsta/cli"
	"github.com/enaml-ops/enaml"
	"github.com/enaml-ops/enaml/enamlbosh"
	"github.com/enaml-ops/enaml/enamlbosh/enamlboshfakes"
	. "github.com/enaml-ops/omg-cli/utils"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("utils", func() {

	Describe("given a GetCloudConfigCommands func", func() {
		Context("when called with a valid plugin dir", func() {
			var commands []cli.Command
			BeforeEach(func() {
				commands = GetCloudConfigCommands("../pluginlib/registry/fixtures/cloudconfig")
			})
			It("then it should return a set of commands for the plugins in the dir", func() {
				Ω(len(commands)).Should(Equal(1))
				Ω(commands[0].Name).Should(Equal("testplugin-linux"))
				Ω(commands[0].Action).ShouldNot(BeNil())
			})
		})
	})

	Describe("given ProcessRemoteBoshAssets", func() {

	})

	Describe("given ProcessRemoteStemcells", func() {
		var doer *enamlboshfakes.FakeHttpClientDoer

		BeforeEach(func() {
			doer = new(enamlboshfakes.FakeHttpClientDoer)
			body, _ := os.Open("fixtures/deployment_tasks.json")
			doer.DoReturns(&http.Response{
				Body: body, //will only support a single call
			}, nil)
		})

		Context("when called with a valid list of remote stemcells", func() {
			var err error
			var myStemcells = []enaml.Stemcell{
				enaml.Stemcell{URL: "someurl.com", SHA1: "lkasdgklhasdglakshdgasdg"},
			}

			BeforeEach(func() {
				err = ProcessRemoteStemcells(
					myStemcells,
					enamlbosh.NewClient("user", "pass", "bosh.com", 25555),
					doer,
				)
			})

			It("then it should upload the given stemcell to bosh", func() {
				Ω(doer.DoCallCount()).Should(Equal(len(myStemcells)))
				Ω(err).ShouldNot(HaveOccurred())
			})
		})

		Context("when called with a valid list of NON remote stemcells", func() {
			var err error
			var myStemcells = []enaml.Stemcell{
				enaml.Stemcell{Name: "hi", Version: "1.2"},
				enaml.Stemcell{Name: "hello", Version: "0.333.4"},
			}

			BeforeEach(func() {
				err = ProcessRemoteStemcells(
					myStemcells,
					enamlbosh.NewClient("user", "pass", "bosh.com", 25555),
					doer,
				)
			})

			It("then it should pass over all and exit successfully", func() {
				Ω(doer.DoCallCount()).Should(Equal(0))
				Ω(err).ShouldNot(HaveOccurred())
			})
		})

		Context("when called with a mixed list of NON remote stemcells and remote stemcells", func() {
			var err error
			var myStemcells = []enaml.Stemcell{
				enaml.Stemcell{Name: "hi", Version: "1.2"},
				enaml.Stemcell{Name: "hello", Version: "0.333.4"},
			}
			var remoteStemcell = enaml.Stemcell{URL: "boshstuff.com", SHA1: "kljaslhdg9ashdgklahsdgklasdg"}

			BeforeEach(func() {
				err = ProcessRemoteStemcells(
					append(myStemcells, remoteStemcell),
					enamlbosh.NewClient("user", "pass", "bosh.com", 25555),
					doer,
				)
			})
			It("then it only upload the remote stemcells and exit successfully", func() {
				Ω(doer.DoCallCount()).Should(Equal(1))
				Ω(err).ShouldNot(HaveOccurred())
			})
		})
	})

	Describe("given ProcessRemoteReleases", func() {
		var doer *enamlboshfakes.FakeHttpClientDoer

		BeforeEach(func() {
			doer = new(enamlboshfakes.FakeHttpClientDoer)
			body, _ := os.Open("fixtures/deployment_tasks.json")
			doer.DoReturns(&http.Response{
				Body: body, //will only support a single call
			}, nil)
		})

		Context("when called with a valid list of remote stemcells", func() {
			var err error
			var myReleases = []enaml.Release{
				enaml.Release{URL: "someurl.com", SHA1: "lkasdgklhasdglakshdgasdg"},
			}

			BeforeEach(func() {
				err = ProcessRemoteReleases(
					myReleases,
					enamlbosh.NewClient("user", "pass", "bosh.com", 25555),
					doer,
				)
			})

			It("then it should upload the given stemcell to bosh", func() {
				Ω(doer.DoCallCount()).Should(Equal(len(myReleases)))
				Ω(err).ShouldNot(HaveOccurred())
			})
		})

		Context("when called with a valid list of NON remote stemcells", func() {
			var err error
			var myReleases = []enaml.Release{
				enaml.Release{Name: "hi", Version: "1.2"},
				enaml.Release{Name: "hello", Version: "0.333.4"},
			}

			BeforeEach(func() {
				err = ProcessRemoteReleases(
					myReleases,
					enamlbosh.NewClient("user", "pass", "bosh.com", 25555),
					doer,
				)
			})

			It("then it should pass over all and exit successfully", func() {
				Ω(doer.DoCallCount()).Should(Equal(0))
				Ω(err).ShouldNot(HaveOccurred())
			})
		})

		Context("when called with a mixed list of NON remote stemcells and remote stemcells", func() {
			var err error
			var myReleases = []enaml.Release{
				enaml.Release{Name: "hi", Version: "1.2"},
				enaml.Release{Name: "hello", Version: "0.333.4"},
			}
			var remoteRelease = enaml.Release{URL: "boshstuff.com", SHA1: "kljaslhdg9ashdgklahsdgklasdg"}

			BeforeEach(func() {
				err = ProcessRemoteReleases(
					append(myReleases, remoteRelease),
					enamlbosh.NewClient("user", "pass", "bosh.com", 25555),
					doer,
				)
			})
			It("then it only upload the remote stemcells and exit successfully", func() {
				Ω(doer.DoCallCount()).Should(Equal(1))
				Ω(err).ShouldNot(HaveOccurred())
			})
		})
	})

	Describe("given a ProcessProductBytes function", func() {
		var (
			printManifest  = true
			deployManifest = false
		)

		Context("when called with valid arguments for deploying the manifest", func() {
			var err error
			var task []enamlbosh.BoshTask
			BeforeEach(func() {
				doer := new(enamlboshfakes.FakeHttpClientDoer)

				body, _ := os.Open("fixtures/deployment_tasks.json")
				doer.DoReturns(&http.Response{
					Body: body,
				}, nil)
				task, err = ProcessProductBytes(new(enaml.DeploymentManifest).Bytes(), deployManifest, "user", "pass", "https://192.168.1.1", 25555, doer)
			})
			It("Then it should deploy the given manifest bytes", func() {
				Ω(err).ShouldNot(HaveOccurred())
				Ω(task).ShouldNot(BeNil())
				Ω(len(task)).Should(Equal(1))
			})
		})

		Context("when called with valid arguments for printing", func() {
			var callCount *int64
			var err error
			BeforeEach(func() {
				var z int64 = 0
				callCount = &z
				UIPrint = func(a ...interface{}) (n int, err error) {
					atomic.AddInt64(callCount, 1)
					return
				}
				doer := new(enamlboshfakes.FakeHttpClientDoer)
				body, _ := os.Open("fixtures/deployment_tasks.json")
				doer.DoReturns(&http.Response{
					Body: body,
				}, nil)
				_, err = ProcessProductBytes(new(enaml.DeploymentManifest).Bytes(), printManifest, "", "", "", 25555, doer)
			})

			AfterEach(func() {
				UIPrint = fmt.Println
			})

			It("Then it should print the yaml of the manifest", func() {
				Ω(err).ShouldNot(HaveOccurred())
				Ω(*callCount).Should(BeNumerically(">", 0))
			})
		})
	})

	Describe("given a GetProductCommands func", func() {
		Context("when called with a valid plugin dir", func() {
			var commands []cli.Command
			BeforeEach(func() {
				commands = GetProductCommands("../pluginlib/registry/fixtures/product")
			})
			It("then it should return a set of commands for the plugins in the dir", func() {
				Ω(len(commands)).Should(Equal(1))
				Ω(commands[0].Name).Should(Equal("testproductplugin-linux"))
				Ω(commands[0].Action).ShouldNot(BeNil())
			})
		})
	})

	Describe("given ClearDefaultStringSliceValue", func() {
		Context("when called on a stringslice containing a default & added values", func() {
			It("then it should clear the default value from the list", func() {
				stringSlice := []string{"default", "useradded1", "useradded2"}
				clearSlice := ClearDefaultStringSliceValue(stringSlice...)
				Ω(clearSlice).Should(ConsistOf("useradded1", "useradded2"))
				Ω(clearSlice).ShouldNot(ContainElement("default"))
			})
		})

		Context("when called on a stringslice only containing a default value", func() {
			It("then it should simply pass through the default value", func() {
				stringSlice := []string{"default"}
				Ω(ClearDefaultStringSliceValue(stringSlice...)).Should(Equal(stringSlice))
			})
		})
	})
})
