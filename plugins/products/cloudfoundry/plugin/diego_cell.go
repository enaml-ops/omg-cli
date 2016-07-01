package cloudfoundry

import (
	"github.com/codegangsta/cli"
	"github.com/enaml-ops/enaml"
)

func NewDiegoCellPartition(c *cli.Context) InstanceGrouper {
	return &diegoCell{
		AZs:                c.StringSlice("az"),
		StemcellName:       c.String("stemcell-name"),
		VMTypeName:         c.String("diego-cell-vm-type"),
		PersistentDiskType: c.String("diego-cell-disk-type"),
		NetworkName:        c.String("network"),
		NetworkIPs:         c.StringSlice("diego-cell-ip"),
	}
}

func (s *diegoCell) ToInstanceGroup() (ig *enaml.InstanceGroup) {
	ig = &enaml.InstanceGroup{
		Name:               "diego_cell-partition",
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
		},
	}
	ig.AddJob(&enaml.InstanceJob{Name: "rep", Release: DiegoReleaseName})
	ig.AddJob(&enaml.InstanceJob{Name: "consul_agent", Release: CFReleaseName})
	ig.AddJob(&enaml.InstanceJob{Name: "cflinuxfs2-rootfs-setup", Release: CFLinuxFSReleaseName})
	ig.AddJob(&enaml.InstanceJob{Name: "garden", Release: GardenReleaseName})
	ig.AddJob(&enaml.InstanceJob{Name: "statsd-injector", Release: CFReleaseName})
	ig.AddJob(&enaml.InstanceJob{Name: "metron_agent", Release: CFReleaseName})
	return
}

func (s *diegoCell) HasValidValues() bool {
	return false
}
