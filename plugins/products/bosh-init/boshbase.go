package boshinit

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/enaml-ops/enaml"
	"github.com/enaml-ops/omg-cli/plugins/products/bosh-init/enaml-gen/director"
	"github.com/enaml-ops/omg-cli/plugins/products/bosh-init/enaml-gen/health_monitor"
	"github.com/enaml-ops/omg-cli/plugins/products/bosh-init/enaml-gen/postgres"
	"github.com/enaml-ops/omg-cli/plugins/products/bosh-init/enaml-gen/registry"
	"github.com/enaml-ops/omg-cli/plugins/products/bosh-init/enaml-gen/uaa"
	"github.com/enaml-ops/omg-cli/utils"
	"github.com/xchapter7x/lo"
)

const (
	directorDBName = "bosh"
	dbUser         = "postgres"
	dbAdapter      = "postgres"
	dbHost         = "127.0.0.1"
	dbPort         = 5432
)

func (s *BoshBase) InitializePasswords() {
	s.DirectorPassword = utils.NewPassword(20)
	s.LoginSecret = utils.NewPassword(20)
	s.RegistryPassword = utils.NewPassword(20)
	s.HealthMonitorSecret = utils.NewPassword(20)
	s.DBPassword = utils.NewPassword(20)
	if s.NatsPassword == "" {
		s.NatsPassword = utils.NewPassword(20)
	}
	s.MBusPassword = utils.NewPassword(20)
}

func (s *BoshBase) HandleDeployment(provider IAASManifestProvider, boshInitDeploy func(string)) error {
	var yamlString string
	var err error
	manifest := provider.CreateDeploymentManifest()

	lo.G.Debug("Got manifest", manifest)
	if yamlString, err = enaml.Paint(manifest); err != nil {
		lo.G.Error(err.Error())
		return err
	}

	if err = s.CreateAuthenticationFiles(); err != nil {
		return err
	}

	if s.PrintManifest {
		fmt.Println(yamlString)
		return nil
	}

	if err = s.deployYaml(yamlString, boshInitDeploy); err != nil {
		lo.G.Error(err.Error())
		return err
	}

	return nil
}

func (s *BoshBase) deployYaml(myYaml string, boshInitDeploy func(string)) error {
	var err error
	var tmpfile *os.File
	content := []byte(myYaml)
	boshdeploypath := utils.GetBoshDeployPath()
	os.Remove(boshdeploypath)
	if tmpfile, err = os.Create(boshdeploypath); err != nil {
		lo.G.Error(err.Error())
		return err
	}
	defer tmpfile.Close()
	defer os.Remove(tmpfile.Name())
	if _, err = tmpfile.Write(content); err != nil {
		lo.G.Error(err.Error())
		return err
	}
	if err = tmpfile.Close(); err != nil {
		lo.G.Error(err.Error())
		return err
	}
	boshInitDeploy(tmpfile.Name())
	return nil
}
func (s *BoshBase) CreateAuthenticationFiles() error {
	if s.CACert != "" {
		if err := ioutil.WriteFile("./rootCA.pem", []byte(s.CACert), 0666); err != nil {
			lo.G.Error(err.Error())
			return err
		}
	}
	if err := ioutil.WriteFile("./director.pwd", []byte(s.DirectorPassword), 0666); err != nil {
		lo.G.Error(err.Error())
		return err
	}
	return nil
}

//IsBasic - is this a basic Bosh director
func (s *BoshBase) IsBasic() bool {
	return strings.ToUpper(s.Mode) == "BASIC"
}

//IsUAA - is this a UAA enabled bosh director
func (s *BoshBase) IsUAA() bool {
	return strings.ToUpper(s.Mode) == "UAA"
}

//InitializeCerts - initializes certs needed for UAA and health monitor
func (s *BoshBase) InitializeCerts() (err error) {
	var cert, key, caCert string
	if caCert, cert, key, err = utils.GenerateCert([]string{s.GetRoutableIP()}); err == nil {
		s.SSLCert = cert
		s.SSLKey = key
		s.CACert = caCert
	}
	return
}

//InitializeKeys - initializes public/private keys
func (s *BoshBase) InitializeKeys() (err error) {
	var publicKey, privateKey string
	if publicKey, privateKey, err = utils.GenerateKeys(); err == nil {
		s.PublicKey = publicKey
		s.PrivateKey = privateKey
	}
	return
}

func (s *BoshBase) CreateDeploymentManifest() *enaml.DeploymentManifest {
	manifest := &enaml.DeploymentManifest{}
	manifest.SetName(s.DirectorName)
	manifest.AddRelease(enaml.Release{
		Name: "bosh",
		URL:  s.BoshReleaseURL,
		SHA1: s.BoshReleaseSHA,
	})

	if s.IsUAA() {
		manifest.AddRelease(enaml.Release{
			Name: "uaa",
			URL:  s.UAAReleaseURL,
			SHA1: s.UAAReleaseSHA,
		})
	}
	manifest.AddJob(s.CreateJob())
	return manifest
}

func (s *BoshBase) CreateJob() enaml.Job {
	boshJob := enaml.Job{
		Name:               "bosh",
		Instances:          1,
		ResourcePool:       "vms",
		PersistentDiskPool: "disks",
		Properties:         enaml.Properties{},
	}
	if s.IsUAA() {
		boshJob.AddTemplate(enaml.Template{Name: "uaa", Release: "uaa"})
		boshJob.AddProperty("uaa", s.createUAAProperties())
		boshJob.AddProperty("uaadb", s.createUAADBProperties())
		boshJob.AddProperty("login", s.createUAALoginProperties())
	}
	boshJob.AddTemplate(enaml.Template{Name: "nats", Release: "bosh"})
	boshJob.AddProperty("nats", s.createNatsJobProperties())

	boshJob.AddTemplate(enaml.Template{Name: "postgres", Release: "bosh"})
	boshJob.AddProperty("postgres", s.createPostgresJobProperties())

	boshJob.AddTemplate(enaml.Template{Name: "registry", Release: "bosh"})
	boshJob.AddProperty("registry", s.createRegistryJobProperties())

	boshJob.AddTemplate(enaml.Template{Name: "director", Release: "bosh"})
	if s.IsUAA() {
		boshJob.AddProperty("director", s.createDirectorUAAProperties())
	} else {
		boshJob.AddProperty("director", s.createDirectorProperties())
	}
	boshJob.AddProperty("ntp", s.NtpServers)

	boshJob.AddTemplate(enaml.Template{Name: "blobstore", Release: "bosh"})
	boshJob.AddProperty("blobstore", s.createBlobStoreJobProperties())

	boshJob.AddTemplate(enaml.Template{Name: "health_monitor", Release: "bosh"})
	hm := s.createHealthMonitorJobProperties()
	if s.IsUAA() {
		s.addHealthMonitorUAA(hm)
	} else {
		s.addHealthMonitorBasicAuth(hm)
	}
	boshJob.AddProperty("hm", hm)

	staticIPs := append(s.PrivateStaticIPs, s.PrivateIP)
	boshJob.AddNetwork(enaml.Network{
		Name:      "private",
		StaticIPs: staticIPs,
		Default:   []interface{}{"dns", "gateway"},
	})
	return boshJob
}

func (s *BoshBase) createHealthMonitorJobProperties() *health_monitor.Hm {
	hm := &health_monitor.Hm{
		ResurrectorEnabled: true,
		Resurrector:        &health_monitor.Resurrector{},
	}
	if s.GraphiteAddress != "" {
		hm.GraphiteEnabled = true
		hm.Graphite = &health_monitor.Graphite{
			Address: s.GraphiteAddress,
			Port:    s.GraphitePort,
		}
	}
	if s.SyslogAddress != "" {
		hm.SyslogEventForwarderEnabled = true
		hm.SyslogEventForwarder = &health_monitor.SyslogEventForwarder{
			Address:   s.SyslogAddress,
			Port:      s.SyslogPort,
			Transport: s.SyslogTransport,
		}
	}
	return hm
}

func (s *BoshBase) addHealthMonitorUAA(hm *health_monitor.Hm) {
	hm.DirectorAccount = &health_monitor.DirectorAccount{
		CaCert:       s.CACert,
		ClientId:     "health_monitor",
		ClientSecret: s.HealthMonitorSecret,
	}
}

func (s *BoshBase) addHealthMonitorBasicAuth(hm *health_monitor.Hm) {
	hm.DirectorAccount = &health_monitor.DirectorAccount{
		User:     "hm",
		Password: s.HealthMonitorSecret,
	}
}

func (s *BoshBase) createBlobStoreJobProperties() *director.Blobstore {
	return &director.Blobstore{
		Port:    25250,
		Address: s.PrivateIP,
		Director: &director.BlobstoreDirector{
			User:     "director",
			Password: s.DirectorPassword,
		},
		Agent: &director.BlobstoreAgent{
			User:     "agent",
			Password: s.NatsPassword,
		},
	}
}

func (s *BoshBase) createRegistryJobProperties() *registry.Registry {
	return &registry.Registry{
		Username: "admin",
		Password: s.RegistryPassword,
		Host:     s.PrivateIP,
		Address:  s.PrivateIP,
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
	}
}

func (s *BoshBase) createPostgresJobProperties() *postgres.Postgres {
	return &postgres.Postgres{
		ListenAddress:       dbHost,
		User:                dbUser,
		Password:            s.DBPassword,
		Database:            directorDBName,
		AdditionalDatabases: []string{"uaa", "registry"},
	}
}

func (s *BoshBase) createUAADBProperties() *uaa.Uaadb {
	return &uaa.Uaadb{
		Address:  dbHost,
		DbScheme: "postgresql",
		Port:     dbPort,
		Databases: []interface{}{
			map[string]string{
				"name": "uaa",
				"tag":  "uaa",
			},
		},
		Roles: []interface{}{
			map[string]string{
				"name":     dbUser,
				"password": s.DBPassword,
				"tag":      "admin",
			},
		},
	}
}
func (s *BoshBase) createUAALoginProperties() *uaa.Login {
	return &uaa.Login{
		Protocol: "https",
		Saml: &uaa.Saml{
			ServiceProviderKey:         s.SSLKey,
			ServiceProviderCertificate: s.SSLCert,
		},
	}
}

func (s *BoshBase) createUAAProperties() *uaa.Uaa {

	return &uaa.Uaa{
		Admin: &uaa.Admin{
			ClientSecret: s.DirectorPassword,
		},
		DisableInternalAuth: false,
		SslCertificate:      s.SSLCert,
		SslPrivateKey:       s.SSLKey,
		RequireHttps:        true,
		Url:                 fmt.Sprintf("https://%s:8443", s.GetRoutableIP()),
		Jwt: &uaa.Jwt{
			SigningKey:      s.PrivateKey,
			VerificationKey: s.PublicKey,
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
				fmt.Sprintf("director|%s|bosh.admin", s.DirectorPassword),
				fmt.Sprintf("admin|%s|bosh.admin,scim.write,clients.write,scim.read,clients.read", s.DirectorPassword),
			},
		},
	}
}

func (s *BoshBase) createDirectorUAAProperties() *director.Director {
	return &director.Director{
		Address:      s.GetRoutableIP(),
		Name:         s.DirectorName,
		CpiJob:       s.CPIJobName,
		MaxThreads:   10,
		TrustedCerts: s.TrustedCerts,
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
				PublicKey: s.PublicKey,
				Url:       fmt.Sprintf("https://%s:8443", s.GetRoutableIP()),
			},
		},
	}
}
func (s *BoshBase) createDirectorProperties() *director.Director {
	return &director.Director{
		Address:      s.GetRoutableIP(),
		Name:         s.DirectorName,
		CpiJob:       s.CPIJobName,
		MaxThreads:   10,
		TrustedCerts: s.TrustedCerts,
		Db: &director.DirectorDb{
			User:     dbUser,
			Password: s.DBPassword,
			Adapter:  dbAdapter,
			Port:     dbPort,
			Host:     dbHost,
		},
		UserManagement: &director.UserManagement{
			Provider: "local",
			Local: &director.Local{
				Users: []user{
					user{
						Name:     "director",
						Password: s.DirectorPassword,
					},
					user{
						Name:     "hm",
						Password: s.HealthMonitorSecret,
					},
				},
			},
		},
	}
}

func (s *BoshBase) createNatsJobProperties() *director.Nats {
	return &director.Nats{
		User:     "nats",
		Password: s.NatsPassword,
		Address:  s.PrivateIP,
	}
}

func (s *BoshBase) GetRoutableIP() string {
	if s.PublicIP != "" {
		return s.PublicIP
	}
	return s.PrivateIP
}
