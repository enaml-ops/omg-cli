package cloudfoundry

import (
	"github.com/codegangsta/cli"
	"github.com/enaml-ops/enaml"
	consullib "github.com/enaml-ops/omg-cli/plugins/products/cloudfoundry/enaml-gen/consul_agent"
)

//NewConsulAgent -
func NewConsulAgent(c *cli.Context, services []string) *ConsulAgent {
	return &ConsulAgent{
		EncryptKeys: c.StringSlice("consul-encryption-key"),
		CaCert:      c.String("consul-ca-cert"),
		AgentCert:   c.String("consul-agent-cert"),
		AgentKey:    c.String("consul-agent-key"),
		ServerCert:  c.String("consul-server-cert"),
		ServerKey:   c.String("consul-server-key"),
		NetworkIPs:  c.StringSlice("consul-ip"),
		Services:    services,
	}
}

//NewConsulAgentServer -
func NewConsulAgentServer(c *cli.Context) *ConsulAgent {
	return &ConsulAgent{
		EncryptKeys: c.StringSlice("consul-encryption-key"),
		CaCert:      c.String("consul-ca-cert"),
		AgentCert:   c.String("consul-agent-cert"),
		AgentKey:    c.String("consul-agent-key"),
		ServerCert:  c.String("consul-server-cert"),
		ServerKey:   c.String("consul-server-key"),
		NetworkIPs:  c.StringSlice("consul-ip"),
		Mode:        "server",
	}
}

//CreateJob - Create the yaml job structure
func (s *ConsulAgent) CreateJob() enaml.InstanceJob {

	serviceMap := make(map[string]map[string]string)
	for _, serviceName := range s.Services {
		serviceMap[serviceName] = make(map[string]string)
	}

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
				Mode:   s.getMode(),
				Servers: &consullib.Servers{
					Lan: s.NetworkIPs,
				},
				Services: serviceMap,
			},
		},
	}
}

func (s *ConsulAgent) getMode() interface{} {
	if s.Mode != "" {
		return s.Mode
	}
	return nil
}

//HasValidValues -
func (s *ConsulAgent) HasValidValues() bool {
	return len(s.NetworkIPs) > 0 &&
		len(s.EncryptKeys) > 0 &&
		s.CaCert != "" &&
		s.AgentCert != "" &&
		s.AgentKey != "" &&
		s.ServerCert != "" &&
		s.ServerKey != ""
}
