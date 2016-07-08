package cloudfoundry

import (
	"fmt"
	"strings"

	"github.com/codegangsta/cli"
	"github.com/enaml-ops/enaml"
	"github.com/enaml-ops/omg-cli/pluginlib/util"
	"github.com/enaml-ops/omg-cli/plugins/products/cloudfoundry/enaml-gen/auctioneer"
	"github.com/enaml-ops/omg-cli/plugins/products/cloudfoundry/enaml-gen/cc_uploader"
	"github.com/enaml-ops/omg-cli/plugins/products/cloudfoundry/enaml-gen/converger"
	"github.com/enaml-ops/omg-cli/plugins/products/cloudfoundry/enaml-gen/file_server"
	"github.com/enaml-ops/omg-cli/plugins/products/cloudfoundry/enaml-gen/nsync"
	"github.com/enaml-ops/omg-cli/plugins/products/cloudfoundry/enaml-gen/route_emitter"
	"github.com/enaml-ops/omg-cli/plugins/products/cloudfoundry/enaml-gen/ssh_proxy"
	"github.com/enaml-ops/omg-cli/plugins/products/cloudfoundry/enaml-gen/stager"
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
		BBSRequireSSL:             c.BoolT("bbs-require-ssl"),
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
		NATSUser:                  c.String("nats-user"),
		NATSPassword:              c.String("nats-pass"),
		NATSPort:                  c.Int("nats-port"),
		NATSMachines:              c.StringSlice("nats-machine-ip"),
		AllowSSHAccess:            c.Bool("allow-app-ssh-access"),
		SSHProxyClientSecret:      c.String("ssh-proxy-uaa-secret"),
		CCExternalPort:            c.Int("cc-external-port"),
		TrafficControllerURL:      c.String("traffic-controller-url"),
		ConsulAgent:               NewConsulAgent(c, []string{}),
		Metron:                    NewMetron(c),
		Statsd:                    NewStatsdInjector(c),
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
	consulJob := d.ConsulAgent.CreateJob()
	metronJob := d.Metron.CreateJob()
	statsdJob := d.Statsd.CreateJob()

	ig.AddJob(d.newAuctioneer())
	ig.AddJob(d.newCCUploader())
	ig.AddJob(d.newConverger())
	ig.AddJob(d.newFileServer())
	ig.AddJob(d.newNsync())
	ig.AddJob(d.newRouteEmitter())
	ig.AddJob(d.newSSHProxy())
	ig.AddJob(d.newStager())
	ig.AddJob(d.newTPS())
	ig.AddJob(&consulJob)
	ig.AddJob(&metronJob)
	ig.AddJob(&statsdJob)
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
		d.SystemDomain != "" &&
		d.ConsulAgent.HasValidValues() &&
		d.Metron.HasValidValues()
}

func (d *diegoBrain) newAuctioneer() *enaml.InstanceJob {
	return &enaml.InstanceJob{
		Name:    "auctioneer",
		Release: DiegoReleaseName,
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
		Release: DiegoReleaseName,
		Properties: &cc_uploader.CcUploaderJob{
			Diego: &cc_uploader.Diego{
				Ssl: &cc_uploader.Ssl{SkipCertVerify: d.SkipSSLCertVerify},
				CcUploader: &cc_uploader.CcUploader{
					Cc: &cc_uploader.Cc{
						JobPollingIntervalInSeconds: d.CCUploaderJobPollInterval,
					},
				},
			},
		},
	}
}

func (d *diegoBrain) newConverger() *enaml.InstanceJob {
	return &enaml.InstanceJob{
		Name:    "converger",
		Release: DiegoReleaseName,
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
		Release: DiegoReleaseName,
		Properties: &file_server.FileServerJob{
			Diego: &file_server.Diego{
				Ssl: &file_server.Ssl{SkipCertVerify: d.SkipSSLCertVerify},
				FileServer: &file_server.FileServer{
					ListenAddr:      d.FSListenAddr,
					DebugAddr:       d.FSDebugAddr,
					LogLevel:        d.FSLogLevel,
					StaticDirectory: d.FSStaticDirectory,
					DropsondePort:   d.MetronPort,
				},
			},
		},
	}
}

func (d *diegoBrain) newNsync() *enaml.InstanceJob {
	return &enaml.InstanceJob{
		Name:    "nsync",
		Release: DiegoReleaseName,
		Properties: &nsync.NsyncJob{
			Diego: &nsync.Diego{
				Ssl: &nsync.Ssl{SkipCertVerify: d.SkipSSLCertVerify},
				Nsync: &nsync.Nsync{
					Cc: &nsync.Cc{
						BaseUrl:                  prefixSystemDomain(d.SystemDomain, "api"),
						BasicAuthUsername:        d.CCInternalAPIUser,
						BasicAuthPassword:        d.CCInternalAPIPassword,
						BulkBatchSize:            d.CCBulkBatchSize,
						FetchTimeoutInSeconds:    d.CCFetchTimeout,
						PollingIntervalInSeconds: d.CCUploaderJobPollInterval,
					},
					Bbs: &nsync.Bbs{
						ApiLocation: d.BBSAPILocation,
						CaCert:      d.BBSCACert,
						ClientCert:  d.BBSClientCert,
						ClientKey:   d.BBSClientKey,
					},
				},
			},
		},
	}
}

func (d *diegoBrain) newRouteEmitter() *enaml.InstanceJob {
	return &enaml.InstanceJob{
		Name:    "route_emitter",
		Release: DiegoReleaseName,
		Properties: &route_emitter.RouteEmitterJob{
			Diego: &route_emitter.Diego{
				RouteEmitter: &route_emitter.RouteEmitter{
					Bbs: &route_emitter.Bbs{
						ApiLocation: d.BBSAPILocation,
						CaCert:      d.BBSCACert,
						ClientCert:  d.BBSClientCert,
						ClientKey:   d.BBSClientKey,
						RequireSsl:  d.BBSRequireSSL,
					},
					Nats: &route_emitter.Nats{
						User:     d.NATSUser,
						Password: d.NATSPassword,
						Port:     d.NATSPort,
						Machines: d.NATSMachines,
					},
				},
			},
		},
	}
}

func (d *diegoBrain) newSSHProxy() *enaml.InstanceJob {
	return &enaml.InstanceJob{
		Name:    "ssh_proxy",
		Release: DiegoReleaseName,
		Properties: &ssh_proxy.SshProxyJob{

			Diego: &ssh_proxy.Diego{
				Ssl: &ssh_proxy.Ssl{SkipCertVerify: d.SkipSSLCertVerify},
				SshProxy: &ssh_proxy.SshProxy{
					Bbs: &ssh_proxy.Bbs{
						ApiLocation: d.BBSAPILocation,
						CaCert:      d.BBSCACert,
						ClientCert:  d.BBSClientCert,
						ClientKey:   d.BBSClientKey,
						RequireSsl:  d.BBSRequireSSL,
					},
					Cc: &ssh_proxy.Cc{
						ExternalPort: d.CCExternalPort,
					},
					EnableCfAuth:    d.AllowSSHAccess,
					EnableDiegoAuth: d.AllowSSHAccess,
					UaaSecret:       d.SSHProxyClientSecret,
					UaaTokenUrl:     prefixSystemDomain(d.SystemDomain, "uaa") + "/oauth/token",
				},
			},
		},
	}
}

func (d *diegoBrain) newStager() *enaml.InstanceJob {
	return &enaml.InstanceJob{
		Name:    "stager",
		Release: DiegoReleaseName,
		Properties: &stager.StagerJob{
			Diego: &stager.Diego{
				Ssl: &stager.Ssl{SkipCertVerify: d.SkipSSLCertVerify},
				Stager: &stager.Stager{
					Bbs: &stager.Bbs{
						ApiLocation: d.BBSAPILocation,
						CaCert:      d.BBSCACert,
						ClientCert:  d.BBSClientCert,
						ClientKey:   d.BBSClientKey,
						RequireSsl:  d.BBSRequireSSL,
					},
					Cc: &stager.Cc{
						BasicAuthUsername: d.CCInternalAPIUser,
						BasicAuthPassword: d.CCInternalAPIPassword,
						ExternalPort:      d.CCExternalPort,
					},
				},
			},
		},
	}
}

func (d *diegoBrain) newTPS() *enaml.InstanceJob {
	return &enaml.InstanceJob{
		Name:    "tps",
		Release: DiegoReleaseName,
		Properties: &tps.TpsJob{

			Diego: &tps.Diego{
				Ssl: &tps.Ssl{SkipCertVerify: d.SkipSSLCertVerify},
				Tps: &tps.Tps{
					TrafficControllerUrl: d.TrafficControllerURL,
					Bbs: &tps.Bbs{
						ApiLocation: d.BBSAPILocation,
						CaCert:      d.BBSCACert,
						ClientCert:  d.BBSClientCert,
						ClientKey:   d.BBSClientKey,
						RequireSsl:  d.BBSRequireSSL,
					},
					Cc: &tps.Cc{
						BasicAuthUsername: d.CCInternalAPIUser,
						BasicAuthPassword: d.CCInternalAPIPassword,
						ExternalPort:      d.CCExternalPort,
					},
				},
			},
		},
	}
}

// prefixSystemDomain adds a prefix to the system domain.
// For example:
//     prefixSystemDomain("https://sys.yourdomain.com", "uaa")
// would return 'https://uaa.sys.yourdomain.com'.
func prefixSystemDomain(domain, prefix string) string {
	d := domain
	// strip leading https:// if necessary
	if strings.HasPrefix(d, "https://") {
		d = d[len("https://"):]
	}
	return fmt.Sprintf("https://%s.%s", prefix, d)
}
