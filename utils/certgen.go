package utils

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"log"
	"math/big"
	"net"
	"os"
	"time"
)

func pemBlockForKey(priv interface{}) (block *pem.Block, err error) {
	block = nil
	switch k := priv.(type) {
	case *rsa.PrivateKey:
		block = &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(k)}
		return
	case *ecdsa.PrivateKey:
		var b []byte
		b, err = x509.MarshalECPrivateKey(k)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to marshal ECDSA private key: %v", err)
			return
		}
		block = &pem.Block{Type: "EC PRIVATE KEY", Bytes: b}
	default:
		block = nil
	}

	return
}

//GenerateCert - will generate a cert based on hosts and if it is CA
func GenerateCert(hosts []string) (caCert, cert, key string, err error) {
	var caKey, certKey *rsa.PrivateKey
	var caCertBytes, certBytes []byte
	var serialNumber *big.Int
	notBefore := time.Now()
	oneYear := 365 * 24 * time.Hour
	notAfter := notBefore.Add(oneYear)
	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	if serialNumber, err = rand.Int(rand.Reader, serialNumberLimit); err != nil {
		log.Fatalf("failed to generate serial number: %s", err)
		return
	}

	ca := &x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{"enaml-ops"},
		},
		NotBefore:             notBefore,
		NotAfter:              notAfter,
		SubjectKeyId:          []byte{1, 2, 3, 4, 5},
		BasicConstraintsValid: true,
		IsCA:        true,
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:    x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
	}

	if caKey, err = rsa.GenerateKey(rand.Reader, 2048); err != nil {
		return
	}
	if certKey, err = rsa.GenerateKey(rand.Reader, 2048); err != nil {
		return
	}
	caCertBytes, err = x509.CreateCertificate(rand.Reader, ca, ca, &caKey.PublicKey, caKey)

	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{"enaml-ops"},
		},
		NotBefore: notBefore,
		NotAfter:  notAfter,

		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}
	for _, h := range hosts {
		if ip := net.ParseIP(h); ip != nil {
			template.IPAddresses = append(template.IPAddresses, ip)
		} else {
			template.DNSNames = append(template.DNSNames, h)
		}
	}

	certBytes, err = x509.CreateCertificate(rand.Reader, &template, ca, &certKey.PublicKey, certKey)
	if err != nil {
		log.Fatalf("Failed to create certificate: %s", err)
		return
	}

	caCert = string(pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: caCertBytes}))
	cert = string(pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certBytes}))
	var keyBlock *pem.Block
	if keyBlock, err = pemBlockForKey(certKey); err == nil {
		keyBytes := pem.EncodeToMemory(keyBlock)
		key = string(keyBytes)
	}
	return
}
