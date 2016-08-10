package vspherecli

import (
	"fmt"

	"github.com/codegangsta/cli"
	"github.com/enaml-ops/enaml"
	"github.com/enaml-ops/omg-cli/plugins/products/bosh-init"
	"github.com/enaml-ops/omg-cli/utils"
)

// GetFlags returns the available CLI flags
func GetFlags() []cli.Flag {
	boshdefaults := boshinit.GetVSphereDefaults()

	boshFlags := boshinit.BoshFlags(boshdefaults)
	vsphereFlags := []cli.Flag{
		// vsphere specific flags
		cli.StringFlag{Name: "vsphere-address", Value: "", Usage: "IP of the vCenter"},
		cli.StringFlag{Name: "vsphere-user", Value: "", Usage: "vSphere user"},
		cli.StringFlag{Name: "vsphere-password", Value: "", Usage: "vSphere user's password"},
		cli.StringFlag{Name: "vsphere-datacenter-name", Value: "", Usage: "name of the datacenter the Director will use for VM creation"},
		cli.StringFlag{Name: "vsphere-vm-folder", Value: "", Usage: "name of the folder created to hold VMs"},
		cli.StringFlag{Name: "vsphere-template-folder", Value: "", Usage: "the name of the folder created to hold stemcells"},
		cli.StringFlag{Name: "vsphere-datastore", Value: "", Usage: "name of the datastore the Director will use for storing VMs"},
		cli.StringFlag{Name: "vsphere-disk-path", Value: "", Usage: "name of the VMs folder, disk folder will be automatically created in the chosen datastore."},
		cli.StringSliceFlag{Name: "vsphere-clusters", Value: &cli.StringSlice{""}, Usage: "one or more vSphere datacenter cluster names"},
		cli.StringFlag{Name: "vsphere-resource-pool", Value: "", Usage: "Name of resource pool for vsphere cluster"},
		// vsphere subnet1 flags
		cli.StringFlag{Name: "vsphere-subnet1-name", Usage: "name of the vSphere network for subnet1"},
		cli.StringFlag{Name: "vsphere-subnet1-range", Usage: "CIDR range for subnet1"},
		cli.StringFlag{Name: "vsphere-subnet1-gateway", Usage: "IP of the default gateway for subnet1"},
		cli.StringSliceFlag{Name: "vsphere-subnet1-dns", Usage: "IP of the DNS server(s) for subnet1"},
	}
	for _, flag := range vsphereFlags {
		boshFlags = append(boshFlags, flag)
	}
	return boshFlags
}

// GetAction returns a function action that can be registered with the CLI
func GetAction(boshInitDeploy func(string)) func(c *cli.Context) error {
	return func(c *cli.Context) (e error) {
		var boshBase *boshinit.BoshBase
		if boshBase, e = boshinit.NewBoshBase(c); e != nil {
			return
		}
		utils.CheckRequired(c, "vsphere-address", "vsphere-user", "vsphere-password", "vsphere-datacenter-name",
			"vsphere-vm-folder", "vsphere-template-folder", "vsphere-datastore", "vsphere-disk-path",
			"vsphere-clusters", "vsphere-resource-pool", "vsphere-subnet1-name", "vsphere-subnet1-range", "vsphere-subnet1-range", "vsphere-subnet1-dns")

		manifest := boshinit.NewVSphereBosh(boshinit.VSphereInitConfig{
			// vsphere specific
			VSphereAddress:        c.String("vsphere-address"),
			VSphereUser:           c.String("vsphere-user"),
			VSpherePassword:       c.String("vsphere-password"),
			VSphereDatacenterName: c.String("vsphere-datacenter-name"),
			VSphereVMFolder:       c.String("vsphere-vm-folder"),
			VSphereTemplateFolder: c.String("vsphere-template-folder"),
			VSphereDataStore:      c.String("vsphere-datastore"),
			VSphereDiskPath:       c.String("vsphere-disk-path"),
			VSphereClusters:       utils.ClearDefaultStringSliceValue(c.StringSlice("vsphere-clusters")...),
			VSphereResourcePool:   c.String("vsphere-resource-pool"),
			VSphereNetworks: []boshinit.Network{boshinit.Network{
				Name:    c.String("vsphere-subnet1-name"),
				Range:   c.String("vsphere-subnet1-range"),
				Gateway: c.String("vsphere-subnet1-gateway"),
				DNS:     utils.ClearDefaultStringSliceValue(c.StringSlice("vsphere-subnet1-dns")...),
			}},
		}, boshBase)

		if yamlString, err := enaml.Paint(manifest); err == nil {

			if c.Bool("print-manifest") {
				fmt.Println(yamlString)

			} else {
				utils.DeployYaml(yamlString, boshInitDeploy)
			}
		} else {
			e = err
		}
		return
	}
}
