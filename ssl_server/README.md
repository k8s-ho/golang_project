# ssl_server
After generating the CA private key, it is a program that creates a CA certificate with self-signing and  
generates a server certificate signed using the CA private key.

```bash
* Server Side: excute server
go build
./ssl_server

* Client side
curl https://localhost:7777 --cacert ca.crt
```
