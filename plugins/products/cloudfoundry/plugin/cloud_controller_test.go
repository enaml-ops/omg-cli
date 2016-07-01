package cloudfoundry_test

import (
	//"fmt"

	. "github.com/enaml-ops/omg-cli/plugins/products/cloudfoundry/plugin"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gopkg.in/yaml.v2"

	ccnglib "github.com/enaml-ops/omg-cli/plugins/products/cloudfoundry/enaml-gen/cloud_controller_ng"
)

var _ = Describe("Cloud Controller Partition", func() {
	Context("When initialized with a complete set of arguments", func() {
		var cloudController InstanceGrouper

		BeforeEach(func() {
			plugin := new(Plugin)
			c := plugin.GetContext([]string{
				"cloudfoundry",
				"--consul-ip", "1.0.0.1",
				"--consul-ip", "1.0.0.2",
				"--az", "az",
				"--stemcell-name", "stemcell",
				"--consul-encryption-key", "consulencryptionkey",
				"--consul-ca-cert", "consul-ca-cert",
				"--consul-agent-cert", "consul-agent-cert",
				"--consul-agent-key", "consul-agent-key",
				"--consul-server-cert", "consulservercert",
				"--consul-server-key", "consulserverkey",
				"--cc-vm-type", "ccvmtype",
				"--network", "foundry",
				"--cc-staging-upload-user", "staginguser",
				"--cc-staging-upload-password", "stagingpassword",
				"--cc-bulk-api-user", "bulkapiuser",
				"--cc-bulk-api-password", "bulkapipassword",
				"--cc-db-encryption-key", "dbencryptionkey",
				"--cc-internal-api-user", "internalapiuser",
				"--cc-internal-api-password", "internalapipassword",
				"--system-domain", "sys.yourdomain.com",
				"--app-domain", "apps.yourdomain.com",
				"--allow-app-ssh-access",
				"--nfs-server-address", "10.0.0.19",
				"--nfs-share-path", "/var/vcap/nfs",
				"--metron-secret", "metronsecret",
				"--metron-zone", "metronzoneguid",
				"--host-key-fingerprint", "hostkeyfingerprint",
				"--support-address", "http://support.pivotal.io",
				"--min-cli-version", "6.7.0",
			})

			cloudController = NewCloudControllerPartition(c)
		})
		It("then should not be nil", func() {
			Ω(cloudController).ShouldNot(BeNil())
		})
		It("should have valid values", func() {
			Ω(cloudController.HasValidValues()).Should(BeTrue())
		})
		It("should have the name of the Network correctly set", func() {
			igf := cloudController.ToInstanceGroup()

			networks := igf.Networks
			Ω(len(networks)).Should(Equal(1))
			Ω(networks[0].Name).Should(Equal("foundry"))
		})
		It("should have 5 jobs under it", func() {
			igf := cloudController.ToInstanceGroup()
			jobs := igf.Jobs
			Ω(len(jobs)).Should(Equal(5))
		})
		It("should have NFS Mounter set as a job", func() {
			igf := cloudController.ToInstanceGroup()
			nfsMounter := igf.Jobs[2]
			Ω(nfsMounter.Name).Should(Equal("nfs_mounter"))
		})
		It("should have NFS Mounter details set properly", func() {
			igf := cloudController.ToInstanceGroup()

			b, _ := yaml.Marshal(igf)
			//fmt.Print(string(b))

			Ω(string(b)).Should(ContainSubstring("https://login.sys.yourdomain.com"))
		})
		XIt("should account for QuotaDefinitions structure", func() {
			igf := cloudController.ToInstanceGroup()
			Ω(igf.Jobs[0].Name).Should(Equal("cloud_controller_worker"))
			ccNg, typecasted := igf.Jobs[0].Properties.(*ccnglib.CloudControllerNg)
			Ω(typecasted).Should(BeTrue())

			_, quotaTypeCasted := ccNg.Cc.QuotaDefinitions.([]string)
			Ω(quotaTypeCasted).Should(BeFalse())
		})
		XIt("should account for InstallBuildPacks structure", func() {
			igf := cloudController.ToInstanceGroup()
			Ω(igf.Jobs[0].Name).Should(Equal("cloud_controller_worker"))
			ccNg, typecasted := igf.Jobs[0].Properties.(*ccnglib.CloudControllerNg)
			Ω(typecasted).Should(BeTrue())

			_, bpTypecasted := ccNg.Cc.InstallBuildpacks.([]string)
			Ω(bpTypecasted).Should(BeFalse())
		})
		XIt("should account for SecurityGroupDefinitions structure", func() {
			igf := cloudController.ToInstanceGroup()
			Ω(igf.Jobs[0].Name).Should(Equal("cloud_controller_worker"))
			ccNg, typecasted := igf.Jobs[0].Properties.(*ccnglib.CloudControllerNg)
			Ω(typecasted).Should(BeTrue())

			_, securityTypcasted := ccNg.Cc.SecurityGroupDefinitions.([]string)
			Ω(securityTypcasted).Should(BeFalse())
		})

	})
})
