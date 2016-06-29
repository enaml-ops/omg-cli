package cloudfoundry

import (
	"github.com/codegangsta/cli"
	"github.com/enaml-ops/enaml"
	consullib "github.com/enaml-ops/omg-cli/plugins/products/cloudfoundry/enaml-gen/consul_agent"
)

//NewConsulAgent -
func NewConsulAgent(c *cli.Context) *ConsulAgent {
	return &ConsulAgent{
		EncryptKeys: c.StringSlice("consul-encryption-key"),
		CaCert:      c.String("consul-ca-cert"),
		AgentCert:   c.String("consul-agent-cert"),
		AgentKey:    c.String("consul-agent-key"),
		ServerCert:  c.String("consul-server-cert"),
		ServerKey:   c.String("consul-server-key"),
		NetworkIPs:  c.StringSlice("consul-ip"),
	}
}

//CreateJob - Create the yaml job structure
func (s *ConsulAgent) CreateJob() enaml.InstanceJob {
	return enaml.InstanceJob{
		Name:    "consul_agent",
		Release: "cf",
		Properties: &consullib.Consul{
			EncryptKeys: s.EncryptKeys,
			CaCert:      s.CaCert,
			AgentCert:   s.AgentCert,
			AgentKey:    s.AgentKey,
			ServerCert:  s.ServerCert,
			ServerKey:   s.ServerKey,
			Agent: &consullib.Agent{
				Domain: "cf.internal",
				Mode:   "server",
				Servers: &consullib.Servers{
					Lan: s.NetworkIPs,
				},
			},
		},
	}
}
