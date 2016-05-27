package concourse_test

import (
	. "github.com/enaml-ops/omg-cli/plugins/deployments/concourse"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Concourse Deployment", func() {
	var deployment Deployment
	BeforeEach(func() {
		deployment = NewDeployment()
	})
	Describe("Given a CreateCompilation", func() {
		Context("when calling with a networkName", func() {
			It("then we should return a valid enaml.Compilation", func() {
				networkName := "private"
				compilation := deployment.CreateCompilation(networkName)
				Ω(compilation.Network).Should(Equal(networkName))
				Ω(compilation.Workers).Should(Equal(3))
				Ω(compilation.ReuseCompilationVMs).Should(Equal(false))
			})
		})
	})
	Describe("Given a CreateResourcePool", func() {
		Context("when calling with a networkName", func() {
			It("then we should return a valid enaml.ResourcePool", func() {
				networkName := "private"
				resourcePool := deployment.CreateResourcePool(networkName)
				Ω(resourcePool.Network).Should(Equal(networkName))
				Ω(resourcePool.Name).Should(Equal("concourse"))
				stemcell := resourcePool.Stemcell
				Ω(stemcell.Name).Should(Equal("bosh-warden-boshlite-ubuntu-trusty-go_agent"))
				Ω(stemcell.Version).Should(Equal("latest"))
				Ω(stemcell.OS).Should(Equal(""))
				Ω(stemcell.SHA1).Should(Equal(""))
				Ω(stemcell.URL).Should(Equal(""))
				Ω(stemcell.Alias).Should(Equal(""))
			})
		})
	})
	Describe("Given a CreateUpdate", func() {
		Context("when calling", func() {
			It("then we should return a valid enaml.Update", func() {
				update := deployment.CreateUpdate()
				Ω(update.Canaries).Should(Equal(1))
				Ω(update.MaxInFlight).Should(Equal(3))
				Ω(update.Serial).Should(Equal(false))
				Ω(update.UpdateWatchTime).Should(Equal("1000-60000"))
				Ω(update.CanaryWatchTime).Should(Equal("1000-60000"))
			})
		})
	})

	Describe("Given a CreateManualDeploymentNetwork", func() {
		Context("when calling with network name", func() {
			It("then we should return a valid enaml.ManualNetwork", func() {
				networkName := "private"
				networkRange := "10.0.0.0/24"
				networkGateway := "10.244.8.1"
				webIPs := []string{"10.244.8.2"}
				manualNetwork := deployment.CreateManualDeploymentNetwork(networkName, networkRange, networkGateway, webIPs)
				Ω(manualNetwork.Name).Should(Equal(networkName))
				Ω(manualNetwork.Type).Should(Equal("manual"))
				subnets := manualNetwork.Subnets
				Ω(len(subnets)).Should(Equal(1))
				Ω(subnets[0].Range).Should(Equal(networkRange))
				Ω(subnets[0].Gateway).Should(Equal(networkGateway))
				Ω(subnets[0].Static).Should(Equal(webIPs))
				Ω(subnets[0].DNS).Should(BeNil())
				Ω(subnets[0].Reserved).Should(BeEmpty())
				Ω(subnets[0].AZs).Should(BeEmpty())
				Ω(subnets[0].AZ).Should(Equal(""))
			})
		})
	})
	Describe("Given a CreateWebInstanceGroup", func() {
		Context("when calling with resourcePoolName on deployment", func() {
			It("then we should return a valid *enaml.InstanceGroup", func() {
				deployment.ResourcePoolName = "concourse"
				deployment.WebInstances = 1
				deployment.WebVMType = "small"
				web, err := deployment.CreateWebInstanceGroup()
				Ω(err).Should(BeNil())
				Ω(web.Name).Should(Equal("web"))
				Ω(web.Instances).Should(Equal(1))
				Ω(web.ResourcePool).Should(Equal("concourse"))
				Ω(web.AZs).Should(BeEmpty())
				Ω(web.PersistentDisk).Should(Equal(0))
				Ω(web.Stemcell).Should(Equal(""))
				Ω(web.VMType).Should(Equal("small"))
				Ω(len(web.Networks)).Should(Equal(1))
				Ω(len(web.Jobs)).Should(Equal(2))
			})
		})
		Context("when calling with WebAzs and StemcellAlias on deployment", func() {
			It("then we should return a valid *enaml.InstanceGroup", func() {
				deployment.WebInstances = 1
				deployment.WebAZs = []string{"z1"}
				deployment.StemcellAlias = "trusty"
				deployment.WebVMType = "small"
				web, err := deployment.CreateWebInstanceGroup()
				Ω(err).Should(BeNil())
				Ω(web.Name).Should(Equal("web"))
				Ω(web.Instances).Should(Equal(1))
				Ω(web.ResourcePool).Should(Equal(""))
				Ω(web.AZs).Should(Equal([]string{"z1"}))
				Ω(web.PersistentDisk).Should(Equal(0))
				Ω(web.Stemcell).Should(Equal("trusty"))
				Ω(web.VMType).Should(Equal("small"))
				Ω(len(web.Networks)).Should(Equal(1))
				Ω(len(web.Jobs)).Should(Equal(2))
			})
		})
		Context("when calling with ResourcePoolName, WebAzs and StemcellAlias on deployment", func() {
			It("then we should return an error", func() {
				deployment.ResourcePoolName = "concourse"
				deployment.WebAZs = []string{"z1"}
				deployment.StemcellAlias = "trusty"
				_, err := deployment.CreateWebInstanceGroup()
				Ω(err).ShouldNot(BeNil())
			})
		})
		Context("when calling with blank ResourcePoolName, WebAzs and StemcellAlias on deployment", func() {
			It("then we should return an error", func() {
				_, err := deployment.CreateWebInstanceGroup()
				Ω(err).ShouldNot(BeNil())
			})
		})
	})
	Describe("Given a CreateDatabaseInstanceGroup", func() {
		Context("when calling with resourcePoolName on deployment", func() {
			It("then we should return a valid *enaml.InstanceGroup", func() {
				deployment.ResourcePoolName = "concourse"
				deployment.DatabaseVMType = "medium"
				database, err := deployment.CreateDatabaseInstanceGroup()
				Ω(err).Should(BeNil())
				Ω(database.Name).Should(Equal("db"))
				Ω(database.Instances).Should(Equal(1))
				Ω(database.ResourcePool).Should(Equal("concourse"))
				Ω(database.AZs).Should(BeEmpty())
				Ω(database.PersistentDisk).Should(Equal(10240))
				Ω(database.Stemcell).Should(Equal(""))
				Ω(database.VMType).Should(Equal("medium"))
				Ω(len(database.Networks)).Should(Equal(1))
				Ω(len(database.Jobs)).Should(Equal(1))
			})
		})
		Context("when calling with Azs and Stemcell on deployment", func() {
			It("then we should return a valid *enaml.InstanceGroup", func() {
				deployment.DatabaseAZs = []string{"z1"}
				deployment.StemcellAlias = "trusty"
				deployment.DatabaseVMType = "medium"
				deployment.DatabaseStorageType = "medium"
				database, err := deployment.CreateDatabaseInstanceGroup()
				Ω(err).Should(BeNil())
				Ω(database.Name).Should(Equal("db"))
				Ω(database.Instances).Should(Equal(1))
				Ω(database.ResourcePool).Should(Equal(""))
				Ω(database.AZs).Should(Equal([]string{"z1"}))
				Ω(database.PersistentDisk).Should(Equal(0))
				Ω(database.PersistentDiskType).Should(Equal("medium"))
				Ω(database.Stemcell).Should(Equal("trusty"))
				Ω(database.VMType).Should(Equal("medium"))
				Ω(len(database.Networks)).Should(Equal(1))
				Ω(len(database.Jobs)).Should(Equal(1))
			})
		})
		Context("when calling with Azs, Stemcell, ResourcePool on deployment", func() {
			It("then we should return an error", func() {
				deployment.DatabaseAZs = []string{"z1"}
				deployment.StemcellAlias = "trusty"
				deployment.ResourcePoolName = "concourse"
				_, err := deployment.CreateDatabaseInstanceGroup()
				Ω(err).ShouldNot(BeNil())
			})
		})
		Context("when calling with blank Azs, Stemcell, ResourcePool on deployment", func() {
			It("then we should return an error", func() {
				_, err := deployment.CreateDatabaseInstanceGroup()
				Ω(err).ShouldNot(BeNil())
			})
		})
	})
	Describe("Given a CreateWorkerInstanceGroup", func() {
		Context("when calling with resourcePoolName on deployment", func() {
			It("then we should return a valid *enaml.InstanceGroup", func() {
				deployment.ResourcePoolName = "concourse"
				deployment.WorkerVMType = "medium"
				worker, err := deployment.CreateWorkerInstanceGroup()
				Ω(err).Should(BeNil())
				Ω(worker.Name).Should(Equal("worker"))
				Ω(worker.Instances).Should(Equal(1))
				Ω(worker.ResourcePool).Should(Equal("concourse"))
				Ω(worker.AZs).Should(BeEmpty())
				Ω(worker.PersistentDisk).Should(Equal(0))
				Ω(worker.Stemcell).Should(Equal(""))
				Ω(worker.VMType).Should(Equal("medium"))
				Ω(len(worker.Networks)).Should(Equal(1))
				Ω(len(worker.Jobs)).Should(Equal(3))
			})
		})
		Context("when calling with Azs and Stemcell on deployment", func() {
			It("then we should return a valid *enaml.InstanceGroup", func() {
				deployment.WorkerAZs = []string{"z1"}
				deployment.StemcellAlias = "trusty"
				deployment.WorkerVMType = "medium"
				worker, err := deployment.CreateWorkerInstanceGroup()
				Ω(err).Should(BeNil())
				Ω(worker.Name).Should(Equal("worker"))
				Ω(worker.Instances).Should(Equal(1))
				Ω(worker.ResourcePool).Should(Equal(""))
				Ω(worker.AZs).Should(Equal([]string{"z1"}))
				Ω(worker.PersistentDisk).Should(Equal(0))
				Ω(worker.Stemcell).Should(Equal("trusty"))
				Ω(worker.VMType).Should(Equal("medium"))
				Ω(len(worker.Networks)).Should(Equal(1))
				Ω(len(worker.Jobs)).Should(Equal(3))
			})
		})
		Context("when calling with ResourcePoolName, Azs and Stemcell on deployment", func() {
			It("then we should return an error", func() {
				deployment.ResourcePoolName = "concourse"
				deployment.WorkerAZs = []string{"z1"}
				deployment.StemcellAlias = "trusty"
				_, err := deployment.CreateWorkerInstanceGroup()
				Ω(err).ShouldNot(BeNil())
			})
		})
		Context("when calling with blank ResourcePoolName, Azs and Stemcell on deployment", func() {
			It("then we should return an error", func() {
				_, err := deployment.CreateWorkerInstanceGroup()
				Ω(err).ShouldNot(BeNil())
			})
		})
	})
	Describe("Given a new deployment", func() {
		XContext("when calling Initialize without a strong password", func() {

			deployment.ConcoursePassword = "test"
			It("then we should error and prompt the user for a better pass", func() {
				err := deployment.Initialize([]byte(""))
				Ω(err).ShouldNot(BeNil())
			})
		})
	})
})
