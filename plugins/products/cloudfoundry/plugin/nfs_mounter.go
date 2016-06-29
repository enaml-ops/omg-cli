package cloudfoundry

import (
	"github.com/codegangsta/cli"
	"github.com/enaml-ops/enaml"
	nfsmounterlib "github.com/enaml-ops/omg-cli/plugins/products/cloudfoundry/enaml-gen/nfs_mounter"
)

//NewNfsMounter -
func NewNFSMounter(c *cli.Context) *NFSMounter {
	return &NFSMounter{
		NFSServerAddress: c.String("nfs-server-address"),
		SharePath:        c.String("nfs-share-path"),
	}
}

//CreateJob - Create the yaml job structure for NFSMounter
func (s *NFSMounter) CreateJob() enaml.InstanceJob {
	return enaml.InstanceJob{
		Name:    "nfs_mounter",
		Release: "cf",
		Properties: &nfsmounterlib.NfsMounter{
			NfsServer: &nfsmounterlib.NfsServer{
				Address: s.NFSServerAddress,
				Share:   s.SharePath,
			},
		},
	}
}
