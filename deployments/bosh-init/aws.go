package boshinit

import (
	"github.com/bosh-ops/bosh-install/deployments/bosh-init/enaml-gen/aws_cpi"
	"github.com/bosh-ops/bosh-install/deployments/bosh-init/enaml-gen/director"
	"github.com/bosh-ops/bosh-install/deployments/bosh-init/enaml-gen/health_monitor"
	"github.com/xchapter7x/enaml"
	"github.com/xchapter7x/enaml/cloudproperties/aws"
)

func NewAWSBosh(cfg BoshInitConfig) *enaml.DeploymentManifest {
	pgsql := NewPostgres("postgres", "127.0.0.1", "postgres-password", "bosh", "postgres")
	var directorProperty = NewDirectorProperty(cfg.BoshDirectorName, "aws_cpi", pgsql.GetDirectorDB())
	var registryProperty = GetRegistry(cfg, pgsql)

	var NTPProperty = []string{
		"0.pool.ntp.org",
		"1.pool.ntp.org",
	}

	var HMProperty = health_monitor.Hm{
		DirectorAccount: &health_monitor.DirectorAccount{
			User:     "hm",
			Password: "hm-password",
		},
		ResurrectorEnabled: true,
	}

	var AWSProperty = aws_cpi.Aws{
		AccessKeyId:           cfg.AWSAccessKeyID,
		SecretAccessKey:       cfg.AWSSecretKey,
		DefaultKeyName:        "bosh",
		DefaultSecurityGroups: []string{"bosh"},
		Region:                cfg.AWSRegion,
	}

	var BlobstoreProperty = director.Blobstore{
		Address:  cfg.BoshPrivateIP,
		Port:     25250,
		Provider: "dav",
		Director: &director.Director{
			User:     "director",
			Password: "director-password",
		},
		Agent: &director.Agent{
			User:     "agent",
			Password: "agent-password",
		},
	}

	var NatsProperty = director.Nats{
		Address:  "127.0.0.1",
		User:     "nats",
		Password: "nats-password",
	}

	var AgentProperty = aws_cpi.Agent{
		Mbus: "nats://nats:nats-password@10.0.0.6:4222",
	}

	manifest := &enaml.DeploymentManifest{}
	manifest.SetName(cfg.Name)
	manifest.AddRelease(enaml.Release{
		Name: "bosh",
		URL:  "https://bosh.io/d/github.com/cloudfoundry/bosh?v=" + cfg.BoshReleaseVersion,
		SHA1: cfg.BoshReleaseSHA,
	})

	manifest.AddRelease(enaml.Release{
		Name: "bosh-aws-cpi",
		URL:  "https://bosh.io/d/github.com/cloudfoundry-incubator/bosh-aws-cpi-release?v=" + cfg.BoshCPIReleaseVersion,
		SHA1: cfg.BoshCPIReleaseSHA,
	})

	resourcePool := enaml.ResourcePool{
		Name:    "vms",
		Network: "private",
	}
	resourcePool.Stemcell = enaml.Stemcell{
		URL:  "https://bosh.io/d/stemcells/bosh-aws-xen-hvm-ubuntu-trusty-go_agent?v=" + cfg.GoAgentVersion,
		SHA1: cfg.GoAgentSHA,
	}
	resourcePool.CloudProperties = awscloudproperties.ResourcePool{
		InstanceType: cfg.BoshInstanceSize,
		EphemeralDisk: awscloudproperties.EphemeralDisk{
			Size:     25000,
			DiskType: "gp2",
		},
		AvailabilityZone: cfg.BoshAvailabilityZone,
	}
	manifest.AddResourcePool(resourcePool)
	manifest.AddDiskPool(enaml.DiskPool{
		Name:     "disks",
		DiskSize: 20000,
		CloudProperties: awscloudproperties.RootDisk{
			DiskType: "gp2",
		},
	})
	net := enaml.NewManualNetwork("private")
	net.AddSubnet(enaml.Subnet{
		Range:   "10.0.0.0/24",
		Gateway: "10.0.0.1",
		DNS:     []string{"10.0.0.2"},
		CloudProperties: awscloudproperties.Network{
			Subnet: cfg.AWSSubnet,
		},
	})
	manifest.AddNetwork(net)
	manifest.AddNetwork(enaml.NewVIPNetwork("public"))
	boshJob := &enaml.Job{
		Name:               "bosh",
		Instances:          1,
		ResourcePool:       "vms",
		PersistentDiskPool: "disks",
	}

	for _, v := range []string{"nats", "postgres", "blobstore", "director", "health_monitor", "registry"} {
		boshJob.AddTemplate(enaml.Template{Name: v, Release: "bosh"})
	}
	boshJob.AddTemplate(enaml.Template{Name: "aws_cpi", Release: "bosh-aws-cpi"})

	boshJob.AddNetwork(enaml.Network{
		Name:      "private",
		StaticIPs: []string{cfg.BoshPrivateIP},
		Default:   []interface{}{"dns", "gateway"},
	})

	boshJob.AddNetwork(enaml.Network{
		Name:      "public",
		StaticIPs: []string{cfg.AWSElasticIP},
	})
	boshJob.AddProperty("director", directorProperty)
	boshJob.AddProperty("nats", NatsProperty)
	boshJob.AddProperty("registry", registryProperty)
	boshJob.AddProperty("hm", HMProperty)
	boshJob.AddProperty("ntp", NTPProperty)
	boshJob.AddProperty("agent", AgentProperty)
	boshJob.AddProperty("postgres", pgsql.GetPostgresDB())
	boshJob.AddProperty("blobstore", BlobstoreProperty)
	boshJob.AddProperty("aws", AWSProperty)
	manifest.AddJob(*boshJob)
	manifest.SetCloudProvider(NewAWSCloudProvider(cfg.AWSElasticIP, cfg.AWSPEMFilePath, AWSProperty, NTPProperty))
	return manifest
}
