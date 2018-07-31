package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/tamarakaufler/go-build-servers/url-shorty-proxy/server/handler"
)

func main() {
	r := mux.NewRouter()

	logger := log.New(os.Stdout, "url-shorty-proxy: ", log.LstdFlags)
	ds, err := handler.NewDatastore(logger)
	if err != nil {
		log.Fatalln(err)
	}

	r.HandleFunc("/", ds.NoopHandler).Methods("GET")
	r.HandleFunc("/all", ds.GetAllHandler).Methods("GET")
	r.HandleFunc("/{shorty}", ds.MappingHandler).Methods("GET")
	r.HandleFunc("/{shorty}/{url}", ds.CreateHandler).Methods("POST")
	r.HandleFunc("/{shorty}", ds.DeleteHandler).Methods("DELETE")
	r.HandleFunc("/info/{shorty}", ds.InfoHandler).Methods("GET")

	log.Println("*** Starting proxy on 127.0.0.1:8888")
	log.Fatal(http.ListenAndServe(":8888", r))
}
