package cloudfoundry

import (
	"fmt"

	"github.com/codegangsta/cli"
	"github.com/enaml-ops/enaml"
	"github.com/enaml-ops/omg-cli/pluginlib/util"
	"github.com/enaml-ops/omg-cli/plugins/products/cloudfoundry/enaml-gen/bbs"
	"github.com/enaml-ops/omg-cli/plugins/products/cloudfoundry/enaml-gen/etcd"
	"github.com/xchapter7x/lo"
)

const diegoDatabaseIGName = "diego_database-partition"

func NewDiegoDatabasePartition(c *cli.Context) InstanceGrouper {

	caCert, err := pluginutil.LoadResourceFromContext(c, "bbs-ca-cert")
	if err != nil {
		lo.G.Panicf("ca cert: %s\n", err.Error())
	}

	bbsServerCert, err := pluginutil.LoadResourceFromContext(c, "bbs-server-cert")
	if err != nil {
		lo.G.Panicf("bbs server cert: %s\n", err.Error())
	}

	bbsServerKey, err := pluginutil.LoadResourceFromContext(c, "bbs-server-key")
	if err != nil {
		lo.G.Panicf("bbs server key: %s\n", err.Error())
	}

	etcdServerCert, err := pluginutil.LoadResourceFromContext(c, "etcd-server-cert")
	if err != nil {
		lo.G.Panicf("etcd server cert: %s\n", err.Error())
	}

	etcdServerKey, err := pluginutil.LoadResourceFromContext(c, "etcd-server-key")
	if err != nil {
		lo.G.Panicf("etcd server key: %s\n", err.Error())
	}

	etcdClientCert, err := pluginutil.LoadResourceFromContext(c, "etcd-client-cert")
	if err != nil {
		lo.G.Panicf("etcd client cert: %s\n", err.Error())
	}

	etcdClientKey, err := pluginutil.LoadResourceFromContext(c, "etcd-client-key")
	if err != nil {
		lo.G.Panicf("etcd client key: %s\n", err.Error())
	}

	etcdPeerCert, err := pluginutil.LoadResourceFromContext(c, "etcd-peer-cert")
	if err != nil {
		lo.G.Panicf("etcd peer cert: %s\n", err.Error())
	}

	etcdPeerKey, err := pluginutil.LoadResourceFromContext(c, "etcd-peer-key")
	if err != nil {
		lo.G.Panicf("etcd peer key: %s\n", err.Error())
	}

	return &diegoDatabase{
		context:            c,
		AZs:                c.StringSlice("az"),
		CACert:             caCert,
		SystemDomain:       c.String("system-domain"),
		StemcellName:       c.String("stemcell-name"),
		VMTypeName:         c.String("diego-db-vm-type"),
		PersistentDiskType: c.String("diego-db-disk-type"),
		NetworkName:        c.String("network"),
		NetworkIPs:         c.StringSlice("diego-db-ip"),
		Passphrase:         c.String("diego-db-passphrase"),
		BBSServerCert:      bbsServerCert,
		BBSServerKey:       bbsServerKey,
		EtcdServerCert:     etcdServerCert,
		EtcdServerKey:      etcdServerKey,
		EtcdClientCert:     etcdClientCert,
		EtcdClientKey:      etcdClientKey,
		EtcdPeerCert:       etcdPeerCert,
		EtcdPeerKey:        etcdPeerKey,
		ConsulAgent:        NewConsulAgentServer(c),
		Metron:             NewMetron(c),
		StatsdInjector:     NewStatsdInjector(c),
		DiegoBrain:         NewDiegoBrainPartition(c).(*diegoBrain),
	}
}

func (s *diegoDatabase) ToInstanceGroup() (ig *enaml.InstanceGroup) {
	ig = &enaml.InstanceGroup{
		Name:               diegoDatabaseIGName,
		Lifecycle:          "service",
		Instances:          len(s.NetworkIPs),
		VMType:             s.VMTypeName,
		AZs:                s.AZs,
		PersistentDiskType: s.PersistentDiskType,
		Stemcell:           s.StemcellName,
		Networks: []enaml.Network{
			{Name: s.NetworkName, StaticIPs: s.NetworkIPs},
		},
		Update: enaml.Update{
			MaxInFlight: 1,
			Serial:      true,
		},
	}

	ig.AddJob(&enaml.InstanceJob{
		Name:       "etcd",
		Release:    EtcdReleaseName,
		Properties: s.newEtcd(),
	})
	ig.AddJob(&enaml.InstanceJob{
		Name:       "bbs",
		Release:    DiegoReleaseName,
		Properties: s.newBBS(),
	})
	ig.AddJob(&enaml.InstanceJob{
		Name:       "consul_agent",
		Release:    CFReleaseName,
		Properties: s.ConsulAgent.CreateJob().Properties,
	})
	ig.AddJob(&enaml.InstanceJob{
		Name:       "statsd-injector",
		Release:    CFReleaseName,
		Properties: s.StatsdInjector.CreateJob().Properties,
	})
	ig.AddJob(&enaml.InstanceJob{
		Name:       "metron_agent",
		Release:    CFReleaseName,
		Properties: s.Metron.CreateJob().Properties,
	})
	return
}

func (s *diegoDatabase) HasValidValues() bool {
	validStrings := hasValidStringFlags(s.context, []string{
		"bbs-ca-cert",
		"bbs-server-cert",
		"bbs-server-key",
		"etcd-server-cert",
		"etcd-server-key",
		"etcd-client-cert",
		"etcd-client-key",
		"etcd-peer-cert",
		"etcd-peer-key",
		"system-domain",
		"stemcell-name",
		"diego-db-vm-type",
		"diego-db-disk-type",
		"network",
		"diego-db-passphrase",
	})
	validSlices := hasValidStringSliceFlags(s.context, []string{"az"})
	return validStrings && validSlices
}

func (s *diegoDatabase) newBBS() (dbdiego *bbs.Diego) {
	var keyname = "key1"
	dbdiego = &bbs.Diego{
		Bbs: &bbs.Bbs{
			RequireSsl:     false,
			CaCert:         s.DiegoBrain.BBSCACert,
			ServerCert:     s.BBSServerCert,
			ServerKey:      s.BBSServerKey,
			ActiveKeyLabel: keyname,
			EncryptionKeys: map[string]string{
				"label":      keyname,
				"passphrase": s.Passphrase,
			},
			Auctioneer: &bbs.Auctioneer{
				ApiUrl: fmt.Sprintf("http://auctioneer.%s:9016", s.SystemDomain),
			},
			Etcd: s.newBBSEtcd(),
		},
	}
	return
}

func (s *diegoDatabase) newEtcd() (dbetcd *etcd.Etcd) {
	dbetcd = &etcd.Etcd{
		CaCert:                 s.CACert,
		ServerCert:             s.EtcdServerCert,
		ServerKey:              s.EtcdServerKey,
		ClientCert:             s.EtcdClientCert,
		ClientKey:              s.EtcdClientKey,
		PeerCaCert:             s.CACert,
		PeerCert:               s.EtcdPeerCert,
		PeerKey:                s.EtcdPeerKey,
		AdvertiseUrlsDnsSuffix: fmt.Sprintf("etcd.%s", s.SystemDomain),
		Machines: []string{
			fmt.Sprintf("etcd.%s", s.SystemDomain),
		},
		Cluster: map[string]interface{}{
			"name":      diegoDatabaseIGName,
			"instances": len(s.NetworkIPs),
		},
	}
	return
}

func (s *diegoDatabase) newBBSEtcd() (dbetcd *bbs.Etcd) {
	dbetcd = &bbs.Etcd{
		CaCert:     s.CACert,
		ClientCert: s.EtcdClientCert,
		ClientKey:  s.EtcdClientKey,
		Machines: []string{
			fmt.Sprintf("etcd.", s.SystemDomain),
		},
	}
	return
}
