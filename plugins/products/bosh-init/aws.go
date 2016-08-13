package boshinit

import (
	"fmt"

	"github.com/enaml-ops/enaml"
	"github.com/enaml-ops/enaml/cloudproperties/aws"
	"github.com/enaml-ops/omg-cli/plugins/products/bosh-init/enaml-gen/aws_cpi"
)

const (
	awsCPIJobName     = "aws_cpi"
	awsCPIReleaseName = "bosh-aws-cpi"
)

type AWSInitConfig struct {
	AWSAvailabilityZone string
	AWSInstanceSize     string
	AWSSubnet           string
	AWSPEMFilePath      string
	AWSAccessKeyID      string
	AWSSecretKey        string
	AWSRegion           string
	AWSSecurityGroups   []string
	AWSKeyName          string
}

type AWSBosh struct {
	cfg      AWSInitConfig
	boshbase *BoshBase
}

func NewAWSIaaSProvider(cfg AWSInitConfig, boshBase *BoshBase) IAASManifestProvider {
	boshBase.CPIJobName = awsCPIJobName
	return &AWSBosh{
		cfg:      cfg,
		boshbase: boshBase,
	}
}

func GetAWSBoshBase() *BoshBase {
	return &BoshBase{
		NetworkCIDR:       "10.0.0.0/24",
		NetworkGateway:    "10.0.0.1",
		NetworkDNS:        []string{"10.0.0.2"},
		BoshReleaseURL:    "https://bosh.io/d/github.com/cloudfoundry/bosh?v=257.3",
		BoshReleaseSHA:    "e4442afcc64123e11f2b33cc2be799a0b59207d0",
		CPIReleaseURL:     "https://bosh.io/d/github.com/cloudfoundry-incubator/bosh-aws-cpi-release?v=57",
		CPIReleaseSHA:     "cbc7ed758f4a41063e9aee881bfc164292664b84",
		GOAgentReleaseURL: "https://bosh.io/d/stemcells/bosh-aws-xen-hvm-ubuntu-trusty-go_agent?v=3262.7",
		GOAgentSHA:        "bf44a5f81d29346af6a309199d2e012237dd222c",
		PrivateIP:         "10.0.0.6",
		NtpServers:        []string{"0.pool.ntp.org", "1.pool.ntp.org"},
		CPIJobName:        awsCPIJobName,
	}
}

func (s *AWSBosh) CreateDeploymentManifest() *enaml.DeploymentManifest {
	manifest := s.boshbase.CreateDeploymentManifest()
	manifest.AddRelease(s.CreateCPIRelease())
	manifest.AddResourcePool(s.CreateResourcePool())
	manifest.AddDiskPool(s.CreateDiskPool())
	manifest.AddNetwork(s.CreateManualNetwork())
	manifest.AddNetwork(s.CreateVIPNetwork())
	boshJob := manifest.Jobs[0]
	boshJob.AddTemplate(s.CreateCPITemplate())
	n := s.CreateJobNetwork()
	if n != nil {
		boshJob.AddNetwork(*n)
	}
	for name, val := range s.CreateCPIJobProperties() {
		boshJob.AddProperty(name, val)
	}
	manifest.Jobs[0] = boshJob
	manifest.SetCloudProvider(s.CreateCloudProvider())
	return manifest
}

func (s *AWSBosh) CreateResourcePool() (resourcePool enaml.ResourcePool) {
	resourcePool = enaml.ResourcePool{
		Name:    "vms",
		Network: "private",
	}
	resourcePool.Stemcell = enaml.Stemcell{
		URL:  s.boshbase.GOAgentReleaseURL,
		SHA1: s.boshbase.GOAgentSHA,
	}
	resourcePool.CloudProperties = awscloudproperties.ResourcePool{
		InstanceType: s.cfg.AWSInstanceSize,
		EphemeralDisk: awscloudproperties.EphemeralDisk{
			Size:     25000,
			DiskType: "gp2",
		},
		AvailabilityZone: s.cfg.AWSAvailabilityZone,
	}
	return
}

func (s *AWSBosh) CreateCPIRelease() enaml.Release {
	return enaml.Release{
		Name: awsCPIReleaseName,
		URL:  s.boshbase.CPIReleaseURL,
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

func (s *AWSBosh) CreateJobNetwork() *enaml.Network {
	if s.boshbase.PublicIP != "" {
		return &enaml.Network{
			Name:      "public",
			StaticIPs: []string{s.boshbase.PublicIP},
		}
	}
	return &enaml.Network{
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
		Name:    s.boshbase.CPIJobName,
		Release: awsCPIReleaseName,
	}
}

func (s *AWSBosh) CreateCloudProvider() (provider enaml.CloudProvider) {
	return enaml.CloudProvider{
		Template: enaml.Template{
			Name:    s.boshbase.CPIJobName,
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
