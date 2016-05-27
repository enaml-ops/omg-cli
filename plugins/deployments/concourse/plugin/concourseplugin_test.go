package concourseplugin_test

import (
	"io/ioutil"

	. "github.com/enaml-ops/omg-cli/plugins/deployments/concourse/plugin"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("given ConcoursePlugin Plugin", func() {
	Context("when plugin is properly initialized", func() {
		var myplugin *ConcoursePlugin
		BeforeEach(func() {
			myplugin = new(ConcoursePlugin)
		})
		Context("when GetProduct is called with valid args", func() {
			var myconcourse []byte
			BeforeEach(func() {
				cloudBytes, _ := ioutil.ReadFile("../fixtures/cloudconfig.yml")
				myconcourse = myplugin.GetProduct([]string{
					"test",
					"--output-filename", "concourse-cloudconfig.yml",
					"--bosh-cloud-config", "true",
					"--bosh-director-uuid", "49525dcc-ecbb-47c2-ad4a-bfddbea27cc7",
					"--network-name", "private",
					"--url", "http://concourse.caleb-washburn.com",
					"--username", "concourse",
					"--password", "concourse",
					"--web-instances", "2",
					"--web-azs", "z1",
					"--worker-azs", "z1",
					"--database-azs", "z1",
					"--bosh-stemcell-alias", "trusty",
					"--postgresql-db-pwd", "secret",
					"--web-vm-type", "small",
					"--worker-vm-type", "medium",
					"--database-vm-type", "medium",
					"--database-storage-type", "large",
				}, cloudBytes)
			})
			It("then it should return the bytes representation of the object", func() {
				Ω(myconcourse).ShouldNot(BeEmpty())
			})
		})
	})
})
