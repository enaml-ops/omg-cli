package boshinitaws

import "github.com/bosh-ops/bosh-install/deployments/bosh-init-aws/enaml-gen/registry"

func GetRegistry(cfg BoshInitConfig, postgresDB *pg) AWSRegistryProperty {
	return AWSRegistryProperty{
		Address: cfg.BoshPrivateIP,
		Ar: Ar{
			Host:     cfg.BoshPrivateIP,
			Username: "admin",
			Password: "admin",
			Port:     25777,
		},
		Rr: Rr{
			Db: postgresDB.GetRegistryDB(),
			Http: &registry.Http{
				User:     "admin",
				Password: "admin",
				Port:     25777,
			},
		},
	}
}
