// build with:
//		go get -v github.com/streadway/amqp
//		go build -v -ldflags="-s -w"
package main

import (
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/streadway/amqp"
)

var (
	address = flag.String("url", "", "e.g. amqps://username:password@example.com/")
)

func main() {
	log.SetOutput(os.Stdout) // for not disturbing PowerShell...

	flag.Parse()

	if *address == "" {
		log.Printf("Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
		return
	}

	tlsConfig := new(tls.Config)
	tlsConfig.InsecureSkipVerify = true
	tlsConfig.VerifyPeerCertificate = func(rawCerts [][]byte, verifiedChains [][]*x509.Certificate) error {
		for i, crt := range rawCerts {
			var subject pkix.Name
			var issuer pkix.Name
			certificates, err := x509.ParseCertificates(crt)
			if err != nil {
				log.Printf("WARN failed to parse %s chain link #%d: %v\n", tlsConfig.ServerName, i, err)
			} else {
				subject = certificates[0].Subject
				issuer = certificates[0].Issuer
			}
			path := fmt.Sprintf("%s-%d.der", tlsConfig.ServerName, i)
			log.Printf(
				"Saving %s certificate chain link #%d (subject=%s; issuer=%s) to %s...",
				tlsConfig.ServerName,
				i,
				subject.CommonName,
				issuer.CommonName,
				path)
			ioutil.WriteFile(
				path,
				crt,
				0644)
		}
		return nil
	}

	c, err := amqp.DialTLS(*address, tlsConfig)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v\n", err)
	}
	defer c.Close()

	logProperties(c.Properties, "")
}

func logProperties(properties amqp.Table, prefix string) {
	for key, value := range properties {
		if innerProperties, ok := value.(amqp.Table); ok {
			logProperties(innerProperties, prefix+key+".")
		} else {
			log.Printf("property %s%s = %v\n", prefix, key, value)
		}
	}
}
