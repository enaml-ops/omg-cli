package cloudfoundry

import (
	"github.com/enaml-ops/enaml"
	grtrlib "github.com/enaml-ops/omg-cli/plugins/products/cloudfoundry/enaml-gen/gorouter"
	"github.com/enaml-ops/omg-cli/plugins/products/cloudfoundry/enaml-gen/metron_agent"
)

//InstanceGroupFactory -
type InstanceGroupFactory interface {
	ToInstanceGroup() (ig *enaml.InstanceGroup)
	hasValidValues() bool
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
	AZs          []string
	StemcellName string
	VMTypeName   string
	NetworkName  string
	NetworkIPs   []string
	EncryptKeys  []string
	CaCert       string
	AgentCert    string
	AgentKey     string
	ServerCert   string
	ServerKey    string
}
type Plugin struct{}
