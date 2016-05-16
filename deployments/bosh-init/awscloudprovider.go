package boshinit

import (
	"fmt"

	"github.com/bosh-ops/bosh-install/deployments/bosh-init/enaml-gen/aws_cpi"
	"github.com/xchapter7x/enaml"
)

func NewAWSCloudProvider(awsElasticIP, awsPEMFilePath string, awsProperty aws_cpi.Aws, ntpProperty []string) enaml.CloudProvider {
	var mbusUserPass = "mbus:mbus-password"

	return enaml.CloudProvider{
		Template: enaml.Template{
			Name:    "aws_cpi",
			Release: "bosh-aws-cpi",
		},
		MBus: fmt.Sprintf("https://%s@%s:6868", mbusUserPass, awsElasticIP),
		SSHTunnel: enaml.SSHTunnel{
			Host:           awsElasticIP,
			Port:           22,
			User:           "vcap",
			PrivateKeyPath: awsPEMFilePath,
		},
		Properties: map[string]interface{}{
			"aws": awsProperty,
			"ntp": ntpProperty,
			"agent": map[string]string{
				"mbus": fmt.Sprintf("https://%s@0.0.0.0:6868", mbusUserPass),
			},
			"blobstore": map[string]string{
				"provider": "local",
				"path":     "/var/vcap/micro_bosh/data/cache",
			},
		},
	}
}
