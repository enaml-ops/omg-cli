package utils

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
)

func GenerateKeys() (publicKeyPem, privateKeyPem string, err error) {
	var privateKey *rsa.PrivateKey
	var publicKeyDer []byte
	if privateKey, err = rsa.GenerateKey(rand.Reader, 2014); err != nil {
		return
	}

	privateKeyDer := x509.MarshalPKCS1PrivateKey(privateKey)
	privateKeyBlock := pem.Block{
		Type:    "RSA PRIVATE KEY",
		Headers: nil,
		Bytes:   privateKeyDer,
	}
	privateKeyPem = string(pem.EncodeToMemory(&privateKeyBlock))

	publicKey := privateKey.PublicKey
	if publicKeyDer, err = x509.MarshalPKIXPublicKey(&publicKey); err != nil {
		return
	}

	publicKeyBlock := pem.Block{
		Type:    "PUBLIC KEY",
		Headers: nil,
		Bytes:   publicKeyDer,
	}
	publicKeyPem = string(pem.EncodeToMemory(&publicKeyBlock))
	return
}
