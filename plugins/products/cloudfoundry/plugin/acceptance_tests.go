package cloudfoundry

import (
	"github.com/codegangsta/cli"
	"github.com/enaml-ops/enaml"
	"github.com/enaml-ops/omg-cli/plugins/products/cloudfoundry/enaml-gen/acceptance-tests"
)

func NewAcceptanceTestsPartition(c *cli.Context, internet bool) InstanceGrouper {
	return &acceptanceTests{
		AZs:                      c.StringSlice("az"),
		StemcellName:             c.String("stemcell-name"),
		NetworkName:              c.String("network"),
		AppsDomain:               c.StringSlice("app-domain"),
		SystemDomain:             c.String("system-domain"),
		AdminPassword:            c.String("admin-password"),
		SkipCertVerify:           c.BoolT("skip-cert-verify"),
		IncludeInternetDependent: internet,
	}
}

func (a *acceptanceTests) ToInstanceGroup() *enaml.InstanceGroup {
	return &enaml.InstanceGroup{
		Name:      "acceptance-tests",
		Instances: 1,
		VMType:    "errand",
		Lifecycle: "errand",
		AZs:       a.AZs,
		Stemcell:  a.StemcellName,
		Networks: []enaml.Network{
			{Name: a.NetworkName},
		},
		Update: enaml.Update{
			MaxInFlight: 1,
		},
		Jobs: []enaml.InstanceJob{
			{
				Name:       "acceptance-tests",
				Release:    CFReleaseName,
				Properties: a.newAcceptanceTestsProperties(a.IncludeInternetDependent),
			},
		},
	}
}

func (a *acceptanceTests) newAcceptanceTestsProperties(internet bool) *acceptance_tests.AcceptanceTests {
	var ad string
	if len(a.AppsDomain) > 0 {
		ad = a.AppsDomain[0]
	}
	return &acceptance_tests.AcceptanceTests{
		Api:                      prefixSystemDomain(a.SystemDomain, "api"),
		AppsDomain:               ad,
		AdminUser:                "admin",
		AdminPassword:            a.AdminPassword,
		IncludeLogging:           true,
		IncludeInternetDependent: internet,
		IncludeOperator:          true,
		IncludeServices:          true,
		IncludeSecurityGroups:    true,
		SkipSslValidation:        a.SkipCertVerify,
		SkipRegex:                "lucid64",
		JavaBuildpackName:        "java_buildpack_offline",
	}
}

func (a *acceptanceTests) HasValidValues() bool {
	return len(a.AZs) > 0 &&
		a.StemcellName != "" &&
		a.NetworkName != "" &&
		len(a.AppsDomain) > 0 &&
		a.SystemDomain != "" &&
		a.AdminPassword != ""
}
