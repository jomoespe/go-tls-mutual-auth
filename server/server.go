package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"golang.org/x/net/http2"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

func SampleHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	io.WriteString(w, `{"protocol": "`+r.Proto+`","common name": "`+r.TLS.PeerCertificates[0].Subject.CommonName+`"}`)
	fmt.Printf("remote: %s, request uri: %s, protocol: %s, subject name: %s\n",
		r.RemoteAddr, r.RequestURI, r.Proto, r.TLS.PeerCertificates[0].Subject.CommonName)

	// https://golang.org/pkg/net/http/#Request
	// https://golang.org/pkg/crypto/tls/#ConnectionState
	// https://golang.org/pkg/crypto/x509/#Certificate
	//fmt.Fprintf(w, "\tpeer certificates=%d\n", len(req.TLS.PeerCertificates))
	// Get the peer get the subject common name from tls request peer certificate
	//fmt.Fprintf(w, "\tCommon name=%s\n", req.TLS.PeerCertificates[0].Subject.CommonName)
}

func main() {
	http.HandleFunc("/sample", SampleHandler)

	caCert, err := ioutil.ReadFile("client.crt")
	if err != nil {
		log.Fatal(err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	// setup HTTPS client (¿client? ¿wtf?)
	tlsConfig := &tls.Config{
		ClientCAs: caCertPool,
		// NoClientCent
		// RequestClientCert
		// RequiredAnyClientCert
		// VerifyClientCartIfGiven
		// RequireAndVerifyClientCert
		ClientAuth: tls.RequireAndVerifyClientCert,
	}
	tlsConfig.BuildNameToCertificate()

	server := &http.Server{
		Addr:      ":8443",
		TLSConfig: tlsConfig,
	}
	http2.ConfigureServer(server, nil)
	server.ListenAndServeTLS("server.crt", "server.key")
}
