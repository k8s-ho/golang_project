package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"log"
	"math/big"
	"os"
	"time"
)

func main() {
	var certfile_name, keyfile_name, common string
	var date int
	fmt.Print("생성하실 KEY 파일명을 입력하세요 = ")
	fmt.Scan(&keyfile_name)
	fmt.Print("생성하실 인증서 파일명을 입력하세요 = ")
	fmt.Scan(&certfile_name)
	fmt.Print("만료기간을 입력하세요(연 단위) = ")
	fmt.Scan(&date)
	fmt.Print("CN을 입력하세요(Common Name) = ")
	fmt.Scan(&common)
	privatekey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Println(err)
	}

	template := x509.Certificate{
		SerialNumber: big.NewInt(time.Now().Unix()),
		Subject:      pkix.Name{CommonName: common},
		NotBefore:    time.Now(),
		NotAfter:     time.Now().AddDate(date, 0, 0),
		KeyUsage:     x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
	}

	data, err := x509.CreateCertificate(rand.Reader, &template, &template, &privatekey.PublicKey, privatekey)
	if err != nil {
		panic(err)
	}

	// 인증서 파일로 생성
	certfile, err := os.OpenFile(certfile_name+".pem", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644) // flag: 쓰기전용, 644 권한부여
	if err != nil {
		panic(err)
	}
	// pem 파일로 변환
	pem.Encode(certfile, &pem.Block{Type: "CERTIFICATE", Bytes: data})
	certfile.Close()

	keyfile, err := os.OpenFile(keyfile_name+".pem", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600) // flag: 쓰기전용 | 없으면생성 | 기존에 존재할 경우 내용을 모두 삭제하고 빈 파일, 600 권한부여
	if err != nil {
		panic(err)
	}
	pem.Encode(keyfile, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(privatekey)})
	keyfile.Close()
	fmt.Println("\n[+] Certificate: " + certfile_name + "\n[+] Key file: " + keyfile_name + "\n 파일이 생성되었습니다!\n")
}
