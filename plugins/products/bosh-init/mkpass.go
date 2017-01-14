package boshinit

import (
	"github.com/tredoe/osutil/user/crypt/sha512_crypt"
)

//SHA512Pass creates a sha-512 password
func SHA512Pass(password string) (string, error) {
	c := sha512_crypt.New()
	return c.Generate([]byte(password), nil)
}
