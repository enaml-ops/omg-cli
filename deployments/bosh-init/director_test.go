package boshinit_test

import (
	. "github.com/bosh-ops/bosh-install/deployments/bosh-init"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("NewDirectorProperty", func() {
	Describe("given the function", func() {
		Context("when called w/ valid parameters", func() {
			var pgsql *PgSql
			var controlName = "my-bosh"
			var controlCpiJob = "cpijob"
			var mydirector DirectorProperty
			BeforeEach(func() {
				pgsql = NewPostgres("theuser", "thehost", "thepass", "thedb", "theadapter")
				mydirector = NewDirectorProperty(controlName, controlCpiJob, pgsql.GetDirectorDB())
			})
			It("then it should return a valid directorProperty object", func() {
				Ω(mydirector.Director.Name).Should(Equal(controlName))
				Ω(mydirector.Director.CpiJob).Should(Equal(controlCpiJob))
				Ω(mydirector.Director.Db).Should(Equal(pgsql.GetDirectorDB()))
			})
		})
	})
})
