package awsccplugin_test

import (
	"fmt"

	"github.com/codegangsta/cli"
	"github.com/enaml-ops/enaml"
	. "github.com/enaml-ops/omg-cli/plugins/cloudconfigs/aws/plugin"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gopkg.in/yaml.v2"
)

var _ = Describe("given AWSCloudConfig Plugin", func() {
	Context("when calling GetFlags", func() {
		var flags []cli.Flag
		BeforeEach(func() {
			cfg := new(AWSCloudConfig)
			flags = cfg.GetFlags()
		})
		It("then it should yield back a bosh-az-name flag", func() {
			testMultiAZFlagExists("bosh-az-name", flags)
		})
		It("then it should yield back a cidr flag", func() {
			testMultiAZFlagExists("cidr", flags)
		})
		It("then it should yield back a gateway flag", func() {
			testMultiAZFlagExists("gateway", flags)
		})
		It("then it should yield back a dns flag", func() {
			testMultiAZFlagExists("dns", flags)
		})
		It("then it should yield back a aws-az-name flag", func() {
			testMultiAZFlagExists("aws-az-name", flags)
		})
		It("then it should yield back a aws-subnet-name flag", func() {
			testMultiAZFlagExists("aws-subnet-name", flags)
		})
		It("then it should yield back a aws-region flag", func() {
			testFlagExists("aws-region", flags)
		})
		It("then it should yield back a bosh bosh-reserve-range flag", func() {
			testMultiAZFlagExists("bosh-reserve-range", flags)
		})

	})

	Context("when plugin is properly initialized", func() {
		var myplugin *AWSCloudConfig
		BeforeEach(func() {
			myplugin = new(AWSCloudConfig)
		})
		Context("when GetCloudConfig is called with valid args for a single az & network", func() {
			var mycloud []byte
			BeforeEach(func() {
				mycloud = myplugin.GetCloudConfig([]string{
					"test",
					"--aws-region", "us-east-1",
					"--aws-security-group", "bosh",
					"--bosh-az-name-1", "bosh-az1",
					"--cidr-1", "10.0.0.0/24",
					"--gateway-1", "10.0.0.1",
					"--dns-1", "10.0.0.240",
					"--aws-az-name-1", "aws-az1-blah",
					"--aws-subnet-name-1", "my-aws-subnet-13857298354792835",
					"--bosh-reserve-range-1", "10.0.0.1-10.0.0.10",
					"--bosh-reserve-range-1", "10.0.0.20-10.0.0.30",
				})
			})
			It("then it should return the bytes representation of the object", func() {
				Ω(mycloud).ShouldNot(BeEmpty())
			})
			It("then should contain the correct number of networks and azs", func() {
				var mynetwork = new(enaml.ManualNetwork)
				fmt.Println(string(mycloud))
				ccManifest := enaml.NewCloudConfigManifest(mycloud)
				testNetwork, _ := yaml.Marshal(ccManifest.Networks[0])
				yaml.Unmarshal(testNetwork, mynetwork)
				subnetCount := len(mynetwork.Subnets)
				azCount := len(ccManifest.AZs)
				Ω(azCount).Should(Equal(1))
				Ω(subnetCount).Should(Equal(1))
			})
		})
		Context("when GetCloudConfig is called with valid args for a multi az & network", func() {
			var mycloud []byte
			BeforeEach(func() {
				mycloud = myplugin.GetCloudConfig([]string{
					"test",
					"--aws-region", "us-east-1",
					"--aws-security-group", "bosh",
					"--bosh-az-name-1", "bosh-az1",
					"--cidr-1", "10.0.0.0/24",
					"--gateway-1", "10.0.0.1",
					"--dns-1", "10.0.0.240",
					"--aws-az-name-1", "aws-az1-blah",
					"--aws-subnet-name-1", "my-aws-subnet-13857298354792835",
					"--bosh-reserve-range-1", "10.0.0.1-10.0.0.10",
					"--bosh-reserve-range-1", "10.0.0.20-10.0.0.30",
					"--bosh-az-name-2", "bosh-az2",
					"--cidr-2", "10.1.0.0/24",
					"--gateway-2", "10.1.0.1",
					"--dns-2", "10.1.0.240",
					"--aws-az-name-2", "aws-az2-blah",
					"--aws-subnet-name-2", "my-aws-subnet2-13857298354792835",
					"--bosh-reserve-range-2", "10.1.0.1-10.1.0.10",
					"--bosh-reserve-range-2", "10.1.0.20-10.1.0.30",
				})
			})
			It("then it should return the bytes representation of the object", func() {
				Ω(mycloud).ShouldNot(BeEmpty())
			})
			It("then should contain the correct number of networks and azs", func() {
				var mynetwork = new(enaml.ManualNetwork)
				fmt.Println(string(mycloud))
				ccManifest := enaml.NewCloudConfigManifest(mycloud)
				testNetwork, _ := yaml.Marshal(ccManifest.Networks[0])
				yaml.Unmarshal(testNetwork, mynetwork)
				subnetCount := len(mynetwork.Subnets)
				azCount := len(ccManifest.AZs)
				Ω(azCount).Should(BeNumerically(">", 1))
				Ω(subnetCount).Should(BeNumerically(">", 1))
			})
		})
	})
})

func testMultiAZFlagExists(flagname string, flags []cli.Flag) {
	for i := 1; i <= AZCountSupported; i++ {
		fn := CreateFlagnameWithSuffix(flagname, i)
		testFlagExists(fn, flags)
	}
}

func testFlagExists(flagname string, flags []cli.Flag) {
	var err = fmt.Errorf("could not find flag %s", flagname)
	for _, flg := range flags {
		if flg.GetName() == flagname {
			err = nil
		}
	}
	Ω(err).ShouldNot(HaveOccurred())
}
