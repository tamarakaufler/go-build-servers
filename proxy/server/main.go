package main

import (
	"crypto/tls"
	"flag"
	"io"
	"log"
	"net"
	"net/http"
	"time"
)

var pemPath string
var keyPath string
var proto string

func init() {
	flag.StringVar(&pemPath, "pem", "server.pem", "path to pem file")
	flag.StringVar(&keyPath, "key", "server.key", "path to key file")
	flag.StringVar(&proto, "proto", "https", "Proxy protocol (http or https)")
}

func main() {
	flag.Parse()

	if proto != "http" && proto != "https" {
		log.Fatal("Protocol must be either http or https")
	}

	server := &http.Server{
		Addr: ":8888",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodConnect {
				log.Println("Want a tunnel ...")
				tunnelHandler(w, r)
			} else {
				log.Println("Want a direct connection ...")
				directHandler(w, r)
			}
		}),
		// Disable HTTP/2.
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler)),
	}

	// start server
	switch proto {
	case "http":
		log.Println("Starting HTTP proxy server ...")
		log.Fatal(server.ListenAndServe())
	default:
		log.Println("Starting HTTPS proxy server ...")
		log.Fatal(server.ListenAndServeTLS(pemPath, keyPath))
	}
}

//------------------------------------------------------------------
func directHandler(w http.ResponseWriter, req *http.Request) {
	log.Println("Handling request ditectly ...")

	resp, err := http.DefaultTransport.RoundTrip(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	defer resp.Body.Close()
	copyHeader(w.Header(), resp.Header)
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

func copyHeader(dst, src http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}

//------------------------------------------------------------------
func tunnelHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Creating a tunnel ...")

	// set up destination connection
	destConn, err := net.DialTimeout("tcp", r.Host, 10*time.Second)
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}

	// set up client connection
	w.WriteHeader(http.StatusOK)
	hijacker, ok := w.(http.Hijacker)
	if !ok {
		http.Error(w, "Hijacking not supported", http.StatusInternalServerError)
		return
	}
	clientConn, _, err := hijacker.Hijack()
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
	}

	go transfer(destConn, clientConn)
	go transfer(clientConn, destConn)
}

func transfer(destination io.WriteCloser, source io.ReadCloser) {
	defer destination.Close()
	defer source.Close()
	io.Copy(destination, source)
}
