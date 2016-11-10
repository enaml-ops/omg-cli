package boshinit

import (
	"github.com/enaml-ops/enaml"
	"github.com/enaml-ops/omg-cli/plugins/products/bosh-init/enaml-gen/aws_cpi"
	"github.com/enaml-ops/omg-cli/plugins/products/bosh-init/enaml-gen/director"
	"github.com/enaml-ops/omg-cli/plugins/products/bosh-init/enaml-gen/postgres"
	"github.com/enaml-ops/omg-cli/plugins/products/bosh-init/enaml-gen/registry"
)

type Rr registry.Registry
type Ar aws_cpi.Registry

type RegistryProperty struct {
	Rr      `yaml:",inline"`
	Ar      `yaml:",inline"`
	Address string `yaml:"address"`
}
type user struct {
	Name     string `yaml:"name"`
	Password string `yaml:"password"`
}

type DirectorProperty struct {
	director.Director `yaml:",inline"`
	Address           string
}

type PgSql struct {
	User     string
	Host     string
	Password string
	Database string
	Adapter  string
}

type IAASManifestProvider interface {
	CreateCPIRelease() enaml.Release
	CreateCPITemplate() enaml.Template
	CreateDiskPool() enaml.DiskPool
	CreateResourcePool() (*enaml.ResourcePool, error)
	CreateManualNetwork() enaml.ManualNetwork
	CreateVIPNetwork() enaml.VIPNetwork
	CreateJobNetwork() *enaml.Network
	CreateCloudProvider() enaml.CloudProvider
	CreateCPIJobProperties() map[string]interface{}
	CreateDeploymentManifest() (*enaml.DeploymentManifest, error)
}

type Postgres interface {
	GetDirectorDB() *director.DirectorDb
	GetRegistryDB() *registry.Db
	GetPostgresDB() postgres.Postgres
}

type BoshBase struct {
	Mode                 string
	NetworkCIDR          string
	NetworkGateway       string
	NetworkDNS           []string
	DirectorName         string
	DirectorPassword     string
	DBPassword           string
	CPIJobName           string
	NtpServers           []string
	PrivateStaticIPs     []string
	PrivateReservedRange string
	NatsPassword         string
	MBusPassword         string
	PrivateIP            string
	PublicIP             string
	SSLCert              string
	SSLKey               string
	PrivateKey           string
	PublicKey            string
	HealthMonitorSecret  string
	LoginSecret          string
	RegistryPassword     string
	CACert               string
	BoshReleaseSHA       string
	BoshReleaseURL       string
	CPIReleaseSHA        string
	CPIReleaseURL        string
	GOAgentSHA           string
	GOAgentReleaseURL    string
	UAAReleaseSHA        string
	UAAReleaseURL        string
	TrustedCerts         string
	PersistentDiskSize   int
	PrintManifest        bool
	GraphiteAddress      string
	GraphitePort         int
	SyslogAddress        string
	SyslogPort           int
	SyslogTransport      string
}

//UAAClient - Structure to represent map of client priviledges
type UAAClient struct {
	ID                   string      `yaml:"id,omitempty"`
	Secret               string      `yaml:"secret,omitempty"`
	Scope                string      `yaml:"scope,omitempty"`
	AuthorizedGrantTypes string      `yaml:"authorized-grant-types,omitempty"`
	Authorities          string      `yaml:"authorities,omitempty"`
	AutoApprove          interface{} `yaml:"autoapprove,omitempty"`
	Override             bool        `yaml:"override,omitempty"`
	RedirectURI          string      `yaml:"redirect-uri,omitempty"`
	AccessTokenValidity  int         `yaml:"access-token-validity,omitempty"`
	RefreshTokenValidity int         `yaml:"refresh-token-validity,omitempty"`
	ResourceIDs          string      `yaml:"resource_ids,omitempty"`
	Name                 string      `yaml:"name,omitempty"`
	AppLaunchURL         string      `yaml:"app-launch-url,omitempty"`
	ShowOnHomepage       bool        `yaml:"show-on-homepage,omitempty"`
	AppIcon              string      `yaml:"app-icon,omitempty"`
}
