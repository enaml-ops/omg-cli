package boshinit

import "github.com/bosh-ops/bosh-install/deployments/bosh-init/enaml-gen/registry"

func GetRegistry(cfg BoshInitConfig, postgresDB *PgSql) RegistryProperty {
	return RegistryProperty{
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
