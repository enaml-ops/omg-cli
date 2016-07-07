package cloudfoundry

import (
	"github.com/codegangsta/cli"
	"github.com/enaml-ops/enaml"
)

//NewHaProxyPartition -
func NewHaProxyPartition(c *cli.Context) InstanceGrouper {
	return &HAProxy{
		AZs:          c.StringSlice("az"),
		StemcellName: c.String("stemcell-name"),
		NetworkIPs:   c.StringSlice("haproxy-ip"),
		NetworkName:  c.String("network"),
		VMTypeName:   c.String("haproxy-vm-type"),
	}
}

//ToInstanceGroup -
func (s *HAProxy) ToInstanceGroup() (ig *enaml.InstanceGroup) {
	ig = &enaml.InstanceGroup{
		Name:      "ha_proxy-partition",
		Instances: len(s.NetworkIPs),
		VMType:    s.VMTypeName,
		AZs:       s.AZs,
		Stemcell:  s.StemcellName,
		Jobs: []enaml.InstanceJob{
			s.createHAProxyJob(),
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

func (s *HAProxy) createHAProxyJob() enaml.InstanceJob {
	return enaml.InstanceJob{
		Name:       "haproxy",
		Release:    "cf",
		Properties: "",
	}
}

//HasValidValues - Check if the datastructure has valid fields
func (s *HAProxy) HasValidValues() bool {
	return (len(s.AZs) > 0 &&
		s.StemcellName != "" &&
		s.VMTypeName != "" &&
		s.NetworkName != "" &&
		len(s.NetworkIPs) > 0)
}
