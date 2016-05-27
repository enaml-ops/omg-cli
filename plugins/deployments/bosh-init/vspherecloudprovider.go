package boshinit

import (
	"fmt"

	"github.com/enaml-ops/enaml"
	"github.com/enaml-ops/omg-cli/plugins/deployments/bosh-init/enaml-gen/vsphere_cpi"
)

// NewVSphereCloudProvider creates a new cloud_provider instance for a vSphere bosh deployment
func NewVSphereCloudProvider(directorIP string, vcenterProperty vsphere_cpi.Vcenter, ntpProperty []string) enaml.CloudProvider {
	var mbusUserPass = "mbus:mbus-password"

	return enaml.CloudProvider{
		Template: enaml.Template{
			Name:    "vsphere_cpi",
			Release: "bosh-vsphere-cpi",
		},
		MBus: fmt.Sprintf("https://%s@%s:6868", mbusUserPass, directorIP),
		Properties: map[string]interface{}{
			"vcenter": vcenterProperty,
			"agent": map[string]string{
				"mbus": fmt.Sprintf("https://%s@0.0.0.0:6868", mbusUserPass),
			},
			"blobstore": map[string]string{
				"provider": "local",
				"path":     "/var/vcap/micro_bosh/data/cache",
			},
			"ntp": ntpProperty,
		},
	}
}
