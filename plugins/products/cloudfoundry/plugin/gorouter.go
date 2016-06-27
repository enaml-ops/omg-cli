package cloudfoundry

import (
	"encoding/json"
	"fmt"

	"github.com/codegangsta/cli"
	"github.com/enaml-ops/enaml"
	grtrlib "github.com/enaml-ops/omg-cli/plugins/products/cloudfoundry/enaml-gen/gorouter"
	"github.com/enaml-ops/omg-cli/plugins/products/cloudfoundry/enaml-gen/loggregator_trafficcontroller"
)

const natsPort = 4222

func NewGoRouterPartition(c *cli.Context) (grtr *gorouter, err error) {
	grtr = &gorouter{
		Instances:    len(c.StringSlice("router-ip")),
		AZs:          c.StringSlice("az"),
		EnableSSL:    c.Bool("router-enable-ssl"),
		StemcellName: c.String("stemcell-name"),
		NetworkIPs:   c.StringSlice("router-ip"),
		NetworkName:  c.String("router-network"),
		VMTypeName:   c.String("router-vm-type"),
		SSLCert:      c.String("router-ssl-cert"),
		SSLKey:       c.String("router-ssl-key"),
		Nats: grtrlib.Nats{
			User:     c.String("nats-user"),
			Password: c.String("nats-pass"),
			Machines: c.StringSlice("nats-machine-ip"),
		},
		Loggregator: loggregator_trafficcontroller.Loggregator{
			Etcd: &loggregator_trafficcontroller.Etcd{
				Machines: c.StringSlice("etcd-machine-ip"),
			},
		},
	}

	if !grtr.hasValidValues() {
		b, _ := json.Marshal(grtr)
		err = fmt.Errorf("invalid values in GoRouter: %v", string(b))
		grtr = nil
	}
	return
}

func (s *gorouter) newNats() *grtrlib.Nats {
	s.Nats.Port = natsPort
	return &s.Nats
}

func (s *gorouter) newRouter() *grtrlib.Router {
	return &grtrlib.Router{
		EnableSsl:     s.EnableSSL,
		SecureCookies: false,
		Status: &grtrlib.Status{
			User:     "router_status",
			Password: "router-status-pass-ouaoihsdgoihasdg",
		},
	}
}

func (s *gorouter) ToInstanceGroup() (ig *enaml.InstanceGroup) {
	ig = &enaml.InstanceGroup{
		Name:      "router-partition",
		Instances: s.Instances,
		VMType:    s.VMTypeName,
		AZs:       s.AZs,
		Stemcell:  s.StemcellName,
		Jobs: []enaml.InstanceJob{
			enaml.InstanceJob{
				Name:    "gorouter",
				Release: "cf",
				Properties: &grtrlib.Gorouter{
					RequestTimeoutInSeconds: 180,
					Nats:   s.newNats(),
					Router: s.newRouter(),
				},
			},
			enaml.InstanceJob{
				Name:    "metron_agent",
				Release: "cf",
			},
		},
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
