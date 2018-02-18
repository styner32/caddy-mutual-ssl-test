package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"golang.org/x/net/http2"
)

func main() {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	pkiPath := filepath.Join(cwd, "pki")

	// load client cert
	cert, err := tls.LoadX509KeyPair(filepath.Join(pkiPath, "issued", "client0.crt"), filepath.Join(pkiPath, "private", "client0.key"))
	if err != nil {
		log.Fatal(err)
	}

	// load CA cert
	caCert, err := ioutil.ReadFile(filepath.Join(pkiPath, "ca.crt"))
	if err != nil {
		log.Fatal(err)
	}

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	// setup HTTPS client
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      caCertPool,
	}

	tlsConfig.BuildNameToCertificate()

	client := &http.Client{
		Transport: &http2.Transport{TLSClientConfig: tlsConfig},
	}

	resp, err := client.Get("https://caddy.local:8443/sample")
	if err != nil {
		log.Fatal(err)
	}

	contents, err := ioutil.ReadAll(resp.Body)
	fmt.Printf("%s\n", string(contents))
}
