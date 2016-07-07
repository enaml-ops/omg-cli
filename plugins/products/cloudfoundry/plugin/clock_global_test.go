package cloudfoundry_test

import (
	"github.com/enaml-ops/enaml"
	"github.com/enaml-ops/omg-cli/plugins/products/cloudfoundry/enaml-gen/cloud_controller_clock"
	. "github.com/enaml-ops/omg-cli/plugins/products/cloudfoundry/plugin"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gopkg.in/yaml.v2"
)

var _ = Describe("given a clock_global partition", func() {
	Context("when initialized WITHOUT a complete set of arguments", func() {
		var ig InstanceGrouper
		BeforeEach(func() {
			cf := new(Plugin)
			c := cf.GetContext([]string{
				"cloudfoundry",
			})
			ig = NewClockGlobalPartition(c)
		})

		It("should contain the appropriate jobs", func() {
			group := ig.ToInstanceGroup()
			Ω(group.GetJobByName("cloud_controller_clock")).ShouldNot(BeNil())
			Ω(group.GetJobByName("metron_agent")).ShouldNot(BeNil())
			Ω(group.GetJobByName("nfs_mounter")).ShouldNot(BeNil())
			Ω(group.GetJobByName("statsd-injector")).ShouldNot(BeNil())
		})

		It("should not have valid values", func() {
			Ω(ig.HasValidValues()).Should(BeFalse())
		})
	})

	Context("when initialized with a complete set of arguments", func() {
		var ig InstanceGrouper
		var deploymentManifest *enaml.DeploymentManifest

		BeforeEach(func() {
			cf := new(Plugin)
			c := cf.GetContext([]string{
				"cloudfoundry",
				"--skip-cert-verify=false",
				"--az", "eastprod-1",
				"--stemcell-name", "cool-ubuntu-animal",
				"--network", "foundry-net",
				"--clock-global-vm-type", "vmtype",
				"--allow-app-ssh-access",
				"--system-domain", "sys.test.com",
				"--app-domain", "apps.test.com",
				"--cc-vm-type", "ccvmtype",
				"--cc-staging-upload-user", "staginguser",
				"--cc-staging-upload-password", "stagingpassword",
				"--cc-bulk-api-user", "bulkapiuser",
				"--cc-bulk-api-password", "bulkapipassword",
				"--cc-db-encryption-key", "dbencryptionkey",
				"--cc-internal-api-user", "internalapiuser",
				"--cc-internal-api-password", "internalapipassword",
				"--cc-service-dashboards-client-secret", "ccsecret",
				"--metron-secret", "metronsecret",
				"--metron-zone", "metronzoneguid",
				"--nfs-server-address", "10.0.0.19",
				"--nfs-share-path", "/var/vcap/nfs",
				"--consul-ca-cert", "consul-ca-cert",
				"--consul-agent-cert", "consul-agent-cert",
				"--consul-agent-key", "consul-agent-key",
				"--consul-server-cert", "consulservercert",
				"--consul-server-key", "consulserverkey",
				"--consul-ip", "1.0.0.1",
				"--consul-ip", "1.0.0.2",
				"--consul-encryption-key", "consulencryptionkey",
				"--mysql-proxy-ip", "1.0.10.3",
				"--mysql-proxy-ip", "1.0.10.4",
				"--db-ccdb-username", "ccdb-user",
				"--db-ccdb-password", "ccdb-password",
				"--uaa-jwt-verification-key", "jwt-verificationkey",
				"--nats-user", "nats",
				"--nats-pass", "pass",
				"--nats-machine-ip", "1.0.0.5",
				"--nats-machine-ip", "1.0.0.6",
				"--nats-port", "4333",
			})
			ig = NewClockGlobalPartition(c)
			deploymentManifest = new(enaml.DeploymentManifest)
			deploymentManifest.AddInstanceGroup(ig.ToInstanceGroup())
		})

		It("should have valid values", func() {
			Ω(ig.HasValidValues()).Should(BeTrue())
		})

		It("then it should allow the user to configure the AZs", func() {
			group := deploymentManifest.GetInstanceGroupByName("clock_global-partition")
			Ω(len(group.AZs)).Should(Equal(1))
			Ω(group.AZs[0]).Should(Equal("eastprod-1"))
		})

		It("then it should allow the user to configure the used stemcell", func() {
			group := deploymentManifest.GetInstanceGroupByName("clock_global-partition")
			Ω(group.Stemcell).Should(Equal("cool-ubuntu-animal"))
		})

		It("then it should allow the user to configure the network to use", func() {
			group := deploymentManifest.GetInstanceGroupByName("clock_global-partition")
			Ω(len(group.Networks)).Should(Equal(1))
			Ω(group.Networks[0].Name).Should(Equal("foundry-net"))
		})

		It("then it should allow the user to configure the VM type", func() {
			group := deploymentManifest.GetInstanceGroupByName("clock_global-partition")
			Ω(group.VMType).Should(Equal("vmtype"))
		})

		It("then it should have a single instance", func() {
			group := deploymentManifest.GetInstanceGroupByName("clock_global-partition")
			Ω(len(group.AZs)).Should(Equal(1))
		})

		It("should have correctly configured the cloud controller clock", func() {
			group := deploymentManifest.GetInstanceGroupByName("clock_global-partition")
			job := group.GetJobByName("cloud_controller_clock")
			Ω(job.Release).Should(Equal(CFReleaseName))
			props := job.Properties.(*cloud_controller_clock.CloudControllerClock)
			Ω(props.Domain).Should(Equal("sys.test.com"))
			Ω(props.SystemDomain).Should(Equal("sys.test.com"))
			Ω(props.SystemDomainOrganization).Should(Equal("system"))

			ad := props.AppDomains.([]string)
			Ω(len(ad)).Should(Equal(1))
			Ω(ad[0]).Should(Equal("apps.test.com"))

			Ω(props.Cc.AllowAppSshAccess).Should(BeTrue())
			Ω(props.Cc.Buildpacks.BlobstoreType).Should(Equal("fog"))
			Ω(props.Cc.Droplets.BlobstoreType).Should(Equal("fog"))
			Ω(props.Cc.Packages.BlobstoreType).Should(Equal("fog"))
			Ω(props.Cc.ResourcePool.BlobstoreType).Should(Equal("fog"))

			Ω(props.Cc.Droplets.FogConnection).Should(HaveKeyWithValue("provider", "Local"))
			Ω(props.Cc.Packages.FogConnection).Should(HaveKeyWithValue("provider", "Local"))
			Ω(props.Cc.ResourcePool.FogConnection).Should(HaveKeyWithValue("provider", "Local"))
			Ω(props.Cc.Droplets.FogConnection).Should(HaveKeyWithValue("local_root", "/var/vcap/nfs/shared"))
			Ω(props.Cc.Packages.FogConnection).Should(HaveKeyWithValue("local_root", "/var/vcap/nfs/shared"))
			Ω(props.Cc.ResourcePool.FogConnection).Should(HaveKeyWithValue("local_root", "/var/vcap/nfs/shared"))

			Ω(props.Cc.LoggingLevel).Should(Equal("debug"))
			Ω(props.Cc.MaximumHealthCheckTimeout).Should(Equal(600))
			Ω(props.Cc.StagingUploadUser).Should(Equal("staginguser"))
			Ω(props.Cc.StagingUploadPassword).Should(Equal("stagingpassword"))
			Ω(props.Cc.BulkApiUser).Should(Equal("bulkapiuser"))
			Ω(props.Cc.BulkApiPassword).Should(Equal("bulkapipassword"))
			Ω(props.Cc.DbEncryptionKey).Should(Equal("dbencryptionkey"))
			Ω(props.Cc.InternalApiUser).Should(Equal("internalapiuser"))
			Ω(props.Cc.InternalApiPassword).Should(Equal("internalapipassword"))

			quotaDefs := `default:
  memory_limit: 10240
  total_services: 100
  non_basic_services_allowed: true
  total_routes: 1000
  trial_db_allowed: true
runaway:
  memory_limit: 102400
  total_services: -1
  total_routes: 1000
  non_basic_services_allowed: true
`
			b, _ := yaml.Marshal(props.Cc.QuotaDefinitions)
			Ω(string(b)).Should(MatchYAML(quotaDefs))

			sgDefs := `- name: all_open
  rules:
    - protocol: all
      destination: 0.0.0.0-255.255.255.255
`
			b, _ = yaml.Marshal(props.Cc.SecurityGroupDefinitions)
			Ω(string(b)).Should(MatchYAML(sgDefs))

			Ω(props.Ccdb.Address).Should(Equal("1.0.10.3"))
			Ω(props.Ccdb.Port).Should(Equal(3306))
			Ω(props.Ccdb.DbScheme).Should(Equal("mysql"))

			Ω(props.Ccdb.Roles).Should(HaveKeyWithValue("tag", "admin"))
			Ω(props.Ccdb.Roles).Should(HaveKeyWithValue("name", "ccdb-user"))
			Ω(props.Ccdb.Roles).Should(HaveKeyWithValue("password", "ccdb-password"))

			Ω(props.Ccdb.Databases).Should(HaveKeyWithValue("tag", "cc"))
			Ω(props.Ccdb.Databases).Should(HaveKeyWithValue("name", "ccdb"))
			Ω(props.Ccdb.Databases).Should(HaveKeyWithValue("citext", "true"))

			Ω(props.Uaa.Url).Should(Equal("https://uaa.sys.test.com"))
			Ω(props.Uaa.Jwt).ShouldNot(BeNil())
			Ω(props.Uaa.Jwt.VerificationKey).Should(Equal("jwt-verificationkey"))

			Ω(props.Uaa.Clients).ShouldNot(BeNil())
			Ω(props.Uaa.Clients.CcServiceDashboards.Secret).Should(Equal("ccsecret"))

			Ω(props.LoggerEndpoint.Port).Should(Equal("443"))
			Ω(props.Ssl.SkipCertVerify).Should(BeFalse())

			Ω(props.Nats.User).Should(Equal("nats"))
			Ω(props.Nats.Password).Should(Equal("pass"))
			Ω(props.Nats.Port).Should(Equal(4333))
			Ω(props.Nats.Machines).Should(ConsistOf("1.0.0.5", "1.0.0.6"))
		})
	})
})
