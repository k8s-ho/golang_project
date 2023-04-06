# ssl_server
After generating the CA private key, it is a program that creates a CA certificate with self-signing and  
generates a server certificate signed using the CA private key.

---
### Usage
```bash
* Server Side: excute server
go build
./ssl_server

* Client side
curl https://localhost:7777 --cacert ca.crt
```  
  
<br>

### Dev information
```bash
go 1.20.2 darwin/arm64
```
<br>

### Reference
```bash
https://pkg.go.dev/crypto/x509
```
---
<img width="855" alt="스크린샷 2023-04-06 오후 1 54 11" src="https://user-images.githubusercontent.com/118821939/230275287-bced8e01-abaa-4ef1-a0da-5d2c8138feb9.png">
