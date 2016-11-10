package boshinit_test

import (
	"os"

	"github.com/enaml-ops/enaml"
	boshinit "github.com/enaml-ops/omg-cli/plugins/products/bosh-init"
	"github.com/enaml-ops/omg-cli/plugins/products/bosh-init/enaml-gen/director"
	"github.com/enaml-ops/omg-cli/plugins/products/bosh-init/enaml-gen/health_monitor"
	"github.com/enaml-ops/omg-cli/plugins/products/bosh-init/enaml-gen/registry"
	"github.com/enaml-ops/omg-cli/plugins/products/bosh-init/enaml-gen/uaa"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type nopIaaSProvider struct{}

func (nopIaaSProvider) CreateCPIRelease() (r enaml.Release)                    { return }
func (nopIaaSProvider) CreateCPITemplate() (r enaml.Template)                  { return }
func (nopIaaSProvider) CreateDiskPool() (r enaml.DiskPool)                     { return }
func (nopIaaSProvider) CreateResourcePool() (r *enaml.ResourcePool, err error) { return }
func (nopIaaSProvider) CreateManualNetwork() (r enaml.ManualNetwork)           { return }
func (nopIaaSProvider) CreateVIPNetwork() (r enaml.VIPNetwork)                 { return }
func (nopIaaSProvider) CreateJobNetwork() (r *enaml.Network)                   { return }
func (nopIaaSProvider) CreateCloudProvider() (r enaml.CloudProvider)           { return }
func (nopIaaSProvider) CreateCPIJobProperties() (r map[string]interface{})     { return }
func (nopIaaSProvider) CreateDeploymentManifest() (*enaml.DeploymentManifest, error) {
	return &enaml.DeploymentManifest{}, nil
}

var _ = Describe("given boshbase", func() {

	const (
		controlSecret          = "health-monitor-secret"
		controlCACert          = "health-monitor-ca-cert"
		controlGraphiteAddress = "graphite.your.org"
		controlSyslogAddress   = "syslog.your.org"
	)
	Context("when configured for Internal Postgresql DB", func() {
		var bb *boshinit.BoshBase
		var job enaml.Job
		BeforeEach(func() {
			bb = &boshinit.BoshBase{
				UseExternalDB: false,
				Mode:          "uaa",
			}
			bb.InitializeDBDefaults()
			job = bb.CreateJob()
			Ω(job).ShouldNot(BeNil())
		})

		It("should configure bosh to use internal postgresql for director", func() {
			Ω(job.Properties).Should(HaveKey("postgres"))
			Ω(job.Properties).Should(HaveKey("director"))
			director := job.Properties["director"].(*director.Director)
			Ω(director.Db.Adapter).Should(Equal("postgres"))
			Ω(director.Db.Host).Should(Equal("127.0.0.1"))
			Ω(director.Db.Port).Should(Equal(5432))
			Ω(director.Db.User).Should(Equal("postgres"))
			Ω(director.Db.Password).ShouldNot(BeNil())
			Ω(director.Db.Database).Should(Equal("bosh"))
		})

		It("should configure bosh to use internal postgresql for registry", func() {
			Ω(job.Properties).Should(HaveKey("postgres"))
			Ω(job.Properties).Should(HaveKey("registry"))
			registry := job.Properties["registry"].(*registry.Registry)
			Ω(registry.Db.Adapter).Should(Equal("postgres"))
			Ω(registry.Db.Host).Should(Equal("127.0.0.1"))
			Ω(registry.Db.Port).Should(Equal(5432))
			Ω(registry.Db.User).Should(Equal("postgres"))
			Ω(registry.Db.Password).ShouldNot(BeNil())
			Ω(registry.Db.Database).Should(Equal("registry"))
		})

		It("should configure uaa to use internal postgresql for uaa db", func() {
			Ω(job.Properties).Should(HaveKey("postgres"))
			Ω(job.Properties).Should(HaveKey("uaadb"))
			uaaDB := job.Properties["uaadb"].(*uaa.Uaadb)
			Ω(uaaDB.DbScheme).Should(Equal("postgresql"))
			Ω(uaaDB.Address).Should(Equal("127.0.0.1"))
			Ω(uaaDB.Port).Should(Equal(5432))
			roles := uaaDB.Roles.([]interface{})[0].(map[string]string)
			Ω(roles["name"]).Should(Equal("postgres"))
			Ω(roles["password"]).ShouldNot(BeNil())
			dbs := uaaDB.Databases.([]interface{})[0].(map[string]string)
			Ω(dbs["name"]).Should(Equal("uaa"))
		})
	})
	Context("when configured for External DB", func() {
		var bb *boshinit.BoshBase
		var job enaml.Job
		BeforeEach(func() {
			bb = &boshinit.BoshBase{
				Mode:                 "uaa",
				UseExternalDB:        true,
				DatabaseDriver:       "mysql2",
				DatabaseHost:         "random.com",
				DatabasePort:         3306,
				DirectorDatabaseName: "bosh",
				RegistryDatabaseName: "registry",
				UAADatabaseName:      "uaa",
				DatabasePassword:     "db.password",
				DatabaseUsername:     "db.user",
				DatabaseScheme:       "mysql",
			}
			job = bb.CreateJob()
			Ω(job).ShouldNot(BeNil())
		})

		It("should configure bosh to use exernal db for bosh db", func() {
			Ω(job.Properties).ShouldNot(HaveKey("postgres"))
			Ω(job.Properties).Should(HaveKey("director"))
			director := job.Properties["director"].(*director.Director)
			Ω(director.Db.Adapter).Should(Equal("mysql2"))
			Ω(director.Db.Host).Should(Equal("random.com"))
			Ω(director.Db.Port).Should(Equal(3306))
			Ω(director.Db.User).Should(Equal("db.user"))
			Ω(director.Db.Password).Should(Equal("db.password"))
			Ω(director.Db.Database).Should(Equal("bosh"))
		})

		It("should configure bosh to use exernal db for registry db", func() {
			Ω(job.Properties).ShouldNot(HaveKey("postgres"))
			Ω(job.Properties).Should(HaveKey("registry"))
			registry := job.Properties["registry"].(*registry.Registry)
			Ω(registry.Db.Adapter).Should(Equal("mysql2"))
			Ω(registry.Db.Port).Should(Equal(3306))
			Ω(registry.Db.User).Should(Equal("db.user"))
			Ω(registry.Db.Password).Should(Equal("db.password"))
			Ω(registry.Db.Database).Should(Equal("registry"))
		})

		It("should configure uaa to use external db for uaa db", func() {
			Ω(job.Properties).Should(HaveKey("uaadb"))
			uaaDB := job.Properties["uaadb"].(*uaa.Uaadb)
			Ω(uaaDB.DbScheme).Should(Equal("mysql"))
			Ω(uaaDB.Address).Should(Equal("random.com"))
			Ω(uaaDB.Port).Should(Equal(3306))
			roles := uaaDB.Roles.([]interface{})[0].(map[string]string)
			Ω(roles["name"]).Should(Equal("db.user"))
			Ω(roles["password"]).Should(Equal("db.password"))
			dbs := uaaDB.Databases.([]interface{})[0].(map[string]string)
			Ω(dbs["name"]).Should(Equal("uaa"))
		})

	})
	Context("when configured for UAA", func() {
		var bb *boshinit.BoshBase
		var job enaml.Job

		BeforeEach(func() {
			bb = &boshinit.BoshBase{
				Mode:                "uaa",
				HealthMonitorSecret: controlSecret,
				CACert:              controlCACert,
				GraphiteAddress:     controlGraphiteAddress,
				GraphitePort:        2003,
			}
			job = bb.CreateJob()
			Ω(bb.IsUAA()).Should(BeTrue())
		})

		It("should create a proper list of clients", func() {
			Ω(job.Properties).Should(HaveKey("uaa"))
			uaa := job.Properties["uaa"].(*uaa.Uaa)
			Ω(uaa.Clients).Should(HaveKey("bosh_cli"))
			Ω(uaa.Clients).Should(HaveKey("health_monitor"))
			Ω(uaa.Clients).Should(HaveKey("director"))
			Ω(uaa.Clients).Should(HaveKey("login"))
		})

		It("configures health monitor", func() {
			Ω(job.Properties).Should(HaveKey("hm"))
			hm := job.Properties["hm"].(*health_monitor.Hm)
			Ω(hm.ResurrectorEnabled).Should(BeTrue())

			Ω(hm.DirectorAccount.CaCert).ShouldNot(BeEmpty())
			Ω(hm.DirectorAccount.ClientId).Should(Equal("health_monitor"))
			Ω(hm.DirectorAccount.ClientSecret).Should(Equal(controlSecret))

			Ω(hm.DirectorAccount.User).Should(BeNil())
			Ω(hm.DirectorAccount.Password).Should(BeNil())

			Ω(hm.GraphiteEnabled).Should(BeTrue())
			Ω(hm.Graphite.Address).Should(Equal(controlGraphiteAddress))
			Ω(hm.Graphite.Port).Should(Equal(2003))

			Ω(hm.SyslogEventForwarderEnabled).Should(BeNil())
			Ω(hm.SyslogEventForwarder).Should(BeNil())
		})

		It("configures director", func() {
			Ω(job.Properties).Should(HaveKey("director"))
			director := job.Properties["director"].(*director.Director)
			Ω(director.GenerateVmPasswords).Should(BeTrue())

		})

		It("configures password for vm", func() {
			rp, err := bb.CreateResourcePool(func() interface{} {
				return nil
			})
			Ω(err).ShouldNot(HaveOccurred())
			Ω(rp).ShouldNot(BeNil())
			Ω(rp.Env).ShouldNot(BeNil())
			Ω(rp.Env["bosh"]).ShouldNot(BeNil())
			pwd := rp.Env["bosh"].(boshinit.BoshPassword)
			Ω(pwd.Password).ShouldNot(BeNil())
		})
	})

	Context("when configured for basic auth", func() {
		var bb *boshinit.BoshBase
		var job enaml.Job

		BeforeEach(func() {
			bb = &boshinit.BoshBase{
				Mode:                "basic",
				HealthMonitorSecret: controlSecret,
				SyslogAddress:       controlSyslogAddress,
				SyslogPort:          5514,
				SyslogTransport:     "tcp",
			}
			job = bb.CreateJob()
			Ω(bb.IsUAA()).Should(BeFalse())
		})

		It("configures health monitor", func() {
			Ω(job.Properties).Should(HaveKey("hm"))
			hm := job.Properties["hm"].(*health_monitor.Hm)
			Ω(hm.ResurrectorEnabled).Should(BeTrue())

			Ω(hm.DirectorAccount.User).Should(Equal("hm"))
			Ω(hm.DirectorAccount.Password).Should(Equal(controlSecret))

			Ω(hm.DirectorAccount.CaCert).Should(BeNil())
			Ω(hm.DirectorAccount.ClientId).Should(BeNil())
			Ω(hm.DirectorAccount.ClientSecret).Should(BeNil())

			Ω(hm.GraphiteEnabled).Should(BeNil())
			Ω(hm.Graphite).Should(BeNil())

			Ω(hm.SyslogEventForwarderEnabled).Should(BeTrue())
			Ω(hm.SyslogEventForwarder).ShouldNot(BeNil())
			Ω(hm.SyslogEventForwarder.Address).Should(Equal(controlSyslogAddress))
			Ω(hm.SyslogEventForwarder.Port).Should(Equal(5514))
			Ω(hm.SyslogEventForwarder.Transport).Should(Equal("tcp"))
		})

		It("configures director", func() {
			Ω(job.Properties).Should(HaveKey("director"))
			director := job.Properties["director"].(*director.Director)
			Ω(director.GenerateVmPasswords).Should(BeTrue())
		})
		It("configures password for vm", func() {
			rp, err := bb.CreateResourcePool(func() interface{} {
				return nil
			})
			Ω(err).ShouldNot(HaveOccurred())
			Ω(rp).ShouldNot(BeNil())
			Ω(rp.Env).ShouldNot(BeNil())
			Ω(rp.Env["bosh"]).ShouldNot(BeNil())
			pwd := rp.Env["bosh"].(boshinit.BoshPassword)
			Ω(pwd.Password).ShouldNot(BeNil())
		})
	})

	Context("handle deployment", func() {
		var bb *boshinit.BoshBase

		const (
			controlPassword = "director-password"
			controlCACert   = "ca-cert"
		)

		nopDeploy := func(s string) {}

		BeforeEach(func() {
			bb = &boshinit.BoshBase{
				DirectorPassword: controlPassword,
				CACert:           controlCACert,
			}

			Ω("./rootCA.pem").ShouldNot(BeAnExistingFile())
			Ω("./director.pwd").ShouldNot(BeAnExistingFile())
			Ω("./nats.pwd").ShouldNot(BeAnExistingFile())
		})

		AfterEach(func() {
			os.Remove("./rootCA.pem")
			os.Remove("./director.pwd")
			os.Remove("./nats.pwd")
		})

		It("creates authentication files when configured to print manfiest", func() {
			bb.PrintManifest = true
			bb.HandleDeployment(nopIaaSProvider{}, nopDeploy)

			Ω("./rootCA.pem").Should(BeAnExistingFile())
			Ω("./director.pwd").Should(BeAnExistingFile())
			Ω("./nats.pwd").Should(BeAnExistingFile())
		})

		It("creates authentication files when configured to deploy", func() {
			bb.PrintManifest = false
			bb.HandleDeployment(nopIaaSProvider{}, nopDeploy)

			Ω("./rootCA.pem").Should(BeAnExistingFile())
			Ω("./director.pwd").Should(BeAnExistingFile())
			Ω("./nats.pwd").Should(BeAnExistingFile())
		})
	})
})
