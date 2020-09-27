package controllers

import (
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"
)

// PublisherAPIHandler manages publisher database records
func (c *Controller) PublisherAPIHandler(w http.ResponseWriter, r *http.Request) {

	var p Publisher

	w.Header().Add("Content-Type", "application/json")

	if r.Method == "GET" {
		name, ok := r.URL.Query()["name"]
		if !ok || len(name[0]) < 1 {
			// err := errors.New("Missing Parameter: 'name'")
			// http.Error(w, err.Error(), http.StatusBadRequest)
			publishers, err := c.getAllPublisher()
			if err != nil {
				log.Debug(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			content, err := json.Marshal(publishers)
			if err != nil {
				log.Debug(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			w.Write(content)
			return
		}
		// URL.Query() returns a []string
		n := name[0]
		p, err := c.getPublisher(n)
		if err != nil {
			log.Debug(err)
			w.WriteHeader(http.StatusNotFound)
			return
		}
		content, err := json.Marshal(p)
		if err != nil {
			log.Debug(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write(content)
		return
	}
	if r.Method == "POST" {
		err := json.NewDecoder(r.Body).Decode(&p)
		if err != nil {
			log.Debug(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		err = p.isValid()
		if err != nil {
			log.Debug(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		err = c.updatePublisher(p)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
		return
	}
	if r.Method == "DELETE" {
		err := json.NewDecoder(r.Body).Decode(&p)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		err = c.deletePublisher(p.Name)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
		return
	}
	w.WriteHeader(http.StatusNotImplemented)
}
