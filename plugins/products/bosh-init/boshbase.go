package boshinit

import (
	"fmt"

	"github.com/enaml-ops/enaml"
	"github.com/enaml-ops/omg-cli/plugins/products/bosh-init/enaml-gen/blobstore"
	"github.com/enaml-ops/omg-cli/plugins/products/bosh-init/enaml-gen/director"
	"github.com/enaml-ops/omg-cli/plugins/products/bosh-init/enaml-gen/health_monitor"
	"github.com/enaml-ops/omg-cli/plugins/products/bosh-init/enaml-gen/nats"
	"github.com/enaml-ops/omg-cli/plugins/products/bosh-init/enaml-gen/postgres"
	"github.com/enaml-ops/omg-cli/plugins/products/bosh-init/enaml-gen/registry"
	"github.com/enaml-ops/omg-cli/plugins/products/bosh-init/enaml-gen/uaa"
)

const (
	directorDBName = "bosh"
	dbUser         = "postgres"
	dbAdapter      = "postgres"
	dbHost         = "127.0.0.1"
	dbPort         = 5432
)

func (s *BoshBase) CreateDeploymentManifest() *enaml.DeploymentManifest {
	manifest := &enaml.DeploymentManifest{}
	manifest.SetName(s.DirectorName)
	manifest.AddRelease(enaml.Release{
		Name: "bosh",
		URL:  "https://bosh.io/d/github.com/cloudfoundry/bosh?v=" + s.BoshReleaseVersion,
		SHA1: s.BoshReleaseSHA,
	})
	manifest.AddRelease(enaml.Release{
		Name: "uaa",
		URL:  "https://bosh.io/d/github.com/cloudfoundry/uaa-release?v=" + s.UAAReleaseVersion,
		SHA1: s.UAAReleaseSHA,
	})
	manifest.AddJob(s.CreateJob())
	return manifest
}

func (s *BoshBase) CreateJob() enaml.Job {
	boshJob := enaml.Job{
		Name:               "bosh",
		Instances:          1,
		ResourcePool:       "vms",
		PersistentDiskPool: "disks",
	}
	boshJob.AddTemplate(enaml.Template{Name: "uaa", Release: "uaa"})
	boshJob.AddProperty(s.createUAAProperties())

	boshJob.AddTemplate(enaml.Template{Name: "nats", Release: "bosh"})
	boshJob.AddProperty(s.createNatsJobProperties())

	boshJob.AddTemplate(enaml.Template{Name: "postgres", Release: "bosh"})
	boshJob.AddProperty(s.createPostgresJobProperties())

	boshJob.AddTemplate(enaml.Template{Name: "registry", Release: "bosh"})
	boshJob.AddProperty(s.createRegistryJobProperties())

	boshJob.AddTemplate(enaml.Template{Name: "director", Release: "bosh"})
	boshJob.AddProperty(s.createDirectorProperties())

	boshJob.AddTemplate(enaml.Template{Name: "blobstore", Release: "bosh"})
	boshJob.AddProperty(s.createBlobStoreJobProperties())

	boshJob.AddTemplate(enaml.Template{Name: "health_monitor", Release: "bosh"})
	boshJob.AddProperty(s.createHeathMonitoryJobProperties())

	boshJob.AddNetwork(enaml.Network{
		Name:      "private",
		StaticIPs: []string{s.PrivateIP},
	})
	return boshJob
}

func (s *BoshBase) createHeathMonitoryJobProperties() *health_monitor.HealthMonitorJob {
	return &health_monitor.HealthMonitorJob{
		Hm: &health_monitor.Hm{
			DirectorAccount: &health_monitor.DirectorAccount{
				CaCert:       s.CACert,
				ClientId:     "health_monitor",
				ClientSecret: s.HealthMonitorSecret,
			},
			ResurrectorEnabled: true,
			Resurrector:        &health_monitor.Resurrector{},
		},
	}
}
func (s *BoshBase) createBlobStoreJobProperties() *blobstore.BlobstoreJob {
	return &blobstore.BlobstoreJob{
		Blobstore: &blobstore.Blobstore{
			Port: 25250,
			Director: &blobstore.Director{
				User:     "director",
				Password: s.DirectorPassword,
			},
			Agent: &blobstore.Agent{
				User:     "agent",
				Password: s.AgentPassword,
			},
		},
	}
}

func (s *BoshBase) createRegistryJobProperties() *registry.RegistryJob {
	return &registry.RegistryJob{
		Registry: &registry.Registry{
			Http: &registry.Http{
				User:     "admin",
				Password: s.RegistryPassword,
				Port:     25777,
			},
			Db: &registry.Db{
				User:     dbUser,
				Password: s.DBPassword,
				Port:     dbPort,
				Adapter:  dbAdapter,
				Database: "registry",
			},
		},
	}
}

func (s *BoshBase) createPostgresJobProperties() *postgres.PostgresJob {
	return &postgres.PostgresJob{
		Postgres: &postgres.Postgres{
			ListenAddress:       dbHost,
			User:                dbUser,
			Password:            s.DBPassword,
			Database:            directorDBName,
			AdditionalDatabases: []string{"uaa", "registry"},
		},
	}
}

func (s *BoshBase) createUAAProperties() *uaa.UaaJob {
	return &uaa.UaaJob{
		Login: &uaa.Login{
			Protocol: "https",
			Saml: &uaa.Saml{
				ServiceProviderKey:         s.SSLKey,
				ServiceProviderCertificate: s.SSLCert,
			},
		},
		Uaa: &uaa.Uaa{
			Admin: &uaa.Admin{
				ClientSecret: s.DirectorPassword,
			},
			DisableInternalAuth: false,
			SslCertificate:      s.SSLCert,
			SslPrivateKey:       s.SSLKey,
			RequireHttps:        true,
			Url:                 fmt.Sprintf("https://%s:8443", s.PublicIP),
			Jwt: &uaa.Jwt{
				SigningKey:      s.SigningKey,
				VerificationKey: s.VerificationKey,
			},
			User: &uaa.UaaUser{
				Authorities: []string{
					"openid",
					"scim.me",
					"password.write",
					"uaa.user",
					"profile",
					"roles",
					"user_attributes",
					"bosh.admin",
					"bosh.read",
					"bosh.*.admin",
					"bosh.*.read",
					"clients.admin"},
			},
			Clients: map[string]UAAClient{
				"bosh_cli": UAAClient{
					Override:             true,
					AuthorizedGrantTypes: "password,refresh_token",
					Scope:                "openid,bosh.admin,bosh.read,bosh.*.admin,bosh.*.read",
					Authorities:          "uaa.none",
					AccessTokenValidity:  120,   // 2 minutes
					RefreshTokenValidity: 86400, //re-login required once a day
					Secret:               "",    //CLI expects this secret to be empty
				},
				"health_monitor": UAAClient{
					AuthorizedGrantTypes: "client_credentials",
					Override:             true,
					Scope:                "",
					Authorities:          "bosh.admin",
					RefreshTokenValidity: 86400,
					AccessTokenValidity:  600,
					Secret:               s.HealthMonitorSecret,
				},
				"login": UAAClient{
					AuthorizedGrantTypes: "password,authorization_code",
					AutoApprove:          true,
					Override:             true,
					Scope:                "bosh.admin,scim.write,scim.read,clients.admin",
					Authorities:          "",
					RefreshTokenValidity: 86400,
					AccessTokenValidity:  600,
					Secret:               s.LoginSecret,
				},
			},
			Scim: &uaa.Scim{
				Users: []string{
					"director|bosh-director-scim-password|bosh.admin",
					fmt.Sprintf("admin|%s|bosh.admin,scim.write,clients.write,scim.read,clients.read", s.DirectorPassword),
				},
			},
		},
		Uaadb: &uaa.Uaadb{
			Address:  dbHost,
			DbScheme: "postgresql",
			Port:     dbPort,
			Databases: map[string]string{
				"name": "uaa",
				"tag":  "uaa",
			},
			Roles: map[string]string{
				"name":     dbUser,
				"password": s.DBPassword,
				"tag":      "admin",
			},
		},
	}
}

func (s *BoshBase) createDirectorProperties() *director.DirectorJob {
	return &director.DirectorJob{
		Director: &director.Director{
			Name:       s.DirectorName,
			CpiJob:     s.CPIName,
			MaxThreads: 10,
			Db: &director.DirectorDb{
				User:     dbUser,
				Password: s.DBPassword,
				Adapter:  dbAdapter,
				Port:     dbPort,
				Host:     dbHost,
			},
			Ssl: &director.Ssl{
				Cert: s.SSLCert,
				Key:  s.SSLKey,
			},
			UserManagement: &director.UserManagement{
				Provider: "uaa",
				Uaa: &director.Uaa{
					PublicKey: s.UAAPublicKey,
					Url:       fmt.Sprintf("https://%s:8443", s.PublicIP),
				},
			},
		},
		Ntp: s.NtpServers,
	}
}

func (s *BoshBase) createNatsJobProperties() *nats.NatsJob {
	return &nats.NatsJob{
		Nats: &nats.Nats{
			User:          "nats",
			Password:      s.NatsPassword,
			ListenAddress: "127.0.0.1",
		},
	}
}
