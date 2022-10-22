package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"golang.org/x/net/http2"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Receided request from %s\n", r.RemoteAddr)
	fmt.Printf("  Number of client certificates: %d\n", len(r.TLS.PeerCertificates))
	for i, c := range r.TLS.PeerCertificates {
		fmt.Printf("  Cert number %d\n", i)
		printCert(c)
	}

	w.Header().Set("Content-Type", "application/json")
	io.WriteString(w, fmt.Sprintf("{\"protocol\": \"%s\", \"commonName\": \"%s\"}", r.Proto, r.TLS.PeerCertificates[0].Subject.CommonName))

	// https://golang.org/pkg/net/http/#Request
	// https://golang.org/pkg/crypto/tls/#ConnectionState
	// https://golang.org/pkg/crypto/x509/#Certificate
	//fmt.Fprintf(w, "\tpeer certificates=%d\n", len(req.TLS.PeerCertificates))
	// Get the peer get the subject common name from tls request peer certificate
	//fmt.Fprintf(w, "\tCommon name=%s\n", req.TLS.PeerCertificates[0].Subject.CommonName)
}

func livenessHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func main() {
	http.HandleFunc("/hello", helloHandler)
	http.HandleFunc("/healthy", livenessHandler)

	caCert, err := os.ReadFile("client.crt")
	if err != nil {
		log.Fatal(err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	// setup HTTPS client (¿client? ¿wtf?)
	tlsConfig := &tls.Config{
		ClientCAs:  caCertPool,
		ClientAuth: tls.RequireAndVerifyClientCert,
		// other possible values:
		//   tls.NoClientCent
		//   tls.RequestClientCert
		//   tls.RequiredAnyClientCert
		//   tls.VerifyClientCartIfGiven
		//   tls.RequireAndVerifyClientCert
	}

	server := &http.Server{
		Addr:      ":8443",
		TLSConfig: tlsConfig,
	}
	http2.ConfigureServer(server, nil)
	fmt.Printf("TLS Muthual auth server listening on port %s\n", server.Addr)
	server.ListenAndServeTLS("server.crt", "server.key")
}

func printCert(c *x509.Certificate) {
	fmt.Printf("    Subject: %v\n", c.Subject)
	fmt.Printf("    Issuer: %v\n", c.Issuer)
	fmt.Printf("    IPAddresses: %v\n", c.IPAddresses)
	fmt.Printf("    Email addresses: %v\n", c.EmailAddresses)
	fmt.Printf("    Total extensions: %d\n", len(c.Extensions))
	for _, ext := range c.Extensions {
		fmt.Printf("      %s: %x\n", ext.Id, ext.Value[:])
	}
	fmt.Printf("    Total extra extensions: %d\n", len(c.ExtraExtensions))
	for _, ext := range c.ExtraExtensions {
		fmt.Printf("      %s: %x\n", ext.Id, ext.Value[:])
	}
	fmt.Printf("    Is CA: %v\n", c.IsCA)
	fmt.Printf("    Issuing certificate URL: %v\n", c.IssuingCertificateURL)
	fmt.Printf("    Serial number: %v\n", c.SerialNumber)
	fmt.Printf("    URIs: %v\n", c.URIs)
	fmt.Printf("    Version: %v\n", c.Version)
}
