package vspherecli_test

import (
	vspherecli "github.com/enaml-ops/omg-cli/vsphere-cli"
	"github.com/enaml-ops/pluginlib/pluginutil"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gopkg.in/urfave/cli.v2"
)

var _ = Describe("given the photon cli", func() {
	Context("when called with a complete set of flags", func() {
		It("then it should NOT panic", func() {
			action := vspherecli.GetAction(func(s string) {})
			var ctx *cli.Context
			ctx = pluginutil.NewContext([]string{"someapp",
				"--vsphere-address", "some",
				"--vsphere-user", "stuff",
				"--vsphere-password", "to",
				"--vsphere-datacenter-name", "do",
				"--vsphere-vm-folder", "asdfadf",
				"--vsphere-template-folder", "asdfasdf",
				"--vsphere-datastore", "asdfasdf",
				"--vsphere-disk-path", "asdfasdf",
				"--vsphere-clusters", "asdfasdf",
				"--vsphere-resource-pool", "asdfasdf",
				"--vsphere-subnet1-name", "asdfasdf",
				"--vsphere-subnet1-range", "asdfasdf",
				"--vsphere-subnet1-gateway", "asdfasdf",
				"--vsphere-subnet1-dns", "asdfasdf",
				"--gateway", "10.0.0.254",
				"--cidr", "10.0.0.1/24",
			}, pluginutil.ToCliFlagArray(vspherecli.GetFlags()))
			err := action(ctx)
			Ω(err).ShouldNot(HaveOccurred())
		})
	})

	Context("when called with an incomplete set of flags", func() {
		It("then it should panic and exit", func() {
			action := vspherecli.GetAction(func(s string) {})
			var ctx *cli.Context
			ctx = pluginutil.NewContext([]string{"someapp"}, pluginutil.ToCliFlagArray(vspherecli.GetFlags()))
			err := action(ctx)
			Ω(err).Should(HaveOccurred())
		})
	})
	Context("when called unbalance clusters and resource pools", func() {
		It("then it should panic and exit", func() {
			action := vspherecli.GetAction(func(s string) {})
			var ctx *cli.Context
			ctx = pluginutil.NewContext([]string{"someapp",
				"--vsphere-address", "some",
				"--vsphere-user", "stuff",
				"--vsphere-password", "to",
				"--vsphere-datacenter-name", "do",
				"--vsphere-vm-folder", "asdfadf",
				"--vsphere-template-folder", "asdfasdf",
				"--vsphere-datastore", "asdfasdf",
				"--vsphere-disk-path", "asdfasdf",
				"--vsphere-clusters", "asdfasdf",
				"--vsphere-subnet1-name", "asdfasdf",
				"--vsphere-subnet1-range", "asdfasdf",
				"--vsphere-subnet1-gateway", "asdfasdf",
				"--vsphere-subnet1-dns", "asdfasdf",
				"--gateway", "10.0.0.254",
				"--cidr", "10.0.0.1/24",
			}, pluginutil.ToCliFlagArray(vspherecli.GetFlags()))
			err := action(ctx)
			Ω(err).Should(HaveOccurred())
		})
	})
})
