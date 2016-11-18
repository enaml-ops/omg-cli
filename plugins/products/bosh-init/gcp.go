package boshinit

import (
	"fmt"

	"github.com/enaml-ops/enaml"
	"github.com/enaml-ops/omg-cli/plugins/products/bosh-init/enaml-gen/google_cpi"
)

const (
	GCPCPIReleaseName = "bosh-google-cpi"
	GCPCPIURL         = "https://storage.googleapis.com/bosh-cpi-artifacts/bosh-google-cpi-24.2.0.tgz"
	GCPCPISHA         = "80d3ef039cb0ed014e97eeea10569598804659d3"
	GCPCPIJobName     = "google_cpi"

	GCPStemcellURL = "https://storage.googleapis.com/bosh-cpi-artifacts/light-bosh-stemcell-3262.4-google-kvm-ubuntu-trusty-go_agent.tgz"
	GCPStemcellSHA = "1f44ee6fc5fd495113694aa772d636bf1a8d645a"
)

type GCPBoshInitConfig struct {
	NetworkName       string
	SubnetworkName    string
	DefaultZone       string
	Project           string
	MachineType       string
	DiskType          string
	DNSRecursor       string
	BlobstoreProvider string
}

type GCPBosh struct {
	BoshInitConfig *GCPBoshInitConfig
	Base           *BoshBase
}

func NewGCPBoshBase() *BoshBase {
	return &BoshBase{
		CPIJobName:        GCPCPIJobName,
		CPIReleaseURL:     GCPCPIURL,
		CPIReleaseSHA:     GCPCPISHA,
		NetworkCIDR:       "10.0.0.0/24",
		NetworkGateway:    "10.0.0.1",
		NetworkDNS:        []string{"10.0.0.2"},
		PrivateIP:         "10.0.0.4",
		BoshReleaseURL:    "https://bosh.io/d/github.com/cloudfoundry/bosh?v=256.2",
		BoshReleaseSHA:    "ff2f4e16e02f66b31c595196052a809100cfd5a8",
		GOAgentReleaseURL: "https://storage.googleapis.com/bosh-cpi-artifacts/light-bosh-stemcell-3262.4-google-kvm-ubuntu-trusty-go_agent.tgz",
		GOAgentSHA:        "1f44ee6fc5fd495113694aa772d636bf1a8d645a",
		NtpServers:        []string{"0.pool.ntp.org", "1.pool.ntp.org"},
	}
}
func NewGCPIaaSProvider(cfg *GCPBoshInitConfig, boshBase *BoshBase) IAASManifestProvider {
	boshBase.CPIJobName = GCPCPIJobName
	return &GCPBosh{
		BoshInitConfig: cfg,
		Base:           boshBase,
	}
}

func (g *GCPBosh) CreateCPIRelease() enaml.Release {
	return enaml.Release{
		Name: GCPCPIReleaseName,
		URL:  g.Base.CPIReleaseURL,
		SHA1: g.Base.CPIReleaseSHA,
	}
}

func (g *GCPBosh) CreateCPITemplate() enaml.Template {
	return enaml.Template{
		Name:    g.Base.CPIJobName,
		Release: GCPCPIReleaseName,
	}
}

func (g *GCPBosh) CreateDiskPool() enaml.DiskPool {
	return enaml.DiskPool{
		Name:     "disks",
		DiskSize: 32768,
		CloudProperties: map[string]interface{}{
			"type": "pd-standard",
		},
	}
}

func (g *GCPBosh) resourcePoolCloudProperties() interface{} {
	return map[string]interface{}{
		"machine_type":      g.BoshInitConfig.MachineType,
		"root_disk_size_gb": 50,
		"root_disk_type":    "pd-standard",
		"service_scopes":    []string{"compute", "devstorage.full_control"},
	}
}
func (g *GCPBosh) CreateResourcePool() (*enaml.ResourcePool, error) {
	return g.Base.CreateResourcePool(g.resourcePoolCloudProperties)
}

func (g *GCPBosh) CreateManualNetwork() enaml.ManualNetwork {
	net := enaml.NewManualNetwork("private")
	net.AddSubnet(enaml.Subnet{
		Range:   g.Base.NetworkCIDR,
		Gateway: g.Base.NetworkGateway,
		DNS:     g.Base.NetworkDNS,
		Static: []string{
			g.Base.PrivateIP,
		},
		CloudProperties: map[string]interface{}{
			"network_name":          g.BoshInitConfig.NetworkName,
			"subnetwork_name":       g.BoshInitConfig.SubnetworkName,
			"ephemeral_external_ip": false,
			"tags":                  []string{"nat-traverse", "no-ip"},
		},
	})
	return net
}

func (g *GCPBosh) CreateVIPNetwork() enaml.VIPNetwork {
	return enaml.VIPNetwork{
		Name: "vip",
		Type: "vip",
	}
}

func (g *GCPBosh) CreateJobNetwork() *enaml.Network {
	return &enaml.Network{
		Name:      "private",
		StaticIPs: []string{g.Base.PrivateIP},
		Default:   []interface{}{"dns", "gateway"},
	}
}

func (g *GCPBosh) CreateCloudProvider() enaml.CloudProvider {
	return enaml.CloudProvider{
		Template: enaml.Template{
			Name:    g.Base.CPIJobName,
			Release: GCPCPIReleaseName,
		},
		SSHTunnel: enaml.SSHTunnel{
			Host:           g.Base.PrivateIP,
			Port:           22,
			User:           "bosh",
			PrivateKeyPath: "~/.ssh/bosh",
		},
		MBus: fmt.Sprintf("https://mbus:%s@%s:6868", g.Base.MBusPassword, g.Base.PrivateIP),
		Properties: &google_cpi.GoogleCpiJob{
			Google: g.createGoogleProperties(),
			Agent: &google_cpi.Agent{
				Mbus: fmt.Sprintf("https://mbus:%s@0.0.0.0:6868", g.Base.MBusPassword),
			},
			Blobstore: &google_cpi.Blobstore{
				Provider: "local",
				Path:     "/var/vcap/micro_bosh/data/cache",
			},
			Ntp: g.createNTP(),
		},
	}
}

func (g *GCPBosh) createGoogleProperties() *google_cpi.Google {
	return &google_cpi.Google{
		Project:     g.BoshInitConfig.Project,
		DefaultZone: g.BoshInitConfig.DefaultZone,
	}
}

func (g *GCPBosh) createNTP() interface{} {
	return g.Base.NtpServers
}

func (g *GCPBosh) CreateCPIJobProperties() map[string]interface{} {
	return map[string]interface{}{
		"google": g.createGoogleProperties(),
		"agent": &google_cpi.Agent{
			Mbus: fmt.Sprintf("nats://nats:%s@%s:4222", g.Base.NatsPassword, g.Base.PrivateIP),
		},
	}
}

func (g *GCPBosh) CreateDeploymentManifest() (*enaml.DeploymentManifest, error) {
	var manifest = g.Base.CreateDeploymentManifest()
	manifest.AddRelease(g.CreateCPIRelease())
	if rp, err := g.CreateResourcePool(); err != nil {
		return nil, err
	} else {
		manifest.AddResourcePool(*rp)
	}
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
	return manifest, nil
}
