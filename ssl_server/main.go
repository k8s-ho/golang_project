package main

import (
	"crypto/x509"
	"crypto/x509/pkix"
	"math/big"
	"time"
)

func main() {
	caCN := "kubernetes"
	serverCN := "localhost"
	caKey := generateKey("ca")
	orgs := []string{"k8sho.com"}

	ca_template := x509.Certificate{
		SerialNumber:          big.NewInt(time.Now().Unix()),
		Subject:               pkix.Name{CommonName: caCN, Organization: orgs},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(1, 0, 0),
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign | x509.KeyUsageKeyEncipherment,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		IsCA:                  true,
		//DNSNames:     []string{"k8s-ho.com"}, //tmp
	}
	generateCrt(caKey, caKey, "ca", ca_template, ca_template)

	serverKey := generateKey("server")
	server_template := x509.Certificate{
		SerialNumber:          big.NewInt(time.Now().Unix()),
		Subject:               pkix.Name{CommonName: serverCN, Organization: orgs}, 
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(1, 0, 0),
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		IsCA:                  false,
		// DNSNames:              []string{"aaa.bbb.com"}, //tmp
	}
	generateCrt(caKey, serverKey, "server", server_template, ca_template)
	server()
}
