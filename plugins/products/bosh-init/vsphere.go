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
	VSphereResourcePool   string
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

// NewVSphereBosh creates a new enaml deployment manifest for vSphere
func NewVSphereBosh(cfg VSphereInitConfig, boshbase *BoshBase) *enaml.DeploymentManifest {
	boshbase.CPIJobName = vSphereCPIJobName
	manifest := boshbase.CreateDeploymentManifest()

	var vcenterProperty = vsphere_cpi.Vcenter{
		Address:     cfg.VSphereAddress,
		User:        cfg.VSphereUser,
		Password:    cfg.VSpherePassword,
		Datacenters: getDataCenters(cfg),
	}

	var agentProperty = vsphere_cpi.Agent{
		Mbus: fmt.Sprintf("nats://nats:%s@%s:4222", boshbase.NatsPassword, boshbase.PrivateIP),
	}

	manifest.AddRelease(enaml.Release{
		Name: vSphereCPIReleaseName,
		URL:  boshbase.CPIReleaseURL,
		SHA1: boshbase.CPIReleaseSHA,
	})

	resourcePool := enaml.ResourcePool{
		Name:    "vms",
		Network: "private",
	}
	resourcePool.Stemcell = enaml.Stemcell{
		URL:  boshbase.GOAgentReleaseURL,
		SHA1: boshbase.GOAgentSHA,
	}
	resourcePool.CloudProperties = VSpherecloudpropertiesResourcePool{
		CPU:         2,
		Disk:        boshbase.PersistentDiskSize,
		RAM:         4096,
		Datacenters: getDataCenters(cfg),
	}
	// c1oudc0w is a default password for vcap user
	/*resourcePool.Env = map[string]interface{}{
		"bosh": map[string]string{
			"password": "$6$4gDD3aV0rdqlrKC$2axHCxGKIObs6tAmMTqYCspcdvQXh3JJcvWOY2WGb4SrdXtnCyNaWlrf3WEqvYR2MYizEGp3kMmbpwBC6jsHt0",
		},
	}*/
	manifest.AddResourcePool(resourcePool)
	manifest.AddDiskPool(enaml.DiskPool{
		Name:     "disks",
		DiskSize: 20000,
	})

	net := enaml.NewManualNetwork("private")
	net.AddSubnet(enaml.Subnet{
		Range:   cfg.VSphereNetworks[0].Range,
		Gateway: cfg.VSphereNetworks[0].Gateway,
		DNS:     cfg.VSphereNetworks[0].DNS,
		CloudProperties: VSpherecloudpropertiesNetwork{
			Name: cfg.VSphereNetworks[0].Name,
		},
	})
	manifest.AddNetwork(net)

	boshJob := manifest.Jobs[0]
	boshJob.AddTemplate(enaml.Template{Name: boshbase.CPIJobName, Release: vSphereCPIReleaseName})
	boshJob.AddProperty("agent", agentProperty)
	boshJob.AddProperty("vcenter", vcenterProperty)
	manifest.Jobs[0] = boshJob
	manifest.SetCloudProvider(createCloudProvider(cfg, boshbase))
	return manifest
}

type VSpherecloudpropertiesResourcePool struct {
	CPU         int                `yaml:"cpu,omitempty"`  // [Integer, required]: Number of CPUs.
	RAM         int                `yaml:"ram,omitempty"`  // [Integer, required]: Specified the amount of RAM in megabytes.
	Disk        int                `yaml:"disk,omitempty"` // [Integer, required]: Specifies the disk size in megabytes.
	Datacenters VSphereDatacenters `yaml:"datacenters,omitempty"`
}

type VSpherecloudpropertiesNetwork struct {
	Name string `yaml:"name,omitempty"` // [String, required]: vSphere network name.
}

type VSphereDatacenters []VSphereDatacenter

type VSphereDatacenter struct {
	Name                       string                    `yaml:"name"`                         // [String, required]: vSphere datacenter name.
	VMFolder                   string                    `yaml:"vm_folder"`                    // [String, required]: The folder to create PCF VMs in.
	TemplateFolder             string                    `yaml:"template_folder"`              // [String, required]: The folder to store stemcells in.
	DatastorePattern           string                    `yaml:"datastore_pattern"`            // [String, required]: The pattern to the vSphere datastore.
	PersistentDatastorePattern string                    `yaml:"persistent_datastore_pattern"` // [String, required]: The pattern to the vSphere datastore for persistent disks.
	DiskPath                   string                    `yaml:"disk_path"`                    // [String, required]: The disk path.
	Clusters                   []map[string]ResourcePool `yaml:"clusters"`                     // [[]String], required]: The vSphere cluster(s).
}

type ResourcePool struct {
	ResourcePool string `yaml:"resource_pool"`
}

func clusterConfig(cfg VSphereInitConfig) (clusters []map[string]ResourcePool) {
	clusters = make([]map[string]ResourcePool, 0)
	for _, clusterName := range cfg.VSphereClusters {
		cluster := make(map[string]ResourcePool)
		cluster[clusterName] = ResourcePool{
			ResourcePool: cfg.VSphereResourcePool,
		}
		clusters = append(clusters, cluster)

	}
	return
}

func createCloudProvider(cfg VSphereInitConfig, boshbase *BoshBase) (provider enaml.CloudProvider) {

	return enaml.CloudProvider{
		Template: enaml.Template{
			Name:    boshbase.CPIJobName,
			Release: vSphereCPIReleaseName,
		},
		MBus: fmt.Sprintf("https://mbus:%s@%s:6868", boshbase.MBusPassword, boshbase.GetRoutableIP()),
		Properties: &vsphere_cpi.VsphereCpiJob{
			Vcenter: &vsphere_cpi.Vcenter{
				Address:     cfg.VSphereAddress,
				User:        cfg.VSphereUser,
				Password:    cfg.VSpherePassword,
				Datacenters: getDataCenters(cfg),
			},
			Ntp: boshbase.NtpServers,
			Agent: &vsphere_cpi.Agent{
				Mbus: fmt.Sprintf("https://mbus:%s@0.0.0.0:6868", boshbase.MBusPassword),
			},
			Blobstore: &vsphere_cpi.Blobstore{
				Provider: "local",
				Path:     "/var/vcap/micro_bosh/data/cache",
			},
		},
	}
}

func getDataCenters(cfg VSphereInitConfig) VSphereDatacenters {
	return VSphereDatacenters{VSphereDatacenter{
		Name:                       cfg.VSphereDatacenterName,
		VMFolder:                   cfg.VSphereVMFolder,
		TemplateFolder:             cfg.VSphereTemplateFolder,
		DatastorePattern:           getDataStorePattern(cfg),
		PersistentDatastorePattern: getDataStorePattern(cfg),
		DiskPath:                   cfg.VSphereDiskPath,
		Clusters:                   clusterConfig(cfg),
	}}
}

func getDataStorePattern(cfg VSphereInitConfig) (pattern string) {
	return fmt.Sprintf("^(%s)$", cfg.VSphereDataStore)
}
