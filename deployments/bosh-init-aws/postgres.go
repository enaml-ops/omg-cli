package boshinitaws

import (
	"github.com/bosh-ops/bosh-install/deployments/bosh-init-aws/enaml-gen/director"
	"github.com/bosh-ops/bosh-install/deployments/bosh-init-aws/enaml-gen/postgres"
	"github.com/bosh-ops/bosh-install/deployments/bosh-init-aws/enaml-gen/registry"
)

func NewPostgres(user, host, pass, database, adapter string) (psql *pg) {
	return &pg{
		User:     user,
		Host:     host,
		Password: pass,
		Database: database,
		Adapter:  adapter,
	}
}

func (s *pg) GetDirectorDB() *director.Db {
	return &director.Db{
		User:     s.User,
		Host:     s.Host,
		Password: s.Password,
		Database: s.Database,
		Adapter:  s.Adapter,
	}
}
func (s *pg) GetRegistryDB() *registry.Db {
	return &registry.Db{
		User:     s.User,
		Host:     s.Host,
		Password: s.Password,
		Database: s.Database,
		Adapter:  s.Adapter,
	}
}
func (s *pg) GetPostgresDB() postgres.Postgres {
	return postgres.Postgres{
		User:          s.User,
		ListenAddress: s.Host,
		Password:      s.Password,
		Database:      s.Database,
	}
}
