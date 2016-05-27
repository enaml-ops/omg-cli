package boshinit

import (
	"github.com/enaml-ops/omg-cli/plugins/products/bosh-init/enaml-gen/director"
	"github.com/enaml-ops/omg-cli/plugins/products/bosh-init/enaml-gen/postgres"
	"github.com/enaml-ops/omg-cli/plugins/products/bosh-init/enaml-gen/registry"
)

func NewPostgres(user, host, pass, database, adapter string) (psql *PgSql) {
	return &PgSql{
		User:     user,
		Host:     host,
		Password: pass,
		Database: database,
		Adapter:  adapter,
	}
}

func (s *PgSql) GetDirectorDB() *director.Db {
	return &director.Db{
		User:     s.User,
		Host:     s.Host,
		Password: s.Password,
		Database: s.Database,
		Adapter:  s.Adapter,
	}
}
func (s *PgSql) GetRegistryDB() *registry.Db {
	return &registry.Db{
		User:     s.User,
		Host:     s.Host,
		Password: s.Password,
		Database: s.Database,
		Adapter:  s.Adapter,
	}
}
func (s *PgSql) GetPostgresDB() postgres.Postgres {
	return postgres.Postgres{
		User:          s.User,
		ListenAddress: s.Host,
		Password:      s.Password,
		Database:      s.Database,
	}
}
