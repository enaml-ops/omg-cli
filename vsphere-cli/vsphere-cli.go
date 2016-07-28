package vspherecli

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/codegangsta/cli"
	"github.com/enaml-ops/enaml"
	"github.com/enaml-ops/omg-cli/plugins/products/bosh-init"
	"github.com/enaml-ops/omg-cli/utils"
)

func deployYaml(myYaml string, boshInitDeploy func(string)) {
	fmt.Println("deploying your bosh")
	content := []byte(myYaml)
	tmpfile, err := ioutil.TempFile("", "bosh-init-deployment")
	defer os.Remove(tmpfile.Name())

	if err != nil {
		log.Fatal(err)
	}

	if _, err := tmpfile.Write(content); err != nil {
		log.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		log.Fatal(err)
	}
	boshInitDeploy(tmpfile.Name())
}

func checkRequired(name string, c *cli.Context) {
	if c.String(name) == "" {
		fmt.Println("Sorry you need to provide " + name)
		os.Exit(1)
	}
}

// GetFlags returns the available CLI flags
func GetFlags() []cli.Flag {
	boshdefaults := boshinit.BoshDefaults{
		CIDR:               "10.0.0.0/24",
		Gateway:            "10.0.0.1",
		DNS:                &cli.StringSlice{"10.0.0.2"},
		BoshReleaseVersion: "256.2",
		BoshReleaseSHA:     "ff2f4e16e02f66b31c595196052a809100cfd5a8",
		CPIReleaseVersion:  "22",
		CPIReleaseSHA:      "dd1827e5f4dfc37656017c9f6e48441f51a7ab73",
		GOAgentVersion:     "3232.4",
		GOAgentSHA:         "27ec32ddbdea13e3025700206388ae5882a23c67",
		PrivateIP:          "10.0.0.6",
		CPIName:            "vsphere_cpi",
	}

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
		cli.StringFlag{Name: "vsphere-subnet1-name", Value: "", Usage: "name of the vSphere network for subnet1"},
		cli.StringFlag{Name: "vsphere-subnet1-range", Value: "10.0.0.0/24", Usage: "CIDR range for subnet1"},
		cli.StringFlag{Name: "vsphere-subnet1-gateway", Value: "10.0.0.1", Usage: "IP of the default gateway for subnet1"},
		cli.StringSliceFlag{Name: "vsphere-subnet1-dns", Value: &cli.StringSlice{"10.0.0.2"}, Usage: "IP of the DNS server(s) for subnet1"},
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
		checkRequired("vsphere-address", c)
		checkRequired("vsphere-user", c)
		checkRequired("vsphere-password", c)
		checkRequired("vsphere-datacenter-name", c)
		checkRequired("vsphere-vm-folder", c)
		checkRequired("vsphere-template-folder", c)
		checkRequired("vsphere-datastore", c)
		checkRequired("vsphere-disk-path", c)
		checkRequired("vsphere-clusters", c)
		checkRequired("vsphere-resource-pool", c)
		checkRequired("vsphere-subnet1-name", c)

		manifest := boshinit.NewVSphereBosh(boshinit.BoshInitConfig{
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
				deployYaml(yamlString, boshInitDeploy)
			}
		} else {
			e = err
		}
		return
	}
}
