package boshinit

import (
	"crypto/sha512"
	"fmt"
)

//SHA512Pass creates a sha-512 password
func SHA512Pass(password string) (string, error) {
	h := sha512.New()
	h.Write([]byte(password))
	bytes := h.Sum(nil)
	return fmt.Sprintf("%x", bytes), nil
}
