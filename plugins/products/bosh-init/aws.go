package boshinit

import (
	"fmt"

	"github.com/enaml-ops/enaml"
	"github.com/enaml-ops/enaml/cloudproperties/aws"
	"github.com/enaml-ops/omg-cli/plugins/products/bosh-init/enaml-gen/aws_cpi"
)

const (
	awsCPIReleaseName = "bosh-aws-cpi"
)

type AWSBosh struct {
	cfg      BoshInitConfig
	boshbase *BoshBase
}

func NewAWSIaaSProvider(cfg BoshInitConfig, boshBase *BoshBase) IAASManifestProvider {

	return &AWSBosh{
		cfg:      cfg,
		boshbase: boshBase,
	}
}

func NewAWSBosh(cfg BoshInitConfig, boshbase *BoshBase) *enaml.DeploymentManifest {
	aws := NewAWSIaaSProvider(cfg, boshbase)
	var manifest = aws.CreateDeploymentManifest()
	manifest.AddRelease(aws.CreateCPIRelease())
	manifest.AddResourcePool(aws.CreateResourcePool())
	manifest.AddDiskPool(aws.CreateDiskPool())
	manifest.AddNetwork(aws.CreateManualNetwork())
	manifest.AddNetwork(aws.CreateVIPNetwork())
	boshJob := manifest.Jobs[0]
	boshJob.AddTemplate(aws.CreateCPITemplate())
	boshJob.AddNetwork(aws.CreateJobNetwork())
	for name, val := range aws.CreateCPIJobProperties() {
		boshJob.AddProperty(name, val)
	}
	manifest.Jobs[0] = boshJob
	manifest.SetCloudProvider(aws.CreateCloudProvider())
	return manifest
}

func (s *AWSBosh) CreateDeploymentManifest() *enaml.DeploymentManifest {
	return s.boshbase.CreateDeploymentManifest()
}

func (s *AWSBosh) CreateResourcePool() (resourcePool enaml.ResourcePool) {
	resourcePool = enaml.ResourcePool{
		Name:    "vms",
		Network: "private",
	}
	resourcePool.Stemcell = enaml.Stemcell{
		URL:  "https://bosh.io/d/stemcells/bosh-aws-xen-hvm-ubuntu-trusty-go_agent?v=" + s.boshbase.GOAgentVersion,
		SHA1: s.boshbase.GOAgentSHA,
	}
	resourcePool.CloudProperties = awscloudproperties.ResourcePool{
		InstanceType: s.cfg.BoshInstanceSize,
		EphemeralDisk: awscloudproperties.EphemeralDisk{
			Size:     25000,
			DiskType: "gp2",
		},
		AvailabilityZone: s.cfg.BoshAvailabilityZone,
	}
	return
}

func (s *AWSBosh) CreateCPIRelease() enaml.Release {
	return enaml.Release{
		Name: awsCPIReleaseName,
		URL:  "https://bosh.io/d/github.com/cloudfoundry-incubator/bosh-aws-cpi-release?v=" + s.boshbase.CPIReleaseVersion,
		SHA1: s.boshbase.CPIReleaseSHA,
	}
}

func (s *AWSBosh) CreateDiskPool() enaml.DiskPool {
	return enaml.DiskPool{
		Name:     "disks",
		DiskSize: 20000,
		CloudProperties: awscloudproperties.RootDisk{
			DiskType: "gp2",
		},
	}
}

func (s *AWSBosh) CreateManualNetwork() (net enaml.ManualNetwork) {
	net = enaml.NewManualNetwork("private")
	net.AddSubnet(enaml.Subnet{
		Range:   s.boshbase.NetworkCIDR,
		Gateway: s.boshbase.NetworkGateway,
		DNS:     s.boshbase.NetworkDNS,
		Static: []string{
			s.boshbase.PrivateIP,
		},
		CloudProperties: awscloudproperties.Network{
			Subnet: s.cfg.AWSSubnet,
		},
	})
	return
}

func (s *AWSBosh) CreateJobNetwork() (net enaml.Network) {
	if s.boshbase.PublicIP != "" {
		return enaml.Network{
			Name:      "public",
			StaticIPs: []string{s.boshbase.PublicIP},
		}
	}
	return enaml.Network{
		Name: "public",
	}

}

func (s *AWSBosh) CreateVIPNetwork() (net enaml.VIPNetwork) {
	return enaml.NewVIPNetwork("public")
}

func (s *AWSBosh) CreateCPIJobProperties() map[string]interface{} {
	return map[string]interface{}{
		"aws": &aws_cpi.Aws{
			AccessKeyId:           s.cfg.AWSAccessKeyID,
			SecretAccessKey:       s.cfg.AWSSecretKey,
			DefaultKeyName:        s.cfg.AWSKeyName,
			DefaultSecurityGroups: s.cfg.AWSSecurityGroups,
			Region:                s.cfg.AWSRegion,
		},
		"agent": &aws_cpi.Agent{
			Mbus: fmt.Sprintf("nats://nats:%s@%s:4222", s.boshbase.NatsPassword, s.boshbase.PrivateIP),
		},
	}
}

func (s *AWSBosh) CreateCPITemplate() (template enaml.Template) {
	return enaml.Template{
		Name:    s.boshbase.CPIName,
		Release: awsCPIReleaseName,
	}
}

func (s *AWSBosh) CreateCloudProvider() (provider enaml.CloudProvider) {
	return enaml.CloudProvider{
		Template: enaml.Template{
			Name:    s.boshbase.CPIName,
			Release: awsCPIReleaseName,
		},
		MBus: fmt.Sprintf("https://mbus:%s@%s:6868", s.boshbase.MBusPassword, s.boshbase.GetRoutableIP()),
		SSHTunnel: enaml.SSHTunnel{
			Host:           s.boshbase.GetRoutableIP(),
			Port:           22,
			User:           "vcap",
			PrivateKeyPath: s.cfg.AWSPEMFilePath,
		},
		Properties: &aws_cpi.AwsCpiJob{
			Aws: &aws_cpi.Aws{
				AccessKeyId:           s.cfg.AWSAccessKeyID,
				SecretAccessKey:       s.cfg.AWSSecretKey,
				DefaultKeyName:        s.cfg.AWSKeyName,
				DefaultSecurityGroups: s.cfg.AWSSecurityGroups,
				Region:                s.cfg.AWSRegion,
			},
			Ntp: s.boshbase.NtpServers,
			Agent: &aws_cpi.Agent{
				Mbus: fmt.Sprintf("https://mbus:%s@0.0.0.0:6868", s.boshbase.MBusPassword),
			},
			Blobstore: &aws_cpi.Blobstore{
				Provider: "local",
				Path:     "/var/vcap/micro_bosh/data/cache",
			},
		},
	}
}
