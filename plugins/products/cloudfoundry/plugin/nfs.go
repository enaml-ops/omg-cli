package cloudfoundry

import (
	"github.com/codegangsta/cli"
	"github.com/enaml-ops/enaml"
	nfslib "github.com/enaml-ops/omg-cli/plugins/products/cloudfoundry/enaml-gen/debian_nfs_server"
)

//NewNFSPartition -
func NewNFSPartition(c *cli.Context) (igf InstanceGrouper) {
	igf = &NFS{
		AZs:                  c.StringSlice("az"),
		StemcellName:         c.String("stemcell-name"),
		NetworkIPs:           c.StringSlice("nfs-ip"),
		NetworkName:          c.String("nfs-network"),
		VMTypeName:           c.String("nfs-vm-type"),
		PersistentDiskType:   c.String("nfs-disk-type"),
		AllowFromNetworkCIDR: c.StringSlice("nfs-allow-from-network-cidr"),
		Metron:               NewMetron(c),
		StatsdInjector:       NewStatsdInjector(c),
	}
	return
}

//ToInstanceGroup -
func (s *NFS) ToInstanceGroup() (ig *enaml.InstanceGroup) {
	ig = &enaml.InstanceGroup{
		Name:      "nfs_server-partition",
		Instances: len(s.NetworkIPs),
		VMType:    s.VMTypeName,
		AZs:       s.AZs,
		Stemcell:  s.StemcellName,
		Jobs: []enaml.InstanceJob{
			s.newNFSJob(),
			s.Metron.CreateJob(),
			s.StatsdInjector.CreateJob(),
		},
		Networks: []enaml.Network{
			enaml.Network{Name: s.NetworkName, StaticIPs: s.NetworkIPs},
		},
		PersistentDiskType: s.PersistentDiskType,
		Update: enaml.Update{
			MaxInFlight: 1,
			Serial:      true,
		},
	}
	return
}

func (s *NFS) newNFSJob() enaml.InstanceJob {

	return enaml.InstanceJob{
		Name:    "debian_nfs_server",
		Release: "cf",
		Properties: &nfslib.NfsServer{
			AllowFromEntries: s.AllowFromNetworkCIDR,
		},
	}
}

//HasValidValues -
func (s *NFS) HasValidValues() bool {
	return (len(s.AZs) > 0 &&
		s.StemcellName != "" &&
		s.VMTypeName != "" &&
		s.NetworkName != "" &&
		len(s.NetworkIPs) > 0 &&
		s.PersistentDiskType != "" &&
		len(s.AllowFromNetworkCIDR) > 0)
}
