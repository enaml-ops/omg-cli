package cloudfoundry

import (
	"fmt"

	"gopkg.in/yaml.v2"

	"github.com/codegangsta/cli"
	"github.com/enaml-ops/enaml"
	consullib "github.com/enaml-ops/omg-cli/plugins/products/cloudfoundry/enaml-gen/consul_agent"
)

//NewConsulPartition -
func NewConsulPartition(c *cli.Context) (igf InstanceGroupFactory, err error) {
	var metron *Metron
	if metron, err = NewMetron(c); err == nil {
		igf = &Consul{
			AZs:          c.StringSlice("az"),
			StemcellName: c.String("stemcell-name"),
			NetworkIPs:   c.StringSlice("consul-ip"),
			NetworkName:  c.String("consul-network"),
			VMTypeName:   c.String("consul-vm-type"),
			EncryptKeys:  c.StringSlice("consul-encryption-key"),
			CaCert:       c.String("consul-ca-cert"),
			AgentCert:    c.String("consul-agent-cert"),
			AgentKey:     c.String("consul-agent-key"),
			ServerCert:   c.String("consul-server-cert"),
			ServerKey:    c.String("consul-server-key"),
			Metron:       metron,
		}
	}

	if !igf.hasValidValues() {
		b, _ := yaml.Marshal(igf)
		err = fmt.Errorf("invalid values in Consul: %v", string(b))
		igf = nil
	}
	return
}

//ToInstanceGroup -
func (s *Consul) ToInstanceGroup() (ig *enaml.InstanceGroup) {
	ig = &enaml.InstanceGroup{
		Name:      "consul-partition",
		Instances: len(s.NetworkIPs),
		VMType:    s.VMTypeName,
		AZs:       s.AZs,
		Stemcell:  s.StemcellName,
		Jobs: []enaml.InstanceJob{
			s.newConsulAgentJob(),
			s.Metron.CreateMetronJob(),
		},
		Networks: []enaml.Network{
			enaml.Network{Name: s.NetworkName, StaticIPs: s.NetworkIPs},
		},
	}
	return
}

func (s *Consul) newConsulAgentJob() enaml.InstanceJob {
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
func (s *Consul) hasValidValues() bool {
	return (len(s.AZs) > 0 &&
		s.StemcellName != "" &&
		s.VMTypeName != "" &&
		s.NetworkName != "" &&
		len(s.NetworkIPs) > 0 &&
		len(s.EncryptKeys) > 0 &&
		s.CaCert != "" &&
		s.AgentCert != "" &&
		s.AgentKey != "" &&
		s.ServerCert != "" &&
		s.ServerKey != "")
}
