package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"

	"github.com/gorilla/mux"
	"github.com/tamarakaufler/go-build-servers/url-shorty-proxy/server/datastore"
)

type Datastore struct {
	Store *datastore.Store
}

type Response struct {
	Status, Content string
}

var httpRegex = regexp.MustCompile("^http(s)?://")

func NewDatastore() (*Datastore, error) {
	st, err := datastore.NewPSQLStore()
	if err != nil {
		return nil, err
	}
	return &Datastore{Store: st}, nil
}

func (d *Datastore) MappingHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	vars := mux.Vars(r)
	shorty := vars["shorty"]

	abbr, err := d.Store.Conn.GetByAbbr(shorty)
	if err != nil {
		log.Printf("ERROR in GetByAbb: %v\n", err)
		writeResponse(w, http.StatusNotFound, "Record not found")
		return
	}

	log.Printf("shorty: %s => url: %s\n", shorty, abbr.Url)

	if !httpRegex.MatchString(abbr.Url) {
		abbr.Url = "http://" + abbr.Url
	}

	//log.Printf("abbr: %#v\n", abbr)
	log.Printf("shorty: %s => url: %s\n", shorty, abbr.Url)

	http.Redirect(w, r, abbr.Url, http.StatusMovedPermanently)
}

func NoopHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("In NoopHandler")
	writeResponse(w, http.StatusOK, "Nothing to do here")
}

func (d *Datastore) CreateHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	vars := mux.Vars(r)
	shorty := vars["shorty"]
	url := vars["url"]

	newShorty := datastore.Shorty{
		Shorty: shorty,
		Url:    url,
	}

	err := d.Store.Conn.Create(newShorty)
	if err != nil {
		log.Printf("ERROR in Create: %v\n", err)
		writeResponse(w, http.StatusOK, fmt.Sprint(err))
		return
	}
	writeResponse(w, http.StatusOK, fmt.Sprintf("%s ... %s", shorty, url))
}

func writeResponse(w http.ResponseWriter, status int, content string) {
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	resp := Response{
		Status:  http.StatusText(status),
		Content: content,
	}

	err := json.NewEncoder(w).Encode(resp)
	if err != nil {
		log.Printf("could not encode response: %v", err)
	}
}
