package vspherecli

import (
	"gopkg.in/urfave/cli.v2"

	"github.com/enaml-ops/omg-cli/plugins/products/bosh-init"
	"github.com/enaml-ops/omg-cli/utils"
	"github.com/enaml-ops/pluginlib/pcli"
	"github.com/xchapter7x/lo"
)

// GetFlags returns the available CLI flags
func GetFlags() []pcli.Flag {
	boshdefaults := boshinit.GetVSphereDefaults()

	boshFlags := boshinit.BoshFlags(boshdefaults)
	vsphereFlags := []pcli.Flag{
		// vsphere specific flags
		pcli.CreateStringFlag("vsphere-address", "IP of the vCenter"),
		pcli.CreateStringFlag("vsphere-user", "vSphere user"),
		pcli.CreateStringFlag("vsphere-password", "vSphere user's password"),
		pcli.CreateStringFlag("vsphere-datacenter-name", "name of the datacenter the Director will use for VM creation"),
		pcli.CreateStringFlag("vsphere-vm-folder", "name of the folder created to hold VMs"),
		pcli.CreateStringFlag("vsphere-template-folder", "the name of the folder created to hold stemcells"),
		pcli.CreateStringFlag("vsphere-datastore", "name of the datastore the Director will use for storing VMs"),
		pcli.CreateStringFlag("vsphere-disk-path", "name of the VMs folder, disk folder will be automatically created in the chosen datastore."),
		pcli.CreateStringSliceFlag("vsphere-clusters", "one or more vSphere datacenter cluster names"),
		pcli.CreateStringSliceFlag("vsphere-resource-pool", "Name of resource pool for vsphere cluster"),
		// vsphere subnet1 flags
		pcli.CreateStringFlag("vsphere-subnet1-name", "name of the vSphere network for subnet1"),
		pcli.CreateStringFlag("vsphere-subnet1-range", "CIDR range for subnet1"),
		pcli.CreateStringFlag("vsphere-subnet1-gateway", "IP of the default gateway for subnet1"),
		pcli.CreateStringSliceFlag("vsphere-subnet1-dns", "IP of the DNS server(s) for subnet1"),
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
			lo.G.Error(e.Error())
			return e
		}
		if err := utils.CheckRequired(c, "vsphere-address", "vsphere-user", "vsphere-password", "vsphere-datacenter-name",
			"vsphere-vm-folder", "vsphere-template-folder", "vsphere-datastore", "vsphere-disk-path",
			"vsphere-clusters", "vsphere-subnet1-name", "vsphere-subnet1-range", "vsphere-subnet1-range", "vsphere-subnet1-dns"); err != nil {
			lo.G.Error(err.Error())
			return err
		}
		vSphereClusters := c.StringSlice("vsphere-clusters")
		vSphereResourcePools := c.StringSlice("vsphere-resource-pool")
		vsphereConfig := boshinit.VSphereInitConfig{
			VSphereAddress:        c.String("vsphere-address"),
			VSphereUser:           c.String("vsphere-user"),
			VSpherePassword:       c.String("vsphere-password"),
			VSphereDatacenterName: c.String("vsphere-datacenter-name"),
			VSphereVMFolder:       c.String("vsphere-vm-folder"),
			VSphereTemplateFolder: c.String("vsphere-template-folder"),
			VSphereDataStore:      c.String("vsphere-datastore"),
			VSphereDiskPath:       c.String("vsphere-disk-path"),
			VSphereClusters:       vSphereClusters,
			VSphereNetworks: []boshinit.Network{boshinit.Network{
				Name:    c.String("vsphere-subnet1-name"),
				Range:   c.String("vsphere-subnet1-range"),
				Gateway: c.String("vsphere-subnet1-gateway"),
				DNS:     c.StringSlice("vsphere-subnet1-dns"),
			}},
		}

		if len(vSphereResourcePools) > 0 {
			vsphereConfig.VSphereResourcePool = vSphereResourcePools
		}
		provider := boshinit.NewVSphereIaaSProvider(vsphereConfig, boshBase)

		if err := boshBase.HandleDeployment(provider, boshInitDeploy); err != nil {
			return err
		}

		return nil
	}
}
