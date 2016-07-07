package cloudfoundry

import (
	"github.com/codegangsta/cli"
	"github.com/enaml-ops/enaml"
	bstraplib "github.com/enaml-ops/omg-cli/plugins/products/cf-mysql/enaml-gen/bootstrap"
)

func NewBootstrapPartition(c *cli.Context) InstanceGrouper {
	return &bootstrap{
		AZs:           c.StringSlice("az"),
		StemcellName:  c.String("stemcell-name"),
		NetworkName:   c.String("network"),
		MySQLIPs:      c.StringSlice("mysql-ip"),
		MySQLUser:     c.String("mysql-bootstrap-username"),
		MySQLPassword: c.String("mysql-bootstrap-password"),
	}
}

func (b *bootstrap) ToInstanceGroup() *enaml.InstanceGroup {
	return &enaml.InstanceGroup{
		Name:      "bootstrap",
		Instances: 1,
		VMType:    "errand",
		Lifecycle: "errand",
		AZs:       b.AZs,
		Stemcell:  b.StemcellName,
		Networks: []enaml.Network{
			{Name: b.NetworkName},
		},
		Update: enaml.Update{
			MaxInFlight: 1,
		},
		Jobs: []enaml.InstanceJob{
			{
				Name:    "bootstrap",
				Release: CFMysqlReleaseName,
				Properties: &bstraplib.Bootstrap{
					ClusterIps:             b.MySQLIPs,
					DatabaseStartupTimeout: 1200,
					BootstrapEndpoint: &bstraplib.BootstrapEndpoint{
						Username: b.MySQLUser,
						Password: b.MySQLPassword,
					},
				},
			},
		},
	}
}

func (b *bootstrap) HasValidValues() bool {
	return len(b.AZs) > 0 &&
		b.StemcellName != "" &&
		b.NetworkName != "" &&
		len(b.MySQLIPs) > 0 &&
		b.MySQLUser != "" &&
		b.MySQLPassword != ""
}
