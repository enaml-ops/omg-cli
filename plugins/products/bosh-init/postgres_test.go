package boshinit_test

import (
	. "github.com/enaml-ops/omg-cli/plugins/products/bosh-init"
	"github.com/enaml-ops/omg-cli/plugins/products/bosh-init/enaml-gen/director"
	"github.com/enaml-ops/omg-cli/plugins/products/bosh-init/enaml-gen/postgres"
	"github.com/enaml-ops/omg-cli/plugins/products/bosh-init/enaml-gen/registry"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("NewPostgres func", func() {
	Describe("given the function pgsql object returned", func() {
		var controluser = "user"
		var controlhost = "host"
		var controlpass = "pass"
		var controldatabase = "db"
		var controladapter = "adapt"
		var pgsql *PgSql

		BeforeEach(func() {
			pgsql = NewPostgres(controluser, controlhost, controlpass, controldatabase, controladapter)
		})

		Context("when calling GetDirectorDB", func() {
			var ddb *director.Db
			BeforeEach(func() {
				ddb = pgsql.GetDirectorDB()
			})
			It("then it should return a properly initialized director.Db object", func() {
				Ω(ddb.User).Should(Equal(controluser))
				Ω(ddb.Host).Should(Equal(controlhost))
				Ω(ddb.Password).Should(Equal(controlpass))
				Ω(ddb.Database).Should(Equal(controldatabase))
				Ω(ddb.Adapter).Should(Equal(controladapter))
			})
		})

		Context("when calling GetRegistryDB", func() {
			var rdb *registry.Db
			BeforeEach(func() {
				rdb = pgsql.GetRegistryDB()
			})
			It("then it should return a properly initialized registry.Db object", func() {
				Ω(rdb.User).Should(Equal(controluser))
				Ω(rdb.Host).Should(Equal(controlhost))
				Ω(rdb.Password).Should(Equal(controlpass))
				Ω(rdb.Database).Should(Equal(controldatabase))
				Ω(rdb.Adapter).Should(Equal(controladapter))
			})
		})

		Context("when calling GetPostgresDB", func() {
			var pdb postgres.Postgres
			BeforeEach(func() {
				pdb = pgsql.GetPostgresDB()
			})
			It("then it should return a properly initialized postgres.Postgres object", func() {
				Ω(pdb.User).Should(Equal(controluser))
				Ω(pdb.ListenAddress).Should(Equal(controlhost))
				Ω(pdb.Password).Should(Equal(controlpass))
				Ω(pdb.Database).Should(Equal(controldatabase))
			})
		})
	})
})
