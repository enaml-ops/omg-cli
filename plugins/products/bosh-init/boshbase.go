package boshinit

import (
	"fmt"

	"github.com/enaml-ops/enaml"
	"github.com/enaml-ops/omg-cli/plugins/products/bosh-init/enaml-gen/blobstore"
	"github.com/enaml-ops/omg-cli/plugins/products/bosh-init/enaml-gen/director"
	"github.com/enaml-ops/omg-cli/plugins/products/bosh-init/enaml-gen/health_monitor"
	"github.com/enaml-ops/omg-cli/plugins/products/bosh-init/enaml-gen/nats"
	"github.com/enaml-ops/omg-cli/plugins/products/bosh-init/enaml-gen/postgres"
	"github.com/enaml-ops/omg-cli/plugins/products/bosh-init/enaml-gen/uaa"
)

func NewBoshDeploymentBase(cfg BoshInitConfig, cpiname string, ntpProperty []string) *enaml.DeploymentManifest {
	var pgsql = NewPostgres("postgres", "127.0.0.1", "postgres-password", "bosh", "postgres")

	manifest := &enaml.DeploymentManifest{}
	manifest.SetName(cfg.Name)
	manifest.AddRelease(enaml.Release{
		Name: "bosh",
		URL:  "https://bosh.io/d/github.com/cloudfoundry/bosh?v=" + cfg.BoshReleaseVersion,
		SHA1: cfg.BoshReleaseSHA,
	})
	manifest.AddRelease(enaml.Release{
		Name: "uaa",
		URL:  "https://bosh.io/d/github.com/cloudfoundry/uaa-release?v=" + cfg.UAAReleaseVersion,
		SHA1: cfg.UAAReleaseSHA,
	})
	directorPassword := "password"
	boshJob := &enaml.Job{
		Name:               "bosh",
		Instances:          1,
		ResourcePool:       "vms",
		PersistentDiskPool: "disks",
		Properties: []interface{}{
			&director.DirectorJob{
				Director: NewDirector(cfg.BoshDirectorName, cpiname, pgsql.GetDirectorDB()),
				Ntp:      ntpProperty,
			},
			&postgres.PostgresJob{
				Postgres: &postgres.Postgres{
					ListenAddress:       pgsql.Host,
					User:                pgsql.User,
					Password:            pgsql.Password,
					Database:            pgsql.GetDirectorDB().Database,
					AdditionalDatabases: []string{"uaa"},
				},
			},
			&nats.NatsJob{
				Nats: &nats.Nats{
					User:          "nats",
					Password:      "nats-password",
					ListenAddress: "127.0.0.1",
				},
			},
			&blobstore.BlobstoreJob{
				Blobstore: &blobstore.Blobstore{
					Port: 25250,
					Director: &blobstore.Director{
						User:     "director",
						Password: "director-password",
					},
					Agent: &blobstore.Agent{
						User:     "agent",
						Password: "agent-password",
					},
				},
			},
			&health_monitor.HealthMonitorJob{
				Hm: &health_monitor.Hm{
					DirectorAccount: &health_monitor.DirectorAccount{
						CaCert:       "TODO",
						ClientId:     "health_monitor",
						ClientSecret: "health-monitor-uaa-secret",
					},
				},
			},
			&uaa.UaaJob{
				Login: &uaa.Login{
					Protocol: "https",
					Branding: &uaa.Branding{},
					Saml: &uaa.Saml{
						ServiceProviderKey:         "TODO",
						ServiceProviderCertificate: "TODO",
					},
				},
				Uaa: &uaa.Uaa{
					Admin: &uaa.Admin{
						ClientSecret: directorPassword,
					},
					DisableInternalAuth: false,
					SslCertificate:      "TODO",
					SslPrivateKey:       "TODO",
					RequireHttps:        true,
					Url:                 "https://<director-ip>:8443",
					Jwt: &uaa.Jwt{
						SigningKey:      "TODO",
						VerificationKey: "TODO",
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
					Clients: nil,
					Scim: &uaa.Scim{
						Users: []string{
							"director|bosh-director-scim-password|bosh.admin",
							fmt.Sprintf("admin|%s|bosh.admin,scim.write,clients.write,scim.read,clients.read", directorPassword),
						},
					},
				},
				Uaadb: &uaa.Uaadb{
					Address:  pgsql.Host,
					DbScheme: "postgresql",
					Port:     5432,
					Databases: map[string]string{
						"name": "uaa",
						"tag":  "uaa",
					},
					Roles: map[string]string{
						"name":     "postgres",
						"password": pgsql.Password,
						"tag":      "admin",
					},
				},
			},
		},
	}

	for _, v := range []string{"nats", "postgres", "blobstore", "director", "health_monitor"} {
		boshJob.AddTemplate(enaml.Template{Name: v, Release: "bosh"})
	}
	for _, v := range []string{"uaa", "uaa_postgres"} {
		boshJob.AddTemplate(enaml.Template{Name: v, Release: "uaa"})
	}

	boshJob.AddNetwork(enaml.Network{
		Name:      "private",
		StaticIPs: []string{cfg.BoshPrivateIP},
		Default:   []interface{}{"dns", "gateway"},
	})

	manifest.AddJob(*boshJob)
	return manifest
}
