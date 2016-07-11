package cloudfoundry_test

import (
	"github.com/enaml-ops/enaml"
	ltc "github.com/enaml-ops/omg-cli/plugins/products/cloudfoundry/enaml-gen/loggregator_trafficcontroller"
	"github.com/enaml-ops/omg-cli/plugins/products/cloudfoundry/enaml-gen/metron_agent"
	"github.com/enaml-ops/omg-cli/plugins/products/cloudfoundry/enaml-gen/route_registrar"
	. "github.com/enaml-ops/omg-cli/plugins/products/cloudfoundry/plugin"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("given the loggregator traffic controller partition", func() {

	Context("when initialized WITHOUT a complete set of arguments", func() {
		var grouper InstanceGrouper
		BeforeEach(func() {
			cf := new(Plugin)
			c := cf.GetContext([]string{
				"cloudfoundry",
			})
			grouper = NewLoggregatorTrafficController(c)
		})

		It("should not have valid values", func() {
			Ω(grouper.HasValidValues()).Should(BeFalse())
		})

		It("should have the loggregator traffic controller job", func() {
			group := grouper.ToInstanceGroup()
			job := group.GetJobByName("loggregator_trafficcontroller")
			Ω(job).ShouldNot(BeNil())
		})
	})

	Context("when initialized with a complete set of arguments", func() {
		var grouper InstanceGrouper
		var dm *enaml.DeploymentManifest
		BeforeEach(func() {
			cf := new(Plugin)
			c := cf.GetContext([]string{
				"cloudfoundry",
				"--az", "eastprod-1",
				"--stemcell-name", "cool-ubuntu-animal",
				"--network", "foundry-net",
				"--system-domain", "sys.yourdomain.com",
				"--skip-cert-verify=false",
				"--loggregator-traffic-controller-ip", "10.0.0.39",
				"--loggregator-traffic-controller-ip", "10.0.0.40",
				"--loggregator-traffic-controller-vmtype", "vmtype",
				"--etcd-machine-ip", "10.0.1.2",
				"--etcd-machine-ip", "10.0.1.3",
				"--etcd-machine-ip", "10.0.1.4",
				"--doppler-client-secret", "dopplersecret",
				"--metron-secret", "metronsecret",
				"--metron-zone", "metronzoneguid",
				"--syslog-address", "syslog-server",
				"--syslog-port", "10601",
				"--syslog-transport", "tcp",
				"--nats-user", "nats",
				"--nats-pass", "pass",
				"--nats-port", "4222",
				"--nats-machine-ip", "1.0.0.5",
				"--nats-machine-ip", "1.0.0.6",
			})
			grouper = NewLoggregatorTrafficController(c)
			dm = new(enaml.DeploymentManifest)
			dm.AddInstanceGroup(grouper.ToInstanceGroup())
		})

		It("then it should allow the user to configure the AZs", func() {
			ig := dm.GetInstanceGroupByName("loggregator_trafficcontroller-partition")
			Ω(len(ig.AZs)).Should(Equal(1))
			Ω(ig.AZs[0]).Should(Equal("eastprod-1"))
		})

		It("then it should allow the user to configure the used stemcell", func() {
			ig := dm.GetInstanceGroupByName("loggregator_trafficcontroller-partition")
			Ω(ig.Stemcell).Should(Equal("cool-ubuntu-animal"))
		})

		It("then it should allow the user to configure the network to use", func() {
			ig := dm.GetInstanceGroupByName("loggregator_trafficcontroller-partition")
			network := ig.GetNetworkByName("foundry-net")
			Ω(network).ShouldNot(BeNil())
			Ω(len(network.StaticIPs)).Should(Equal(2))
			Ω(network.StaticIPs[0]).Should(Equal("10.0.0.39"))
			Ω(network.StaticIPs[1]).Should(Equal("10.0.0.40"))
		})

		It("should use the correct number of instances based on the network IPs", func() {
			ig := dm.GetInstanceGroupByName("loggregator_trafficcontroller-partition")
			network := ig.Networks[0]
			Ω(ig.Instances).Should(Equal(len(network.StaticIPs)))
		})

		It("then it should allow the user to configure the VM type", func() {
			ig := dm.GetInstanceGroupByName("loggregator_trafficcontroller-partition")
			Ω(ig.VMType).Should(Equal("vmtype"))
		})

		It("then it should have update max in flight 1", func() {
			ig := dm.GetInstanceGroupByName("loggregator_trafficcontroller-partition")
			Ω(ig.Update.MaxInFlight).Should(Equal(1))
		})

		It("then it should have correctly configured the loggregator traffic controller job", func() {
			ig := dm.GetInstanceGroupByName("loggregator_trafficcontroller-partition")
			job := ig.GetJobByName("loggregator_trafficcontroller")
			Ω(job).ShouldNot(BeNil())
			Ω(job.Release).Should(Equal(CFReleaseName))

			props := job.Properties.(*ltc.LoggregatorTrafficcontroller)
			Ω(props.SystemDomain).Should(Equal("sys.yourdomain.com"))
			Ω(props.Cc.SrvApiUri).Should(Equal("https://api.sys.yourdomain.com"))
			Ω(props.Ssl.SkipCertVerify).Should(BeFalse())
			Ω(props.TrafficController.Zone).Should(Equal("metronzoneguid"))
			Ω(props.Doppler.UaaClientId).Should(Equal("doppler"))
			Ω(props.Uaa.Clients.Doppler.Secret).Should(Equal("dopplersecret"))
			Ω(props.Loggregator.Etcd.Machines).Should(ConsistOf("10.0.1.2", "10.0.1.3", "10.0.1.4"))
		})

		It("then it should have the metron_agent job", func() {
			ig := dm.GetInstanceGroupByName("loggregator_trafficcontroller-partition")
			job := ig.GetJobByName("metron_agent")
			Ω(job).ShouldNot(BeNil())
			Ω(job.Release).Should(Equal(CFReleaseName))

			props := job.Properties.(*metron_agent.MetronAgent)
			Ω(props.MetronAgent.Zone).Should(Equal("metronzoneguid"))
			Ω(props.MetronAgent.Deployment).Should(Equal(CFReleaseName))
			Ω(props.MetronEndpoint.SharedSecret).Should(Equal("metronsecret"))
			Ω(props.Loggregator.Etcd.Machines).Should(ConsistOf("10.0.1.2", "10.0.1.3", "10.0.1.4"))
		})

		It("then it should have the route_registrar job", func() {
			ig := dm.GetInstanceGroupByName("loggregator_trafficcontroller-partition")
			job := ig.GetJobByName("route_registrar")
			Ω(job).ShouldNot(BeNil())

			props, _ := job.Properties.(*route_registrar.RouteRegistrar)
			Ω(props.Nats).ShouldNot(BeNil())
			Ω(props.Nats.User).Should(Equal("nats"))
			Ω(props.Nats.Password).Should(Equal("pass"))
			Ω(props.Nats.Port).Should(Equal(4222))
			Ω(props.Nats.Machines).Should(ConsistOf("1.0.0.5", "1.0.0.6"))

			Ω(props.Routes).ShouldNot(BeNil())
			routes := props.Routes.([]map[string]interface{})
			Ω(len(routes)).Should(Equal(2))
			Ω(routes[0]["name"]).Should(Equal("doppler"))
			Ω(routes[0]["port"]).Should(Equal(8081))
			Ω(routes[0]["registration_interval"]).Should(Equal("20s"))
			Ω(routes[0]["uris"]).Should(ConsistOf("doppler.sys.yourdomain.com"))
			Ω(routes[1]["name"]).Should(Equal("loggregator"))
			Ω(routes[1]["port"]).Should(Equal(8080))
			Ω(routes[1]["registration_interval"]).Should(Equal("20s"))
			Ω(routes[1]["uris"]).Should(ConsistOf("loggregator.sys.yourdomain.com"))
		})

		It("then it should have the statsd-injector job", func() {
			ig := dm.GetInstanceGroupByName("loggregator_trafficcontroller-partition")
			job := ig.GetJobByName("statsd-injector")
			Ω(job).ShouldNot(BeNil())
			Ω(job.Release).Should(Equal(CFReleaseName))
		})
	})
})
