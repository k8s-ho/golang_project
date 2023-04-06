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
<img width="779" alt="스크린샷 2023-04-06 오후 1 46 48" src="https://user-images.githubusercontent.com/118821939/230274394-26b73889-8d57-47c2-9664-3fe22a5acb1e.png">
