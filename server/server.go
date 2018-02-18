package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"golang.org/x/net/http2"
)

func SampleHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("GET /sample requested")
	w.Header().Set("Content-Type", "application/json")
	io.WriteString(w, `{"protocol": "`+r.Proto+`","common name": "`+r.TLS.PeerCertificates[0].Subject.CommonName+`"}`)
	fmt.Printf("remote: %s, request uri: %s, protocol: %s, subject name: %s\n",
		r.RemoteAddr, r.RequestURI, r.Proto, r.TLS.PeerCertificates[0].Subject.CommonName)

	fmt.Println("GET /sample responded")
}

func verifyPeerCertificate(rawCerts [][]byte, verifiedChains [][]*x509.Certificate) error {
	fmt.Println("verifying now")
	return nil
}

func main() {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	pkiPath := filepath.Join(cwd, "pki")
	rootCertPath := filepath.Join(pkiPath, "ca.crt")

	http.HandleFunc("/sample", SampleHandler)

	caCert, err := ioutil.ReadFile(rootCertPath)
	if err != nil {
		log.Fatal(err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	tlsConfig := &tls.Config{
		ClientCAs:             caCertPool,
		ClientAuth:            tls.RequireAndVerifyClientCert,
		VerifyPeerCertificate: verifyPeerCertificate,
	}

	tlsConfig.BuildNameToCertificate()

	server := &http.Server{
		Addr:      "0.0.0.0:8080",
		TLSConfig: tlsConfig,
	}

	http2.ConfigureServer(server, nil)
	fmt.Println("Server Started", os.Getenv("GODEBUG"))

	certPath := filepath.Join(pkiPath, "issued", "sunjin.local.crt")
	keyPath := filepath.Join(pkiPath, "private", "sunjin.local.key")
	if err := server.ListenAndServeTLS(certPath, keyPath); err != nil {
		log.Fatal(err)
	}
}
