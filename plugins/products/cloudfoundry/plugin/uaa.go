package cloudfoundry

import (
	"fmt"

	"github.com/codegangsta/cli"
	"github.com/enaml-ops/enaml"
	"github.com/enaml-ops/omg-cli/plugins/products/cloudfoundry/enaml-gen/uaa"
)

//NewUAAPartition -
func NewUAAPartition(c *cli.Context) InstanceGrouper {
	return &UAA{
		AZs:            c.StringSlice("az"),
		StemcellName:   c.String("stemcell-name"),
		NetworkName:    c.String("network"),
		VMTypeName:     c.String("uaa-vm-type"),
		Instances:      c.Int("uaa-instances"),
		SystemDomain:   c.String("system-domain"),
		Metron:         NewMetron(c),
		ConsulAgent:    NewConsulAgent(c, []string{"uaa"}),
		StatsdInjector: NewStatsdInjector(c),
	}
}

//ToInstanceGroup -
func (s *UAA) ToInstanceGroup() (ig *enaml.InstanceGroup) {
	ig = &enaml.InstanceGroup{
		Name:      "uaa-partition",
		Instances: s.Instances,
		VMType:    s.VMTypeName,
		AZs:       s.AZs,
		Stemcell:  s.StemcellName,
		Jobs: []enaml.InstanceJob{
			s.createUAAJob(),
			s.Metron.CreateJob(),
			s.ConsulAgent.CreateJob(),
			s.StatsdInjector.CreateJob(),
			s.createRouteRegistrarJob(),
		},
		Networks: []enaml.Network{
			enaml.Network{Name: s.NetworkName},
		},
		Update: enaml.Update{
			MaxInFlight: 1,
		},
	}
	return
}

func (s *UAA) createRouteRegistrarJob() enaml.InstanceJob {
	return enaml.InstanceJob{
		Name:    "route_registrar",
		Release: "cf",
	}
}

func (s *UAA) createUAAJob() enaml.InstanceJob {
	protocol := "https"
	return enaml.InstanceJob{
		Name:    "uaa",
		Release: "cf",
		Properties: &uaa.Uaa{
			Login: &uaa.Login{
				Branding: &uaa.Branding{
					CompanyName:     "",
					ProductLogo:     "",
					SquareLogo:      "",
					FooterLegalText: "",
				},
				SelfServiceLinksEnabled: true,
				SignupsEnabled:          true,
				Protocol:                protocol,
				Links: &uaa.Links{
					Signup: "",
					Passwd: "",
				},
				UaaBase: fmt.Sprintf("%s://uaa.%s", protocol, s.SystemDomain),
				Notifications: &uaa.Notifications{
					Url: fmt.Sprintf("%s://notifications.%s", protocol, s.SystemDomain),
				},
				Saml: &uaa.Saml{
					Entityid:                   fmt.Sprintf("%s://login.%s", protocol, s.SystemDomain),
					ServiceProviderKey:         "",
					ServiceProviderCertificate: "",
					SignRequest:                true,
					WantAssertionSigned:        false,
				},
				Logout: &uaa.Logout{
					Redirect: &uaa.Redirect{
						Parameter: &uaa.Parameter{
							Disable:   false,
							Whitelist: []string{fmt.Sprintf("%s://console.%s", protocol, s.SystemDomain), fmt.Sprintf("%s://apps.%s", protocol, s.SystemDomain)},
						},
						Url: "/login",
					},
				},
			},
			Uaa: &uaa.Uaa{
				RequireHttps: true,
				Ssl: &uaa.Ssl{
					Port: -1,
				},
				Ldap: &uaa.Ldap{
					ProfileType:         "search-and-bind",
					Url:                 "",
					UserDN:              "",
					UserPassword:        "",
					SearchBase:          "",
					SearchFilter:        "",
					SslCertificate:      "",
					SslCertificateAlias: "",
					MailAttributeName:   "",
					Enabled:             false,
					Groups: &uaa.Groups{
						ProfileType:       "no-groups",
						SearchBase:        "",
						GroupSearchFilter: "",
					},
				},
				CatalinaOpts: "-Xmx768m -XX:MaxPermSize=256m",
				Url:          fmt.Sprintf("%s://uaa.%s", protocol, s.SystemDomain),
				Jwt: &uaa.Jwt{
					SigningKey:      "",
					VerificationKey: "",
				},
			},
			Admin: &uaa.Admin{
				ClientSecret: "",
			},
			Proxy: &uaa.Proxy{
				Servers: s.RouterMachines,
			},
			//TODO create map of clients
			Clients: "",
			Scim: &uaa.Scim{
				User: &uaa.User{
					Override: true,
				},
				UseridsEnabled: true,
				//TODO | delimited list of strings...
				//- admin|aa7abaa063a34c7269ba|scim.write,scim.read,openid,cloud_controller.admin,dashboard.user,console.admin,console.support,doppler.firehose,notification_preferences.read,notification_preferences.write,notifications.manage,notification_templates.read,notification_templates.write,emails.write,notifications.write,zones.read,zones.write
				//- push_apps_manager|a5687a3153a58b9fb491|cloud_controller.admin
				//- smoke_tests|b9781aeee53b0f933591|cloud_controller.admin
				//- system_services|efee313efd2ad0134548|cloud_controller.admin
				//- system_verification|bfbeed4cc362c1ae6ec4|scim.write,scim.read,openid,cloud_controller.admin,dashboard.user,console.admin,console.support
				Users: "",
			},
			Domain: s.SystemDomain,
			Uaadb: &uaa.Uaadb{
				Address:  "",
				Port:     3306,
				DbScheme: "mysql",
				//TODO define stuct for
				//- tag: admin
				//  name: c040224fe68cf6e00b52
				//  password: fcfccfef6d0c542bacbd
				Roles: nil,
				//TODO define stuct for
				//- tag: uaa
				//  name: uaa
				Databases: nil,
			},
		},
	}
}

//HasValidValues - Check if the datastructure has valid fields
func (s *UAA) HasValidValues() bool {
	return (len(s.AZs) > 0 &&
		s.StemcellName != "" &&
		s.VMTypeName != "" &&
		s.NetworkName != "" &&
		s.Instances > 0 &&
		s.SystemDomain != "" &&
		s.Metron.HasValidValues() &&
		s.StatsdInjector.HasValidValues() &&
		s.ConsulAgent.HasValidValues())
}
