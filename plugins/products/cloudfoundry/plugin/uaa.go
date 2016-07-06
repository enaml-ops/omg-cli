package cloudfoundry

import (
	"fmt"

	"github.com/codegangsta/cli"
	"github.com/enaml-ops/enaml"
	"github.com/enaml-ops/omg-cli/plugins/products/cloudfoundry/enaml-gen/route_registrar"
	"github.com/enaml-ops/omg-cli/plugins/products/cloudfoundry/enaml-gen/uaa"
)

//NewUAAPartition -
func NewUAAPartition(c *cli.Context) InstanceGrouper {
	protocol := "https"
	if c.IsSet("uaa-login-protocol") {
		protocol = c.String("uaa-login-protocol")
	}
	UAA := &UAA{
		AZs:            c.StringSlice("az"),
		StemcellName:   c.String("stemcell-name"),
		NetworkName:    c.String("network"),
		VMTypeName:     c.String("uaa-vm-type"),
		Instances:      c.Int("uaa-instances"),
		SystemDomain:   c.String("system-domain"),
		Metron:         NewMetron(c),
		ConsulAgent:    NewConsulAgent(c, []string{"uaa"}),
		StatsdInjector: NewStatsdInjector(c),
		Nats: &route_registrar.Nats{
			User:     c.String("nats-user"),
			Password: c.String("nats-pass"),
			Machines: c.StringSlice("nats-machine-ip"),
			Port:     4222,
		},
		Protocol:                       protocol,
		SAMLServiceProviderKey:         c.String("uaa-saml-service-provider-key"),
		SAMLServiceProviderCertificate: c.String("uaa-saml-service-provider-certificate"),
		JWTSigningKey:                  c.String("uaa-jwt-signing-key"),
		JWTVerificationKey:             c.String("uaa-jwt-verification-key"),
		AdminSecret:                    c.String("uaa-admin-secret"),
		RouterMachines:                 c.StringSlice("router-ip"),
		MySQLProxyHost:                 c.String("mysql-proxy-external-host"),
		DBUserName:                     c.String("db-uaa-username"),
		DBPassword:                     c.String("db-uaa-password"),
		AdminPassword:                  c.String("admin-password"),
		PushAppsManagerPassword:        c.String("push-apps-manager-password"),
		SmokeTestsPassword:             c.String("smoke-tests-password"),
		SystemServicesPassword:         c.String("system-services-password"),
		SystemVerificationPassword:     c.String("system-verification-password"),
	}
	UAA.Login = UAA.CreateLogin(c)
	UAA.UAA = UAA.CreateUAA(c)
	return UAA
}

//CreateUAA - Helper method to create uaa structure
func (s *UAA) CreateUAA(c *cli.Context) (login *uaa.Uaa) {
	return &uaa.Uaa{
		RequireHttps: true,
		Ssl: &uaa.Ssl{
			Port: -1,
		},
		Authentication: &uaa.Authentication{
			Policy: &uaa.Policy{
				LockoutAfterFailures: 5,
			},
		},
		Password: &uaa.Password{
			Policy: &uaa.Policy{
				MinLength:                 0,
				RequireLowerCaseCharacter: 0,
				RequireUpperCaseCharacter: 0,
				RequireDigit:              0,
				RequireSpecialCharacter:   0,
				ExpirePasswordInMonths:    0,
			},
		},

		Ldap: &uaa.Ldap{
			ProfileType:         "search-and-bind",
			Url:                 c.String("uaa-ldap-url"),
			UserDN:              c.String("uaa-ldap-user-dn"),
			UserPassword:        c.String("uaa-ldap-user-password"),
			SearchBase:          c.String("uaa-ldap-search-base"),
			SearchFilter:        c.String("uaa-ldap-search-filter"),
			SslCertificate:      "",
			SslCertificateAlias: "",
			MailAttributeName:   c.String("uaa-ldap-mail-attributename"),
			Enabled:             c.BoolT("uaa-ldap-enabled"),
			Groups: &uaa.Groups{
				ProfileType:       "no-groups",
				SearchBase:        "",
				GroupSearchFilter: "",
			},
		},
		CatalinaOpts: "-Xmx768m -XX:MaxPermSize=256m",
		Url:          fmt.Sprintf("%s://uaa.%s", s.Protocol, s.SystemDomain),
		Jwt: &uaa.Jwt{
			SigningKey:      s.JWTSigningKey,
			VerificationKey: s.JWTVerificationKey,
		},
	}
}

//CreateLogin - Helper method to create login structure
func (s *UAA) CreateLogin(c *cli.Context) (login *uaa.Login) {
	return &uaa.Login{
		Branding:                CreateBranding(c),
		SelfServiceLinksEnabled: c.BoolT("uaa-enable-selfservice-links"),
		SignupsEnabled:          c.BoolT("uaa-signups-enabled"),
		Protocol:                s.Protocol,
		Links: &uaa.Links{
			Signup: fmt.Sprintf("%s://login.%s/create_account", s.Protocol, s.SystemDomain),
			Passwd: fmt.Sprintf("%s://login.%s/forgot_password", s.Protocol, s.SystemDomain),
		},
		UaaBase: fmt.Sprintf("%s://uaa.%s", s.Protocol, s.SystemDomain),
		Notifications: &uaa.Notifications{
			Url: fmt.Sprintf("%s://notifications.%s", s.Protocol, s.SystemDomain),
		},
		Saml: &uaa.Saml{
			Entityid:                   fmt.Sprintf("%s://login.%s", s.Protocol, s.SystemDomain),
			ServiceProviderKey:         s.SAMLServiceProviderKey,
			ServiceProviderCertificate: s.SAMLServiceProviderCertificate,
			SignRequest:                true,
			WantAssertionSigned:        false,
		},
		Logout: &uaa.Logout{
			Redirect: &uaa.Redirect{
				Parameter: &uaa.Parameter{
					Disable:   false,
					Whitelist: []string{fmt.Sprintf("%s://console.%s", s.Protocol, s.SystemDomain), fmt.Sprintf("%s://apps.%s", s.Protocol, s.SystemDomain)},
				},
				Url: "/login",
			},
		},
	}

}

func CreateBranding(c *cli.Context) (branding *uaa.Branding) {
	branding = &uaa.Branding{
		CompanyName:     c.String("uaa-company-name"),
		ProductLogo:     c.String("uaa-product-logo"),
		SquareLogo:      c.String("uaa-square-logo"),
		FooterLegalText: c.String("uaa-footer-legal-txt"),
	}
	return
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
	routes := make(map[string]interface{})
	routes["name"] = "uaa"
	routes["port"] = 8080
	routes["registration_interval"] = "40s"
	routes["uris"] = []string{fmt.Sprintf("uaa.%s", s.SystemDomain), fmt.Sprintf("*.uaa.%s", s.SystemDomain), fmt.Sprintf("login.%s", s.SystemDomain), fmt.Sprintf("*.login.%s", s.SystemDomain)}
	return enaml.InstanceJob{
		Name:    "route_registrar",
		Release: "cf",
		Properties: &route_registrar.RouteRegistrar{
			Routes: routes,
			Nats:   s.Nats,
		},
	}
}

func (s *UAA) createUAAJob() enaml.InstanceJob {
	return enaml.InstanceJob{
		Name:    "uaa",
		Release: "cf",
		Properties: &uaa.Uaa{
			Login: s.Login,
			Uaa:   s.UAA,
			Admin: &uaa.Admin{
				ClientSecret: s.AdminSecret,
			},
			Proxy: &uaa.Proxy{
				Servers: s.RouterMachines,
			},
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
				Users: []string{
					fmt.Sprintf("admin|%s|scim.write,scim.read,openid,cloud_controller.admin,dashboard.user,console.admin,console.support,doppler.firehose,notification_preferences.read,notification_preferences.write,notifications.manage,notification_templates.read,notification_templates.write,emails.write,notifications.write,zones.read,zones.write", s.AdminPassword),
					fmt.Sprintf("push_apps_manager|%s|cloud_controller.admin", s.PushAppsManagerPassword),
					fmt.Sprintf("smoke_tests|%s|cloud_controller.admin", s.SmokeTestsPassword),
					fmt.Sprintf("system_services|%s|cloud_controller.admin", s.SystemServicesPassword),
					fmt.Sprintf("system_verification|%s|scim.write,scim.read,openid,cloud_controller.admin,dashboard.user,console.admin,console.support", s.SystemVerificationPassword),
				},
			},
			Domain: s.SystemDomain,
			Uaadb:  s.createUAADB(),
		},
	}
}

func (s *UAA) createUAADB() (uaadb *uaa.Uaadb) {
	const uaaVal = "uaa"
	roles := make(map[string]string)
	roles["tag"] = "admin"
	roles["name"] = s.DBUserName
	roles["password"] = s.DBPassword

	dbs := make(map[string]string)
	dbs["tag"] = uaaVal
	dbs["name"] = uaaVal
	return &uaa.Uaadb{
		Address:   s.MySQLProxyHost,
		Port:      3306,
		DbScheme:  "mysql",
		Roles:     roles,
		Databases: dbs,
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
		s.ConsulAgent.HasValidValues() &&
		s.Nats.User != "" &&
		s.Nats.Password != "" &&
		len(s.Nats.Machines.([]string)) > 0 &&
		s.SAMLServiceProviderKey != "" &&
		s.SAMLServiceProviderCertificate != "" &&
		s.JWTSigningKey != "" &&
		s.JWTVerificationKey != "" &&
		s.AdminSecret != "" &&
		len(s.RouterMachines) > 0 &&
		s.MySQLProxyHost != "" &&
		s.DBUserName != "" &&
		s.DBPassword != "" && s.AdminPassword != "" &&
		s.PushAppsManagerPassword != "" &&
		s.SmokeTestsPassword != "" &&
		s.SystemServicesPassword != "" &&
		s.SystemVerificationPassword != "")
}
