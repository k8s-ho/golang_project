# ssl_server
After generating the CA private key, it is a program that creates a CA certificate with self-signing and  
generates a server certificate signed using the CA private key.

```bash
go build
./ssl_server
curl https://localhost:7777 --cacert ca.crt
```