package cloudfoundry

import (
	"fmt"

	"github.com/codegangsta/cli"
	"github.com/enaml-ops/enaml"
)

func NewGoRouterPartition(c *cli.Context) (grtr *gorouter, err error) {
	grtr = &gorouter{
		NetworkName: c.String("router-network"),
		NetworkIPs:  c.StringSlice("router-ip"),
	}
	if !grtr.hasValidValues() {
		err = fmt.Errorf("invalid values in GoRouter: %v", grtr)
		grtr = nil
	}
	return
}

func (s *gorouter) ToInstanceGroup() (ig *enaml.InstanceGroup) {
	ig = &enaml.InstanceGroup{
		Networks: []enaml.Network{
			enaml.Network{Name: s.NetworkName, StaticIPs: s.NetworkIPs},
		},
	}
	return
}

func (s *gorouter) hasValidValues() bool {
	return (len(s.AZs) > 0 &&
		s.StemcellName != "" &&
		s.VMTypeName != "" &&
		s.NetworkName != "" &&
		len(s.NetworkIPs) > 0 &&
		s.SSLCert != "" &&
		s.SSLKey != "" &&
		s.Nats.User != "" &&
		s.Nats.Password != "" &&
		s.Nats.Machines != nil &&
		s.Loggregator.Etcd.Machines != nil)
}
