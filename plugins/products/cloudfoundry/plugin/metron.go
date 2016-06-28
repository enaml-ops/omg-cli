package cloudfoundry

import (
	"fmt"

	"github.com/codegangsta/cli"
	"github.com/enaml-ops/enaml"
	"github.com/enaml-ops/omg-cli/plugins/products/cloudfoundry/enaml-gen/metron_agent"
	"gopkg.in/yaml.v2"
)

//NewMetron -
func NewMetron(c *cli.Context) (metron *Metron, err error) {
	metron = &Metron{
		Zone:            c.String("metron-zone"),
		Secret:          c.String("metron-secret"),
		SyslogAddress:   c.String("syslog-address"),
		SyslogPort:      c.Int("syslog-port"),
		SyslogTransport: c.String("syslog-transport"),
		Loggregator: metron_agent.Loggregator{
			Etcd: &metron_agent.Etcd{
				Machines: c.StringSlice("etcd-machine-ip"),
			},
		},
	}
	if metron.SyslogTransport == "" {
		metron.SyslogTransport = "tcp"
	}

	if metron.Zone == "" || metron.Secret == "" {
		b, _ := yaml.Marshal(metron)
		err = fmt.Errorf("invalid values in Metron: %v", string(b))
		metron = nil
	}
	return
}

//CreateMetronJob -
func (s *Metron) CreateMetronJob() enaml.InstanceJob {
	return enaml.InstanceJob{
		Name:    "metron_agent",
		Release: "cf",
		Properties: &metron_agent.MetronAgent{
			SyslogDaemonConfig: &metron_agent.SyslogDaemonConfig{
				Transport: s.SyslogTransport,
				Address:   s.SyslogAddress,
				Port:      s.SyslogPort,
			},
			MetronAgent: &metron_agent.MetronAgent{
				Zone:       s.Zone,
				Deployment: DeploymentName,
			},
			MetronEndpoint: &metron_agent.MetronEndpoint{
				SharedSecret: s.Secret,
			},
			Loggregator: &s.Loggregator,
		},
	}
}
