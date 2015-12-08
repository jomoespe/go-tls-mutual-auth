package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

func HelloServer(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, "hello world!\n")
	// https://golang.org/pkg/net/http/#Request
	// https://golang.org/pkg/crypto/tls/#ConnectionState
	// https://golang.org/pkg/crypto/x509/#Certificate
	//fmt.Fprintf(w, "\tpeer certificates=%d\n", len(req.TLS.PeerCertificates))
	// Get the peer get the subject common name from tls request peer certificate
	fmt.Fprintf(w, "\tCommon name=%s\n", req.TLS.PeerCertificates[0].Subject.CommonName)
}

func main() {
	http.HandleFunc("/hello", HelloServer)

	caCert, err := ioutil.ReadFile("cert/client/selfsigned.crt")
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
		Addr: ":8080",
		TLSConfig: tlsConfig,
	}

	server.ListenAndServeTLS("cert/server/selfsigned.crt", "cert/server/selfsigned.key") // private cert
}