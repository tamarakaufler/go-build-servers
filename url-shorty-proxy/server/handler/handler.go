package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/tamarakaufler/go-build-servers/url-shorty-proxy/server/datastore"
)

type Datastore struct {
	Store *datastore.Store
}

type Response struct {
	Status, Content string
}

func NewDatastore() (*Datastore, error) {
	st, err := datastore.NewPSQLStore()
	if err != nil {
		return nil, err
	}
	return &Datastore{Store: st}, nil
}

func (d *Datastore) MappingHandler(w http.ResponseWriter, r *http.Request) {
	log.Println(">>> In MappingHandler")

	fmt.Printf("%#v\n", r.URL)

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	vars := mux.Vars(r)
	shorty := vars["shorty"]
	log.Printf("shorty: %s\n", shorty)

	abbr, err := d.Store.Conn.GetByAbbr(shorty)
	if err != nil {
		log.Printf("GetByAbb: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.Printf("abbr: %v\n", abbr)
	log.Printf("shorty: %s => url: %s\n", shorty, abbr.Url)

	// if shorty == "ggl.c" {
	// 	url = "https://www.google.com"
	// } else {
	// 	url = "https://amazon.co.uk"
	// }

	http.Redirect(w, r, abbr.Url, http.StatusMovedPermanently)
}

func NoopHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("In NoopHandler")

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	resp := Response{
		Status:  http.StatusText(http.StatusOK),
		Content: "Nothing to do here",
	}

	log.Printf("%v\n", resp)

	err := json.NewEncoder(w).Encode(resp)
	if err != nil {
		log.Printf("could not encode response: %v", err)
	}
}
