package boshinit

import (
	"fmt"

	"github.com/enaml-ops/enaml"
	"github.com/enaml-ops/omg-cli/plugins/products/bosh-init/enaml-gen/vsphere_cpi"
)

const (
	vSphereCPIJobName     = "vsphere_cpi"
	vSphereCPIReleaseName = "bosh-vsphere-cpi"
)

type VSphereInitConfig struct {
	VSphereAddress        string
	VSphereUser           string
	VSpherePassword       string
	VSphereDatacenterName string
	VSphereVMFolder       string
	VSphereTemplateFolder string
	VSphereDataStore      string
	VSphereDiskPath       string
	VSphereResourcePool   []string
	VSphereClusters       []string
	VSphereNetworks       []Network
}

type Network struct {
	Name    string
	Range   string
	Gateway string
	DNS     []string
}

func GetVSphereDefaults() *BoshBase {
	return &BoshBase{
		BoshReleaseURL:     "https://bosh.io/d/github.com/cloudfoundry/bosh?v=256.2",
		BoshReleaseSHA:     "ff2f4e16e02f66b31c595196052a809100cfd5a8",
		CPIReleaseURL:      "https://bosh.io/d/github.com/cloudfoundry-incubator/bosh-vsphere-cpi-release?v=22",
		CPIReleaseSHA:      "dd1827e5f4dfc37656017c9f6e48441f51a7ab73",
		GOAgentReleaseURL:  "https://bosh.io/d/stemcells/bosh-vsphere-esxi-ubuntu-trusty-go_agent?v=3232.4",
		GOAgentSHA:         "27ec32ddbdea13e3025700206388ae5882a23c67",
		CPIJobName:         vSphereCPIJobName,
		PersistentDiskSize: 20000,
	}
}

type VSpherecloudpropertiesResourcePool struct {
	CPU         int                `yaml:"cpu,omitempty"`  // [Integer, required]: Number of CPUv.
	RAM         int                `yaml:"ram,omitempty"`  // [Integer, required]: Specified the amount of RAM in megabytev.
	Disk        int                `yaml:"disk,omitempty"` // [Integer, required]: Specifies the disk size in megabytev.
	Datacenters VSphereDatacenters `yaml:"datacenters,omitempty"`
}

type VSpherecloudpropertiesNetwork struct {
	Name string `yaml:"name,omitempty"` // [String, required]: vSphere network name.
}

type VSphereDatacenters []VSphereDatacenter

type VSphereDatacenter struct {
	Name                       string      `yaml:"name"`                         // [String, required]: vSphere datacenter name.
	VMFolder                   string      `yaml:"vm_folder"`                    // [String, required]: The folder to create PCF VMs in.
	TemplateFolder             string      `yaml:"template_folder"`              // [String, required]: The folder to store stemcells in.
	DatastorePattern           string      `yaml:"datastore_pattern"`            // [String, required]: The pattern to the vSphere datastore.
	PersistentDatastorePattern string      `yaml:"persistent_datastore_pattern"` // [String, required]: The pattern to the vSphere datastore for persistent diskv.
	DiskPath                   string      `yaml:"disk_path"`                    // [String, required]: The disk path.
	Clusters                   interface{} `yaml:"clusters"`                     // [[]String], required]: The vSphere cluster(s).
}

type ResourcePool struct {
	ResourcePool string `yaml:"resource_pool"`
}

type VSphereBosh struct {
	cfg      VSphereInitConfig
	boshbase *BoshBase
}

func NewVSphereIaaSProvider(cfg VSphereInitConfig, boshBase *BoshBase) IAASManifestProvider {
	boshBase.CPIJobName = vSphereCPIJobName
	return &VSphereBosh{
		cfg:      cfg,
		boshbase: boshBase,
	}
}

func (v *VSphereBosh) CreateCPIRelease() enaml.Release {
	return enaml.Release{
		Name: vSphereCPIReleaseName,
		URL:  v.boshbase.CPIReleaseURL,
		SHA1: v.boshbase.CPIReleaseSHA,
	}
}
func (v *VSphereBosh) CreateCPITemplate() enaml.Template {
	return enaml.Template{
		Name:    v.boshbase.CPIJobName,
		Release: vSphereCPIReleaseName,
	}
}
func (v *VSphereBosh) CreateDiskPool() enaml.DiskPool {
	return enaml.DiskPool{
		Name:     "disks",
		DiskSize: v.boshbase.PersistentDiskSize,
	}
}

func (v *VSphereBosh) resourcePoolCloudProperties() interface{} {
	return VSpherecloudpropertiesResourcePool{
		CPU:         2,
		Disk:        v.boshbase.PersistentDiskSize,
		RAM:         4096,
		Datacenters: v.getDataCenters(),
	}
}
func (v *VSphereBosh) CreateResourcePool() (*enaml.ResourcePool, error) {
	return v.boshbase.CreateResourcePool(v.resourcePoolCloudProperties)

}
func (v *VSphereBosh) CreateManualNetwork() enaml.ManualNetwork {
	net := enaml.NewManualNetwork("private")
	net.AddSubnet(enaml.Subnet{
		Range:   v.cfg.VSphereNetworks[0].Range,
		Gateway: v.cfg.VSphereNetworks[0].Gateway,
		DNS:     v.cfg.VSphereNetworks[0].DNS,
		CloudProperties: VSpherecloudpropertiesNetwork{
			Name: v.cfg.VSphereNetworks[0].Name,
		},
	})
	return net
}
func (v *VSphereBosh) CreateVIPNetwork() enaml.VIPNetwork {
	return enaml.NewVIPNetwork("public")
}
func (v *VSphereBosh) CreateJobNetwork() *enaml.Network {
	return nil
}
func (v *VSphereBosh) CreateCloudProvider() enaml.CloudProvider {
	return enaml.CloudProvider{
		Template: enaml.Template{
			Name:    v.boshbase.CPIJobName,
			Release: vSphereCPIReleaseName,
		},
		MBus: fmt.Sprintf("https://mbus:%s@%s:6868", v.boshbase.MBusPassword, v.boshbase.GetRoutableIP()),
		Properties: &vsphere_cpi.VsphereCpiJob{
			Vcenter: &vsphere_cpi.Vcenter{
				Address:     v.cfg.VSphereAddress,
				User:        v.cfg.VSphereUser,
				Password:    v.cfg.VSpherePassword,
				Datacenters: v.getDataCenters(),
			},
			Ntp: v.boshbase.NtpServers,
			Agent: &vsphere_cpi.Agent{
				Mbus: fmt.Sprintf("https://mbus:%s@0.0.0.0:6868", v.boshbase.MBusPassword),
			},
			Blobstore: &vsphere_cpi.Blobstore{
				Provider: "local",
				Path:     "/var/vcap/micro_bosh/data/cache",
			},
		},
	}
}
func (v *VSphereBosh) CreateCPIJobProperties() map[string]interface{} {
	return map[string]interface{}{
		"vcenter": &vsphere_cpi.Vcenter{
			Address:     v.cfg.VSphereAddress,
			User:        v.cfg.VSphereUser,
			Password:    v.cfg.VSpherePassword,
			Datacenters: v.getDataCenters(),
		},
		"agent": &vsphere_cpi.Agent{
			Mbus: fmt.Sprintf("nats://nats:%s@%s:4222", v.boshbase.NatsPassword, v.boshbase.PrivateIP),
		},
	}
}

func (v *VSphereBosh) CreateDeploymentManifest() (*enaml.DeploymentManifest, error) {
	manifest := v.boshbase.CreateDeploymentManifest()
	manifest.AddRelease(v.CreateCPIRelease())
	if rp, err := v.CreateResourcePool(); err != nil {
		return nil, err
	} else {
		manifest.AddResourcePool(*rp)
	}
	manifest.AddDiskPool(v.CreateDiskPool())
	manifest.AddNetwork(v.CreateManualNetwork())
	//manifest.AddNetwork(v.CreateVIPNetwork())
	boshJob := manifest.Jobs[0]
	boshJob.AddTemplate(v.CreateCPITemplate())
	n := v.CreateJobNetwork()
	if n != nil {
		boshJob.AddNetwork(*n)
	}
	for name, val := range v.CreateCPIJobProperties() {
		boshJob.AddProperty(name, val)
	}
	manifest.Jobs[0] = boshJob
	manifest.SetCloudProvider(v.CreateCloudProvider())
	return manifest, nil
}

func (v *VSphereBosh) getDataCenters() VSphereDatacenters {
	return VSphereDatacenters{VSphereDatacenter{
		Name:                       v.cfg.VSphereDatacenterName,
		VMFolder:                   v.cfg.VSphereVMFolder,
		TemplateFolder:             v.cfg.VSphereTemplateFolder,
		DatastorePattern:           v.getDataStorePattern(),
		PersistentDatastorePattern: v.getDataStorePattern(),
		DiskPath:                   v.cfg.VSphereDiskPath,
		Clusters:                   v.clusterConfig(),
	}}
}

func (v *VSphereBosh) getDataStorePattern() (pattern string) {
	return fmt.Sprintf("^(%s)$", v.cfg.VSphereDataStore)
}
func (v *VSphereBosh) clusterConfig() interface{} {

	clusters := make([]map[string]interface{}, 0)

	for index, clusterName := range v.cfg.VSphereClusters {
		cluster := make(map[string]interface{})

		if len(v.cfg.VSphereResourcePool) > index {
			cluster[clusterName] = &ResourcePool{
				ResourcePool: v.cfg.VSphereResourcePool[index],
			}
		} else {
			cluster[clusterName] = make(map[string]string, 0)
		}
		clusters = append(clusters, cluster)
	}
	return clusters

}
