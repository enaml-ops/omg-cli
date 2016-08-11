package boshinit

import (
	"fmt"

	"github.com/enaml-ops/enaml"
	"github.com/enaml-ops/omg-cli/plugins/products/bosh-init/enaml-gen/photoncpi"
)

const (
	PhotonCPIReleaseName = "bosh-photon-cpi"
	PhotonCPIURL         = "https://github.com/enaml-ops/bosh-photon-cpi-release/releases/download/v0.9.0/bosh-photon-cpi-0.9.0.dev.1.tgz"
	PhotonCPISHA         = "dc90202ee3981087a237fa0eb249d3345157a1e4"
	PhotonCPIJobName     = "cpi"

	PhotonStemcellURL = "https://bosh.io/d/stemcells/bosh-vsphere-esxi-ubuntu-trusty-go_agent?v=3232.1"
	PhotonStemcellSHA = "169df93e3e344cd84ac6ef16d76dd0276e321a25"
)

type PhotonBoshInitConfig struct {
	photoncpi.Photon
	MachineType string
	NetworkName string
}

type PhotonBosh struct {
	BoshInitConfig *PhotonBoshInitConfig
	Base           *BoshBase
}

func NewPhotonBoshBase(boshBase *BoshBase) *BoshBase {
	boshBase.CPIReleaseURL = PhotonCPIURL
	boshBase.CPIReleaseSHA = PhotonCPISHA
	boshBase.CPIJobName = "cpi"
	boshBase.NetworkCIDR = "10.0.0.0/24"
	boshBase.NetworkGateway = "10.0.0.1"
	boshBase.NetworkDNS = []string{"10.0.0.2"}
	boshBase.PrivateIP = "10.0.0.4"
	boshBase.BoshReleaseURL = "https//bosh.io/d/github.com/cloudfoundry/bosh?v=256.2"
	boshBase.BoshReleaseSHA = "ff2f4e16e02f66b31c595196052a809100cfd5a8"
	boshBase.GOAgentReleaseURL = PhotonStemcellURL
	boshBase.GOAgentSHA = PhotonStemcellSHA
	boshBase.NtpServers = []string{"0.pool.ntp.org", "1.pool.ntp.org"}
	return boshBase
}
func NewPhotonIaaSProvider(cfg *PhotonBoshInitConfig, boshBase *BoshBase) IAASManifestProvider {
	return &PhotonBosh{
		BoshInitConfig: cfg,
		Base:           boshBase,
	}
}

func (g *PhotonBosh) CreateCPIRelease() enaml.Release {
	return enaml.Release{
		Name: PhotonCPIReleaseName,
		URL:  g.Base.CPIReleaseURL,
		SHA1: g.Base.CPIReleaseSHA,
	}
}

func (g *PhotonBosh) CreateCPITemplate() enaml.Template {
	return enaml.Template{
		Name:    PhotonCPIJobName,
		Release: PhotonCPIReleaseName,
	}
}

func (g *PhotonBosh) CreateDiskPool() enaml.DiskPool {
	return enaml.DiskPool{
		Name:     "disks",
		DiskSize: 32768,
		CloudProperties: map[string]interface{}{
			"disk_flavor": "core-200",
		},
	}
}

func (g *PhotonBosh) CreateResourcePool() enaml.ResourcePool {
	return enaml.ResourcePool{
		Name:    "vms",
		Network: "private",
		Stemcell: enaml.Stemcell{
			URL:  PhotonStemcellURL,
			SHA1: PhotonStemcellSHA,
		},
		CloudProperties: map[string]interface{}{
			"vm_flavor":   g.BoshInitConfig.MachineType,
			"disk_flavor": "core-200",
		},
	}
}

func (g *PhotonBosh) CreateManualNetwork() enaml.ManualNetwork {
	net := enaml.NewManualNetwork("private")
	net.AddSubnet(enaml.Subnet{
		Range:   g.Base.NetworkCIDR,
		Gateway: g.Base.NetworkGateway,
		DNS:     g.Base.NetworkDNS,
		Static: []string{
			g.Base.PrivateIP,
		},
		CloudProperties: map[string]interface{}{
			"network_id": g.BoshInitConfig.NetworkName,
		},
	})
	return net
}

func (g *PhotonBosh) CreateVIPNetwork() enaml.VIPNetwork {
	return enaml.VIPNetwork{
		Name: "vip",
		Type: "vip",
	}
}

func (g *PhotonBosh) CreateJobNetwork() *enaml.Network {
	// photon just needs the default private network provided by boshbase
	return nil
}

func (g *PhotonBosh) CreateCloudProvider() enaml.CloudProvider {
	return enaml.CloudProvider{
		Template: g.CreateCPITemplate(),
		MBus:     fmt.Sprintf("https://mbus:%s@%s:6868", g.Base.MBusPassword, g.Base.PrivateIP),
		Properties: &photoncpi.PhotoncpiJob{
			Photon: &g.BoshInitConfig.Photon,
			Agent: &photoncpi.Agent{
				Mbus: fmt.Sprintf("https://mbus:%s@0.0.0.0:6868", g.Base.MBusPassword),
			},
			Blobstore: &photoncpi.Blobstore{
				Provider: "local",
				Options: map[string]string{
					"blobstore_path": "/var/vcap/micro_bosh/data/cache",
				},
			},
			Ntp: g.createNTP(),
		},
	}
}

func (g *PhotonBosh) createNTP() interface{} {
	return g.Base.NtpServers
}

func (g *PhotonBosh) CreateCPIJobProperties() map[string]interface{} {
	return map[string]interface{}{
		"photon": g.BoshInitConfig.Photon,
		"agent": &photoncpi.Agent{
			Mbus: fmt.Sprintf("nats://nats:%s@%s:4222", g.Base.NatsPassword, g.Base.PrivateIP),
		},
	}
}

func (g *PhotonBosh) CreateDeploymentManifest() *enaml.DeploymentManifest {
	var manifest = g.Base.CreateDeploymentManifest()
	manifest.AddRelease(g.CreateCPIRelease())
	manifest.AddResourcePool(g.CreateResourcePool())
	manifest.AddDiskPool(g.CreateDiskPool())
	manifest.AddNetwork(g.CreateManualNetwork())
	manifest.AddNetwork(g.CreateVIPNetwork())
	boshJob := manifest.Jobs[0]
	boshJob.AddTemplate(g.CreateCPITemplate())
	n := g.CreateJobNetwork()
	if n != nil {
		boshJob.AddNetwork(*n)
	}
	for name, val := range g.CreateCPIJobProperties() {
		boshJob.AddProperty(name, val)
	}
	manifest.Jobs[0] = boshJob
	manifest.SetCloudProvider(g.CreateCloudProvider())
	return manifest
}
