package boshinit_test

import (
	"os"

	"github.com/enaml-ops/enaml"
	boshinit "github.com/enaml-ops/omg-cli/plugins/products/bosh-init"
	"github.com/enaml-ops/omg-cli/plugins/products/bosh-init/enaml-gen/health_monitor"
	"github.com/enaml-ops/omg-cli/plugins/products/bosh-init/enaml-gen/uaa"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type nopIaaSProvider struct{}

func (nopIaaSProvider) CreateCPIRelease() (r enaml.Release)                { return }
func (nopIaaSProvider) CreateCPITemplate() (r enaml.Template)              { return }
func (nopIaaSProvider) CreateDiskPool() (r enaml.DiskPool)                 { return }
func (nopIaaSProvider) CreateResourcePool() (r enaml.ResourcePool)         { return }
func (nopIaaSProvider) CreateManualNetwork() (r enaml.ManualNetwork)       { return }
func (nopIaaSProvider) CreateVIPNetwork() (r enaml.VIPNetwork)             { return }
func (nopIaaSProvider) CreateJobNetwork() (r *enaml.Network)               { return }
func (nopIaaSProvider) CreateCloudProvider() (r enaml.CloudProvider)       { return }
func (nopIaaSProvider) CreateCPIJobProperties() (r map[string]interface{}) { return }
func (nopIaaSProvider) CreateDeploymentManifest() *enaml.DeploymentManifest {
	return new(enaml.DeploymentManifest)
}

var _ = Describe("given boshbase", func() {

	const (
		controlSecret          = "health-monitor-secret"
		controlCACert          = "health-monitor-ca-cert"
		controlGraphiteAddress = "graphite.your.org"
		controlSyslogAddress   = "syslog.your.org"
	)

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

		FIt("should create a proper list of clients", func() {
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
		})

		AfterEach(func() {
			os.Remove("./rootCA.pem")
			os.Remove("./director.pwd")
		})

		It("creates authentication files when configured to print manfiest", func() {
			bb.PrintManifest = true
			bb.HandleDeployment(nopIaaSProvider{}, nopDeploy)

			Ω("./rootCA.pem").Should(BeAnExistingFile())
			Ω("./director.pwd").Should(BeAnExistingFile())
		})

		It("creates authentication files when configured to deploy", func() {
			bb.PrintManifest = false
			bb.HandleDeployment(nopIaaSProvider{}, nopDeploy)

			Ω("./rootCA.pem").Should(BeAnExistingFile())
			Ω("./director.pwd").Should(BeAnExistingFile())
		})
	})
})
