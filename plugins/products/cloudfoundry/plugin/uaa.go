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
	var mysqlProxyIP string
	mysqlProxys := c.StringSlice("mysql-proxy-ip")
	if len(mysqlProxys) > 0 {
		mysqlProxyIP = mysqlProxys[0]
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
		Protocol:                                  protocol,
		SAMLServiceProviderKey:                    c.String("uaa-saml-service-provider-key"),
		SAMLServiceProviderCertificate:            c.String("uaa-saml-service-provider-certificate"),
		JWTSigningKey:                             c.String("uaa-jwt-signing-key"),
		JWTVerificationKey:                        c.String("uaa-jwt-verification-key"),
		AdminSecret:                               c.String("uaa-admin-secret"),
		RouterMachines:                            c.StringSlice("router-ip"),
		MySQLProxyHost:                            mysqlProxyIP,
		DBUserName:                                c.String("db-uaa-username"),
		DBPassword:                                c.String("db-uaa-password"),
		AdminPassword:                             c.String("admin-password"),
		PushAppsManagerPassword:                   c.String("push-apps-manager-password"),
		SmokeTestsPassword:                        c.String("smoke-tests-password"),
		SystemServicesPassword:                    c.String("system-services-password"),
		SystemVerificationPassword:                c.String("system-verification-password"),
		OpentsdbFirehoseNozzleClientSecret:        c.String("opentsdb-firehose-nozzle-client-secret"),
		IdentityClientSecret:                      c.String("identity-client-secret"),
		LoginClientSecret:                         c.String("login-client-secret"),
		PortalClientSecret:                        c.String("portal-client-secret"),
		AutoScalingServiceClientSecret:            c.String("autoscaling-service-client-secret"),
		SystemPasswordsClientSecret:               c.String("system-passwords-client-secret"),
		CCServiceDashboardsClientSecret:           c.String("cc-service-dashboards-client-secret"),
		DopplerClientSecret:                       c.String("doppler-client-secret"),
		GoRouterClientSecret:                      c.String("gorouter-client-secret"),
		NotificationsClientSecret:                 c.String("notifications-client-secret"),
		NotificationsUIClientSecret:               c.String("notifications-ui-client-secret"),
		CloudControllerUsernameLookupClientSecret: c.String("cloud-controller-username-lookup-client-secret"),
		CCRoutingClientSecret:                     c.String("cc-routing-client-secret"),
		SSHProxyClientSecret:                      c.String("ssh-proxy-client-secret"),
		AppsMetricsClientSecret:                   c.String("apps-metrics-client-secret"),
		AppsMetricsProcessingClientSecret:         c.String("apps-metrics-processing-client-secret"),
	}
	UAA.Login = UAA.CreateLogin(c)
	UAA.UAA = UAA.CreateUAA(c)
	return UAA
}

//CreateUAA - Helper method to create uaa structure
func (s *UAA) CreateUAA(c *cli.Context) (login *uaa.Uaa) {
	clientMap := make(map[string]UAAClient)
	clientMap["opentsdb-firehose-nozzle"] = UAAClient{
		AccessTokenValidity:  1209600,
		AuthorizedGrantTypes: "authorization_code,client_credentials,refresh_token",
		Override:             true,
		Secret:               s.OpentsdbFirehoseNozzleClientSecret,
		Scope:                "openid,oauth.approvals,doppler.firehose",
		Authorities:          "oauth.login,doppler.firehose",
	}
	clientMap["identity"] = UAAClient{
		ID:                   "identity",
		Secret:               s.IdentityClientSecret,
		Scope:                "cloud_controller.admin,cloud_controller.read,cloud_controller.write,openid,zones.*.*,zones.*.*.*,zones.read,zones.write,scim.read",
		ResourceIDs:          "none",
		AuthorizedGrantTypes: "authorization_code,client_credentials,refresh_token",
		AutoApprove:          true,
		Authorities:          "scim.zones,zones.read,uaa.resource,zones.write,cloud_controller.admin",
		RedirectURI:          fmt.Sprintf("%s://p-identity.%s/dashboard/,%s://p-identity.%s/dashboard/**", s.Protocol, s.SystemDomain, s.Protocol, s.SystemDomain),
	}
	clientMap["login"] = UAAClient{
		ID:                   "login",
		Secret:               s.LoginClientSecret,
		AutoApprove:          true,
		Override:             true,
		Authorities:          "oauth.login,scim.write,clients.read,notifications.write,critical_notifications.write,emails.write,scim.userids,password.write",
		AuthorizedGrantTypes: "authorization_code,client_credentials,refresh_token",
		Scope:                "openid,oauth.approvals",
	}
	clientMap["portal"] = UAAClient{
		ID:                   "portal",
		Secret:               s.PortalClientSecret,
		Override:             true,
		AutoApprove:          true,
		Authorities:          "scim.write,scim.read,cloud_controller.read,cloud_controller.write,password.write,uaa.admin,uaa.resource,cloud_controller.admin,emails.write,notifications.write",
		Scope:                "openid,cloud_controller.read,cloud_controller.write,password.write,console.admin,console.support,cloud_controller.admin",
		AccessTokenValidity:  1209600,
		RefreshTokenValidity: 1209600,
		Name:                 "Pivotal Apps Manager",
		AppLaunchURL:         fmt.Sprintf("%s://apps.%s", s.Protocol, s.SystemDomain),
		ShowOnHomepage:       true,
		AppIcon:              "iVBORw0KGgoAAAANSUhEUgAAAGwAAABsCAYAAACPZlfNAAAAAXNSR0IArs4c6QAABYtJREFUeAHtnVtsFFUYx7/d3ruWotUKVIkNaCw02YgJGBRTMd4CokUejD4QH4gxQcIDeHnBmPjkhSghUYLGe3ywPtAHNCo0QgkWwi2tXG2V1kIpLXTbLt1tS9dzlmzSJssZhv32zDk7/2km2znn7Pd9+/vt2Z2dmW0D9Obat4gCiwiLBQQSLflSViAQeN6Can1fYiJBFPQ9BcsAQBiEWUbAsnIxwyDMMgKWlYsZBmGWEbCsXMwwCLOMgGXlYoZBmGUELCsXMwzCLCNgWbmYYRBmGQHLysUMgzDLCFhWLmYYhFlGwLJyMcMgzDIClpWLGQZhlhGwrFzMMAizjIBl5WKGQZhlBCwrV1xbb96y59V1VFJQmLawQNrWa43x8XEaHo1fW+Oj1H8lSqf6eulEbw+dvNhLvcNDinvb0WWksAdm3UWhwiJ2gt2RAWo80UY7jrdSU8cZGrt6lT1HtgMaKSxbD7qqfDq99tAjyTUSG6FP9v1BH+3dTUPxeLZSssf17U5HeXEJbXr8aerY+A6tf7iOxFeu2OFmI6BvhaVgVoRCtHl5PTW8/AoV5xekmo299b2wlJn6+WFqWrOWKkpDqSYjbyFskpZFs++hL1e9NKnFvF+t3OmQOwzdkcgUmnnBABXm5Ys1j8qKisVadFPvS8tramn1goX09eEDU+KbsmGlsMbjbbT6x++UDOVORGXoFppXOYMerLqbVsyrpcWzqykYdH5R+fjZlcnd/8sjV5Q5vOh0rt6LqhhyJsQ3uC+ID8ry89aHYtf90W1bKLzlffr19EnH6HIP8oXasOM4LwbkrLB0MP+6cJ6e+eoz+vTP5nTdU9peDC+Ysm3Khq+ESehy5r3e2ECHu7uUDuqq59Id4iXVtMV3wqSACSHt3V2/KF3I97qayjuVY7zo9KUwCfq3M6coNjamZD6zrFzZ70Wnb4XFxseoK3JZyXzWtGnKfi86fStMwu6LRpXMZ5RBmBKQ7k75XqZa8gLmPZ/Nq0hFkLnvttJSZUT5Oc60xbfC5CGs6lsrlD56hgaV/V50+lbYkuo5VFygPp3SMwxhXjwp0+bcsGRp2vZU48TEBB09153aNObWlzNMHo1/6r4apYTmsx10MTqsHONFp5VH6zMBtWbhYtq6YpVjiJ/ajjmO8WKAL4QFxamWZffPT1678dicex05D4jTKj8cO+Q4zosBOSXs7bonktci5ovjgPIUye3ieo3wzKrk+TC5faPLGz83On6ovtFY3ONySth7Ty67qbPMk6Hu+edv+vzg/slNRv3uy52O6xk40HWW6r/94nrdRrTn1AzLhOju9tP03DfbKTo6mkmYrN/X98L6xQHgTb/vpG0t+5LnybJOPMMEvhXWOXCJvj9yiD7Yu4sGRkYyxKjv7r4RJi+Na+05Rwf/66SG1qO0v/NffZQZM+WUsI07d1BC/MTE144GYzHxJYcYDYq1vb/f8WQlI9OshsopYZubm7IKy4Tg2K03wYKLGiDMBSwThkKYCRZc1ABhLmCZMBTCTLDgogYIcwHLhKEQZoIFFzVAmAtYJgyFMBMsuKgBwlzAMmEohJlgwUUNEOYClglDIcwECy5qgDAXsEwYCmEmWHBRA4S5gGXCUAgzwYKLGow84yyvuyhR/GW19kt9Lh5ibg01UtjS7VtzizLjo8FLIiNMHaEgTAdlxhwQxghTRygI00GZMQeEMcLUEQrCdFBmzAFhjDB1hIIwHZQZc0AYI0wdoSBMB2XGHBDGCFNHKAjTQZkxB4QxwtQRCsJ0UGbMAWGMMHWEgjAdlBlzQBgjTB2hIEwHZcYcEMYIU0coCNNBmTEHhDHC1BEKwnRQZswBYYwwdYSCMB2UGXNAGCNMHaEgTAdlxhziUu1Ei8M/+WFMh1CZEUi0/A+j7hNSB5Wo2wAAAABJRU5ErkJggg==",
	}
	clientMap["apps_manager_js"] = UAAClient{
		Override:             true,
		AutoApprove:          []string{"cloud_controller.read", "cloud_controller.write", "cloud_controller.admin"},
		Scope:                "cloud_controller.read,cloud_controller.write,cloud_controller.admin",
		AuthorizedGrantTypes: "implicit",
		AccessTokenValidity:  28800,
	}
	clientMap["cf"] = UAAClient{
		ID:                   "cf",
		Override:             true,
		Authorities:          "uaa.none",
		AuthorizedGrantTypes: "password,refresh_token",
		Scope:                "cloud_controller.read,cloud_controller.write,openid,password.write,cloud_controller.admin,scim.read,scim.write,doppler.firehose,uaa.user",
		AccessTokenValidity:  7200,
		RefreshTokenValidity: 1209600,
	}
	clientMap["autoscaling_service"] = UAAClient{
		ID:                   "autoscaling_service",
		Secret:               s.AutoScalingServiceClientSecret,
		Override:             true,
		AutoApprove:          true,
		Authorities:          "cloud_controller.write,cloud_controller.read,cloud_controller.admin,notifications.write,critical_notifications.write,emails.write",
		AuthorizedGrantTypes: "client_credentials,authorization_code,refresh_token",
		Scope:                "openid,cloud_controller.permissions,cloud_controller.read,cloud_controller.write",
		AccessTokenValidity:  3600,
	}
	clientMap["system_passwords"] = UAAClient{
		ID:                   "system_passwords",
		Secret:               s.SystemPasswordsClientSecret,
		Override:             true,
		AutoApprove:          true,
		Authorities:          "uaa.admin,scim.read,scim.write,password.write",
		AuthorizedGrantTypes: "client_credentials",
	}
	clientMap["cc-service-dashboards"] = UAAClient{
		ID:                   "cc-service-dashboards",
		Secret:               s.CCServiceDashboardsClientSecret,
		Override:             true,
		Authorities:          "clients.read,clients.write,clients.admin",
		AuthorizedGrantTypes: "client_credentials",
		Scope:                "cloud_controller.write,openid,cloud_controller.read,cloud_controller_service_permissions.read",
	}
	clientMap["doppler"] = UAAClient{
		ID:          "doppler",
		Secret:      s.DopplerClientSecret,
		Authorities: "uaa.resource",
	}
	clientMap["gorouter"] = UAAClient{
		ID:                   "gorouter",
		Secret:               s.GoRouterClientSecret,
		Authorities:          "clients.read,clients.write,clients.admin,routing.routes.write,routing.routes.read",
		AuthorizedGrantTypes: "client_credentials,refresh_token",
		Scope:                "openid,cloud_controller_service_permissions.read",
	}
	clientMap["notifications"] = UAAClient{
		ID:                   "notifications",
		Secret:               s.NotificationsClientSecret,
		Authorities:          "cloud_controller.admin,scim.read,notifications.write,critical_notifications.write,emails.write",
		AuthorizedGrantTypes: "client_credentials",
	}
	clientMap["notifications_template"] = UAAClient{
		ID:                   "notifications_template",
		Secret:               "bb6be96896c5ab64c897",
		Scope:                "openid,clients.read,clients.write,clients.secret",
		Authorities:          "clients.read,clients.write,clients.secret,notification_templates.write,notification_templates.read,notifications.manage",
		AuthorizedGrantTypes: "client_credentials",
	}
	clientMap["notifications_ui_client"] = UAAClient{
		ID:                   "notifications_ui_client",
		Secret:               s.NotificationsUIClientSecret,
		Scope:                "notification_preferences.read,notification_preferences.write,openid",
		AuthorizedGrantTypes: "authorization_code,client_credentials,refresh_token",
		Authorities:          "notification_preferences.admin",
		AutoApprove:          true,
		Override:             true,
		RedirectURI:          fmt.Sprintf("%s://notifications-ui.%s/sessions/create", s.Protocol, s.SystemDomain),
	}
	clientMap["cloud_controller_username_lookup"] = UAAClient{
		ID:                   "cloud_controller_username_lookup",
		Secret:               s.CloudControllerUsernameLookupClientSecret,
		AuthorizedGrantTypes: "client_credentials",
		Authorities:          "scim.userids",
	}
	clientMap["cc_routing"] = UAAClient{
		Authorities:          "routing.router_groups.read",
		AuthorizedGrantTypes: "client_credentials",
		Secret:               s.CCRoutingClientSecret,
	}
	clientMap["ssh-proxy"] = UAAClient{
		AuthorizedGrantTypes: "authorization_code",
		AutoApprove:          true,
		Override:             true,
		RedirectURI:          "/login",
		Scope:                "openid,cloud_controller.read,cloud_controller.write",
		Secret:               s.SSHProxyClientSecret,
	}
	clientMap["apps_metrics"] = UAAClient{
		ID:                   "apps_metrics",
		Secret:               s.AppsMetricsClientSecret,
		Override:             true,
		AuthorizedGrantTypes: "authorization_code,refresh_token",
		RedirectURI:          fmt.Sprintf("%s://apm.%s,%s://apm.%s/,%s://apm.%s/*,%s://metrics.%s,%s://metrics.%s/,%s://metrics.%s/*", s.Protocol, s.SystemDomain, s.Protocol, s.SystemDomain, s.Protocol, s.SystemDomain, s.Protocol, s.SystemDomain, s.Protocol, s.SystemDomain, s.Protocol, s.SystemDomain),
		Scope:                "cloud_controller.admin,cloud_controller.read,metrics.read",
		AccessTokenValidity:  900,
		RefreshTokenValidity: 2628000,
	}
	clientMap["apps_metrics_processing"] = UAAClient{
		ID:                   "apps_metrics_processing",
		Secret:               s.AppsMetricsProcessingClientSecret,
		Override:             true,
		AuthorizedGrantTypes: "authorization_code,client_credentials,refresh_token",
		Authorities:          "oauth.login,doppler.firehose,cloud_controller.admin",
		Scope:                "openid,oauth.approvals,doppler.firehose,cloud_controller.admin",
		AccessTokenValidity:  1209600,
	}
	return &uaa.Uaa{
		RequireHttps: true,
		Ssl: &uaa.Ssl{
			Port: -1,
		},
		Admin: &uaa.Admin{
			ClientSecret: s.AdminSecret,
		},
		Authentication: &uaa.Authentication{
			Policy: &uaa.AuthenticationPolicy{
				LockoutAfterFailures: 5,
			},
		},
		Password: &uaa.UaaPassword{
			Policy: &uaa.PasswordPolicy{
				MinLength:                 0,
				RequireLowerCaseCharacter: 0,
				RequireUpperCaseCharacter: 0,
				RequireDigit:              0,
				RequireSpecialCharacter:   0,
				ExpirePasswordInMonths:    0,
			},
		},

		Ldap: &uaa.UaaLdap{
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
			Groups: &uaa.LdapGroups{
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
		Proxy: &uaa.Proxy{
			Servers: s.RouterMachines,
		},
		Clients: clientMap,
		Scim: &uaa.Scim{
			User: &uaa.ScimUser{
				Override: true,
			},
			UseridsEnabled: true,
			Users: []string{
				fmt.Sprintf("admin|%s|scim.write,scim.read,openid,cloud_controller.admin,dashboard.user,console.admin,console.support,doppler.firehose,notification_preferences.read,notification_preferences.write,notifications.manage,notification_templates.read,notification_templates.write,emails.write,notifications.write,zones.read,zones.write", s.AdminPassword),
				fmt.Sprintf("push_apps_manager|%s|cloud_controller.admin", s.PushAppsManagerPassword),
				fmt.Sprintf("smoke_tests|%s|cloud_controller.admin", s.SmokeTestsPassword),
				fmt.Sprintf("system_services|%s|cloud_controller.admin", s.SystemServicesPassword),
				fmt.Sprintf("system_verification|%s|scim.write,scim.read,openid,cloud_controller.admin,dashboard.user,console.admin,console.support", s.SystemVerificationPassword),
			},
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
		Properties: &route_registrar.RouteRegistrarJob{
			RouteRegistrar: &route_registrar.RouteRegistrar{
				Routes: routes,
			},
			Nats: s.Nats,
		},
	}
}

func (s *UAA) createUAAJob() enaml.InstanceJob {
	return enaml.InstanceJob{
		Name:    "uaa",
		Release: "cf",
		Properties: &uaa.UaaJob{
			Login:  s.Login,
			Uaa:    s.UAA,
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
		s.SystemVerificationPassword != "" &&
		s.OpentsdbFirehoseNozzleClientSecret != "" &&
		s.IdentityClientSecret != "" &&
		s.LoginClientSecret != "" &&
		s.PortalClientSecret != "" &&
		s.AutoScalingServiceClientSecret != "" &&
		s.SystemPasswordsClientSecret != "" &&
		s.CCServiceDashboardsClientSecret != "" &&
		s.DopplerClientSecret != "" &&
		s.GoRouterClientSecret != "" &&
		s.NotificationsClientSecret != "" &&
		s.NotificationsUIClientSecret != "" &&
		s.CloudControllerUsernameLookupClientSecret != "" &&
		s.CCRoutingClientSecret != "" &&
		s.SSHProxyClientSecret != "" &&
		s.AppsMetricsClientSecret != "" &&
		s.AppsMetricsProcessingClientSecret != "")
}
