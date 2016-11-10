package boshinit

import sha512 "github.com/kless/osutil/user/crypt/sha512_crypt"

//SHA512Pass creates a sha-512 password
func SHA512Pass(password string) (string, error) {
	c := sha512.New()
	if shadowHash, err := c.Generate([]byte(password), []byte("")); err != nil {
		return "", err
	} else {
		return shadowHash, nil
	}
}
