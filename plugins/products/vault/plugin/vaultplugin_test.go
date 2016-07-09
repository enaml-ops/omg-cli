package vault_test

import (
	"fmt"
	"io/ioutil"

	"github.com/codegangsta/cli"
	"github.com/enaml-ops/enaml"
	. "github.com/enaml-ops/omg-cli/plugins/products/vault/plugin"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("given vault Plugin", func() {
	var plgn *Plugin

	BeforeEach(func() {
		plgn = new(Plugin)
	})

	Context("when calling GetProduct while targeting an un-compatible cloud config'd bosh", func() {
		var cloudConfigBytes []byte
		var controlNetName = "hello"
		var controlDisk = "medium"
		var controlVM = "large"
		var controlIP = "1.2.3.4"

		BeforeEach(func() {
			cloudConfigBytes, _ = ioutil.ReadFile("./fixtures/sample-aws.yml")
		})

		It("then we should fail fast and give the user guidance on what is wrong", func() {
			Ω(func() {
				plgn.GetProduct([]string{
					"appname",
					"--disk-type", controlDisk,
					"--network", controlNetName,
					"--vm-type", controlVM,
					"--ip", controlIP,
					"--az", "z1",
					"--stemcell-url", "something",
					"--stemcell-ver", "12.3.44",
					"--stemcell-sha", "ilkjag09dhsg90ahsd09gsadg9",
				}, cloudConfigBytes)
			}).Should(Panic())
		})
	})

	Context("when calling plugin without all required flags", func() {
		It("then it should fail fast and give the user guidance on what is wrong", func() {
			Ω(func() {
				plgn.GetProduct([]string{"appname"}, []byte(``))
			}).Should(Panic())
		})
	})

	Context("when calling GetProduct w/ valid flags and matching cloud config", func() {
		var deployment *enaml.DeploymentManifest
		var controlNetName = "private"
		var controlDisk = "medium"
		var controlVM = "medium"
		var controlIP = "1.2.3.4"

		BeforeEach(func() {
			cloudConfigBytes, _ := ioutil.ReadFile("./fixtures/sample-aws.yml")
			dmBytes := plgn.GetProduct([]string{
				"appname",
				"--network", controlNetName,
				"--vm-type", controlVM,
				"--disk-type", controlDisk,
				"--ip", controlIP,
				"--az", "z1",
				"--stemcell-url", "something",
				"--stemcell-ver", "12.3.44",
				"--stemcell-sha", "ilkjag09dhsg90ahsd09gsadg9",
			}, cloudConfigBytes)
			deployment = enaml.NewDeploymentManifest(dmBytes)
		})
		It("then we should have a properly initialized deployment set", func() {
			Ω(deployment.Update).ShouldNot(BeNil())
			Ω(len(deployment.Releases)).Should(Equal(2))
			Ω(len(deployment.Stemcells)).Should(Equal(1))
			Ω(len(deployment.InstanceGroups)).Should(Equal(1))
		})
	})

	Context("when calling the plugin", func() {
		var flags []cli.Flag

		BeforeEach(func() {
			flags = plgn.GetFlags()
		})
		It("then there should be valid flags available", func() {
			for _, flagname := range []string{
				"ip",
				"az",
				"network",
				"vm-type",
				"disk-type",
				"stemcell-url",
				"stemcell-ver",
				"stemcell-sha",
				"stemcell-name",
			} {
				Ω(checkFlags(flags, flagname)).ShouldNot(HaveOccurred())
			}
		})
	})
})

func checkFlags(flags []cli.Flag, flagName string) error {
	var err = fmt.Errorf("could not find an flag %s in plugin", flagName)
	for _, f := range flags {
		if f.GetName() == flagName {
			err = nil
		}
	}
	return err
}
