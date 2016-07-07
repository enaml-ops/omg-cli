package cloudfoundry_test

import (
	"github.com/enaml-ops/omg-cli/plugins/products/cloudfoundry/enaml-gen/route_registrar"
	"github.com/enaml-ops/omg-cli/plugins/products/cloudfoundry/enaml-gen/uaa"
	. "github.com/enaml-ops/omg-cli/plugins/products/cloudfoundry/plugin"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("UAA Partition", func() {
	Context("when initialized WITHOUT a complete set of arguments", func() {
		It("then HasValidValues returns false", func() {
			plugin := new(Plugin)
			c := plugin.GetContext([]string{
				"cloudfoundry",
			})
			uaaPartition := NewUAAPartition(c)
			Ω(uaaPartition.HasValidValues()).Should(BeFalse())
		})
	})
	Context("when initialized WITH a complete set of arguments", func() {
		var uaaPartition InstanceGrouper
		BeforeEach(func() {
			plugin := new(Plugin)
			c := plugin.GetContext([]string{
				"cloudfoundry",
				"--network", "foundry-net",
				"--stemcell-name", "cool-ubuntu-animal",
				"--az", "eastprod-1",
				"--network", "foundry-net",
				"--uaa-vm-type", "blah",
				"--uaa-instances", "1",
				"--system-domain", "sys.test.com",
				"--consul-ip", "1.0.0.1",
				"--consul-ip", "1.0.0.2",
				"--consul-encryption-key", "consulencryptionkey",
				"--consul-ca-cert", "consul-ca-cert",
				"--consul-agent-cert", "consul-agent-cert",
				"--consul-agent-key", "consul-agent-key",
				"--consul-server-cert", "consulservercert",
				"--consul-server-key", "consulserverkey",
				"--metron-secret", "metronsecret",
				"--metron-zone", "metronzoneguid",
				"--syslog-address", "syslog-server",
				"--syslog-port", "10601",
				"--syslog-transport", "tcp",
				"--etcd-machine-ip", "1.0.0.7",
				"--etcd-machine-ip", "1.0.0.8",
				"--nats-user", "nats",
				"--nats-pass", "pass",
				"--nats-machine-ip", "1.0.0.5",
				"--nats-machine-ip", "1.0.0.6",
				"--uaa-saml-service-provider-key", "saml-key",
				"--uaa-saml-service-provider-certificate", "saml-cert",
				"--uaa-jwt-verification-key", "jwt-verificationkey",
				"--uaa-jwt-signing-key", "jwt-signingkey",
				"--uaa-ldap-enabled",
				"--uaa-ldap-url", "ldap://ldap.test.com",
				"--uaa-ldap-user-dn", "userdn",
				"--uaa-ldap-user-password", "userpwd",
				"--uaa-ldap-search-filter", "filter",
				"--uaa-ldap-search-base", "base",
				"--uaa-ldap-mail-attributename", "mail",
				"--uaa-admin-secret", "adminclientsecret",
				"--router-ip", "1.0.0.1",
				"--router-ip", "1.0.0.2",
				"--mysql-proxy-ip", "1.0.10.3",
				"--mysql-proxy-ip", "1.0.10.4",
				"--db-uaa-username", "uaa-db-user",
				"--db-uaa-password", "uaa-db-pwd",
				"--admin-password", "admin",
				"--push-apps-manager-password", "appsman",
				"--smoke-tests-password", "smoke",
				"--system-services-password", "sysservices",
				"--system-verification-password", "sysverification",
				"--opentsdb-firehose-nozzle-client-secret", "opentsdb-firehose-nozzle-client-secret",
				"--identity-client-secret", "identity-client-secret",
				"--login-client-secret", "login-client-secret",
				"--portal-client-secret", "portal-client-secret",
				"--autoscaling-service-client-secret", "autoscaling-service-client-secret",
				"--system-passwords-client-secret", "system-passwords-client-secret",
				"--cc-service-dashboards-client-secret", "cc-service-dashboards-client-secret",
				"--doppler-client-secret", "doppler-client-secret",
				"--gorouter-client-secret", "gorouter-client-secret",
				"--notifications-client-secret", "notifications-client-secret",
				"--notifications-ui-client-secret", "notifications-ui-client-secret",
				"--cloud-controller-username-lookup-client-secret", "cloud-controller-username-lookup-client-secret",
				"--cc-routing-client-secret", "cc-routing-client-secret",
				"--ssh-proxy-client-secret", "ssh-proxy-client-secret",
				"--apps-metrics-client-secret", "apps-metrics-client-secret",
				"--apps-metrics-processing-client-secret", "apps-metrics-processing-client-secret",
			})
			uaaPartition = NewUAAPartition(c)
		})
		It("then HasValidValues should return true", func() {
			Ω(uaaPartition.HasValidValues()).Should(Equal(true))
		})
		It("then it should not configure static ips for uaaPartition", func() {
			ig := uaaPartition.ToInstanceGroup()
			network := ig.Networks[0]
			Ω(len(network.StaticIPs)).Should(Equal(0))
		})
		It("then it should have 1 instances", func() {
			ig := uaaPartition.ToInstanceGroup()
			Ω(ig.Instances).Should(Equal(1))
		})
		It("then it should allow the user to configure the AZs", func() {
			ig := uaaPartition.ToInstanceGroup()
			Ω(len(ig.AZs)).Should(Equal(1))
			Ω(ig.AZs[0]).Should(Equal("eastprod-1"))
		})

		It("then it should allow the user to configure vm-type", func() {
			ig := uaaPartition.ToInstanceGroup()
			Ω(ig.VMType).ShouldNot(BeEmpty())
			Ω(ig.VMType).Should(Equal("blah"))
		})

		It("then it should allow the user to configure network to use", func() {
			ig := uaaPartition.ToInstanceGroup()
			network := ig.GetNetworkByName("foundry-net")
			Ω(network).ShouldNot(BeNil())
		})

		It("then it should allow the user to configure the used stemcell", func() {
			ig := uaaPartition.ToInstanceGroup()
			Ω(ig.Stemcell).ShouldNot(BeEmpty())
			Ω(ig.Stemcell).Should(Equal("cool-ubuntu-animal"))
		})
		It("then it should have update max in-flight 1 and serial false", func() {
			ig := uaaPartition.ToInstanceGroup()
			Ω(ig.Update.MaxInFlight).Should(Equal(1))
			Ω(ig.Update.Serial).Should(Equal(false))
		})

		It("then it should then have 5 jobs", func() {
			ig := uaaPartition.ToInstanceGroup()
			Ω(len(ig.Jobs)).Should(Equal(5))
		})
		It("then it should then have uaa job with client secret", func() {
			ig := uaaPartition.ToInstanceGroup()
			job := ig.GetJobByName("uaa")
			Ω(job).ShouldNot(BeNil())
			props, _ := job.Properties.(*uaa.Uaa)
			Ω(props.Admin).ShouldNot(BeNil())
			Ω(props.Admin.ClientSecret).Should(Equal("adminclientsecret"))
		})
		It("then it should then have uaa job with proxy configured", func() {
			ig := uaaPartition.ToInstanceGroup()
			job := ig.GetJobByName("uaa")
			Ω(job).ShouldNot(BeNil())
			props, _ := job.Properties.(*uaa.Uaa)
			Ω(props.Proxy).ShouldNot(BeNil())
			Ω(props.Proxy.Servers).Should(ConsistOf("1.0.0.1", "1.0.0.2"))
		})
		It("then it should then have uaa job with UAADB", func() {
			ig := uaaPartition.ToInstanceGroup()
			job := ig.GetJobByName("uaa")
			Ω(job).ShouldNot(BeNil())
			props, _ := job.Properties.(*uaa.Uaa)
			Ω(props.Uaadb).ShouldNot(BeNil())
			Ω(props.Uaadb.DbScheme).Should(Equal("mysql"))
			Ω(props.Uaadb.Port).Should(Equal(3306))
			Ω(props.Uaadb.Address).Should(Equal("1.0.10.3"))
			Ω(props.Uaadb.Roles).ShouldNot(BeNil())
			roles := props.Uaadb.Roles.(map[string]string)
			Ω(roles["tag"]).Should(Equal("admin"))
			Ω(roles["name"]).Should(Equal("uaa-db-user"))
			Ω(roles["password"]).Should(Equal("uaa-db-pwd"))
			Ω(props.Uaadb.Databases).ShouldNot(BeNil())
			dbs := props.Uaadb.Databases.(map[string]string)
			Ω(dbs["tag"]).Should(Equal("uaa"))
			Ω(dbs["name"]).Should(Equal("uaa"))
		})
		It("then it should then have uaa job with Clients", func() {
			ig := uaaPartition.ToInstanceGroup()
			job := ig.GetJobByName("uaa")
			Ω(job).ShouldNot(BeNil())
			props, _ := job.Properties.(*uaa.Uaa)
			Ω(props.Clients).ShouldNot(BeNil())
			clientMap := props.Clients.(map[string]UAAClient)
			Ω(len(clientMap)).Should(Equal(19))

		})
		It("then it should then have uaa job with SCIM", func() {
			ig := uaaPartition.ToInstanceGroup()
			job := ig.GetJobByName("uaa")
			Ω(job).ShouldNot(BeNil())
			props, _ := job.Properties.(*uaa.Uaa)
			Ω(props.Scim).ShouldNot(BeNil())
			Ω(props.Scim.User).ShouldNot(BeNil())
			Ω(props.Scim.User.Override).Should(BeTrue())
			Ω(props.Scim.UseridsEnabled).Should(BeTrue())
			Ω(props.Scim.Users).ShouldNot(BeNil())
			users := props.Scim.Users.([]string)
			Ω(len(users)).Should(Equal(5))
		})
		It("then it should then have uaa job with valid login information", func() {
			ig := uaaPartition.ToInstanceGroup()
			job := ig.GetJobByName("uaa")
			Ω(job).ShouldNot(BeNil())
			props, _ := job.Properties.(*uaa.Uaa)
			Ω(props.Domain).Should(Equal("sys.test.com"))
		})
		It("then it should then have uaa job with valid uaa information", func() {
			ig := uaaPartition.ToInstanceGroup()
			job := ig.GetJobByName("uaa")
			Ω(job).ShouldNot(BeNil())
			props, _ := job.Properties.(*uaa.Uaa)
			Ω(props.Uaa).ShouldNot(BeNil())
			Ω(props.Uaa.CatalinaOpts).Should(Equal("-Xmx768m -XX:MaxPermSize=256m"))
			Ω(props.Uaa.RequireHttps).Should(BeTrue())
			Ω(props.Uaa.Url).Should(Equal("https://uaa.sys.test.com"))
			Ω(props.Uaa.Ssl).ShouldNot(BeNil())
			Ω(props.Uaa.Ssl.Port).Should(Equal(-1))

			Ω(props.Uaa.Authentication).ShouldNot(BeNil())
			Ω(props.Uaa.Authentication.Policy).ShouldNot(BeNil())
			Ω(props.Uaa.Authentication.Policy.LockoutAfterFailures).Should(Equal(5))

			Ω(props.Uaa.Password).ShouldNot(BeNil())
			Ω(props.Uaa.Password.Policy).ShouldNot(BeNil())
			Ω(props.Uaa.Password.Policy.MinLength).Should(Equal(0))
			Ω(props.Uaa.Password.Policy.RequireLowerCaseCharacter).Should(Equal(0))
			Ω(props.Uaa.Password.Policy.RequireUpperCaseCharacter).Should(Equal(0))
			Ω(props.Uaa.Password.Policy.RequireDigit).Should(Equal(0))
			Ω(props.Uaa.Password.Policy.RequireSpecialCharacter).Should(Equal(0))
			Ω(props.Uaa.Password.Policy.ExpirePasswordInMonths).Should(Equal(0))

			Ω(props.Uaa.Jwt).ShouldNot(BeNil())
			Ω(props.Uaa.Jwt.SigningKey).Should(Equal("jwt-signingkey"))
			Ω(props.Uaa.Jwt.VerificationKey).Should(Equal("jwt-verificationkey"))

			Ω(props.Uaa.Ldap).ShouldNot(BeNil())
			Ω(props.Uaa.Ldap.Enabled).Should(BeTrue())
			Ω(props.Uaa.Ldap.Url).Should(Equal("ldap://ldap.test.com"))
			Ω(props.Uaa.Ldap.UserDN).Should(Equal("userdn"))
			Ω(props.Uaa.Ldap.UserPassword).Should(Equal("userpwd"))
			Ω(props.Uaa.Ldap.SearchBase).Should(Equal("base"))
			Ω(props.Uaa.Ldap.SearchFilter).Should(Equal("filter"))
			Ω(props.Uaa.Ldap.MailAttributeName).Should(Equal("mail"))
			Ω(props.Uaa.Ldap.ProfileType).Should(Equal("search-and-bind"))
			Ω(props.Uaa.Ldap.SslCertificate).Should(Equal(""))
			Ω(props.Uaa.Ldap.SslCertificateAlias).Should(Equal(""))
			Ω(props.Uaa.Ldap.Groups).ShouldNot(BeNil())
			Ω(props.Uaa.Ldap.Groups.ProfileType).Should(Equal("no-groups"))
			Ω(props.Uaa.Ldap.Groups.SearchBase).Should(Equal(""))
			Ω(props.Uaa.Ldap.Groups.GroupSearchFilter).Should(Equal(""))
		})
		It("then it should then have uaa job with valid login information", func() {
			ig := uaaPartition.ToInstanceGroup()
			job := ig.GetJobByName("uaa")
			Ω(job).ShouldNot(BeNil())
			props, _ := job.Properties.(*uaa.Uaa)
			Ω(props.Login).ShouldNot(BeNil())
			Ω(props.Login.SelfServiceLinksEnabled).Should(BeTrue())
			Ω(props.Login.SignupsEnabled).Should(BeTrue())
			Ω(props.Login.Protocol).Should(Equal("https"))
			Ω(props.Login.UaaBase).Should(Equal("https://uaa.sys.test.com"))
			Ω(props.Login.Branding).ShouldNot(BeNil())

			Ω(props.Login.Links).ShouldNot(BeNil())
			links := props.Login.Links.(*uaa.Links)
			Ω(links.Passwd).Should(Equal("https://login.sys.test.com/forgot_password"))
			Ω(links.Signup).Should(Equal("https://login.sys.test.com/create_account"))

			Ω(props.Login.Notifications).ShouldNot(BeNil())
			Ω(props.Login.Notifications.Url).Should(Equal("https://notifications.sys.test.com"))

			Ω(props.Login.Saml).ShouldNot(BeNil())
			Ω(props.Login.Saml.Entityid).Should(Equal("https://login.sys.test.com"))
			Ω(props.Login.Saml.SignRequest).Should(BeTrue())
			Ω(props.Login.Saml.WantAssertionSigned).Should(BeFalse())
			Ω(props.Login.Saml.ServiceProviderKey).Should(Equal("saml-key"))
			Ω(props.Login.Saml.ServiceProviderCertificate).Should(Equal("saml-cert"))

			Ω(props.Login.Logout).ShouldNot(BeNil())
			Ω(props.Login.Logout.Redirect).ShouldNot(BeNil())
			Ω(props.Login.Logout.Redirect.Url).Should(Equal("/login"))
			Ω(props.Login.Logout.Redirect.Parameter).ShouldNot(BeNil())
			Ω(props.Login.Logout.Redirect.Parameter.Disable).Should(BeFalse())
			Ω(props.Login.Logout.Redirect.Parameter.Whitelist).Should(ConsistOf("https://console.sys.test.com", "https://apps.sys.test.com"))
		})
		It("then it should then have route_registrar job", func() {
			ig := uaaPartition.ToInstanceGroup()
			job := ig.GetJobByName("route_registrar")
			Ω(job).ShouldNot(BeNil())
			props, _ := job.Properties.(*route_registrar.RouteRegistrar)
			Ω(props.Nats).ShouldNot(BeNil())
			Ω(props.Nats.User).Should(Equal("nats"))
			Ω(props.Nats.Password).Should(Equal("pass"))
			Ω(props.Nats.Port).Should(Equal(4222))
			Ω(props.Nats.Machines).Should(ConsistOf("1.0.0.5", "1.0.0.6"))
			Ω(props.Routes).ShouldNot(BeNil())
			routes := props.Routes.(map[string]interface{})
			Ω(routes["name"]).Should(Equal("uaa"))
			Ω(routes["port"]).Should(Equal(8080))
			Ω(routes["registration_interval"]).Should(Equal("40s"))
			Ω(routes["uris"]).Should(ConsistOf("uaa.sys.test.com", "*.uaa.sys.test.com", "login.sys.test.com", "*.login.sys.test.com"))
		})
	})
	Context("when Creating Branding with flags", func() {
		var branding *uaa.Branding
		BeforeEach(func() {
			plugin := new(Plugin)
			c := plugin.GetContext([]string{
				"cloudfoundry",
				"--uaa-company-name", "company",
				"--uaa-product-logo", "product-logo",
				"--uaa-square-logo", "square-logo",
				"--uaa-footer-legal-txt", "legal",
			})
			branding = CreateBranding(c)
		})
		It("branding should be initialized", func() {
			Ω(branding).ShouldNot(BeNil())
			Ω(branding.CompanyName).Should(Equal("company"))
			Ω(branding.ProductLogo).Should(Equal("product-logo"))
			Ω(branding.SquareLogo).Should(Equal("square-logo"))
			Ω(branding.FooterLegalText).Should(Equal("legal"))
			Ω(branding.FooterLinks).Should(BeNil())
		})
	})
})
