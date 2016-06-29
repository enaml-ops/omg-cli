package cloudfoundry

import (
	"fmt"

	"gopkg.in/yaml.v2"

	"github.com/codegangsta/cli"
	"github.com/enaml-ops/enaml"
	mysqllib "github.com/enaml-ops/omg-cli/plugins/products/cf-mysql/enaml-gen/mysql"
)

//NewMySQLPartition -
func NewMySQLPartition(c *cli.Context) (igf InstanceGrouper, err error) {
	var seededDBs []MySQLSeededDatabase
	if seededDBs, err = MySQLParseSeededDBs(c); err != nil {
		return
	}
	igf = &MySQL{
		AZs:                    c.StringSlice("az"),
		StemcellName:           c.String("stemcell-name"),
		NetworkIPs:             c.StringSlice("mysql-ip"),
		NetworkName:            c.String("mysql-network"),
		VMTypeName:             c.String("mysql-vm-type"),
		PersistentDiskType:     c.String("mysql-disk-type"),
		AdminPassword:          c.String("mysql-admin-password"),
		BootstrapUsername:      c.String("mysql-bootstrap-username"),
		BootstrapPassword:      c.String("mysql-bootstrap-password"),
		DatabaseStartupTimeout: 1200,
		InnodbBufferPoolSize:   2147483648,
		MaxConnections:         1500,
		SyslogAddress:          c.String("syslog-address"),
		SyslogPort:             c.Int("syslog-port"),
		SyslogTransport:        c.String("syslog-transport"),
		MySQLSeededDatabases:   seededDBs,
	}

	if !igf.HasValidValues() {
		b, _ := yaml.Marshal(igf)
		err = fmt.Errorf("invalid values in MySQL: %v", string(b))
		igf = nil
	}
	return
}

//MySQLParseSeededDBs -
func MySQLParseSeededDBs(c *cli.Context) (dbs []MySQLSeededDatabase, err error) {
	return
}

//ToInstanceGroup -
func (s *MySQL) ToInstanceGroup() (ig *enaml.InstanceGroup) {
	ig = &enaml.InstanceGroup{
		Name:      "mysql-partition",
		Instances: len(s.NetworkIPs),
		VMType:    s.VMTypeName,
		AZs:       s.AZs,
		Stemcell:  s.StemcellName,
		Jobs: []enaml.InstanceJob{
			s.newMySQLJob(),
		},
		Networks: []enaml.Network{
			enaml.Network{Name: s.NetworkName, StaticIPs: s.NetworkIPs},
		},
		Update: enaml.Update{
			MaxInFlight: 1,
		},
	}
	return
}

func (s *MySQL) newMySQLJob() enaml.InstanceJob {
	return enaml.InstanceJob{
		Name:    "mysql",
		Release: "cf-mysql",
		Properties: &mysqllib.Mysql{
			AdminPassword:          s.AdminPassword,
			ClusterIps:             s.NetworkIPs,
			DatabaseStartupTimeout: s.DatabaseStartupTimeout,
			InnodbBufferPoolSize:   s.InnodbBufferPoolSize,
			MaxConnections:         s.MaxConnections,
			BootstrapEndpoint: &mysqllib.BootstrapEndpoint{
				Username: s.BootstrapUsername,
				Password: s.BootstrapPassword,
			},
			SeededDatabases: s.MySQLSeededDatabases,
			SyslogAggregator: &mysqllib.SyslogAggregator{
				Address:   s.SyslogAddress,
				Port:      s.SyslogPort,
				Transport: s.SyslogTransport,
			},
		},
	}
}

//HasValidValues -
func (s *MySQL) HasValidValues() bool {
	return (len(s.AZs) > 0 &&
		s.StemcellName != "" &&
		s.VMTypeName != "" &&
		s.NetworkName != "" &&
		len(s.NetworkIPs) > 0 &&
		s.PersistentDiskType != "" &&
		s.AdminPassword != "" &&
	s.BootstrapUsername != "" &&
	s.BootstrapPassword != "")
}
