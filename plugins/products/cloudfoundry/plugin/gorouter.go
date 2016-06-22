package cloudfoundry

import "github.com/codegangsta/cli"

func NewGoRouterPartition(c *cli.Context) (grtr *gorouter, err error) {
	grtr = new(gorouter)
	if !grtr.hasValidValues() {
		grtr = nil
	}
	return
}

func (s *gorouter) hasValidValues() bool {
	return (len(s.AZs) > 0 &&
		s.StemcellName != "" &&
		s.VMTypeName != "" &&
		s.NetworkName != "" &&
		len(s.NetworkIPs) > 0 &&
		s.SSLCert != "" &&
		s.SSLKey != "" &&
		s.Nats.User != "" &&
		s.Nats.Password != "" &&
		s.Nats.Machines != nil &&
		s.Loggregator.Etcd.Machines != nil)
}
