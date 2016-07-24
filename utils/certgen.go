package utils

import (
	"crypto/x509"
	"fmt"
	"net"

	"github.com/square/certstrap/pkix"
)

const (
	orgUnit     = "enaml-ops"
	caCertYears = 10
	signYears   = 2
	org         = "omg"
	country     = "US"
)

func GenerateCert(hosts []string) (caCert, cert, key string, err error) {
	var cakey, crtKey *pkix.Key
	var cacrt, crt *pkix.Certificate
	if cakey, cacrt, err = initialize("enaml"); err == nil {
		if crt, crtKey, err = createCert(cacrt, cakey, hosts); err != nil {
			return
		}
		crtBytes, _ := crt.Export()
		cert = string(crtBytes[:])

		crkKeyBytes, _ := crtKey.ExportPrivate()
		key = string(crkKeyBytes[:])

		caCrtBytes, _ := cacrt.Export()
		caCert = string(caCrtBytes[:])
	}

	return
}

func initialize(commonName string) (key *pkix.Key, crt *pkix.Certificate, err error) {
	if key, err = pkix.CreateRSAKey(2048); err != nil {
		return
	}
	if crt, err = pkix.CreateCertificateAuthority(key, orgUnit, caCertYears, org, "", "", "", commonName); err != nil {
		return
	}
	return
}

func createCert(cacrt *pkix.Certificate, cakey *pkix.Key, hosts []string) (crt *pkix.Certificate, csrKey *pkix.Key, err error) {
	var csr *pkix.CertificateSigningRequest
	var rawCrt *x509.Certificate
	var name = ""
	var domains []string
	var ips []net.IP

	// Validate that crt is allowed to sign certificates.
	if rawCrt, err = cacrt.GetRawCertificate(); err != nil {
		return
	}

	if !rawCrt.IsCA {
		err = fmt.Errorf("Selected CA certificate is not allowed to sign certificates")
		return
	}
	for _, h := range hosts {
		if ip := net.ParseIP(h); ip != nil {
			ips = append(ips, ip)
		} else {
			domains = append(domains, h)
		}
	}

	switch {
	case len(domains) != 0:
		name = domains[0]
	case len(ips) != 0:
		name = ips[0].String()
	default:
		err = fmt.Errorf("Must provide Common Name or SAN")
		return
	}

	if csrKey, err = pkix.CreateRSAKey(2048); err != nil {
		return
	}
	if csr, err = pkix.CreateCertificateSigningRequest(csrKey, orgUnit, ips, domains, org, country, "", "", name); err != nil {
		return
	}
	if crt, err = pkix.CreateCertificateHost(cacrt, cakey, csr, signYears); err != nil {
		return
	}
	return
}
