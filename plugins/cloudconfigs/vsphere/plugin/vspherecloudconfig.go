package vsphereccplugin

import (
	"github.com/codegangsta/cli"
	"github.com/enaml-ops/omg-cli/pluginlib/cloudconfig"
	"github.com/enaml-ops/omg-cli/pluginlib/util"
	cconfig "github.com/enaml-ops/omg-cli/plugins/cloudconfigs/vsphere/cloud-config"
	"github.com/enaml-ops/omg-cli/utils"
	"github.com/xchapter7x/lo"
)

type VSphereCloudConfig struct{}

func (s *VSphereCloudConfig) GetFlags() (flags []cli.Flag) {
	return []cli.Flag{
		// flags for availability zones
		cli.StringFlag{Name: "vsphere-az1-name", Value: "az1", Usage: "name of az1 availability zone"},
		cli.StringFlag{Name: "vsphere-az1-cluster", Value: "az1", Usage: "name of the vSphere cluster for az1"},
		cli.StringFlag{Name: "vsphere-az1-resource-pool", Value: "az1", Usage: "name of the vSphere resource pool for az1"},
		// flags shared with bosh-init for networking
		cli.StringFlag{Name: "vsphere-subnet1-name", Value: "", Usage: "name of the vSphere network for subnet1"},
		cli.StringFlag{Name: "vsphere-subnet1-range", Value: "10.0.0.0/24", Usage: "CIDR range for subnet1"},
		cli.StringFlag{Name: "vsphere-subnet1-gateway", Value: "10.0.0.1", Usage: "IP of the default gateway for subnet1"},
		cli.StringSliceFlag{Name: "vsphere-subnet1-dns", Value: &cli.StringSlice{"10.0.0.2"}, Usage: "IP of the DNS server(s) for subnet1"},
	}
}

func (s *VSphereCloudConfig) GetMeta() cloudconfig.Meta {
	return cloudconfig.Meta{
		Name: "vsphere",
	}
}

func (s *VSphereCloudConfig) GetCloudConfig(args []string) (b []byte) {
	var err error
	c := pluginutil.NewContext(args, s.GetFlags())

	cfg := &cconfig.VSphereCloudConfig{
		AZs: []cconfig.VSphereAZ{cconfig.VSphereAZ{
			Name: c.String("vsphere-az1-name"),
			Cluster: cconfig.VSphereCluster{
				Name:         c.String("vsphere-az1-cluster"),
				ResourcePool: c.String("vsphere-az1-resource-pool"),
			},
			Network: cconfig.VSphereNetwork{
				Name:    c.String("vsphere-subnet1-name"),
				Range:   c.String("vsphere-subnet1-range"),
				Gateway: c.String("vsphere-subnet1-gateway"),
				DNS:     utils.ClearDefaultStringSliceValue(c.StringSlice("vsphere-subnet1-dns")...),
			},
		}},
	}

	cloud := cconfig.NewVSphereCloudConfig(cfg)
	if b, err = cloud.Bytes(); err != nil {
		lo.G.Error("cloud bytes call yielded error: ", err)
	}
	return b
}
