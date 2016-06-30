package cloudfoundry

import (
	"github.com/codegangsta/cli"
	"github.com/enaml-ops/enaml"
	etcdlib "github.com/enaml-ops/omg-cli/plugins/products/cloudfoundry/enaml-gen/etcd"
	etcdmetricslib "github.com/enaml-ops/omg-cli/plugins/products/cloudfoundry/enaml-gen/etcd_metrics_server"
)

//NewEtcdPartition -
func NewEtcdPartition(c *cli.Context) (igf InstanceGrouper) {
	igf = &Etcd{
		AZs:                c.StringSlice("az"),
		StemcellName:       c.String("stemcell-name"),
		NetworkIPs:         c.StringSlice("etcd-machine-ip"),
		NetworkName:        c.String("network"),
		VMTypeName:         c.String("etcd-vm-type"),
		PersistentDiskType: c.String("etcd-disk-type"),
		Metron:             NewMetron(c),
		StatsdInjector:     NewStatsdInjector(c),
		Nats: &etcdmetricslib.Nats{
			Username: c.String("nats-user"),
			Password: c.String("nats-pass"),
			Machines: c.StringSlice("nats-machine-ip"),
		},
	}
	return
}

//ToInstanceGroup -
func (s *Etcd) ToInstanceGroup() (ig *enaml.InstanceGroup) {
	ig = &enaml.InstanceGroup{
		Name:      "etcd_server-partition",
		Instances: len(s.NetworkIPs),
		VMType:    s.VMTypeName,
		AZs:       s.AZs,
		Stemcell:  s.StemcellName,
		Jobs: []enaml.InstanceJob{
			s.newEtcdJob(),
			s.newEtcdMetricsServerJob(),
			s.Metron.CreateJob(),
			s.StatsdInjector.CreateJob(),
		},
		Networks: []enaml.Network{
			enaml.Network{Name: s.NetworkName, StaticIPs: s.NetworkIPs},
		},
		PersistentDiskType: s.PersistentDiskType,
		Update: enaml.Update{
			MaxInFlight: 1,
			Serial:      false,
		},
	}
	return
}

func (s *Etcd) newEtcdJob() enaml.InstanceJob {
	return enaml.InstanceJob{
		Name:    "etcd",
		Release: "cf",
		Properties: &etcdlib.Etcd{
			PeerRequireSsl: false,
			RequireSsl:     false,
			Machines:       s.NetworkIPs,
		},
	}
}

func (s *Etcd) newEtcdMetricsServerJob() enaml.InstanceJob {
	return enaml.InstanceJob{
		Name:    "etcd_metrics_server",
		Release: "cf",
		Properties: &etcdmetricslib.EtcdMetricsServer{
			Nats: s.Nats,
		},
	}
}

//HasValidValues -
func (s *Etcd) HasValidValues() bool {
	return (len(s.AZs) > 0 &&
		s.StemcellName != "" &&
		s.VMTypeName != "" &&
		s.NetworkName != "" &&
		len(s.NetworkIPs) > 0 &&
		s.PersistentDiskType != "")
}
