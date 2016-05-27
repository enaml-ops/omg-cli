package boshinit

import (
	"github.com/enaml-ops/omg-cli/plugins/deployments/bosh-init/enaml-gen/director"
	"github.com/enaml-ops/enaml"
)

func NewBoshDeploymentBase(cfg BoshInitConfig, cpiname string, ntpProperty []string) *enaml.DeploymentManifest {
	var pgsql = NewPostgres("postgres", "127.0.0.1", "postgres-password", "bosh", "postgres")

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

	manifest := &enaml.DeploymentManifest{}
	manifest.SetName(cfg.Name)
	manifest.AddRelease(enaml.Release{
		Name: "bosh",
		URL:  "https://bosh.io/d/github.com/cloudfoundry/bosh?v=" + cfg.BoshReleaseVersion,
		SHA1: cfg.BoshReleaseSHA,
	})

	boshJob := &enaml.Job{
		Name:               "bosh",
		Instances:          1,
		ResourcePool:       "vms",
		PersistentDiskPool: "disks",
	}

	for _, v := range []string{"nats", "postgres", "blobstore", "director", "health_monitor", "registry"} {
		boshJob.AddTemplate(enaml.Template{Name: v, Release: "bosh"})
	}

	boshJob.AddNetwork(enaml.Network{
		Name:      "private",
		StaticIPs: []string{cfg.BoshPrivateIP},
		Default:   []interface{}{"dns", "gateway"},
	})

	boshJob.AddProperty("director", NewDirectorProperty(cfg.BoshDirectorName, cpiname, pgsql.GetDirectorDB()))
	boshJob.AddProperty("nats", NewNats("nats", "nats-password"))
	boshJob.AddProperty("registry", GetRegistry(cfg, pgsql))
	boshJob.AddProperty("hm", NewHealthMonitor(true))
	boshJob.AddProperty("ntp", ntpProperty)
	boshJob.AddProperty("postgres", pgsql.GetPostgresDB())
	boshJob.AddProperty("blobstore", BlobstoreProperty)
	manifest.AddJob(*boshJob)
	return manifest
}
