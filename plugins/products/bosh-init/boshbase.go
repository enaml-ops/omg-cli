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
	"github.com/enaml-ops/pluginlib/pluginutil"
	"github.com/xchapter7x/lo"
)

type BoshPassword struct {
	Password string `yaml:"password"`
}

func (s *BoshBase) InitializeDBDefaults() {
	s.DatabaseDriver = "postgres"
	s.DatabaseUsername = "postgres"
	s.DatabaseScheme = "postgresql"
	s.DirectorDatabaseName = "bosh"
	s.RegistryDatabaseName = "registry"
	s.UAADatabaseName = "uaa"
	s.DatabasePort = 5432
	s.DatabaseHost = "127.0.0.1"
}

func (s *BoshBase) InitializePasswords() {
	s.DirectorPassword = pluginutil.NewPassword(20)
	s.LoginSecret = pluginutil.NewPassword(20)
	s.RegistryPassword = pluginutil.NewPassword(20)
	s.HealthMonitorSecret = pluginutil.NewPassword(20)
	s.MBusPassword = pluginutil.NewPassword(20)
	if s.NatsPassword == "" {
		s.NatsPassword = pluginutil.NewPassword(20)
	}
	if s.DatabasePassword == "" {
		s.DatabasePassword = pluginutil.NewPassword(20)
	}
}

func (s *BoshBase) HandleDeployment(provider IAASManifestProvider, boshInitDeploy func(string)) error {
	var yamlString string
	if manifest, err := provider.CreateDeploymentManifest(); err != nil {
		lo.G.Error(err.Error())
		return err
	} else {
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
	if err := ioutil.WriteFile("./nats.pwd", []byte(s.NatsPassword), 0666); err != nil {
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
	if caCert, cert, key, err = pluginutil.GenerateCert([]string{s.GetRoutableIP()}); err == nil {
		s.SSLCert = cert
		s.SSLKey = key
		s.CACert = caCert
	}
	return
}

//InitializeKeys - initializes public/private keys
func (s *BoshBase) InitializeKeys() (err error) {
	var publicKey, privateKey string
	if publicKey, privateKey, err = pluginutil.GenerateKeys(); err == nil {
		s.PublicKey = publicKey
		s.PrivateKey = privateKey
	}
	return
}

//CreateResourcePool creates the bosh resource pool
func (s *BoshBase) CreateResourcePool(cloudPropertiesFunction func() interface{}) (*enaml.ResourcePool, error) {
	if passwordHash, err := SHA512Pass(s.DirectorPassword); err != nil {
		return nil, err
	} else {
		resourcePool := &enaml.ResourcePool{
			Name:    "vms",
			Network: "private",
		}
		resourcePool.Stemcell = enaml.Stemcell{
			URL:  s.GOAgentReleaseURL,
			SHA1: s.GOAgentSHA,
		}
		resourcePool.CloudProperties = cloudPropertiesFunction()

		resourcePool.Env = map[string]interface{}{
			"bosh": BoshPassword{
				Password: passwordHash,
			},
		}
		return resourcePool, nil
	}
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
	if s.ConfigureBlobstore == nil {
		s.ConfigureBlobstore = s.configureBlobstore
	}
	boshJob := &enaml.Job{
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

	if !s.UseExternalDB {
		boshJob.AddTemplate(enaml.Template{Name: "postgres", Release: "bosh"})
		boshJob.AddProperty("postgres", s.createPostgresJobProperties())
	}
	boshJob.AddTemplate(enaml.Template{Name: "registry", Release: "bosh"})
	boshJob.AddProperty("registry", s.createRegistryJobProperties())

	boshJob.AddTemplate(enaml.Template{Name: "director", Release: "bosh"})
	if s.IsUAA() {
		boshJob.AddProperty("director", s.createDirectorProperties(s.sslFunction, s.uaaUserManagement))
	} else {
		boshJob.AddProperty("director", s.createDirectorProperties(s.noSSlFunction, s.localUserManagement))
	}
	s.ConfigureBlobstore(boshJob)
	boshJob.AddProperty("ntp", s.NtpServers)
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
	return *boshJob
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

func (s *BoshBase) configureBlobstore(boshJob *enaml.Job) {
	boshJob.AddTemplate(enaml.Template{Name: "blobstore", Release: "bosh"})
	boshJob.AddProperty("blobstore", &director.Blobstore{
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
	})
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
			User:     s.DatabaseUsername,
			Password: s.DatabasePassword,
			Port:     s.DatabasePort,
			Adapter:  s.DatabaseDriver,
			Database: s.RegistryDatabaseName,
			Host:     s.DatabaseHost,
		},
	}
}

func (s *BoshBase) createPostgresJobProperties() *postgres.Postgres {
	return &postgres.Postgres{
		ListenAddress:       s.DatabaseHost,
		User:                s.DatabaseUsername,
		Password:            s.DatabasePassword,
		Database:            s.DirectorDatabaseName,
		AdditionalDatabases: []string{s.UAADatabaseName, s.RegistryDatabaseName},
	}
}

func (s *BoshBase) createUAADBProperties() *uaa.Uaadb {
	return &uaa.Uaadb{
		Address:  s.DatabaseHost,
		DbScheme: s.DatabaseScheme,
		Port:     s.DatabasePort,
		Databases: []interface{}{
			map[string]string{
				"name": s.UAADatabaseName,
				"tag":  "uaa",
			},
		},
		Roles: []interface{}{
			map[string]string{
				"name":     s.DatabaseUsername,
				"password": s.DatabasePassword,
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
			"director": UAAClient{
				AuthorizedGrantTypes: "client_credentials",
				Override:             true,
				Scope:                "",
				Authorities:          "bosh.admin",
				RefreshTokenValidity: 86400,
				AccessTokenValidity:  600,
				Secret:               s.DirectorPassword,
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

func (s *BoshBase) uaaUserManagement() *director.UserManagement {
	return &director.UserManagement{
		Provider: "uaa",
		Uaa: &director.Uaa{
			PublicKey: s.PublicKey,
			Url:       fmt.Sprintf("https://%s:8443", s.GetRoutableIP()),
		},
	}
}

func (s *BoshBase) localUserManagement() *director.UserManagement {
	return &director.UserManagement{
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
	}
}

func (s *BoshBase) sslFunction() *director.Ssl {
	return &director.Ssl{
		Cert: s.SSLCert,
		Key:  s.SSLKey,
	}
}
func (s *BoshBase) noSSlFunction() *director.Ssl {
	return nil
}

func (s *BoshBase) createDirectorProperties(sslFunction func() *director.Ssl, userManagementFunction func() *director.UserManagement) *director.Director {
	return &director.Director{
		Address:      s.GetRoutableIP(),
		Name:         s.DirectorName,
		CpiJob:       s.CPIJobName,
		MaxThreads:   10,
		TrustedCerts: s.TrustedCerts,
		Ssl:          sslFunction(),
		Db: &director.DirectorDb{
			User:     s.DatabaseUsername,
			Password: s.DatabasePassword,
			Adapter:  s.DatabaseDriver,
			Port:     s.DatabasePort,
			Host:     s.DatabaseHost,
			Database: s.DirectorDatabaseName,
		},
		UserManagement:      userManagementFunction(),
		GenerateVmPasswords: true,
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
