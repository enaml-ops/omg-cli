package boshinit

import "github.com/bosh-ops/bosh-install/deployments/bosh-init/enaml-gen/director"

func NewNats(user, pass string) director.Nats {
	return director.Nats{
		Address:  "127.0.0.1",
		User:     user,
		Password: pass,
	}

}
