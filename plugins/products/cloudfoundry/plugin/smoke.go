package cloudfoundry

import (
	"fmt"

	"github.com/codegangsta/cli"
	"github.com/enaml-ops/enaml"
	"github.com/enaml-ops/omg-cli/plugins/products/cloudfoundry/enaml-gen/smoke-tests"
)

//NewSmokeErrand - errand definition for smoke tests
func NewSmokeErrand(c *cli.Context) InstanceGrouper {
	return &SmokeErrand{
		AZs:          c.StringSlice("az"),
		StemcellName: c.String("stemcell-name"),
		NetworkName:  c.String("network"),
		VMTypeName:   c.String("errand-vm-type"),
		Protocol:     c.String("uaa-login-protocol"),
		Password:     c.String("smoke-tests-password"),
		SystemDomain: c.String("system-domain"),
		AppsDomain:   c.StringSlice("app-domain")[0],
	}
}

//ToInstanceGroup -
func (s *SmokeErrand) ToInstanceGroup() (ig *enaml.InstanceGroup) {
	ig = &enaml.InstanceGroup{
		Name:      "smoke-tests",
		Instances: 1,
		VMType:    s.VMTypeName,
		AZs:       s.AZs,
		Stemcell:  s.StemcellName,
		Jobs: []enaml.InstanceJob{
			s.createSmokeJob(),
		},
		Lifecycle: "errand",
		Networks: []enaml.Network{
			enaml.Network{Name: s.NetworkName},
		},
		Update: enaml.Update{
			MaxInFlight: 1,
		},
	}
	return
}

func (s *SmokeErrand) createSmokeJob() enaml.InstanceJob {
	return enaml.InstanceJob{
		Name:    "smoke-tests",
		Release: "cf",
		Properties: &smoke_tests.SmokeTestsJob{
			SmokeTests: &smoke_tests.SmokeTests{
				UseExistingOrg:   false,
				UseExistingSpace: false,
				Space:            "CF_SMOKE_TEST_SPACE",
				Org:              "CF_SMOKE_TEST_ORG",
				Password:         s.Password,
				User:             "smoke_tests",
				Api:              fmt.Sprintf("%s://api.%s", s.Protocol, s.SystemDomain),
				AppsDomain:       s.AppsDomain,
			},
		},
	}
}

//HasValidValues - Check if the datastructure has valid fields
func (s *SmokeErrand) HasValidValues() bool {
	return (len(s.AZs) > 0 &&
		s.StemcellName != "" &&
		s.VMTypeName != "" &&
		s.NetworkName != "" &&
		s.Protocol != "" &&
		s.Password != "" &&
		s.SystemDomain != "" &&
		s.AppsDomain != "")
}
