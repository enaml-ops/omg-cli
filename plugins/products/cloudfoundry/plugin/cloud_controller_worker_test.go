package cloudfoundry_test

import (
	"fmt"

	. "github.com/enaml-ops/omg-cli/plugins/products/cloudfoundry/plugin"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gopkg.in/yaml.v2"
)

var _ = Describe("Cloud Controller Worker Partition", func() {
	Context("When initialized with a complete set of arguments", func() {
		var cloudControllerWorker InstanceGrouper
		BeforeEach(func() {
			plugin := new(Plugin)
			c := plugin.GetContext([]string{
				"cloudfoundry",
				"--consul-ip", "1.0.0.1",
				"--consul-ip", "1.0.0.2",
				"--consul-ca-cert", "consul-ca-cert",
				"--consul-agent-cert", "consul-agent-cert",
				"--consul-server-cert", "consulservercert",
				"--consul-server-key", "consulserverkey",
				"--cc-worker-vm-type", "ccworkervmtype",
				"--cc-worker-network", "foundry",
				"--cc-staging-upload-user", "staginguser",
				"--cc-staging-upload-password", "stagingpassword",
				"--cc-bulk-api-user", "bulkapiuser",
				"--cc-bulk-api-password", "bulkapipassword",
				"--cc-db-encryption-key", "dbencryptionkey",
				"--cc-internal-api-user", "internalapiuser",
				"--cc-internal-api-password", "internalapipassword",
				"--system-domain", "sys.yourdomain.com",
				"--app-domain", "apps.yourdomain.com",
				"--allow-app-ssh-access", "true",
				"--nfs-server-address", "10.0.0.19",
				"--nfs-share-path", "/var/vcap/nfs",
				"--metron-secret", "metronsecret",
				"--metron-zone", "metronzoneguid",
			})

			cloudControllerWorker = NewCloudControllerWorkerPartition(c)
		})

		It("then should not be nil", func() {
			Ω(cloudControllerWorker).ShouldNot(BeNil())
		})
		It("should have the name of the Network correctly set", func() {
			igf := cloudControllerWorker.ToInstanceGroup()

			networks := igf.Networks
			Ω(len(networks)).Should(Equal(1))
			Ω(networks[0].Name).Should(Equal("foundry"))
		})
		It("should have 5 jobs under it", func() {
			igf := cloudControllerWorker.ToInstanceGroup()
			jobs := igf.Jobs
			Ω(len(jobs)).Should(Equal(5))
		})
		It("should have NFS Mounter set as a job", func() {
			igf := cloudControllerWorker.ToInstanceGroup()
			nfsMounter := igf.Jobs[2]
			Ω(nfsMounter.Name).Should(Equal("nfs_mounter"))
		})
		XIt("should have NFS Mounter details set properly", func() {
			igf := cloudControllerWorker.ToInstanceGroup()

			b, _ := yaml.Marshal(igf)
			fmt.Print(string(b))
		})

	})
})
