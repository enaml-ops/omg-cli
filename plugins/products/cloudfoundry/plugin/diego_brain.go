package cloudfoundry

import (
	"strings"

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
	caCert, err := pluginutil.LoadResourceFromContext(c, "bbs-ca-cert")
	if err != nil {
		lo.G.Panicf("bbs ca cert: %s\n", err.Error())
	}

	clientCert, err := pluginutil.LoadResourceFromContext(c, "bbs-client-cert")
	if err != nil {
		lo.G.Panicf("bbs client cert: %s\n", err.Error())
	}

	clientKey, err := pluginutil.LoadResourceFromContext(c, "bbs-client-key")
	if err != nil {
		lo.G.Panicf("bbs client key: %s\n", err.Error())
	}

	return &diegoBrain{
		AZs:                       c.StringSlice("az"),
		StemcellName:              c.String("stemcell-name"),
		VMTypeName:                c.String("diego-brain-vm-type"),
		PersistentDiskType:        c.String("diego-brain-disk-type"),
		NetworkName:               c.String("network"),
		NetworkIPs:                c.StringSlice("diego-brain-ip"),
		BBSCACert:                 caCert,
		BBSClientCert:             clientCert,
		BBSClientKey:              clientKey,
		BBSAPILocation:            c.String("bbs-api"),
		SkipSSLCertVerify:         c.BoolT("skip-cert-verify"),
		CCUploaderJobPollInterval: c.Int("cc-uploader-poll-interval"),
		SystemDomain:              c.String("system-domain"),
		CCInternalAPIUser:         c.String("cc-internal-api-user"),
		CCInternalAPIPassword:     c.String("cc-internal-api-password"),
		CCFetchTimeout:            c.Int("cc-fetch-timeout"),
		CCBulkBatchSize:           c.Int("cc-bulk-batch-size"),
		FSListenAddr:              c.String("fs-listen-addr"),
		FSStaticDirectory:         c.String("fs-static-dir"),
		FSDebugAddr:               c.String("fs-debug-addr"),
		FSLogLevel:                c.String("fs-log-level"),
		MetronPort:                c.Int("metron-port"),
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
		d.PersistentDiskType != "" &&
		d.NetworkName != "" &&
		d.BBSCACert != "" &&
		d.BBSClientCert != "" &&
		d.BBSClientKey != "" &&
		d.BBSAPILocation != "" &&
		d.CCInternalAPIUser != "" &&
		d.CCInternalAPIPassword != "" &&
		d.SystemDomain != ""
}

func (d *diegoBrain) newAuctioneer() *enaml.InstanceJob {
	return &enaml.InstanceJob{
		Name:    "auctioneer",
		Release: "diego",
		Properties: &auctioneer.Auctioneer{
			Bbs: &auctioneer.Bbs{
				ApiLocation: d.BBSAPILocation,
				CaCert:      d.BBSCACert,
				ClientCert:  d.BBSClientCert,
				ClientKey:   d.BBSClientKey,
			},
		},
	}
}

func (d *diegoBrain) newCCUploader() *enaml.InstanceJob {
	return &enaml.InstanceJob{
		Name:    "cc_uploader",
		Release: "diego",
		Properties: &cc_uploader.CcUploader{
			Diego: &cc_uploader.Diego{
				Ssl: &cc_uploader.Ssl{SkipCertVerify: d.SkipSSLCertVerify},
			},
			Cc: &cc_uploader.Cc{
				JobPollingIntervalInSeconds: d.CCUploaderJobPollInterval,
			},
		},
	}
}

func (d *diegoBrain) newConverger() *enaml.InstanceJob {
	return &enaml.InstanceJob{
		Name:    "converger",
		Release: "diego",
		Properties: &converger.Converger{
			Bbs: &converger.Bbs{
				ApiLocation: d.BBSAPILocation,
				CaCert:      d.BBSCACert,
				ClientCert:  d.BBSClientCert,
				ClientKey:   d.BBSClientKey,
			},
		},
	}
}

func (d *diegoBrain) newFileServer() *enaml.InstanceJob {
	return &enaml.InstanceJob{
		Name:    "file_server",
		Release: "diego",
		Properties: &file_server.FileServer{
			Diego: &file_server.Diego{
				Ssl: &file_server.Ssl{SkipCertVerify: d.SkipSSLCertVerify},
			},
			ListenAddr:      d.FSListenAddr,
			DebugAddr:       d.FSDebugAddr,
			LogLevel:        d.FSLogLevel,
			StaticDirectory: d.FSStaticDirectory,
			DropsondePort:   d.MetronPort,
		},
	}
}

func (d *diegoBrain) newNsync() *enaml.InstanceJob {
	// construct API URL from system domain,
	// stripping leading "https://" if necessary
	sys := d.SystemDomain
	if strings.HasPrefix(sys, "https://") {
		sys = sys[len("https://"):]
	}
	api := "https://api." + sys

	return &enaml.InstanceJob{
		Name:    "nsync",
		Release: "diego",
		Properties: &nsync.Nsync{
			Bbs: &nsync.Bbs{
				ApiLocation: d.BBSAPILocation,
				CaCert:      d.BBSCACert,
				ClientCert:  d.BBSClientCert,
				ClientKey:   d.BBSClientKey,
			},
			Cc: &nsync.Cc{
				BaseUrl:                  api,
				BasicAuthUsername:        d.CCInternalAPIUser,
				BasicAuthPassword:        d.CCInternalAPIPassword,
				BulkBatchSize:            d.CCBulkBatchSize,
				FetchTimeoutInSeconds:    d.CCFetchTimeout,
				PollingIntervalInSeconds: d.CCUploaderJobPollInterval,
			},
			Diego: &nsync.Diego{
				Ssl: &nsync.Ssl{SkipCertVerify: d.SkipSSLCertVerify},
			},
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
