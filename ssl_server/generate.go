package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"log"
)

func generateKey(name string) *rsa.PrivateKey {
	privatekey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Println(err)
	}

	privData := x509.MarshalPKCS1PrivateKey(privatekey)
	privatefile := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privData,
	})

	err = ioutil.WriteFile(name+".key", privatefile, 0600)
	if err != nil {
		panic(err)
	}
	fmt.Println("[+] Generated Key file: " + name + ".key")
	return privatekey
}

func generateCrt(caKey *rsa.PrivateKey, privatekey *rsa.PrivateKey, name string, template x509.Certificate, parentTemplate x509.Certificate) {
	certData, err := x509.CreateCertificate(rand.Reader, &template, &parentTemplate, &privatekey.PublicKey, caKey)
	if err != nil {
		panic(err)
	}

	certfile := pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certData,
	})

	err = ioutil.WriteFile(name+".crt", certfile, 0644)
	if err != nil {
		panic(err)
	}
	fmt.Println("[+] Generated Certificate file: " + name + ".crt")
}
