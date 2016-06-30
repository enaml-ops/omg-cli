package cloudfoundry

import (
	"github.com/codegangsta/cli"
	"github.com/enaml-ops/enaml"
)

//NewDopplerPartition -
func NewDopplerPartition(c *cli.Context) InstanceGrouper {
	return &Doppler{
		AZs:                    c.StringSlice("az"),
		StemcellName:           c.String("stemcell-name"),
		NetworkIPs:             c.StringSlice("doppler-ip"),
		NetworkName:            c.String("network"),
		VMTypeName:             c.String("doppler-vm-type"),
		Metron:                 NewMetron(c),
		StatsdInjector:         NewStatsdInjector(c),
		Zone:                   c.String("doppler-zone"),
		StatusUser:             c.String("doppler-status-user"),
		StatusPassword:         c.String("doppler-status-password"),
		StatusPort:             c.Int("doppler-status-port"),
		MessageDrainBufferSize: c.Int("doppler-drain-buffer-size"),
		SharedSecret:           c.String("doppler-shared-secret"),
		SystemDomain:           c.String("system-domain"),
		CCBuilkAPIPassword:     c.String("cc-bulk-api-password"),
		SkipSSLCertify:         c.Bool("skip-cert-verify"),
	}
}

//ToInstanceGroup -
func (s *Doppler) ToInstanceGroup() (ig *enaml.InstanceGroup) {
	ig = &enaml.InstanceGroup{
		Name:      "doppler-partition",
		Instances: len(s.NetworkIPs),
		VMType:    s.VMTypeName,
		AZs:       s.AZs,
		Stemcell:  s.StemcellName,
		Jobs: []enaml.InstanceJob{
			s.Metron.CreateJob(),
			s.StatsdInjector.CreateJob(),
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

//HasValidValues - Check if the datastructure has valid fields
func (s *Doppler) HasValidValues() bool {
	return (len(s.AZs) > 0 &&
		s.StemcellName != "" &&
		s.VMTypeName != "" &&
		s.NetworkName != "" &&
		len(s.NetworkIPs) > 0 &&
		s.Zone != "" &&
		s.StatusUser != "" &&
		s.StatusPassword != "" &&
		s.StatusPort > 0 &&
		s.MessageDrainBufferSize > 0 &&
		s.SharedSecret != "" &&
		s.SystemDomain != "" &&
		s.CCBuilkAPIPassword != "")
}
