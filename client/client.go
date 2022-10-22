package main

import (
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"golang.org/x/net/http2"
)

func main() {
	// process parameters
	totalRequest := flag.Int("request", 1, "How many request perform.")
	flag.Parse()

	// load client cert
	cert, err := tls.LoadX509KeyPair("client.crt", "client.key")
	if err != nil {
		log.Fatal(err)
	}

	// load CA cert
	caCert, err := ioutil.ReadFile("server.crt")
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

	for i := 0; i < *totalRequest; i++ {
		client := &http.Client{
			Transport: &http2.Transport{TLSClientConfig: tlsConfig},
		}

		resp, err := client.Get("https://localhost:8443/hello")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		contents, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println(err)
			os.Exit(2)
		}

		fmt.Printf("%s\n", string(contents))
	}
}
