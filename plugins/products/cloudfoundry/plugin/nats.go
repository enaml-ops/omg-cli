package cloudfoundry

import (
	"github.com/codegangsta/cli"
	"github.com/enaml-ops/enaml"
	natslib "github.com/enaml-ops/omg-cli/plugins/products/cloudfoundry/enaml-gen/nats"
)

//NewNatsPartition --
func NewNatsPartition(c *cli.Context) (igf InstanceGrouper) {
	igf = &NatsPartition{
		AZs:          c.StringSlice("az"),
		StemcellName: c.String("stemcell-name"),
		NetworkIPs:   c.StringSlice("nats-machine-ip"),
		NetworkName:  c.String("network"),
		VMTypeName:   c.String("nats-vm-type"),
		Metron:       NewMetron(c),
		Nats: natslib.Nats{
			User:     c.String("nats-user"),
			Password: c.String("nats-pass"),
			Machines: c.StringSlice("nats-machine-ip"),
			Port:     natsPort,
		},
		StatsdInjector: NewStatsdInjector(c),
	}
	return
}

//ToInstanceGroup --
func (s *NatsPartition) ToInstanceGroup() (ig *enaml.InstanceGroup) {
	ig = &enaml.InstanceGroup{
		Name:      "nats-partition",
		Instances: len(s.NetworkIPs),
		VMType:    s.VMTypeName,
		AZs:       s.AZs,
		Stemcell:  s.StemcellName,
		Jobs: []enaml.InstanceJob{
			s.newNatsJob(),
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

func (s *NatsPartition) newNatsJob() enaml.InstanceJob {
	return enaml.InstanceJob{
		Name:       "nats",
		Release:    "cf",
		Properties: s.Nats,
	}
}

//HasValidValues - Checks that fields in NatsPartition are valid
func (s *NatsPartition) HasValidValues() bool {
	return (len(s.AZs) > 0 &&
		s.StemcellName != "" &&
		s.VMTypeName != "" &&
		s.Metron.Zone != "" &&
		s.Metron.Secret != "" &&
		s.NetworkName != "" &&
		len(s.NetworkIPs) > 0 &&
		s.Nats.User != "" &&
		s.Nats.Password != "" &&
		s.Nats.Machines != nil)
}
