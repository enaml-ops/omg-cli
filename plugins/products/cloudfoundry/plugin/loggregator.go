package cloudfoundry

import (
	"github.com/codegangsta/cli"
	"github.com/enaml-ops/enaml"
	ltc "github.com/enaml-ops/omg-cli/plugins/products/cloudfoundry/enaml-gen/loggregator_trafficcontroller"
	"github.com/enaml-ops/omg-cli/plugins/products/cloudfoundry/enaml-gen/route_registrar"
)

func NewLoggregatorTrafficController(c *cli.Context) InstanceGrouper {
	return &loggregatorTrafficController{
		AZs:               c.StringSlice("az"),
		StemcellName:      c.String("stemcell-name"),
		NetworkName:       c.String("network"),
		NetworkIPs:        c.StringSlice("loggregator-traffic-controller-ip"),
		VMTypeName:        c.String("loggregator-traffic-controller-vmtype"),
		SystemDomain:      c.String("system-domain"),
		DopplerSecret:     c.String("doppler-client-secret"),
		SkipSSLCertVerify: c.BoolT("skip-cert-verify"),
		EtcdMachines:      c.StringSlice("etcd-machine-ip"),
		Metron:            NewMetron(c),
		Nats: &route_registrar.Nats{
			User:     c.String("nats-user"),
			Password: c.String("nats-pass"),
			Machines: c.StringSlice("nats-machine-ip"),
			Port:     4222,
		},
	}
}

func (l *loggregatorTrafficController) HasValidValues() bool {
	return len(l.AZs) > 0 &&
		l.StemcellName != "" &&
		l.NetworkName != "" &&
		len(l.NetworkIPs) > 0 &&
		l.VMTypeName != "" &&
		l.SystemDomain != "" &&
		len(l.EtcdMachines) > 0 &&
		l.DopplerSecret != "" &&
		l.Metron.HasValidValues()
}

func (l *loggregatorTrafficController) ToInstanceGroup() *enaml.InstanceGroup {
	return &enaml.InstanceGroup{
		Name:      "loggregator_trafficcontroller-partition",
		AZs:       l.AZs,
		Stemcell:  l.StemcellName,
		VMType:    l.VMTypeName,
		Instances: len(l.NetworkIPs),

		Networks: []enaml.Network{
			{
				Name:      l.NetworkName,
				StaticIPs: l.NetworkIPs,
			},
		},
		Update: enaml.Update{
			MaxInFlight: 1,
		},
		Jobs: []enaml.InstanceJob{
			l.createLoggregatorTrafficControllerJob(),
			l.Metron.CreateJob(),
			l.createRouteRegistrarJob(),
			l.createStatsdInjectorJob(),
		},
	}
}

func (l *loggregatorTrafficController) createLoggregatorTrafficControllerJob() enaml.InstanceJob {
	return enaml.InstanceJob{
		Name:    "loggregator_trafficcontroller",
		Release: CFReleaseName,
		Properties: &ltc.LoggregatorTrafficcontroller{
			SystemDomain: l.SystemDomain,
			Cc: &ltc.Cc{
				SrvApiUri: prefixSystemDomain(l.SystemDomain, "api"),
			},
			Ssl: &ltc.Ssl{
				SkipCertVerify: l.SkipSSLCertVerify,
			},
			TrafficController: &ltc.TrafficController{
				Zone: l.Metron.Zone,
			},
			Doppler: &ltc.Doppler{
				UaaClientId: "doppler",
			},
			Uaa: &ltc.Uaa{
				Clients: &ltc.Clients{
					Doppler: &ltc.Doppler{
						Secret: l.DopplerSecret,
					},
				},
			},
			Loggregator: &ltc.Loggregator{
				Etcd: &ltc.Etcd{
					Machines: l.EtcdMachines,
				},
			},
		},
	}
}

func (l *loggregatorTrafficController) createRouteRegistrarJob() enaml.InstanceJob {
	routes := make([]map[string]interface{}, 2)

	routes[0] = map[string]interface{}{
		"name":                  "doppler",
		"port":                  8081,
		"registration_interval": "20s",
		"uris":                  []string{"doppler." + l.SystemDomain},
	}
	routes[1] = map[string]interface{}{
		"name":                  "loggregator",
		"port":                  8080,
		"registration_interval": "20s",
		"uris":                  []string{"loggregator." + l.SystemDomain},
	}
	return enaml.InstanceJob{
		Name:    "route_registrar",
		Release: CFReleaseName,
		Properties: &route_registrar.RouteRegistrar{
			Routes: routes,
			Nats:   l.Nats,
		},
	}
}

func (l *loggregatorTrafficController) createStatsdInjectorJob() enaml.InstanceJob {
	return enaml.InstanceJob{
		Name:       "statsd-injector",
		Release:    CFReleaseName,
		Properties: struct{}{},
	}
}
