package main

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
)

func server() {
	cert, err := tls.LoadX509KeyPair("./server.crt", "./server.key")
	if err != nil {
		panic(err)
	}
	tlsConfig := &tls.Config{
		ServerName:         "kubernetes.default.svc",
		Certificates:       []tls.Certificate{cert},
		InsecureSkipVerify: true,
	}

	addr := "0.0.0.0:7777"

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}

	fileServer := http.FileServer(http.Dir("./web"))

	http.Handle("/", fileServer)

	server := &http.Server{
		TLSConfig: tlsConfig,
	}

	err = server.ServeTLS(listener, "./server.crt", "./server.key")
	if err != nil {
		panic(err)
	}
	fmt.Println("Server started")
}
