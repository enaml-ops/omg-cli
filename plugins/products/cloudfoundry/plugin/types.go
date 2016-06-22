package cloudfoundry

import (
	"github.com/enaml-ops/omg-cli/plugins/products/cloudfoundry/enaml-gen/loggregator_trafficcontroller"
	"github.com/enaml-ops/omg-cli/plugins/products/cloudfoundry/enaml-gen/nats"
)

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
	Nats         nats.Nats
	Loggregator  loggregator_trafficcontroller.Loggregator
}
type Plugin struct{}
