package cloudfoundry

import (
	"fmt"

	"github.com/codegangsta/cli"
	"github.com/enaml-ops/enaml"
	"github.com/enaml-ops/omg-cli/plugins/products/cloudfoundry/enaml-gen/doppler"
	"github.com/enaml-ops/omg-cli/plugins/products/cloudfoundry/enaml-gen/syslog_drain_binder"
)

//NewDopplerPartition -
func NewDopplerPartition(c *cli.Context) InstanceGrouper {
	skipSSL := true
	if c.IsSet("skip-cert-verify") {
		skipSSL = c.Bool("skip-cert-verify")
	}
	return &Doppler{
		AZs:            c.StringSlice("az"),
		StemcellName:   c.String("stemcell-name"),
		NetworkIPs:     c.StringSlice("doppler-ip"),
		NetworkName:    c.String("network"),
		VMTypeName:     c.String("doppler-vm-type"),
		Metron:         NewMetron(c),
		StatsdInjector: NewStatsdInjector(c),
		Zone:           c.String("doppler-zone"),
		MessageDrainBufferSize: c.Int("doppler-drain-buffer-size"),
		SharedSecret:           c.String("doppler-shared-secret"),
		SystemDomain:           c.String("system-domain"),
		CCBuilkAPIPassword:     c.String("cc-bulk-api-password"),
		SkipSSLCertify:         skipSSL,
		EtcdMachines:           c.StringSlice("etcd-machine-ip"),
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
			s.createDopplerJob(),
			s.Metron.CreateJob(),
			s.createSyslogDrainBinderJob(),
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

func (s *Doppler) createDopplerJob() enaml.InstanceJob {
	return enaml.InstanceJob{
		Name:    "doppler",
		Release: "cf",
		Properties: &doppler.Doppler{
			Zone: s.Zone,
			MessageDrainBufferSize: s.MessageDrainBufferSize,
			DopplerEndpoint: &doppler.DopplerEndpoint{
				SharedSecret: s.SharedSecret,
			},
			Loggregator: &doppler.Loggregator{
				Etcd: &doppler.Etcd{
					Machines: s.EtcdMachines,
				},
			},
		},
	}
}

func (s *Doppler) createSyslogDrainBinderJob() enaml.InstanceJob {
	return enaml.InstanceJob{
		Name:    "syslog_drain_binder",
		Release: "cf",
		Properties: &syslog_drain_binder.SyslogDrainBinder{
			Ssl: &syslog_drain_binder.Ssl{
				SkipCertVerify: s.SkipSSLCertify,
			},
			SystemDomain: s.SystemDomain,
			Cc: &syslog_drain_binder.Cc{
				BulkApiPassword: s.CCBuilkAPIPassword,
				SrvApiUri:       fmt.Sprintf("https://api.%s", s.SystemDomain),
			},
			Loggregator: &syslog_drain_binder.Loggregator{
				Etcd: &syslog_drain_binder.Etcd{
					Machines: s.EtcdMachines,
				},
			},
		},
	}
}

//HasValidValues - Check if the datastructure has valid fields
func (s *Doppler) HasValidValues() bool {
	return (len(s.AZs) > 0 &&
		s.StemcellName != "" &&
		s.VMTypeName != "" &&
		s.NetworkName != "" &&
		len(s.NetworkIPs) > 0 &&
		s.Zone != "" &&
		s.MessageDrainBufferSize > 0 &&
		s.SharedSecret != "" &&
		s.SystemDomain != "" &&
		s.CCBuilkAPIPassword != "" &&
		len(s.EtcdMachines) > 0)
}
