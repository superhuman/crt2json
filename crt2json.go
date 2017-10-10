package main

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
)

func main() {
	if len(os.Args) <= 1 {
		decodeFromStdin()
		return
	}

	if os.Args[1] == "-h" || os.Args[1] == "--help" {
		fmt.Println("Usage: crt2json https://superhuman.com\n" +
			"       crt2json < superhuman.crt\n" +
			"\n" +
			"Prints out a JSON summary of an SSL certificate.\n" +
			"If an argument is passed it is assumed to be a server to connect to, if not STDIN is scanned for any certificates")

		os.Exit(1)
	}

	u, err := url.Parse(os.Args[1])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	toDial := u.Hostname() + ":443"

	if toDial == ":443" {
		toDial = u.RawPath + ":443"
	}

	c, err := tls.Dial("tcp", toDial, &tls.Config{
		InsecureSkipVerify: true,
	})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	printCert(c.ConnectionState().PeerCertificates[0])

}

func decodeFromStdin() {

	bytes, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for {
		var block *pem.Block
		block, bytes = pem.Decode(bytes)

		if block == nil {
			break
		}

		certs, err := x509.ParseCertificates(block.Bytes)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		for _, c := range certs {
			printCert(c)
		}

	}

}

func printCert(c *x509.Certificate) {

	e := json.NewEncoder(os.Stdout)
	e.SetIndent("", "   ")
	err := e.Encode(map[string]interface{}{
		"not_before":    c.NotBefore,
		"not_after":     c.NotAfter,
		"version":       c.Version,
		"serial_number": fmt.Sprintf("%036x", c.SerialNumber),
		//	"issuer":        c.Issuer,
		"common_name": c.Subject.CommonName,
		"dns_names":   c.DNSNames,
	})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
