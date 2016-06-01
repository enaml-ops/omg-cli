package cloudconfig

import (
	"github.com/enaml-ops/enaml"
	"github.com/xchapter7x/lo"
)

// NewVSphereCloudConfig creates a new ready to execute vSphere cloud config
func NewVSphereCloudConfig(cfg *VSphereCloudConfig) (vsphereCloudConfig *enaml.CloudConfigManifest) {
	if err := cfg.validate(); err != nil {
		lo.G.Error(err)
		return
	}
	manifest := &enaml.CloudConfigManifest{}
	addAZs(manifest, cfg)
	addVMTypes(manifest)
	addDisk(manifest)
	addNetworks(manifest, cfg)
	addCompilation(manifest, cfg)
	return
}

// azs:
// - name: az1
//   cloud_properties:
//     datacenters:
//       clusters:
//       - PCF_CLUSTER_2:
//         resource_pool: SERVICES_RP
func addAZs(manifest *enaml.CloudConfigManifest, cfg *VSphereCloudConfig) {
	for _, az := range cfg.AZs {
		newAZ := enaml.AZ{
			Name: az.Name,
			CloudProperties: []vspherecloudpropertiesDatacenter{
				vspherecloudpropertiesDatacenter{
					Clusters: []map[string]map[string]string{},
				},
			},
		}

		props := newAZ.CloudProperties.([]vspherecloudpropertiesDatacenter)
		dc := props[0]
		dc.Clusters = make([]map[string]map[string]string, 1)
		dc.Clusters[0] = make(map[string]map[string]string)
		dc.Clusters[0][az.Cluster.Name] = map[string]string{
			"resource_pool": az.Cluster.ResourcePool,
		}

		manifest.AddAZ(newAZ)
	}
}

func addVMTypes(manifest *enaml.CloudConfigManifest) {
	manifest.AddVMType(enaml.VMType{
		Name: "small",
		CloudProperties: vspherecloudpropertiesVMType{
			CPU:  1,
			RAM:  1024,
			Disk: 3240,
		},
	})
	manifest.AddVMType(enaml.VMType{
		Name: "medium",
		CloudProperties: vspherecloudpropertiesVMType{
			CPU:  2,
			RAM:  2048,
			Disk: 20000,
		},
	})
	manifest.AddVMType(enaml.VMType{
		Name: "large",
		CloudProperties: vspherecloudpropertiesVMType{
			CPU:  2,
			RAM:  4096,
			Disk: 50000,
		},
	})
}

func addDisk(manifest *enaml.CloudConfigManifest) {
	manifest.AddDiskType(enaml.DiskType{
		Name:     "small",
		DiskSize: 3240,
	})
	manifest.AddDiskType(enaml.DiskType{
		Name:     "medium",
		DiskSize: 20000,
	})
	manifest.AddDiskType(enaml.DiskType{
		Name:     "large",
		DiskSize: 50000,
	})
}

func addNetworks(manifest *enaml.CloudConfigManifest, cfg *VSphereCloudConfig) {
	for _, a := range cfg.AZs {
		net := enaml.NewManualNetwork("private")
		net.AddSubnet(enaml.Subnet{
			AZ:      a.Name,
			Range:   a.Network.Range,
			Gateway: a.Network.Gateway,
			DNS:     a.Network.DNS,
			CloudProperties: vspherecloudpropertiesNetwork{
				Name: a.Network.Name,
			},
		})
		manifest.AddNetwork(net)
	}
}

func addCompilation(manifest *enaml.CloudConfigManifest, cfg *VSphereCloudConfig) {
	az := cfg.AZs[0]
	manifest.SetCompilation(&enaml.Compilation{
		Workers:             5,
		ReuseCompilationVMs: true,
		AZ:                  az.Name,
		VMType:              "medium",
		Network:             "private",
	})
}

type vspherecloudpropertiesDatacenter struct {
	Name     string                         `yaml:"name,omitempty"`
	Clusters []map[string]map[string]string `yaml:"clusters"`
}

type vspherecloudpropertiesVMType struct {
	CPU  int
	RAM  int
	Disk int
}

type vspherecloudpropertiesNetwork struct {
	Name string `yaml:"name,omitempty"`
}
