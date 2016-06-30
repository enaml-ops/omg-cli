package cloudfoundry

import (
	"io/ioutil"

	"github.com/codegangsta/cli"
	"github.com/enaml-ops/enaml"
	grtrlib "github.com/enaml-ops/omg-cli/plugins/products/cloudfoundry/enaml-gen/gorouter"
	"github.com/enaml-ops/omg-cli/plugins/products/cloudfoundry/enaml-gen/metron_agent"
	"github.com/xchapter7x/lo"
)

const natsPort = 4222

func loadSSLFromContext(c *cli.Context, strFlag, fileFlag string) string {
	flag := c.String(strFlag)
	if file := c.String(fileFlag); file != "" {
		b, err := ioutil.ReadFile(file)
		if err != nil {
			lo.G.Panicf("error reading SSL cert/key file (%s): %s", fileFlag, err.Error())
		}
		flag = string(b)
	}
	return flag
}

//NewGoRouterPartition -
func NewGoRouterPartition(c *cli.Context) InstanceGrouper {
	return &gorouter{
		Instances:    len(c.StringSlice("router-ip")),
		AZs:          c.StringSlice("az"),
		EnableSSL:    c.Bool("router-enable-ssl"),
		StemcellName: c.String("stemcell-name"),
		NetworkIPs:   c.StringSlice("router-ip"),
		NetworkName:  c.String("network"),
		VMTypeName:   c.String("router-vm-type"),
		SSLCert:      loadSSLFromContext(c, "router-ssl-cert", "router-ssl-cert-file"),
		SSLKey:       loadSSLFromContext(c, "router-ssl-key", "router-ssl-key-file"),
		RouterUser:   c.String("router-user"),
		RouterPass:   c.String("router-pass"),
		MetronZone:   c.String("metron-zone"),
		MetronSecret: c.String("metron-secret"),
		Nats: grtrlib.Nats{
			User:     c.String("nats-user"),
			Password: c.String("nats-pass"),
			Machines: c.StringSlice("nats-machine-ip"),
		},
		Loggregator: metron_agent.Loggregator{
			Etcd: &metron_agent.Etcd{
				Machines: c.StringSlice("etcd-machine-ip"),
			},
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
			s.newRouterJob(),
			s.newMetronJob(),
			s.newStatsdInjectorJob(),
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

func (s *gorouter) newNats() *grtrlib.Nats {
	s.Nats.Port = natsPort
	return &s.Nats
}

func (s *gorouter) newRouter() *grtrlib.Router {
	return &grtrlib.Router{
		EnableSsl:     s.EnableSSL,
		SecureCookies: false,
		SslKey:        s.SSLKey,
		SslCert:       s.SSLCert,
		Status: &grtrlib.Status{
			User:     s.RouterUser,
			Password: s.RouterPass,
		},
	}
}

func (s *gorouter) newStatsdInjectorJob() enaml.InstanceJob {
	return enaml.InstanceJob{
		Name:       "statsd-injector",
		Release:    "cf",
		Properties: make(map[interface{}]interface{}),
	}
}

func (s *gorouter) newRouterJob() enaml.InstanceJob {
	return enaml.InstanceJob{
		Name:    "gorouter",
		Release: "cf",
		Properties: &grtrlib.Gorouter{
			RequestTimeoutInSeconds: 180,
			Nats:   s.newNats(),
			Router: s.newRouter(),
		},
	}
}

func (s *gorouter) newMetronJob() enaml.InstanceJob {
	return enaml.InstanceJob{
		Name:    "metron_agent",
		Release: "cf",
		Properties: &metron_agent.MetronAgent{
			SyslogDaemonConfig: &metron_agent.SyslogDaemonConfig{
				Transport: "tcp",
			},
			MetronAgent: &metron_agent.MetronAgent{
				Zone:       s.MetronZone,
				Deployment: DeploymentName,
			},
			MetronEndpoint: &metron_agent.MetronEndpoint{
				SharedSecret: s.MetronSecret,
			},
			Loggregator: &s.Loggregator,
		},
	}
}

func (s *gorouter) HasValidValues() bool {
	return (len(s.AZs) > 0 &&
		s.StemcellName != "" &&
		s.VMTypeName != "" &&
		s.MetronZone != "" &&
		s.MetronSecret != "" &&
		s.RouterPass != "" &&
		s.RouterUser != "" &&
		s.NetworkName != "" &&
		len(s.NetworkIPs) > 0 &&
		s.SSLCert != "" &&
		s.SSLKey != "" &&
		s.Nats.User != "" &&
		s.Nats.Password != "" &&
		s.Nats.Machines != nil &&
		s.Loggregator.Etcd.Machines != nil)
}
