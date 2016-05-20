package awsccplugin_test

import (
	. "github.com/bosh-ops/bosh-install/cloudconfigs/aws/plugin"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("given AWSCloudConfig Plugin", func() {
	Context("when plugin is properly initialized", func() {
		var myplugin *AWSCloudConfig
		BeforeEach(func() {
			myplugin = new(AWSCloudConfig)
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
