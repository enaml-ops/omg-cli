package cloudfoundry

import (
	"github.com/codegangsta/cli"
	"github.com/enaml-ops/enaml"
	"github.com/enaml-ops/omg-cli/plugins/products/cloudfoundry/enaml-gen/rep"
)

func NewDiegoCellPartition(c *cli.Context) InstanceGrouper {

	return &diegoCell{
		AZs:                c.StringSlice("az"),
		StemcellName:       c.String("stemcell-name"),
		VMTypeName:         c.String("diego-cell-vm-type"),
		PersistentDiskType: c.String("diego-cell-disk-type"),
		NetworkName:        c.String("network"),
		NetworkIPs:         c.StringSlice("diego-cell-ip"),
		ConsulAgent:        NewConsulAgentServer(c),
		Metron:             NewMetron(c),
		StatsdInjector:     NewStatsdInjector(c),
		DiegoBrain:         NewDiegoBrainPartition(c).(*diegoBrain),
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

	ig.AddJob(&enaml.InstanceJob{
		Name:       "rep",
		Release:    DiegoReleaseName,
		Properties: s.newRDiego(),
	})
	ig.AddJob(&enaml.InstanceJob{
		Name:       "consul_agent",
		Release:    CFReleaseName,
		Properties: s.ConsulAgent.CreateJob().Properties,
	})
	ig.AddJob(&enaml.InstanceJob{
		Name:    "cflinuxfs2-rootfs-setup",
		Release: CFLinuxFSReleaseName,
	})
	ig.AddJob(&enaml.InstanceJob{
		Name:       "garden",
		Release:    GardenReleaseName,
		Properties: nil,
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

func (s *diegoCell) HasValidValues() bool {
	return false
}

func (s *diegoCell) newRDiego() (rdiego *rep.Diego) {
	rdiego = &rep.Diego{
		Executor: &rep.Executor{
			PostSetupHook: `sh -c "rm -f /home/vcap/app/.java-buildpack.log /home/vcap/app/**/.java-buildpack.log"`,
			PostSetupUser: "root",
		},
		Rep: &rep.Rep{
			Bbs: &rep.Bbs{
				ApiLocation: s.DiegoBrain.BBSAPILocation,
				CaCert:      s.DiegoBrain.BBSCACert,
				ClientCert:  s.DiegoBrain.BBSClientCert,
				ClientKey:   s.DiegoBrain.BBSClientKey,
			},
			PreloadedRootfses: map[string]string{
				"cflinuxfs2": "/var/vcap/packages/cflinuxfs2/rootfs",
			},
			Zone: s.Metron.Zone,
		},
	}
	return
}
