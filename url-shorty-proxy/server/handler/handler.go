package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"

	"github.com/gorilla/mux"
	//"github.com/tamarakaufler/go-build-servers/url-shorty-proxy/server/datastore/psql"
	"github.com/tamarakaufler/go-build-servers/url-shorty-proxy/server/datastore/mongo"
)

type Datastore struct {
	Store  *mongo.Store
	logger *log.Logger
}

type Response struct {
	Status, Content string
}

var httpRegex = regexp.MustCompile("^http(s)?://")

func NewDatastore(logger *log.Logger) (*Datastore, error) {
	st, err := mongo.NewMGOStore(logger)
	if err != nil {
		return nil, err
	}
	return &Datastore{
		Store:  st,
		logger: logger,
	}, nil
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
		d.logger.Printf("ERROR in GetByAbb: %v\n", err)
		writeResponse(d.logger, w, http.StatusNotFound, "Record not found")
		return
	}

	d.logger.Printf("shorty: %s => url: %s\n", shorty, abbr.Url)

	if !httpRegex.MatchString(abbr.Url) {
		abbr.Url = "http://" + abbr.Url
	}

	//log.Printf("abbr: %#v\n", abbr)
	d.logger.Printf("shorty: %s => url: %s\n", shorty, abbr.Url)

	http.Redirect(w, r, abbr.Url, http.StatusMovedPermanently)
}

func (d *Datastore) NoopHandler(w http.ResponseWriter, r *http.Request) {
	d.logger.Println("In NoopHandler")
	writeResponse(d.logger, w, http.StatusOK, "Nothing to do here")
}

func (d *Datastore) CreateHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	vars := mux.Vars(r)
	shorty := vars["shorty"]
	url := vars["url"]

	newShorty := mongo.MGOShorty{
		Shorty: shorty,
		Url:    url,
	}

	err := d.Store.Conn.Create(newShorty)
	if err != nil {
		d.logger.Printf("ERROR in Create: %v\n", err)
		writeResponse(d.logger, w, http.StatusOK, fmt.Sprint(err))
		return
	}
	writeResponse(d.logger, w, http.StatusOK, fmt.Sprintf("%s ... %s", shorty, url))
}

func (d *Datastore) DeleteHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodDelete {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	vars := mux.Vars(r)
	shorty := vars["shorty"]

	err := d.Store.Conn.Delete(shorty)
	if err != nil {
		d.logger.Printf("ERROR in Delete: %v\n", err)
		writeResponse(d.logger, w, http.StatusNotFound, fmt.Sprint(err))
		return
	}
	writeResponse(d.logger, w, http.StatusOK, fmt.Sprintf("%s deleted", shorty))
}

func (d *Datastore) InfoHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	vars := mux.Vars(r)
	shorty := vars["shorty"]

	abbr, err := d.Store.Conn.GetByAbbr(shorty)
	if err != nil {
		d.logger.Printf("ERROR in Info: %v\n", err)
		writeResponse(d.logger, w, http.StatusNotFound, fmt.Sprint(err))
		return
	}
	writeResponse(d.logger, w, http.StatusOK, fmt.Sprintf("%s ... %s", shorty, abbr.Url))
}

func (d *Datastore) GetAllHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	abbrs, err := d.Store.Conn.GetAll()
	if err != nil {
		d.logger.Printf("ERROR in Info: %v\n", err)
		writeResponse(d.logger, w, http.StatusNotFound, fmt.Sprint(err))
		return
	}
	writeResponse(d.logger, w, http.StatusOK, fmt.Sprintf("%+v", abbrs))
}

// Helper functions and methods -----------------------------------------------------------

func writeResponse(logger *log.Logger, w http.ResponseWriter, status int, content string) {
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	resp := Response{
		Status:  http.StatusText(status),
		Content: content,
	}

	err := json.NewEncoder(w).Encode(resp)
	if err != nil {
		logger.Printf("could not encode response: %v", err)
	}
}
