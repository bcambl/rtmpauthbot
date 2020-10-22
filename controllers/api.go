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

	// API GET REQUESTS
	if r.Method == "GET" {
		name, ok := r.URL.Query()["name"]
		if !ok || len(name[0]) < 1 {
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
			log.Info("listing all publishers")
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
		log.Infof("listing publisher %s", p.Name)
		w.Write(content)
		return
	}

	// API POST REQUESTS
	if r.Method == "POST" {
		err := json.NewDecoder(r.Body).Decode(&p)
		if err != nil {
			log.Debug(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		err = p.IsValid()
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
		log.Infof("updating publisher %s", p.Name)
		w.WriteHeader(http.StatusCreated)
		return
	}

	// API DELETE REQUESTS
	if r.Method == "DELETE" {
		err := json.NewDecoder(r.Body).Decode(&p)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		_, err = c.getPublisher(p.Name)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		err = c.deletePublisher(p.Name)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		log.Infof("deleting publisher %s", p.Name)
		w.WriteHeader(http.StatusNoContent)
		return
	}
	w.WriteHeader(http.StatusNotImplemented)
	return
}
