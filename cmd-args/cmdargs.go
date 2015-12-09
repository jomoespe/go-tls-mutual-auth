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

func main() {

	organization := *flag.String("org", "My Organization", "CA Organization nane")
	commonName := *flag.String("name", "localhost", "The subject name. Usually the DNS server name")
	duration := *flag.Int("duration", 365, "How log the certificate will be valid.")
	certFilename := *flag.String("certificate", "selfsigned.crt", "Certificate filename.")
	keyFilename := *flag.String("key", "selfsigned.key", "Privake Key filename.")
	isClientCert := flag.Bool("client", false, "If the certificate usage is client. Default is false (server usage)")
	flag.Parse()
	addresses := flag.Args()

	println("\n____ Command arguments ______________________________________")
	fmt.Println("organization:         ", organization)
	fmt.Println("common name:          ", commonName)
	fmt.Println("certificate filename: ", certFilename)
	fmt.Println("key filename:         ", keyFilename)
	fmt.Println("is client cert?       ", *isClientCert)
	fmt.Println("duration:             ", duration)
	fmt.Println("Addresses:            ", addresses)

	template := x509.Certificate{
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
	if *isClientCert {
		template.ExtKeyUsage = []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth}
	}
	//
	for i := 0; i < len(addresses); i++ {
		if ip := net.ParseIP(addresses[i]); ip != nil {
			template.IPAddresses = append(template.IPAddresses, ip)
		} else {
			template.DNSNames = append(template.DNSNames, addresses[i])
		}
	}

	printTemplate(template, certFilename, keyFilename)
	//generate(template, certFilename, keyFilename)
}

func printTemplate(template x509.Certificate, certFilename, keyFilename string) {
	println("\n____ Template info ______________________________________")
	println(template.Subject.Organization[0])
	println(template.Subject.CommonName)
	println(template.ExtKeyUsage[0])
	println(certFilename)
	println(keyFilename)

	println("\n____ Addresses ______________________________________")
	fmt.Println("DNS SANs:")
	if len(template.DNSNames) == 0 {
		fmt.Println("    None")
	} else {
		for _, e := range template.DNSNames {
			fmt.Println("   ", e)
		}
	}
	fmt.Println("IP SANs:")
	if len(template.IPAddresses) == 0 {
		fmt.Println("    None")
	} else {
		for _, e := range template.IPAddresses {
			fmt.Println("   ", e)
		}
	}
	fmt.Println()

}

func generate(template x509.Certificate, certFilename, keyFilename string) {
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		fmt.Println("Failed to generate private key:", err)
		os.Exit(1)
	}

	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	template.SerialNumber, err = rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		fmt.Println("Failed to generate serial number:", err)
		os.Exit(1)
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)
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

	fmt.Println("Successfully generated certificate")
	fmt.Println("\tCertificate: ", certFilename)
	fmt.Println("\tPrivate Key: ", keyFilename)
}
