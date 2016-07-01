package cloudfoundry

import (
	"github.com/codegangsta/cli"
	"github.com/enaml-ops/enaml"
)

//NewConsulPartition -
func NewConsulPartition(c *cli.Context) InstanceGrouper {
	return &Consul{
		AZs:            c.StringSlice("az"),
		StemcellName:   c.String("stemcell-name"),
		NetworkIPs:     c.StringSlice("consul-ip"),
		NetworkName:    c.String("network"),
		VMTypeName:     c.String("consul-vm-type"),
		ConsulAgent:    NewConsulAgent(c, true),
		Metron:         NewMetron(c),
		StatsdInjector: NewStatsdInjector(c),
	}
}

//ToInstanceGroup -
func (s *Consul) ToInstanceGroup() (ig *enaml.InstanceGroup) {
	ig = &enaml.InstanceGroup{
		Name:      "consul-partition",
		Instances: len(s.NetworkIPs),
		VMType:    s.VMTypeName,
		AZs:       s.AZs,
		Stemcell:  s.StemcellName,
		Jobs: []enaml.InstanceJob{
			s.ConsulAgent.CreateJob(),
			s.Metron.CreateJob(),
			s.StatsdInjector.CreateJob(),
		},
		Networks: []enaml.Network{
			enaml.Network{Name: s.NetworkName, StaticIPs: s.NetworkIPs},
		},
		Update: enaml.Update{
			MaxInFlight: 1,
			Serial:      true,
		},
	}
	return
}

//HasValidValues - Check if the datastructure has valid fields
func (s *Consul) HasValidValues() bool {
	return (len(s.AZs) > 0 &&
		s.StemcellName != "" &&
		s.VMTypeName != "" &&
		s.NetworkName != "" &&
		len(s.NetworkIPs) > 0 &&
		len(s.ConsulAgent.EncryptKeys) > 0 &&
		s.ConsulAgent.hasValidValues())
}
