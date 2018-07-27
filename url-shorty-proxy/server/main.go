package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/tamarakaufler/go-build-servers/url-shorty-proxy/server/handler"
)

func main() {
	r := mux.NewRouter()

	ds, err := handler.NewDatastore()
	if err != nil {
		log.Fatalln(err)
	}

	r.HandleFunc("/", handler.NoopHandler)
	r.HandleFunc("/{shorty}", ds.MappingHandler)

	log.Println("*** Starting proxy on 127.0.0.1:8888")
	log.Fatal(http.ListenAndServe(":8888", r))
}