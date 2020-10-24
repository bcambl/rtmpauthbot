package controllers

import (
	"encoding/json"
	"io/ioutil"
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
				log.Debug("error retrieving all publishers: ", err)
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
			log.Debugf("error retrieving publisher '%s': %s\n", p.Name, err)
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

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Debug("error reading POST body: ", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		err = json.Unmarshal(body, &p)
		if err != nil {
			log.Debug("error unmarshaling body json: ", err)
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
			log.Debugf("error updating publisher '%s': %s\n", p.Name, err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		log.Infof("publisher updated: %s", p.Name)
		w.WriteHeader(http.StatusCreated)
		return
	}

	// API DELETE REQUESTS
	if r.Method == "DELETE" {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Debug("error reading DELETE body: ", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		err = json.Unmarshal(body, &p)
		if err != nil {
			log.Debug("error unmarshaling body json: ", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		_, err = c.getPublisher(p.Name)
		if err != nil {
			log.Debugf("error retrieving publisher for deletion '%s': %s\n", p.Name, err)
			w.WriteHeader(http.StatusNotFound)
			return
		}
		err = c.deletePublisher(p.Name)
		if err != nil {
			log.Debugf("error deleting publisher '%s': %s\n", p.Name, err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		log.Infof("publisher deleted: %s", p.Name)
		w.WriteHeader(http.StatusNoContent)
		return
	}
	log.Debug(http.StatusNotImplemented)
	w.WriteHeader(http.StatusNotImplemented)
	return
}
