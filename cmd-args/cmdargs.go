package main

import (
	"crypto/x509"
	"crypto/x509/pkix"
	"flag"
	//	"fmt"
	"time"
)

const ORG_ARG = "org"
const ORG_ARG_DEFAULT = "Jose Moreno"
const ORG_ARG_DESC = "CA Organization nane"
const NAME_ARG = "name"
const NAME_ARG_DEFAULT = "localhost"
const NAME_ARG_DESC = "The subject name. Usually the DNS server name"

var org, name, certFilename, keyFilename string
var isClientCert bool
var duration int
var addresses []string

func main() {
	isClientCert = *flag.Bool("client", false, "If the certificate usage is client. Default is false (server usage)")
	println(isClientCert)
	
	template := x509.Certificate{
		Subject: pkix.Name{
			Organization: []string{*flag.String(ORG_ARG, ORG_ARG_DEFAULT, ORG_ARG_DESC)},
			CommonName:   *flag.String(NAME_ARG, NAME_ARG_DEFAULT, NAME_ARG_DESC),
		},
		NotBefore: time.Now(),
		NotAfter:  time.Now().Add(time.Duration(duration) * time.Hour * 24),

		KeyUsage:    x509.KeyUsageCertSign | x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth}, // for server certificate usage
		//ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},                             // for client certificate usage
		//ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth}, // for both (server & client) certificate usage
		BasicConstraintsValid: true,
		IsCA: true,
	}

	//	org          = *flag.String("org",         "",               "CA Organization nane")
	//	name         = *flag.String("name",        "localhost",      "The subject name. Usually the DNS server name")
	//	certFilename = *flag.String("certificate", "selfsigned.crt", "Certificate filename")
	//	keyFilename  = *flag.String("key",         "selfsigned.key", "Privake Key filename")
	//	isClientCert = *flag.Bool("clientCert",    false,            "If the certificate usage is client. Default is false (server usage)")
	//	duration     = *flag.Int("duration",       365,              "How log the certificate will be valid.")
	//	flag.Parse()
	//	addresses    = flag.Args()

	//fmt.Println("org:                  ", org)
	//fmt.Println("name:                 ", name)
	//fmt.Println("certificate filename: ", certFilename)
	//fmt.Println("key filename:         ", keyFilename)
	//fmt.Println("is client cert:       ", isClientCert)
	//fmt.Println("duration:             ", duration)
	//fmt.Println("Addresses:            ", addresses)

	generate(&template)
}

func generate(template *x509.Certificate) {
	println(template.Subject.Organization[0])
	println(template.Subject.CommonName)
}
