package boshinit_test

import (
	. "github.com/enaml-ops/omg-cli/plugins/products/bosh-init"
	"github.com/enaml-ops/omg-cli/plugins/products/bosh-init/enaml-gen/director"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("NewDirectorProperty", func() {
	Describe("given the function", func() {
		Context("when called w/ valid parameters", func() {
			var pgsql *PgSql
			var controlName = "my-bosh"
			var controlCpiJob = "cpijob"
			var mydirector *director.Director
			BeforeEach(func() {
				pgsql = NewPostgres("theuser", "thehost", "thepass", "thedb", "theadapter")
				mydirector = NewDirector(controlName, controlCpiJob, pgsql.GetDirectorDB())
			})
			It("then it should return a valid directorProperty object", func() {
				Ω(mydirector.Name).Should(Equal(controlName))
				Ω(mydirector.CpiJob).Should(Equal(controlCpiJob))
				Ω(mydirector.Db).Should(Equal(pgsql.GetDirectorDB()))
			})
		})
	})
})
