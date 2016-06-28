package cloudfoundry

import (
	"github.com/enaml-ops/enaml"
	etcdmetricslib "github.com/enaml-ops/omg-cli/plugins/products/cloudfoundry/enaml-gen/etcd_metrics_server"
	grtrlib "github.com/enaml-ops/omg-cli/plugins/products/cloudfoundry/enaml-gen/gorouter"
	natslib "github.com/enaml-ops/omg-cli/plugins/products/cloudfoundry/enaml-gen/nats"
	"github.com/enaml-ops/omg-cli/plugins/products/cloudfoundry/enaml-gen/metron_agent"
)

//InstanceGroupFactory -
type InstanceGroupFactory interface {
	ToInstanceGroup() (ig *enaml.InstanceGroup)
	HasValidValues() bool
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

//Consul -
type Consul struct {
	AZs            []string
	StemcellName   string
	VMTypeName     string
	NetworkName    string
	NetworkIPs     []string
	EncryptKeys    []string
	CaCert         string
	AgentCert      string
	AgentKey       string
	ServerCert     string
	ServerKey      string
	Metron         *Metron
	StatsdInjector *StatsdInjector
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

type Plugin struct{}
