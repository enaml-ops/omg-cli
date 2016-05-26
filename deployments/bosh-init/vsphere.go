package boshinit

import (
	"fmt"

	"github.com/enaml-ops/enaml"
	"github.com/enaml-ops/omg-cli/deployments/bosh-init/enaml-gen/vsphere_cpi"
)

// NewVSphereBosh creates a new enaml deployment manifest for vSphere
func NewVSphereBosh(cfg BoshInitConfig) *enaml.DeploymentManifest {
	var ntpProperty = NewNTP("0.pool.ntp.org", "1.pool.ntp.org")
	var manifest = NewBoshDeploymentBase(cfg, "vsphere_cpi", ntpProperty)
	var vcenterProperty = vsphere_cpi.Vcenter{
		Address:  cfg.VSphereAddress,
		User:     cfg.VSphereUser,
		Password: cfg.VSpherePassword,
		Datacenters: []map[string]interface{}{{
			"name":                         cfg.VSphereDatacenterName,
			"vm_folder":                    cfg.VSphereVMFolder,
			"template_folder":              cfg.VSphereTemplateFolder,
			"datastore_pattern":            cfg.VSphereDatastorePattern,
			"persistent_datastore_pattern": cfg.VSpherePersistentDatastorePattern,
			"disk_path":                    cfg.VSphereDiskPath,
			"clusters":                     cfg.VSphereClusters,
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
	resourcePool.CloudProperties = vspherecloudpropertiesResourcePool{
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
		Range:   "10.0.0.0/24",
		Gateway: "10.0.0.1",
		DNS:     []string{"10.0.0.2"},
		CloudProperties: vspherecloudpropertiesNetwork{
			Name: cfg.VSphereNetworkName,
		},
	})
	manifest.AddNetwork(net)

	boshJob := manifest.Jobs[0]
	boshJob.AddTemplate(enaml.Template{Name: "vsphere_cpi", Release: "bosh-vsphere-cpi"})
	boshJob.AddProperty("agent", agentProperty)
	boshJob.AddProperty("vcenter", vcenterProperty)
	manifest.Jobs[0] = boshJob
	manifest.SetCloudProvider(NewVSphereCloudProvider(cfg.BoshPrivateIP, vcenterProperty, ntpProperty))
	return manifest
}

type vspherecloudpropertiesResourcePool struct {
	CPU  int `yaml:"cpu,omitempty"`  // [Integer, required]: Number of CPUs.
	RAM  int `yaml:"ram,omitempty"`  // [Integer, required]: Specified the amount of RAM in megabytes.
	Disk int `yaml:"disk,omitempty"` // [Integer, required]: Specifies the disk size in megabytes.
}

type vspherecloudpropertiesNetwork struct {
	Name string `yaml:"name,omitempty"` // [String, required]: vSphere network name.
}
