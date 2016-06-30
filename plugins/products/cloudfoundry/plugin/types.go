package cloudfoundry

import (
	"github.com/codegangsta/cli"
	"github.com/enaml-ops/enaml"
	mysqlproxylib "github.com/enaml-ops/omg-cli/plugins/products/cf-mysql/enaml-gen/proxy"
	etcdmetricslib "github.com/enaml-ops/omg-cli/plugins/products/cloudfoundry/enaml-gen/etcd_metrics_server"
	grtrlib "github.com/enaml-ops/omg-cli/plugins/products/cloudfoundry/enaml-gen/gorouter"
	"github.com/enaml-ops/omg-cli/plugins/products/cloudfoundry/enaml-gen/metron_agent"
	natslib "github.com/enaml-ops/omg-cli/plugins/products/cloudfoundry/enaml-gen/nats"
)

// InstanceGrouper creates and validates InstanceGroups.
type InstanceGrouper interface {
	ToInstanceGroup() (ig *enaml.InstanceGroup)
	HasValidValues() bool
}

// InstanceGrouperFactory is a function that creates InstanceGroupers from CLI args.
type InstanceGrouperFactory func(*cli.Context) InstanceGrouper

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

type diegoBrain struct {
	AZs                  []string
	StemcellName         string
	VMTypeName           string
	PersistentDiskType   string
	NetworkName          string
	NetworkIPs           []string
	AuctioneerCACert     string
	AuctioneerServerCert string
	AuctioneerServerKey  string
	AuctioneerClientCert string
	AuctioneerClientKey  string
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
	Nats           natslib.Nats
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

//Plugin -
type Plugin struct{}
