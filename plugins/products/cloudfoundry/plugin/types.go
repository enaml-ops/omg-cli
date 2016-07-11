package cloudfoundry

import (
	"github.com/codegangsta/cli"
	"github.com/enaml-ops/enaml"
	mysqlproxylib "github.com/enaml-ops/omg-cli/plugins/products/cf-mysql/enaml-gen/proxy"
	etcdmetricslib "github.com/enaml-ops/omg-cli/plugins/products/cloudfoundry/enaml-gen/etcd_metrics_server"
	grtrlib "github.com/enaml-ops/omg-cli/plugins/products/cloudfoundry/enaml-gen/gorouter"
	"github.com/enaml-ops/omg-cli/plugins/products/cloudfoundry/enaml-gen/metron_agent"
	natslib "github.com/enaml-ops/omg-cli/plugins/products/cloudfoundry/enaml-gen/nats"
	routereglib "github.com/enaml-ops/omg-cli/plugins/products/cloudfoundry/enaml-gen/route_registrar"
	"github.com/enaml-ops/omg-cli/plugins/products/cloudfoundry/enaml-gen/uaa"
)

// InstanceGrouper creates and validates InstanceGroups.
type InstanceGrouper interface {
	ToInstanceGroup() (ig *enaml.InstanceGroup)
	HasValidValues() bool
}

// InstanceGrouperFactory is a function that creates InstanceGroupers from CLI args.
type InstanceGrouperFactory func(*cli.Context) InstanceGrouper

type acceptanceTests struct {
	AZs                      []string
	StemcellName             string
	NetworkName              string
	AppsDomain               []string
	SystemDomain             string
	AdminPassword            string
	SkipCertVerify           bool
	IncludeInternetDependent bool
}

type bootstrap struct {
	AZs           []string
	StemcellName  string
	NetworkName   string
	MySQLIPs      []string
	MySQLUser     string
	MySQLPassword string
}

type clockGlobal struct {
	Instances                int
	AZs                      []string
	StemcellName             string
	VMTypeName               string
	NetworkName              string
	SystemDomain             string
	AppDomains               []string
	Metron                   *Metron
	Statsd                   *StatsdInjector
	NFS                      *NFSMounter
	AllowSSHAccess           bool
	SkipSSLCertVerify        bool
	NATSUser                 string
	NATSPassword             string
	NATSPort                 int
	NATSMachines             []string
	CloudController          *CloudControllerPartition
	CCDBAddress              string
	CCDBUser                 string
	CCDBPassword             string
	JWTVerificationKey       string
	CCServiceDashboardSecret string
}

type gorouter struct {
	Instances    int
	AZs          []string
	StemcellName string
	VMTypeName   string
	NetworkName  string
	NetworkIPs   []string
	SSLCert      string
	SSLKey       string
	EnableSSL    bool
	Nats         grtrlib.Nats
	Loggregator  metron_agent.Loggregator
	RouterUser   string
	RouterPass   string
	MetronZone   string
	MetronSecret string
}

type diegoDatabase struct {
	AZs                []string
	Passphrase         string
	SystemDomain       string
	StemcellName       string
	VMTypeName         string
	PersistentDiskType string
	NetworkName        string
	NetworkIPs         []string
	CACert             string
	BBSServerCert      string
	BBSServerKey       string
	EtcdServerCert     string
	EtcdServerKey      string
	EtcdClientCert     string
	EtcdClientKey      string
	EtcdPeerCert       string
	EtcdPeerKey        string
	ConsulAgent        *ConsulAgent
	StatsdInjector     *StatsdInjector
	Metron             *Metron
	DiegoBrain         *diegoBrain
}

type diegoCell struct {
	AZs                []string
	StemcellName       string
	VMTypeName         string
	PersistentDiskType string
	NetworkName        string
	NetworkIPs         []string
	ConsulAgent        *ConsulAgent
	StatsdInjector     *StatsdInjector
	Metron             *Metron
	DiegoBrain         *diegoBrain
}

type diegoBrain struct {
	AZs                       []string
	StemcellName              string
	VMTypeName                string
	PersistentDiskType        string
	NetworkName               string
	NetworkIPs                []string
	BBSCACert                 string
	BBSClientCert             string
	BBSClientKey              string
	BBSAPILocation            string
	BBSRequireSSL             bool
	SkipSSLCertVerify         bool
	CCUploaderJobPollInterval int
	CCInternalAPIUser         string
	CCInternalAPIPassword     string
	CCBulkBatchSize           int
	CCFetchTimeout            int
	SystemDomain              string
	FSListenAddr              string
	FSStaticDirectory         string
	FSDebugAddr               string
	FSLogLevel                string
	MetronPort                int
	NATSUser                  string
	NATSPassword              string
	NATSPort                  int
	NATSMachines              []string
	AllowSSHAccess            bool
	SSHProxyClientSecret      string
	CCExternalPort            int
	TrafficControllerURL      string
	ConsulAgent               *ConsulAgent
	Metron                    *Metron
	Statsd                    *StatsdInjector
}

type loggregatorTrafficController struct {
	AZs               []string
	StemcellName      string
	VMTypeName        string
	NetworkName       string
	NetworkIPs        []string
	SystemDomain      string
	SkipSSLCertVerify bool
	EtcdMachines      []string
	DopplerSecret     string
	Metron            *Metron
	Nats              *routereglib.Nats
}

// Consul -
type Consul struct {
	AZs            []string
	StemcellName   string
	VMTypeName     string
	NetworkName    string
	NetworkIPs     []string
	ConsulAgent    *ConsulAgent
	Metron         *Metron
	StatsdInjector *StatsdInjector
}

//ConsulAgent -
type ConsulAgent struct {
	EncryptKeys []string
	CaCert      string
	AgentCert   string
	AgentKey    string
	ServerCert  string
	ServerKey   string
	NetworkIPs  []string
	Mode        string
	Services    []string
}

//Etcd -
type Etcd struct {
	AZs                []string
	StemcellName       string
	VMTypeName         string
	NetworkName        string
	NetworkIPs         []string
	PersistentDiskType string
	Metron             *Metron
	StatsdInjector     *StatsdInjector
	Nats               *etcdmetricslib.Nats
}

//Metron -
type Metron struct {
	Zone            string
	Secret          string
	SyslogAddress   string
	SyslogPort      int
	SyslogTransport string
	Loggregator     metron_agent.Loggregator
}

//StatsdInjector -
type StatsdInjector struct {
}

//NFSMounter -
type NFSMounter struct {
	NFSServerAddress string
	SharePath        string
}

//NatsPartition -
type NatsPartition struct {
	AZs            []string
	StemcellName   string
	VMTypeName     string
	NetworkName    string
	NetworkIPs     []string
	Nats           natslib.NatsJob
	Metron         *Metron
	StatsdInjector *StatsdInjector
}

//NFS -
type NFS struct {
	AZs                  []string
	StemcellName         string
	VMTypeName           string
	NetworkName          string
	NetworkIPs           []string
	PersistentDiskType   string
	AllowFromNetworkCIDR []string
	Metron               *Metron
	StatsdInjector       *StatsdInjector
}

//MySQL -
type MySQL struct {
	AZs                    []string
	StemcellName           string
	VMTypeName             string
	NetworkName            string
	NetworkIPs             []string
	PersistentDiskType     string
	AdminPassword          string
	DatabaseStartupTimeout int
	InnodbBufferPoolSize   int
	MaxConnections         int
	BootstrapUsername      string
	BootstrapPassword      string
	SyslogAddress          string
	SyslogPort             int
	SyslogTransport        string
	MySQLSeededDatabases   []MySQLSeededDatabase
}

//MySQLSeededDatabase -
type MySQLSeededDatabase struct {
	Name     string `yaml:"name"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

//MySQLProxy -
type MySQLProxy struct {
	AZs              []string
	StemcellName     string
	VMTypeName       string
	NetworkName      string
	NetworkIPs       []string
	ExternalHost     string
	APIUsername      string
	APIPassword      string
	ClusterIPs       []string
	Nats             *mysqlproxylib.Nats
	SyslogAggregator *mysqlproxylib.SyslogAggregator
}

//CloudControllerWorkerPartition - Cloud Controller Worker Partition
type CloudControllerWorkerPartition struct {
	AZs                   []string
	VMTypeName            string
	StemcellName          string
	NetworkName           string
	SystemDomain          string
	AppDomains            []string
	AllowedCorsDomains    []string
	AllowAppSSHAccess     bool
	Metron                *Metron
	ConsulAgent           *ConsulAgent
	StatsdInjector        *StatsdInjector
	NFSMounter            *NFSMounter
	StagingUploadUser     string
	StagingUploadPassword string
	BulkAPIUser           string
	BulkAPIPassword       string
	DbEncryptionKey       string
	InternalAPIUser       string
	InternalAPIPassword   string
}

//CloudControllerPartition - Cloud Controller Partition
type CloudControllerPartition struct {
	AZs                   []string
	VMTypeName            string
	StemcellName          string
	NetworkName           string
	SystemDomain          string
	AppDomains            []string
	AllowedCorsDomains    []string
	AllowAppSSHAccess     bool
	Metron                *Metron
	ConsulAgent           *ConsulAgent
	StatsdInjector        *StatsdInjector
	NFSMounter            *NFSMounter
	StagingUploadUser     string
	StagingUploadPassword string
	BulkAPIUser           string
	BulkAPIPassword       string
	DbEncryptionKey       string
	InternalAPIUser       string
	InternalAPIPassword   string
	HostKeyFingerprint    string
	SupportAddress        string
	MinCliVersion         string
}

//Doppler -
type Doppler struct {
	AZs                    []string
	StemcellName           string
	VMTypeName             string
	NetworkName            string
	NetworkIPs             []string
	Metron                 *Metron
	StatsdInjector         *StatsdInjector
	Zone                   string
	MessageDrainBufferSize int
	SharedSecret           string
	SystemDomain           string
	CCBuilkAPIPassword     string
	SkipSSLCertify         bool
	EtcdMachines           []string
}

//UAA -
type UAA struct {
	AZs                                       []string
	StemcellName                              string
	VMTypeName                                string
	NetworkName                               string
	Instances                                 int
	SystemDomain                              string
	RouterMachines                            []string
	Metron                                    *Metron
	StatsdInjector                            *StatsdInjector
	ConsulAgent                               *ConsulAgent
	Nats                                      *routereglib.Nats
	Login                                     *uaa.Login
	UAA                                       *uaa.Uaa
	SAMLServiceProviderKey                    string
	SAMLServiceProviderCertificate            string
	JWTSigningKey                             string
	JWTVerificationKey                        string
	Protocol                                  string
	AdminSecret                               string
	MySQLProxyHost                            string
	DBUserName                                string
	DBPassword                                string
	AdminPassword                             string
	PushAppsManagerPassword                   string
	SmokeTestsPassword                        string
	SystemServicesPassword                    string
	SystemVerificationPassword                string
	OpentsdbFirehoseNozzleClientSecret        string
	IdentityClientSecret                      string
	LoginClientSecret                         string
	PortalClientSecret                        string
	AutoScalingServiceClientSecret            string
	SystemPasswordsClientSecret               string
	CCServiceDashboardsClientSecret           string
	DopplerClientSecret                       string
	GoRouterClientSecret                      string
	NotificationsClientSecret                 string
	NotificationsUIClientSecret               string
	CloudControllerUsernameLookupClientSecret string
	CCRoutingClientSecret                     string
	SSHProxyClientSecret                      string
	AppsMetricsClientSecret                   string
	AppsMetricsProcessingClientSecret         string
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

// HAProxy -
type HAProxy struct {
	AZs            []string
	StemcellName   string
	VMTypeName     string
	NetworkName    string
	NetworkIPs     []string
	ConsulAgent    *ConsulAgent
	Metron         *Metron
	StatsdInjector *StatsdInjector
	SSLPem         string
	RouterMachines []string
}

type SmokeErrand struct {
	AZs          []string
	StemcellName string
	VMTypeName   string
	NetworkName  string
	Instances    int
	Protocol     string
	SystemDomain string
	AppsDomain   string
	Password     string
}

//Plugin -
type Plugin struct{}
