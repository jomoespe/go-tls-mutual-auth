package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"flag"
	"fmt"
	"math/big"
	"net"
	"os"
	"time"
)

const (
	ORGANIZATION_DEFAULT = "My Organization"
	COMMON_NAME_DEFAULT  = "localhost"
	DURATION_DEFAULT     = 365
	CERTIFICATE_DEFAULT  = "selfsigned.crt"
	KEY_DEFAULT          = "selfsigned.key"
	IS_CLIENT_DEFAULT    = false
)

func main() {
	organization := flag.String("org", ORGANIZATION_DEFAULT, "CA Organization nane")
	commonName := flag.String("name", COMMON_NAME_DEFAULT, "The subject name. Usually the DNS server name")
	duration := flag.Int("duration", DURATION_DEFAULT, "How log the certificate will be valid.")
	certFilename := flag.String("cert", CERTIFICATE_DEFAULT, "Certificate filename.")
	keyFilename := flag.String("key", KEY_DEFAULT, "Privake Key filename.")
	isClientCert := flag.Bool("client", IS_CLIENT_DEFAULT, "If the certificate usage is client. Default is false (server usage)")
	flag.Parse()
	addresses := []string{"localhost", "127.0.0.1"} // default localhost
	if len(flag.Args()) > 0 {
		addresses = flag.Args()
	}

	certificate := newCertificate(*organization, *commonName, *duration, *isClientCert, addresses)
	generate(*certificate, *certFilename, *keyFilename)
}

func newCertificate(organization, commonName string, duration int, isClientCert bool, addresses []string) *x509.Certificate {
	certificate := x509.Certificate{
		Subject: pkix.Name{
			Organization: []string{organization},
			CommonName:   commonName,
		},
		NotBefore: time.Now(),
		NotAfter:  time.Now().Add(time.Duration(duration) * time.Hour * 24),

		KeyUsage:    x509.KeyUsageCertSign | x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth}, // for server certificate usage
		//ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth}, // for client certificate usage
		//ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth}, // for both (server & client) certificate usage
		BasicConstraintsValid: true,
		IsCA: true,
	}
	// change key usage if is client cert
	if isClientCert {
		certificate.ExtKeyUsage = []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth}
	}
	//
	for i := 0; i < len(addresses); i++ {
		if ip := net.ParseIP(addresses[i]); ip != nil {
			certificate.IPAddresses = append(certificate.IPAddresses, ip)
		} else {
			certificate.DNSNames = append(certificate.DNSNames, addresses[i])
		}
	}

	return &certificate
}

func generate(certificate x509.Certificate, certFilename, keyFilename string) {
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		fmt.Println("Failed to generate private key:", err)
		os.Exit(1)
	}

	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	certificate.SerialNumber, err = rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		fmt.Println("Failed to generate serial number:", err)
		os.Exit(1)
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, &certificate, &certificate, &priv.PublicKey, priv)
	if err != nil {
		fmt.Println("Failed to create certificate:", err)
		os.Exit(1)
	}

	certOut, err := os.Create(certFilename)
	if err != nil {
		fmt.Println("Failed to open "+certFilename+" for writing:", err)
		os.Exit(1)
	}
	pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	certOut.Close()

	keyOut, err := os.OpenFile(keyFilename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		fmt.Println("Failed to open key "+keyFilename+" for writing:", err)
		os.Exit(1)
	}
	pem.Encode(keyOut, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(priv)})
	keyOut.Close()

	fmt.Println("Certificate generated successfully")
	fmt.Println("\tCertificate: ", certFilename)
	fmt.Println("\tPrivate Key: ", keyFilename)
}

/*
func printCmdArgs(organization, commonName, certFilename, keyFilename *string, isClientCert *bool, duration *int, addresses []string) {
	println("\n____ Command arguments ______________________________________")
	fmt.Println("organization:         ", *organization)
	fmt.Println("common name:          ", *commonName)
	fmt.Println("certificate filename: ", *certFilename)
	fmt.Println("key filename:         ", *keyFilename)
	fmt.Println("is client cert?       ", *isClientCert)
	fmt.Println("duration:             ", *duration)
	fmt.Println("Addresses:            ", addresses)
}

func printCertificate(certificate x509.Certificate, certFilename, keyFilename string) {
	println("\n____ Template info ______________________________________")
	println(certificate.Subject.Organization[0])
	println(certificate.Subject.CommonName)
	println(certificate.ExtKeyUsage[0])
	println(certFilename)
	println(keyFilename)
	fmt.Println("DNS SANs:")
	for _, e := range certificate.DNSNames {
		fmt.Println("   ", e)
	}
	fmt.Println("IP SANs:")
	for _, e := range certificate.IPAddresses {
		fmt.Println("   ", e)
	}
}
*/
