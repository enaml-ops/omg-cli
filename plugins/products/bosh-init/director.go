package boshinit

import "github.com/enaml-ops/omg-cli/plugins/products/bosh-init/enaml-gen/director"

func NewDirector(name, cpijob string, db *director.DirectorDb) *director.Director {
	return &director.Director{
		Name:       name,
		CpiJob:     cpijob,
		MaxThreads: 10,
		Db:         db,
		UserManagement: &director.UserManagement{
			Provider: "local",
			Local: &director.Local{
				Users: []user{
					user{Name: "admin", Password: "admin"},
					user{Name: "hm", Password: "hm-password"},
				},
			},
		},
	}
}
