package vsphereccplugin_test

import (
	. "github.com/enaml-ops/omg-cli/plugins/cloudconfigs/vsphere/plugin"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("given VSphereCloudConfig Plugin", func() {
	Context("when plugin is properly initialized", func() {
		var myplugin *VSphereCloudConfig
		BeforeEach(func() {
			myplugin = new(VSphereCloudConfig)
		})
		Context("when GetCloudConfig is called with valid args", func() {
			var mycloud []byte
			BeforeEach(func() {
				mycloud = myplugin.GetCloudConfig([]string{
					"test",
					"--az-subnet-map", "us-east-1c:subnet-12345",
					"--region", "us-east-1",
					"--security-group", "bosh",
				})
			})
			It("then it should return the bytes representation of the object", func() {
				Ω(mycloud).ShouldNot(BeEmpty())
			})
		})
	})
})
