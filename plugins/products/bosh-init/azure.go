package boshinit

import (
	"fmt"

	"github.com/enaml-ops/enaml"
	"github.com/enaml-ops/enaml/cloudproperties/azure"
	"github.com/enaml-ops/omg-cli/plugins/products/bosh-init/enaml-gen/azure_cpi"
)

const (
	azureCPIJobName     = "cpi"
	azureCPIReleaseName = "bosh-azure-cpi"
)

type AzureInitConfig struct {
	AzureInstanceSize         string
	AzureVnet                 string
	AzureSubnet               string
	AzureSubscriptionID       string
	AzureTenantID             string
	AzureClientID             string
	AzureClientSecret         string
	AzureResourceGroup        string
	AzureStorageAccount       string
	AzureDefaultSecurityGroup string
	AzureSSHPubKey            string
	AzureSSHUser              string
	AzureEnvironment          string
	AzurePrivateKeyPath       string
}

func GetAzureDefaults() *BoshBase {
	return &BoshBase{
		NetworkCIDR:        "10.0.0.0/24",
		NetworkGateway:     "10.0.0.1",
		NetworkDNS:         []string{"168.63.129.16"},
		BoshReleaseURL:     "https://bosh.io/d/github.com/cloudfoundry/bosh?v=256.2",
		BoshReleaseSHA:     "ff2f4e16e02f66b31c595196052a809100cfd5a8",
		CPIReleaseURL:      "https://bosh.io/d/github.com/cloudfoundry-incubator/bosh-azure-cpi-release?v=11",
		CPIReleaseSHA:      "395fc05c11ead59711188ebd0a684842a03dc93d",
		GOAgentReleaseURL:  "https://bosh.io/d/stemcells/bosh-azure-hyperv-ubuntu-trusty-go_agent?v=3262.4",
		GOAgentSHA:         "1ec76310cd99d4ad2dd2b239b3dfde09c609b292",
		PrivateIP:          "10.0.0.4",
		NtpServers:         []string{"0.pool.ntp.org", "1.pool.ntp.org"},
		CPIJobName:         azureCPIJobName,
		PersistentDiskSize: 20000,
	}
}

type AzureBosh struct {
	cfg      AzureInitConfig
	boshbase *BoshBase
}

func NewAzureIaaSProvider(cfg AzureInitConfig, boshBase *BoshBase) IAASManifestProvider {
	boshBase.CPIJobName = azureCPIJobName
	return &AzureBosh{
		cfg:      cfg,
		boshbase: boshBase,
	}
}

func (a *AzureBosh) CreateCPIRelease() enaml.Release {
	return enaml.Release{
		Name: azureCPIReleaseName,
		URL:  a.boshbase.CPIReleaseURL,
		SHA1: a.boshbase.CPIReleaseSHA,
	}
}
func (a *AzureBosh) CreateCPITemplate() enaml.Template {
	return enaml.Template{
		Name:    a.boshbase.CPIJobName,
		Release: azureCPIReleaseName}
}
func (a *AzureBosh) CreateDiskPool() enaml.DiskPool {
	return enaml.DiskPool{
		Name:     "disks",
		DiskSize: a.boshbase.PersistentDiskSize,
	}
}

func (a *AzureBosh) resourcePoolCloudProperties() interface{} {
	return azurecloudproperties.ResourcePool{
		InstanceType: a.cfg.AzureInstanceSize,
	}
}
func (a *AzureBosh) CreateResourcePool() (*enaml.ResourcePool, error) {
	return a.boshbase.CreateResourcePool(a.resourcePoolCloudProperties)
}

func (a *AzureBosh) CreateManualNetwork() enaml.ManualNetwork {
	net := enaml.NewManualNetwork("private")
	net.AddSubnet(enaml.Subnet{
		Range:   a.boshbase.NetworkCIDR,
		Gateway: a.boshbase.NetworkGateway,
		DNS:     a.boshbase.NetworkDNS,
		CloudProperties: azurecloudproperties.Network{
			VnetName:   a.cfg.AzureVnet,
			SubnetName: a.cfg.AzureSubnet,
		},
	})
	return net
}
func (a *AzureBosh) CreateVIPNetwork() enaml.VIPNetwork {
	return enaml.NewVIPNetwork("public")
}
func (a *AzureBosh) CreateJobNetwork() *enaml.Network {
	if a.boshbase.PublicIP != "" {
		return &enaml.Network{
			Name:      "public",
			StaticIPs: []string{a.boshbase.PublicIP},
		}
	}
	return nil
}
func (a *AzureBosh) CreateCloudProvider() enaml.CloudProvider {
	return enaml.CloudProvider{
		Template: a.CreateCPITemplate(),
		MBus:     fmt.Sprintf("https://mbus:%s@%s:6868", a.boshbase.MBusPassword, a.boshbase.GetRoutableIP()),
		SSHTunnel: enaml.SSHTunnel{
			Host:           a.boshbase.GetRoutableIP(),
			Port:           22,
			User:           "vcap",
			PrivateKeyPath: a.cfg.AzurePrivateKeyPath,
		},
		Properties: azure_cpi.AzureCpiJob{
			Azure: a.createAzure(),
			Agent: &azure_cpi.Agent{
				Mbus: fmt.Sprintf("https://mbus:%s@0.0.0.0:6868", a.boshbase.MBusPassword),
			},
			Blobstore: &azure_cpi.Blobstore{
				Provider: "local",
				Path:     "/var/vcap/micro_bosh/data/cache",
			},
			Ntp: a.boshbase.NtpServers,
		},
	}
}

func (a *AzureBosh) createAzure() *azure_cpi.Azure {
	return &azure_cpi.Azure{
		Environment:          a.cfg.AzureEnvironment,
		SubscriptionId:       a.cfg.AzureSubscriptionID,
		TenantId:             a.cfg.AzureTenantID,
		ClientId:             a.cfg.AzureClientID,
		ClientSecret:         a.cfg.AzureClientSecret,
		ResourceGroupName:    a.cfg.AzureResourceGroup,
		StorageAccountName:   a.cfg.AzureStorageAccount,
		DefaultSecurityGroup: a.cfg.AzureDefaultSecurityGroup,
		SshUser:              a.cfg.AzureSSHUser,
		SshPublicKey:         a.cfg.AzureSSHPubKey,
	}
}
func (a *AzureBosh) CreateCPIJobProperties() map[string]interface{} {
	return map[string]interface{}{
		"azure": a.createAzure(),
		"agent": &azure_cpi.Agent{
			Mbus: fmt.Sprintf("nats://nats:%s@%s:4222", a.boshbase.NatsPassword, a.boshbase.PrivateIP),
		},
	}
}

func (a *AzureBosh) CreateDeploymentManifest() (*enaml.DeploymentManifest, error) {
	manifest := a.boshbase.CreateDeploymentManifest()
	manifest.AddRelease(a.CreateCPIRelease())
	if rp, err := a.CreateResourcePool(); err != nil {
		return nil, err
	} else {
		manifest.AddResourcePool(*rp)
	}
	manifest.AddDiskPool(a.CreateDiskPool())
	manifest.AddNetwork(a.CreateManualNetwork())
	manifest.AddNetwork(a.CreateVIPNetwork())
	boshJob := manifest.Jobs[0]
	boshJob.AddTemplate(a.CreateCPITemplate())
	n := a.CreateJobNetwork()
	if n != nil {
		boshJob.AddNetwork(*n)
	}
	for name, val := range a.CreateCPIJobProperties() {
		boshJob.AddProperty(name, val)
	}
	manifest.Jobs[0] = boshJob
	manifest.SetCloudProvider(a.CreateCloudProvider())
	return manifest, nil
}
