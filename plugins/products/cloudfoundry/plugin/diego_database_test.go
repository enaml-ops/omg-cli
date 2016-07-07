package cloudfoundry_test

import (
	"github.com/enaml-ops/enaml"
	"github.com/enaml-ops/omg-cli/plugins/products/cloudfoundry/enaml-gen/bbs"
	. "github.com/enaml-ops/omg-cli/plugins/products/cloudfoundry/plugin"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("given a Diego Database Partition", func() {
	Describe("given valid flags", func() {

		var instanceGroup *enaml.InstanceGroup
		var grouper InstanceGrouper

		Context("when ToInstanceGroup is called", func() {

			BeforeEach(func() {
				cf := new(Plugin)
				c := cf.GetContext([]string{
					"cloudfoundry",
					"--az", "eastprod-1",
					"--system-domain", "service.cf.domain.com",
					"--stemcell-name", "cool-ubuntu-animal",
					"--network", "foundry-net",
					"--diego-db-ip", "10.0.0.39",
					"--diego-db-ip", "10.0.0.40",
					"--diego-db-vm-type", "dbvmtype",
					"--diego-db-disk-type", "dbdisktype",
					"--diego-db-passphrase", "random-db-encrytionkey",
					"--bbs-ca-cert", "cacert",
					"--etcd-server-cert", "blah-cert",
					"--etcd-server-key", "blah-key",
					"--etcd-client-cert", "bleh-cert",
					"--etcd-client-key", "bleh-key",
					"--etcd-peer-cert", "blee-cert",
					"--etcd-peer-key", "blee-key",
					"--bbs-client-cert", "clientcert",
					"--bbs-client-key", "clientkey",
					"--bbs-server-cert", "clientcert",
					"--bbs-server-key", "clientkey",
					"--bbs-api", "bbs.service.cf.internal:8889",
					"--consul-ip", "1.0.0.1",
					"--consul-ip", "1.0.0.2",
					"--consul-vm-type", "blah",
					"--consul-encryption-key", "encyption-key",
					"--consul-ca-cert", "ca-cert",
					"--consul-agent-cert", "agent-cert",
					"--consul-agent-key", "agent-key",
					"--consul-server-cert", "server-cert",
					"--consul-server-key", "server-key",
					"--metron-secret", "metronsecret",
					"--metron-zone", "metronzoneguid",
					"--syslog-address", "syslog-server",
					"--syslog-port", "10601",
					"--syslog-transport", "tcp",
					"--etcd-machine-ip", "1.0.0.7",
					"--etcd-machine-ip", "1.0.0.8",
				})
				grouper = NewDiegoDatabasePartition(c)
				instanceGroup = grouper.ToInstanceGroup()
			})

			It("then it should be populated with valid network configs", func() {
				ignet := instanceGroup.GetNetworkByName("foundry-net")
				Ω(ignet).ShouldNot(BeNil())
				Ω(ignet.StaticIPs).Should(ConsistOf("10.0.0.39", "10.0.0.40"))
			})

			It("then it should have an instance count in line with given IPs", func() {
				ignet := instanceGroup.GetNetworkByName("foundry-net")
				Ω(len(ignet.StaticIPs)).Should(Equal(instanceGroup.Instances))
			})

			It("then it should be populated the required jobs", func() {
				Ω(instanceGroup.GetJobByName("etcd")).ShouldNot(BeNil())
				Ω(instanceGroup.GetJobByName("bbs")).ShouldNot(BeNil())
				Ω(instanceGroup.GetJobByName("consul_agent")).ShouldNot(BeNil())
				Ω(instanceGroup.GetJobByName("metron_agent")).ShouldNot(BeNil())
				Ω(instanceGroup.GetJobByName("statsd-injector")).ShouldNot(BeNil())
			})
			Describe("given a consul_agent job", func() {
				Context("when defined", func() {
					var job *enaml.InstanceJob
					BeforeEach(func() {
						job = instanceGroup.GetJobByName("consul_agent")
					})
					It("then it should use the correct release", func() {
						Ω(job.Release).Should(Equal(CFReleaseName))
					})

					It("then it should populate my properties", func() {
						Ω(job.Properties).ShouldNot(BeNil())
					})
				})
			})

			Describe("given a statsd-injector job", func() {
				Context("when defined", func() {
					var job *enaml.InstanceJob
					BeforeEach(func() {
						job = instanceGroup.GetJobByName("statsd-injector")
					})
					It("then it should use the correct release", func() {
						Ω(job.Release).Should(Equal(CFReleaseName))
					})

					It("then it should populate my properties", func() {
						Ω(job.Properties).ShouldNot(BeNil())
					})
				})
			})

			Describe("given a metron_agent job", func() {
				Context("when defined", func() {
					var job *enaml.InstanceJob
					BeforeEach(func() {
						job = instanceGroup.GetJobByName("metron_agent")
					})
					It("then it should use the correct release", func() {
						Ω(job.Release).Should(Equal(CFReleaseName))
					})
					It("then it should populate my properties", func() {
						Ω(job.Properties).ShouldNot(BeNil())
					})
				})
			})

			Describe("given a etcd job", func() {
				Context("when defined", func() {
					var job *enaml.InstanceJob
					BeforeEach(func() {
						job = instanceGroup.GetJobByName("etcd")
					})
					It("then it should use the correct release", func() {
						Ω(job.Release).Should(Equal(EtcdReleaseName))
					})
					It("then it should populate my properties", func() {
						Ω(job.Properties).ShouldNot(BeNil())
					})
				})
			})

			Describe("given a bbs job", func() {
				Context("when defined", func() {
					var job *enaml.InstanceJob
					BeforeEach(func() {
						job = instanceGroup.GetJobByName("bbs")
					})
					It("then it should use the correct release", func() {
						Ω(job.Release).Should(Equal(DiegoReleaseName))
					})

					It("should properly set my server key/cert", func() {
						propertiesCasted := job.Properties.(*bbs.Diego)
						Ω(propertiesCasted.Bbs.ServerCert).Should(Equal("clientcert"))
						Ω(propertiesCasted.Bbs.ServerKey).Should(Equal("clientkey"))
					})

					It("should properly set my db passphrase", func() {
						propertiesCasted := job.Properties.(*bbs.Diego)
						Ω(propertiesCasted.Bbs.EncryptionKeys.(map[string]string)["passphrase"]).Should(Equal("random-db-encrytionkey"))
					})

					It("should properly set my bbs.etcd", func() {
						propertiesCasted := job.Properties.(*bbs.Diego)
						Ω(propertiesCasted.Bbs.Etcd).ShouldNot(BeNil())
					})

					It("then it should populate my properties", func() {
						Ω(job.Properties).ShouldNot(BeNil())
					})
				})
			})
		})
	})
})
