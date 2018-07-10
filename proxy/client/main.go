package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	neturl "net/url"
)

var insecure bool
var localCert, server, url string

// points to self-signed certificate
func init() {
	flag.BoolVar(&insecure, "insecure", false, "Accept/Ignore all server SSL certificates")
	flag.StringVar(&localCert, "cert", "/home/tamara/programming/go/src/github.com/tamarakaufler/go-build-servers/certs/localhost.pem", "Local CA certificate")
	flag.StringVar(&server, "server-host", "localhost:8888", "Server hostname")
	flag.StringVar(&url, "url", "https://docs.docker.com", "Url to fetch")
}

func main() {

	flag.Parse()

	rootCAs, err := getCertPool(localCert)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	// Let the client trust the extended cert pool
	tlsConfig := &tls.Config{
		InsecureSkipVerify: insecure,
		RootCAs:            rootCAs,
	}

	proto := "https://"
	if insecure == true {
		proto = "http://"
	}
	server = proto + server

	serverURL, err := neturl.Parse(server)
	if err != nil {
		panic(err)
	}
	tr := &http.Transport{
		TLSClientConfig: tlsConfig,
		Proxy:           http.ProxyURL(serverURL),
		// Disable HTTP/2
		TLSNextProto: make(map[string]func(authority string, c *tls.Conn) http.RoundTripper),
	}
	client := &http.Client{Transport: tr}
	resp, err := client.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	dump, err := httputil.DumpResponse(resp, true)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%q", dump)
}
