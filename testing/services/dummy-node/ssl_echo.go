package main

import (
	"crypto/tls"
	_ "embed"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

// certs/out/localhost.crt
var localhostCert string

// certs/out/localhost.csr
// var localhostCsr string

// certs/out/localhost.key
var localhostKey string

// "Inspired" by https://youngkin.github.io/post/gohttpsclientserver/
func StartSslEchoServer() {
	fmt.Printf("Starting SSL echo server...\n")
	host := "localhost"
	port := "8443"

	server := &http.Server{
		Addr:         ":" + port,
		ReadTimeout:  5 * time.Minute, // 5 min to allow for delays when 'curl' on OSx prompts for username/password
		WriteTimeout: 10 * time.Second,
		TLSConfig:    &tls.Config{ServerName: host},
	}

	server.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Received %s request for host %s from IP address %s and X-FORWARDED-FOR %s",
			r.Method, r.Host, r.RemoteAddr, r.Header.Get("X-FORWARDED-FOR"))
		body, err := io.ReadAll(r.Body)
		if err != nil {
			body = []byte(fmt.Sprintf("error reading request body: %s", err))
		}
		resp := fmt.Sprintf("Hello, %s from Simple Server!", body)
		if _, err := w.Write([]byte(resp)); err != nil {
			log.Printf("Error writing response: %v", err)
			return
		}
		log.Printf("SimpleServer: Sent response %s", resp)
	})

	if err := os.WriteFile("localhost.crt", []byte(localhostCert), 0644); err != nil {
		log.Printf("Error writing localhost.crt: %v", err)
		return
	}

	if err := os.WriteFile("localhost.key", []byte(localhostKey), 0644); err != nil {
		log.Printf("Error writing localhost.key: %v", err)
		return
	}

	log.Printf("Starting HTTPS server on host %s and port %s", host, port)
	if err := server.ListenAndServeTLS("localhost.crt", "localhost.key"); err != nil {
		log.Fatal(err)
	}
}
