package cloudfoundry_test

import (
	"github.com/enaml-ops/omg-cli/plugins/products/cloudfoundry/enaml-gen/metron_agent"
	. "github.com/enaml-ops/omg-cli/plugins/products/cloudfoundry/plugin"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Metron", func() {
	Context("when initialized WITHOUT a complete set of arguments", func() {
		It("then HasValidValues should return false", func() {

			plugin := new(Plugin)
			c := plugin.GetContext([]string{
				"cloudfoundry",
			})

			Ω(NewMetron(c).HasValidValues()).Should(BeFalse())

		})
	})
	Context("when initialized WITH a complete set of arguments", func() {
		var metron *Metron
		BeforeEach(func() {
			plugin := new(Plugin)
			c := plugin.GetContext([]string{
				"cloudfoundry",
				"--metron-secret", "metronsecret",
				"--metron-zone", "metronzoneguid",
				"--syslog-address", "syslog-server",
				"--syslog-port", "10601",
				"--syslog-transport", "tcp",
				"--etcd-machine-ip", "1.0.0.7",
				"--etcd-machine-ip", "1.0.0.8",
			})
			metron = NewMetron(c)
		})
		It("then HasValidValues should return true", func() {
			Ω(metron.HasValidValues()).Should(BeTrue())
		})
		It("then it should allow the user to configure the metron agent", func() {
			job := metron.CreateJob()
			Ω(job).ShouldNot(BeNil())
			props, _ := job.Properties.(*metron_agent.MetronAgent)
			Ω(props.MetronAgent.Zone).Should(Equal("metronzoneguid"))
			Ω(props.SyslogDaemonConfig.Address).Should(Equal("syslog-server"))
			Ω(props.SyslogDaemonConfig.Port).Should(Equal(10601))
			Ω(props.SyslogDaemonConfig.Transport).Should(Equal("tcp"))
			Ω(props.MetronEndpoint.SharedSecret).Should(Equal("metronsecret"))
			Ω(props.Loggregator.Etcd.Machines).Should(Equal([]string{"1.0.0.7", "1.0.0.8"}))
		})
	})
})
