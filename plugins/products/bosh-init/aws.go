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
		NetworkCIDR:        "10.0.0.0/24",
		NetworkGateway:     "10.0.0.1",
		NetworkDNS:         []string{"10.0.0.2"},
		BoshReleaseURL:     "https://bosh.io/d/github.com/cloudfoundry/bosh?v=260",
		BoshReleaseSHA:     "f8f086974d9769263078fb6cb7927655744dacbc",
		CPIReleaseURL:      "https://bosh.io/d/github.com/cloudfoundry-incubator/bosh-aws-cpi-release?v=60",
		CPIReleaseSHA:      "8e40a9ff892204007889037f094a1b0d23777058",
		GOAgentReleaseURL:  "https://bosh.io/d/stemcells/bosh-aws-xen-hvm-ubuntu-trusty-go_agent?v=3263.10",
		GOAgentSHA:         "42ea4577caec4aec463bed951cfdffb935961270",
		PrivateIP:          "10.0.0.6",
		NtpServers:         []string{"0.amazon.pool.ntp.org", "1.amazon.pool.ntp.org", "2.amazon.pool.ntp.org", "3.amazon.pool.ntp.org"},
		CPIJobName:         awsCPIJobName,
		PersistentDiskSize: 51200,
	}
}

func (s *AWSBosh) CreateDeploymentManifest() (*enaml.DeploymentManifest, error) {
	manifest := s.boshbase.CreateDeploymentManifest()
	manifest.AddRelease(s.CreateCPIRelease())
	if rp, err := s.CreateResourcePool(); err != nil {
		return nil, err
	} else {
		manifest.AddResourcePool(*rp)
	}

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
	return manifest, nil
}

func (s *AWSBosh) resourcePoolCloudProperties() interface{} {
	return awscloudproperties.ResourcePool{
		InstanceType: s.cfg.AWSInstanceSize,
		EphemeralDisk: awscloudproperties.EphemeralDisk{
			Size:     s.boshbase.PersistentDiskSize,
			DiskType: "gp2",
		},
		AvailabilityZone: s.cfg.AWSAvailabilityZone,
	}
}
func (s *AWSBosh) CreateResourcePool() (*enaml.ResourcePool, error) {
	return s.boshbase.CreateResourcePool(s.resourcePoolCloudProperties)
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
		DiskSize: s.boshbase.PersistentDiskSize,
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
