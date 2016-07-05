package cloudfoundry_test

import (
	"github.com/enaml-ops/enaml"
	"github.com/enaml-ops/omg-cli/plugins/products/cloudfoundry/enaml-gen/auctioneer"
	"github.com/enaml-ops/omg-cli/plugins/products/cloudfoundry/enaml-gen/cc_uploader"
	"github.com/enaml-ops/omg-cli/plugins/products/cloudfoundry/enaml-gen/converger"
	"github.com/enaml-ops/omg-cli/plugins/products/cloudfoundry/enaml-gen/file_server"
	"github.com/enaml-ops/omg-cli/plugins/products/cloudfoundry/enaml-gen/nsync"
	"github.com/enaml-ops/omg-cli/plugins/products/cloudfoundry/enaml-gen/route_emitter"
	"github.com/enaml-ops/omg-cli/plugins/products/cloudfoundry/enaml-gen/ssh_proxy"
	"github.com/enaml-ops/omg-cli/plugins/products/cloudfoundry/enaml-gen/stager"
	"github.com/enaml-ops/omg-cli/plugins/products/cloudfoundry/enaml-gen/tps"
	. "github.com/enaml-ops/omg-cli/plugins/products/cloudfoundry/plugin"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("given a Diego Brain Partition", func() {
	Context("when initialized WITHOUT a complete set of arguments", func() {
		var ig InstanceGrouper
		BeforeEach(func() {
			cf := new(Plugin)
			c := cf.GetContext([]string{""})
			ig = NewDiegoBrainPartition(c)
		})

		It("then it should contain the appropriate jobs", func() {
			group := ig.ToInstanceGroup()
			Ω(group.GetJobByName("auctioneer")).ShouldNot(BeNil())
			Ω(group.GetJobByName("cc_uploader")).ShouldNot(BeNil())
			Ω(group.GetJobByName("converger")).ShouldNot(BeNil())
			Ω(group.GetJobByName("file_server")).ShouldNot(BeNil())
			Ω(group.GetJobByName("nsync")).ShouldNot(BeNil())
			Ω(group.GetJobByName("route_emitter")).ShouldNot(BeNil())
			Ω(group.GetJobByName("ssh_proxy")).ShouldNot(BeNil())
			Ω(group.GetJobByName("stager")).ShouldNot(BeNil())
			Ω(group.GetJobByName("tps")).ShouldNot(BeNil())
			Ω(group.GetJobByName("consul_agent")).ShouldNot(BeNil())
			Ω(group.GetJobByName("metron_agent")).ShouldNot(BeNil())
			Ω(group.GetJobByName("statsd-injector")).ShouldNot(BeNil())
		})

		It("then it should not validate", func() {
			Ω(ig.HasValidValues()).Should(BeFalse())
		})
	})

	Context("when initialized with a complete set of arguments", func() {
		var deploymentManifest *enaml.DeploymentManifest
		var grouper InstanceGrouper

		BeforeEach(func() {
			cf := new(Plugin)
			c := cf.GetContext([]string{
				"cloudfoundry",
				"--az", "eastprod-1",
				"--stemcell-name", "cool-ubuntu-animal",
				"--network", "foundry-net",
				"--allow-app-ssh-access",
				"--diego-brain-ip", "10.0.0.39",
				"--diego-brain-ip", "10.0.0.40",
				"--diego-brain-vm-type", "brainvmtype",
				"--diego-brain-disk-type", "braindisktype",
				"--bbs-ca-cert", "cacert",
				"--bbs-client-cert", "clientcert",
				"--bbs-client-key", "clientkey",
				"--bbs-api", "bbs.service.cf.internal:8889",
				"--bbs-require-ssl=false",
				"--skip-cert-verify=false",
				"--cc-uploader-poll-interval", "25",
				"--cc-external-port", "9023",
				"--system-domain", "sys.test.com",
				"--cc-internal-api-user", "internaluser",
				"--cc-internal-api-password", "internalpassword",
				"--cc-bulk-batch-size", "5",
				"--cc-fetch-timeout", "30",
				"--fs-listen-addr", "0.0.0.0:12345",
				"--fs-static-dir", "/foo/bar/baz",
				"--fs-debug-addr", "10.0.1.2:22222",
				"--fs-log-level", "debug",
				"--metron-port", "3458",
				"--nats-user", "nats",
				"--nats-port", "1234",
				"--nats-pass", "natspass",
				"--nats-machine-ip", "10.0.0.11",
				"--nats-machine-ip", "10.0.0.12",
				"--ssh-proxy-uaa-secret", "secret",
				"--traffic-controller-url", "wss://doppler.sys.yourdomain.com:443",
			})
			grouper = NewDiegoBrainPartition(c)
			deploymentManifest = new(enaml.DeploymentManifest)
			deploymentManifest.AddInstanceGroup(grouper.ToInstanceGroup())
		})

		It("then it should validate", func() {
			Ω(grouper.HasValidValues()).Should(BeTrue())
		})

		It("then it should allow the user to configure the AZs", func() {
			ig := deploymentManifest.GetInstanceGroupByName("diego_brain-partition")
			Ω(len(ig.AZs)).Should(Equal(1))
			Ω(ig.AZs[0]).Should(Equal("eastprod-1"))
		})

		It("then it should allow the user to configure the used stemcell", func() {
			ig := deploymentManifest.GetInstanceGroupByName("diego_brain-partition")
			Ω(ig.Stemcell).Should(Equal("cool-ubuntu-animal"))
		})

		It("then it should allow the user to configure the network to use", func() {
			ig := deploymentManifest.GetInstanceGroupByName("diego_brain-partition")
			network := ig.GetNetworkByName("foundry-net")
			Ω(network).ShouldNot(BeNil())
			Ω(len(network.StaticIPs)).Should(Equal(2))
			Ω(network.StaticIPs[0]).Should(Equal("10.0.0.39"))
			Ω(network.StaticIPs[1]).Should(Equal("10.0.0.40"))
		})

		It("then it should allow the user to configure the VM type", func() {
			ig := deploymentManifest.GetInstanceGroupByName("diego_brain-partition")
			Ω(ig.VMType).Should(Equal("brainvmtype"))
		})

		It("then it should allow the user to configure the disk type", func() {
			ig := deploymentManifest.GetInstanceGroupByName("diego_brain-partition")
			Ω(ig.PersistentDiskType).Should(Equal("braindisktype"))
		})

		It("then it should configure the correct number of instances automatically from the count of IPs", func() {
			ig := deploymentManifest.GetInstanceGroupByName("diego_brain-partition")
			Ω(len(ig.Networks)).Should(Equal(1))
			Ω(len(ig.Networks[0].StaticIPs)).Should(Equal(ig.Instances))
		})

		It("then it should have update max-in-flight 1", func() {
			ig := deploymentManifest.GetInstanceGroupByName("diego_brain-partition")
			Ω(ig.Update.MaxInFlight).Should(Equal(1))
		})

		It("then it should allow the user to configure the auctioneer BBS", func() {
			ig := deploymentManifest.GetInstanceGroupByName("diego_brain-partition")
			job := ig.GetJobByName("auctioneer")
			a := job.Properties.(*auctioneer.Auctioneer)
			Ω(a.Bbs.CaCert).Should(Equal("cacert"))
			Ω(a.Bbs.ClientCert).Should(Equal("clientcert"))
			Ω(a.Bbs.ClientKey).Should(Equal("clientkey"))
			Ω(a.Bbs.ApiLocation).Should(Equal("bbs.service.cf.internal:8889"))
		})

		It("then it should allow the user to configure the CC uploader", func() {
			ig := deploymentManifest.GetInstanceGroupByName("diego_brain-partition")
			job := ig.GetJobByName("cc_uploader")
			cc := job.Properties.(*cc_uploader.CcUploader)
			Ω(cc.Diego.Ssl.SkipCertVerify).Should(BeFalse())
			Ω(cc.Cc.JobPollingIntervalInSeconds).Should(Equal(25))
		})

		It("then it should allow the user to configure the converger", func() {
			ig := deploymentManifest.GetInstanceGroupByName("diego_brain-partition")
			job := ig.GetJobByName("converger")
			c := job.Properties.(*converger.Converger)
			Ω(c.Bbs.ApiLocation).Should(Equal("bbs.service.cf.internal:8889"))
			Ω(c.Bbs.CaCert).Should(Equal("cacert"))
			Ω(c.Bbs.ClientCert).Should(Equal("clientcert"))
			Ω(c.Bbs.ClientKey).Should(Equal("clientkey"))
		})

		It("then it should allow the user to configure the file server", func() {
			ig := deploymentManifest.GetInstanceGroupByName("diego_brain-partition")
			job := ig.GetJobByName("file_server")
			fs := job.Properties.(*file_server.FileServer)

			Ω(fs.Diego.Ssl.SkipCertVerify).Should(BeFalse())

			Ω(fs.ListenAddr).Should(Equal("0.0.0.0:12345"))
			Ω(fs.StaticDirectory).Should(Equal("/foo/bar/baz"))
			Ω(fs.DebugAddr).Should(Equal("10.0.1.2:22222"))
			Ω(fs.LogLevel).Should(Equal("debug"))
			Ω(fs.DropsondePort).Should(Equal(3458))
		})

		It("then it should allow the user to configure nsync", func() {
			ig := deploymentManifest.GetInstanceGroupByName("diego_brain-partition")
			job := ig.GetJobByName("nsync")
			n := job.Properties.(*nsync.Nsync)
			Ω(n.Diego.Ssl.SkipCertVerify).Should(BeFalse())
			Ω(n.Bbs.ApiLocation).Should(Equal("bbs.service.cf.internal:8889"))
			Ω(n.Bbs.CaCert).Should(Equal("cacert"))
			Ω(n.Bbs.ClientCert).Should(Equal("clientcert"))
			Ω(n.Bbs.ClientKey).Should(Equal("clientkey"))
			Ω(n.Cc.BaseUrl).Should(Equal("https://api.sys.test.com"))
			Ω(n.Cc.BasicAuthUsername).Should(Equal("internaluser"))
			Ω(n.Cc.BasicAuthPassword).Should(Equal("internalpassword"))
			Ω(n.Cc.BulkBatchSize).Should(Equal(5))
			Ω(n.Cc.FetchTimeoutInSeconds).Should(Equal(30))
			Ω(n.Cc.PollingIntervalInSeconds).Should(Equal(25))
		})

		It("then it should allows the user to configure the route emitter", func() {
			ig := deploymentManifest.GetInstanceGroupByName("diego_brain-partition")
			job := ig.GetJobByName("route_emitter")
			r := job.Properties.(*route_emitter.RouteEmitter)
			Ω(r.Bbs.ApiLocation).Should(Equal("bbs.service.cf.internal:8889"))
			Ω(r.Bbs.CaCert).Should(Equal("cacert"))
			Ω(r.Bbs.ClientCert).Should(Equal("clientcert"))
			Ω(r.Bbs.ClientKey).Should(Equal("clientkey"))
			Ω(r.Bbs.RequireSsl).Should(BeFalse())
			Ω(r.Nats.User).Should(Equal("nats"))
			Ω(r.Nats.Password).Should(Equal("natspass"))
			Ω(r.Nats.Port).Should(Equal(1234))
			Ω(r.Nats.Machines).Should(ContainElement("10.0.0.11"))
			Ω(r.Nats.Machines).Should(ContainElement("10.0.0.12"))
		})

		It("then it should allow the user to configure the SSH proxy", func() {
			ig := deploymentManifest.GetInstanceGroupByName("diego_brain-partition")
			job := ig.GetJobByName("ssh_proxy")
			s := job.Properties.(*ssh_proxy.SshProxy)
			Ω(s.Diego.Ssl.SkipCertVerify).Should(BeFalse())
			Ω(s.Bbs.ApiLocation).Should(Equal("bbs.service.cf.internal:8889"))
			Ω(s.Bbs.CaCert).Should(Equal("cacert"))
			Ω(s.Bbs.ClientCert).Should(Equal("clientcert"))
			Ω(s.Bbs.ClientKey).Should(Equal("clientkey"))
			Ω(s.Bbs.RequireSsl).Should(BeFalse())
			Ω(s.EnableCfAuth).Should(BeTrue())    // tied to allow-app-ssh-access
			Ω(s.EnableDiegoAuth).Should(BeTrue()) // tied to allow-app-ssh-access
			Ω(s.Cc.ExternalPort).Should(Equal(9023))
			Ω(s.UaaTokenUrl).Should(Equal("https://uaa.sys.test.com/oauth/token"))
			Ω(s.UaaSecret).Should(Equal("secret"))
		})

		It("then it should allow the user to configure the stager", func() {
			ig := deploymentManifest.GetInstanceGroupByName("diego_brain-partition")
			job := ig.GetJobByName("stager")
			s := job.Properties.(*stager.Stager)
			Ω(s.Diego.Ssl.SkipCertVerify).Should(BeFalse())
			Ω(s.Bbs.ApiLocation).Should(Equal("bbs.service.cf.internal:8889"))
			Ω(s.Bbs.CaCert).Should(Equal("cacert"))
			Ω(s.Bbs.ClientCert).Should(Equal("clientcert"))
			Ω(s.Bbs.ClientKey).Should(Equal("clientkey"))
			Ω(s.Bbs.RequireSsl).Should(BeFalse())
			Ω(s.Cc.ExternalPort).Should(Equal(9023))
			Ω(s.Cc.BasicAuthUsername).Should(Equal("internaluser"))
			Ω(s.Cc.BasicAuthPassword).Should(Equal("internalpassword"))
		})

		It("then it should allow the user to configure the tps", func() {
			ig := deploymentManifest.GetInstanceGroupByName("diego_brain-partition")
			job := ig.GetJobByName("tps")
			t := job.Properties.(*tps.Tps)
			Ω(t.Diego.Ssl.SkipCertVerify).Should(BeFalse())
			Ω(t.TrafficControllerUrl).Should(Equal("wss://doppler.sys.yourdomain.com:443"))
			Ω(t.Bbs.ApiLocation).Should(Equal("bbs.service.cf.internal:8889"))
			Ω(t.Bbs.CaCert).Should(Equal("cacert"))
			Ω(t.Bbs.ClientCert).Should(Equal("clientcert"))
			Ω(t.Bbs.ClientKey).Should(Equal("clientkey"))
			Ω(t.Bbs.RequireSsl).Should(BeFalse())
			Ω(t.Cc.ExternalPort).Should(Equal(9023))
			Ω(t.Cc.BasicAuthUsername).Should(Equal("internaluser"))
			Ω(t.Cc.BasicAuthPassword).Should(Equal("internalpassword"))
		})
	})
})
