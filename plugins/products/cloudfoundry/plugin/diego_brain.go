package cloudfoundry

import (
	"github.com/codegangsta/cli"
	"github.com/enaml-ops/enaml"
	"github.com/enaml-ops/omg-cli/pluginlib/util"
	"github.com/enaml-ops/omg-cli/plugins/products/cloudfoundry/enaml-gen/auctioneer"
	"github.com/enaml-ops/omg-cli/plugins/products/cloudfoundry/enaml-gen/cc_uploader"
	"github.com/enaml-ops/omg-cli/plugins/products/cloudfoundry/enaml-gen/consul_agent"
	"github.com/enaml-ops/omg-cli/plugins/products/cloudfoundry/enaml-gen/converger"
	"github.com/enaml-ops/omg-cli/plugins/products/cloudfoundry/enaml-gen/file_server"
	"github.com/enaml-ops/omg-cli/plugins/products/cloudfoundry/enaml-gen/metron_agent"
	"github.com/enaml-ops/omg-cli/plugins/products/cloudfoundry/enaml-gen/nsync"
	"github.com/enaml-ops/omg-cli/plugins/products/cloudfoundry/enaml-gen/route_emitter"
	"github.com/enaml-ops/omg-cli/plugins/products/cloudfoundry/enaml-gen/ssh_proxy"
	"github.com/enaml-ops/omg-cli/plugins/products/cloudfoundry/enaml-gen/stager"
	"github.com/enaml-ops/omg-cli/plugins/products/cloudfoundry/enaml-gen/statsd-injector"
	"github.com/enaml-ops/omg-cli/plugins/products/cloudfoundry/enaml-gen/tps"
	"github.com/xchapter7x/lo"
)

func NewDiegoBrainPartition(c *cli.Context) InstanceGrouper {
	caCert, err := pluginutil.LoadResourceFromContext(c, "auctioneer-ca-cert")
	if err != nil {
		lo.G.Panicf("auctioneer ca cert: %s\n", err.Error())
	}

	clientCert, err := pluginutil.LoadResourceFromContext(c, "auctioneer-client-cert")
	if err != nil {
		lo.G.Panicf("auctioneer client cert: %s\n", err.Error())
	}

	clientKey, err := pluginutil.LoadResourceFromContext(c, "auctioneer-client-key")
	if err != nil {
		lo.G.Panicf("auctioneer client key: %s\n", err.Error())
	}

	return &diegoBrain{
		AZs:                  c.StringSlice("az"),
		StemcellName:         c.String("stemcell-name"),
		VMTypeName:           c.String("diego-brain-vm-type"),
		PersistentDiskType:   c.String("diego-brain-disk-type"),
		NetworkName:          c.String("network"),
		NetworkIPs:           c.StringSlice("diego-brain-ip"),
		AuctioneerCACert:     caCert,
		AuctioneerClientCert: clientCert,
		AuctioneerClientKey:  clientKey,
		BBSAPILocation:       c.String("bbs-api"),
	}
}

func (d *diegoBrain) ToInstanceGroup() *enaml.InstanceGroup {
	ig := &enaml.InstanceGroup{
		Name:               "diego_brain-partition",
		Instances:          len(d.NetworkIPs),
		VMType:             d.VMTypeName,
		AZs:                d.AZs,
		PersistentDiskType: d.PersistentDiskType,
		Stemcell:           d.StemcellName,
		Networks: []enaml.Network{
			{Name: d.NetworkName, StaticIPs: d.NetworkIPs},
		},
		Update: enaml.Update{
			MaxInFlight: 1,
		},
	}
	ig.AddJob(d.newAuctioneer())
	ig.AddJob(d.newCCUploader())
	ig.AddJob(d.newConverger())
	ig.AddJob(d.newFileServer())
	ig.AddJob(d.newNsync())
	ig.AddJob(d.newRouteEmitter())
	ig.AddJob(d.newSSHProxy())
	ig.AddJob(d.newStager())
	ig.AddJob(d.newTPS())
	ig.AddJob(d.newConsulAgent())
	ig.AddJob(d.newMetronAgent())
	ig.AddJob(d.newStatsdInjector())
	return ig
}

func (d *diegoBrain) HasValidValues() bool {
	return len(d.AZs) > 0 &&
		d.StemcellName != "" &&
		len(d.NetworkIPs) > 0 &&
		d.VMTypeName != "" &&
		d.NetworkName != ""
}

func (d *diegoBrain) newAuctioneer() *enaml.InstanceJob {
	return &enaml.InstanceJob{
		Name:    "auctioneer",
		Release: "diego",
		Properties: &auctioneer.Auctioneer{
			Bbs: &auctioneer.Bbs{
				ApiLocation: d.BBSAPILocation,
				CaCert:      d.AuctioneerCACert,
				ClientCert:  d.AuctioneerClientCert,
				ClientKey:   d.AuctioneerClientKey,
			},
		},
	}
}

func (d *diegoBrain) newCCUploader() *enaml.InstanceJob {
	return &enaml.InstanceJob{
		Name:       "cc_uploader",
		Release:    "diego",
		Properties: &cc_uploader.CcUploader{
		// TODO
		},
	}
}

func (d *diegoBrain) newConverger() *enaml.InstanceJob {
	return &enaml.InstanceJob{
		Name:       "converger",
		Release:    "diego",
		Properties: &converger.Converger{
		// TODO
		},
	}
}

func (d *diegoBrain) newFileServer() *enaml.InstanceJob {
	return &enaml.InstanceJob{
		Name:       "file_server",
		Release:    "diego",
		Properties: &file_server.FileServer{
		// TODO
		},
	}
}

func (d *diegoBrain) newNsync() *enaml.InstanceJob {
	return &enaml.InstanceJob{
		Name:       "nsync",
		Release:    "diego",
		Properties: &nsync.Nsync{
		// TODO
		},
	}
}

func (d *diegoBrain) newRouteEmitter() *enaml.InstanceJob {
	return &enaml.InstanceJob{
		Name:       "route_emitter",
		Release:    "diego",
		Properties: &route_emitter.RouteEmitter{
		// TODO
		},
	}
}

func (d *diegoBrain) newSSHProxy() *enaml.InstanceJob {
	return &enaml.InstanceJob{
		Name:       "ssh_proxy",
		Release:    "diego",
		Properties: &ssh_proxy.SshProxy{
		// TODO
		},
	}
}

func (d *diegoBrain) newStager() *enaml.InstanceJob {
	return &enaml.InstanceJob{
		Name:       "stager",
		Release:    "diego",
		Properties: &stager.Stager{
		// TODO
		},
	}
}

func (d *diegoBrain) newTPS() *enaml.InstanceJob {
	return &enaml.InstanceJob{
		Name:       "tps",
		Release:    "diego",
		Properties: &tps.Tps{
		// TODO
		},
	}
}

// TODO(zbergquist) reuse cloudfoundry.NewConsulAgent() ??
func (d *diegoBrain) newConsulAgent() *enaml.InstanceJob {
	return &enaml.InstanceJob{
		Name:       "consul_agent",
		Release:    "diego",
		Properties: &consul_agent.ConsulAgent{
		// TODO
		},
	}
}

func (d *diegoBrain) newMetronAgent() *enaml.InstanceJob {
	return &enaml.InstanceJob{
		Name:       "metron_agent",
		Release:    "diego",
		Properties: &metron_agent.MetronAgent{
		// TODO
		},
	}
}

func (d *diegoBrain) newStatsdInjector() *enaml.InstanceJob {
	return &enaml.InstanceJob{
		Name:       "statsd-injector",
		Release:    "diego",
		Properties: &statsd_injector.StatsdInjector{
		// TODO
		},
	}
}
