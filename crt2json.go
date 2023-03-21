package main

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"flag"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"strings"
)

func main() {

	var sni string

	flag.StringVar(&sni, "sni", "", "server name indication to instruct the server to return the correct certificate")

	flag.Usage = func() {
		fmt.Println(`Usage: crt2json [ -sni HOSTNAME ] [ SERVER | FILENAME ]

Prints out a JSON summary of an SSL certificate.

If the argument exists on disk then the file is assumed to be a certificate file, otherwise
it is interpreted as a URL or a hostname to connect to.`)

		flag.PrintDefaults()
	}

	flag.Parse()
	// nonsense
	if len(flag.Args()) == 3 && flag.Args()[1] == "-sni" {
		flag.Set("sni", flag.Args()[2])
	} else if len(flag.Args()) != 1 {
		flag.Usage()
		return
	}

	arg := flag.Args()[0]
	if _, err := os.Stat(arg); err == nil {
		decodeFile(arg)
		return
	}

	toDial := arg
	if strings.Contains(arg, "//") {
		u, err := url.Parse(arg)

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		toDial = u.Hostname() + ":443"

		if toDial == ":443" {
			toDial = u.RawPath + ":443"
		}

	} else if !strings.Contains(arg, ":") {
		toDial += ":443"
	}

	c, err := tls.Dial("tcp", toDial, &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         sni,
	})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	printCert(c.ConnectionState().PeerCertificates[0])

}

func decodeFile(name string) {
	bytes, err := ioutil.ReadFile(name)
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
		"not_before":          c.NotBefore,
		"not_after":           c.NotAfter,
		"version":             c.Version,
		"serial_number":       fmt.Sprintf("%036x", c.SerialNumber),
		"issuer_organization": c.Issuer.Organization,
		"common_name":         c.Subject.CommonName,
		"dns_names":           c.DNSNames,
	})

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
