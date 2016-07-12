package boshinit

import (
	"fmt"

	"github.com/enaml-ops/enaml"
	"github.com/enaml-ops/omg-cli/plugins/products/bosh-init/enaml-gen/vsphere_cpi"
)

// NewVSphereBosh creates a new enaml deployment manifest for vSphere
func NewVSphereBosh(cfg BoshInitConfig) *enaml.DeploymentManifest {
	var ntpProperty = NewNTP("0.pool.ntp.org", "1.pool.ntp.org")
	var manifest = NewBoshDeploymentBase(cfg, "vsphere_cpi", ntpProperty)

	persistentDatastorePattern := cfg.VSpherePersistentDatastorePattern
	if len(persistentDatastorePattern) == 0 {
		persistentDatastorePattern = cfg.VSphereDatastorePattern
	}
	var vcenterProperty = vsphere_cpi.Vcenter{
		Address:  cfg.VSphereAddress,
		User:     cfg.VSphereUser,
		Password: cfg.VSpherePassword,
		Datacenters: VSphereDatacenters{VSphereDatacenter{
			Name:                       cfg.VSphereDatacenterName,
			VMFolder:                   cfg.VSphereVMFolder,
			TemplateFolder:             cfg.VSphereTemplateFolder,
			DatastorePattern:           cfg.VSphereDatastorePattern,
			PersistentDatastorePattern: persistentDatastorePattern,
			DiskPath:                   cfg.VSphereDiskPath,
			Clusters:                   cfg.VSphereClusters,
		}},
	}

	var agentProperty = vsphere_cpi.Agent{
		Mbus: fmt.Sprintf("nats://nats:nats-password@%s:4222", cfg.BoshPrivateIP),
	}

	manifest.AddRelease(enaml.Release{
		Name: "bosh-vsphere-cpi",
		URL:  "https://bosh.io/d/github.com/cloudfoundry-incubator/bosh-vsphere-cpi-release?v=" + cfg.BoshCPIReleaseVersion,
		SHA1: cfg.BoshCPIReleaseSHA,
	})

	resourcePool := enaml.ResourcePool{
		Name:    "vms",
		Network: "private",
	}
	resourcePool.Stemcell = enaml.Stemcell{
		URL:  "https://bosh.io/d/stemcells/bosh-vsphere-esxi-ubuntu-trusty-go_agent?v=" + cfg.GoAgentVersion,
		SHA1: cfg.GoAgentSHA,
	}
	resourcePool.CloudProperties = VSpherecloudpropertiesResourcePool{
		CPU:  2,
		Disk: 20000,
		RAM:  4096,
	}
	// c1oudc0w is a default password for vcap user
	resourcePool.Env = map[string]interface{}{
		"bosh": map[string]string{
			"password": "$6$4gDD3aV0rdqlrKC$2axHCxGKIObs6tAmMTqYCspcdvQXh3JJcvWOY2WGb4SrdXtnCyNaWlrf3WEqvYR2MYizEGp3kMmbpwBC6jsHt0",
		},
	}
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
	boshJob.AddTemplate(enaml.Template{Name: "vsphere_cpi", Release: "bosh-vsphere-cpi"})
	boshJob.AddProperty(agentProperty)
	boshJob.AddProperty(vcenterProperty)
	manifest.Jobs[0] = boshJob
	manifest.SetCloudProvider(NewVSphereCloudProvider(cfg.BoshPrivateIP, vcenterProperty, ntpProperty))
	return manifest
}

type VSpherecloudpropertiesResourcePool struct {
	CPU  int `yaml:"cpu,omitempty"`  // [Integer, required]: Number of CPUs.
	RAM  int `yaml:"ram,omitempty"`  // [Integer, required]: Specified the amount of RAM in megabytes.
	Disk int `yaml:"disk,omitempty"` // [Integer, required]: Specifies the disk size in megabytes.
}

type VSpherecloudpropertiesNetwork struct {
	Name string `yaml:"name,omitempty"` // [String, required]: vSphere network name.
}

type VSphereDatacenters []VSphereDatacenter

type VSphereDatacenter struct {
	Name                       string   `yaml:"name"`                         // [String, required]: vSphere datacenter name.
	VMFolder                   string   `yaml:"vm_folder"`                    // [String, required]: The folder to create PCF VMs in.
	TemplateFolder             string   `yaml:"template_folder"`              // [String, required]: The folder to store stemcells in.
	DatastorePattern           string   `yaml:"datastore_pattern"`            // [String, required]: The pattern to the vSphere datastore.
	PersistentDatastorePattern string   `yaml:"persistent_datastore_pattern"` // [String, required]: The pattern to the vSphere datastore for persistent disks.
	DiskPath                   string   `yaml:"disk_path"`                    // [String, required]: The disk path.
	Clusters                   []string `yaml:"clusters"`                     // [[]String], required]: The vSphere cluster(s).
}
