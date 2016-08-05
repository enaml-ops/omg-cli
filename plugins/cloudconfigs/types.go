package cloudconfigs

import "github.com/enaml-ops/enaml"

type CloudConfigProvider interface {
	CreateAZs() ([]enaml.AZ, error)
	CreateNetworks() ([]enaml.DeploymentNetwork, error)
	CreateVMTypes() ([]enaml.VMType, error)
	CreateDiskTypes() ([]enaml.DiskType, error)
	CreateCompilation() (*enaml.Compilation, error)
}
