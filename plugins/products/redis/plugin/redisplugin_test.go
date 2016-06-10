package redis_test

import (
	"fmt"

	"github.com/codegangsta/cli"
	"github.com/enaml-ops/enaml"
	. "github.com/enaml-ops/omg-cli/plugins/products/redis/plugin"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("given redis Plugin", func() {
	var plgn *Plugin

	BeforeEach(func() {
		plgn = new(Plugin)
	})

	XContext("when passing a valid cloud config to GetProduct", func() {
		It("", func() {
			Ω(false).Should(BeTrue())
		})
	})

	Context("when calling GetProduct w/ valid flags", func() {
		var deployment *enaml.DeploymentManifest
		var controlInstances = "1"
		var controlNetName = "hello"
		var controlPass = "pss"
		var controlDisk = "4033"
		var controlVM = "large"
		var controlIP = "1.2.3.4"

		BeforeEach(func() {
			dmBytes := plgn.GetProduct([]string{
				"appname",
				"--disk-size", controlDisk,
				"--leader-instances", controlInstances,
				"--network-name", controlNetName,
				"--redis-pass", controlPass,
				"--vm-size", controlVM,
				"--leader-ip", controlIP,
				"--slave-ip", controlIP,
				"--stemcell-url", "something",
				"--stemcell-ver", "12.3.44",
				"--stemcell-sha", "ilkjag09dhsg90ahsd09gsadg9",
			}, []byte(``))
			deployment = enaml.NewDeploymentManifest(dmBytes)
		})
		It("then we should have a properly initialized deployment set", func() {
			Ω(deployment.Update).ShouldNot(BeNil())
			Ω(len(deployment.Releases)).Should(Equal(1))
			Ω(len(deployment.Stemcells)).Should(Equal(1))
			Ω(len(deployment.Jobs)).Should(Equal(4))
		})
	})

	Context("when calling the plugin", func() {
		var flags []cli.Flag

		BeforeEach(func() {
			flags = plgn.GetFlags()
		})
		It("then there should be valid flags available", func() {
			for _, flagname := range []string{
				"leader-ip",
				"leader-instances",
				"redis-pass",
				"pool-instances",
				"disk-size",
				"slave-instances",
				"slave-ip",
				"network-name",
				"vm-size",
				"errand-instances",
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
	var err = fmt.Errorf("could not find an ip flag %s in plugin", flagName)
	for _, f := range flags {
		if f.GetName() == flagName {
			err = nil
		}
	}
	return err
}
