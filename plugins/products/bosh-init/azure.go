package boshinit

import (
	"fmt"

	"github.com/enaml-ops/enaml"
	"github.com/enaml-ops/enaml/cloudproperties/azure"
	"github.com/enaml-ops/omg-cli/plugins/products/bosh-init/enaml-gen/aws_cpi"
	"github.com/enaml-ops/omg-cli/plugins/products/bosh-init/enaml-gen/cpi"
)

const (
	azureCPIJobName     = "cpi"
	azureCPIReleaseName = "bosh-azure-cpi"
)

func GetAzureDefaults() *BoshBase {
	return &BoshBase{
		NetworkCIDR:       "10.0.0.0/24",
		NetworkGateway:    "10.0.0.1",
		NetworkDNS:        []string{"168.63.129.16"},
		BoshReleaseURL:    "https://bosh.io/d/github.com/cloudfoundry/bosh?v=256.2",
		BoshReleaseSHA:    "ff2f4e16e02f66b31c595196052a809100cfd5a8",
		CPIReleaseURL:     "https://bosh.io/d/github.com/cloudfoundry-incubator/bosh-azure-cpi-release?v=11",
		CPIReleaseSHA:     "395fc05c11ead59711188ebd0a684842a03dc93d",
		GOAgentReleaseURL: "https://bosh.io/d/stemcells/bosh-azure-hyperv-ubuntu-trusty-go_agent?v=3262.4",
		GOAgentSHA:        "1ec76310cd99d4ad2dd2b239b3dfde09c609b292",
		PrivateIP:         "10.0.0.4",
		NtpServers:        []string{"0.pool.ntp.org", "1.pool.ntp.org"},
		CPIJobName:        azureCPIJobName,
	}
}

func NewAzureBosh(cfg BoshInitConfig, boshbase *BoshBase) *enaml.DeploymentManifest {
	boshbase.CPIJobName = azureCPIJobName
	var cpiTemplate = enaml.Template{Name: boshbase.CPIJobName, Release: azureCPIReleaseName}
	manifest := boshbase.CreateDeploymentManifest()

	manifest.AddRelease(enaml.Release{
		Name: azureCPIReleaseName,
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
	resourcePool.CloudProperties = azurecloudproperties.ResourcePool{
		InstanceType: cfg.BoshInstanceSize,
	}
	manifest.AddResourcePool(resourcePool)
	manifest.AddDiskPool(enaml.DiskPool{
		Name:     "disks",
		DiskSize: 20000,
	})
	net := enaml.NewManualNetwork("private")
	net.AddSubnet(enaml.Subnet{
		Range:   boshbase.NetworkCIDR,
		Gateway: boshbase.NetworkGateway,
		DNS:     boshbase.NetworkDNS,
		CloudProperties: azurecloudproperties.Network{
			VnetName:   cfg.AzureVnet,
			SubnetName: cfg.AzureSubnet,
		},
	})
	manifest.AddNetwork(net)
	manifest.AddNetwork(enaml.NewVIPNetwork("public"))
	boshJob := manifest.Jobs[0]
	boshJob.AddTemplate(cpiTemplate)
	if boshbase.PublicIP != "" {
		boshJob.AddNetwork(enaml.Network{
			Name:      "public",
			StaticIPs: []string{boshbase.PublicIP},
		})
	}
	var agentProperty = aws_cpi.Agent{
		Mbus: "nats://nats:nats-password@" + boshbase.PrivateIP + ":4222",
	}
	boshJob.AddProperty("agent", agentProperty)
	azureProperty := NewAzureProperty(
		cfg.AzureEnvironment,
		cfg.AzureSubscriptionID,
		cfg.AzureTenantID,
		cfg.AzureClientID,
		cfg.AzureClientSecret,
		cfg.AzureResourceGroup,
		cfg.AzureStorageAccount,
		cfg.AzureDefaultSecurityGroup,
		cfg.AzureSSHUser,
		cfg.AzureSSHPubKey,
	)
	boshJob.AddProperty("azure", azureProperty)
	manifest.Jobs[0] = boshJob
	manifest.SetCloudProvider(NewAzureCloudProvider(azureProperty, cpiTemplate, boshbase.GetRoutableIP(), cfg.AzurePrivateKeyPath, boshbase.NtpServers))
	return manifest
}

func NewAzureCloudProvider(myazure cpi.Azure, cpiTemplate enaml.Template, pubip, keypath string, ntpProperty []string) enaml.CloudProvider {
	return enaml.CloudProvider{
		Template: cpiTemplate,
		MBus:     fmt.Sprintf("https://mbus:mbus-password@%s:6868", pubip),
		SSHTunnel: enaml.SSHTunnel{
			Host:           pubip,
			Port:           22,
			User:           "vcap",
			PrivateKeyPath: keypath,
		},
		Properties: map[string]interface{}{
			"azure": myazure,
			"agent": map[string]string{
				"mbus": "https://mbus:mbus-password@0.0.0.0:6868",
			},
			"blobstore": map[string]string{
				"provider": "local",
				"path":     "/var/vcap/micro_bosh/data/cache",
			},
			"ntp": ntpProperty,
		},
	}
}

func NewAzureProperty(azureenv, subid, tenantid, clientid, clientsecret, resourcegroup, storageaccount, securitygroup, sshuser, sshkey string) cpi.Azure {
	return cpi.Azure{
		Environment:          azureenv,
		SubscriptionId:       subid,
		TenantId:             tenantid,
		ClientId:             clientid,
		ClientSecret:         clientsecret,
		ResourceGroupName:    resourcegroup,
		StorageAccountName:   storageaccount,
		DefaultSecurityGroup: securitygroup,
		SshUser:              sshuser,
		SshPublicKey:         sshkey,
	}
}
