package boshinit_test

import (
	. "github.com/enaml-ops/omg-cli/plugins/deployments/bosh-init"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("GetRegistry", func() {
	var postgresDB = NewPostgres("postgres", "127.0.0.1", "postgres-password", "bosh", "postgres")
	Context("when called with valid args", func() {
		It("should yield a complete and valid object", func() {
			reg := GetRegistry(BoshInitConfig{}, postgresDB)
			Ω(reg.Http).ShouldNot(BeNil())
			Ω(reg.Db).ShouldNot(BeNil())
			Ω(reg.Username).ShouldNot(BeEmpty())
			Ω(reg.Password).ShouldNot(BeEmpty())
		})
	})
})
