package controllers

import (
	"encoding/json"
	"errors"
	"net/http"

	log "github.com/sirupsen/logrus"
	bolt "go.etcd.io/bbolt"
)

// Publisher struct contains streamer names and stream keys
type Publisher struct {
	Name string `json:"name"`
	Key  string `json:"key"`
}

func (c *Controller) getPublisher(name string) (Publisher, error) {
	var key []byte
	p := Publisher{}
	c.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("PublisherBucket"))
		key = b.Get([]byte(name))
		return nil
	})
	if len(key) < 1 {
		return p, errors.New("publisher not found")
	}
	p.Name = name
	p.Key = string(key)

	return p, nil
}

func (c *Controller) updatePublisher(p Publisher) error {
	var err error
	c.DB.Update(func(tx *bolt.Tx) error {
		log.Debug(p)
		b := tx.Bucket([]byte("PublisherBucket"))
		err = b.Put([]byte(p.Name), []byte(p.Key))
		return err
	})
	return nil
}

func (c *Controller) deletePublisher(name string) error {
	log.Debug("deleting ", name)
	c.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("PublisherBucket"))
		err := b.Delete([]byte(name))
		log.Debug(err)
		return err
	})
	return nil
}

// PublisherhHandler adds a publisher to the database
func (c *Controller) PublisherhHandler(w http.ResponseWriter, r *http.Request) {

	var p Publisher

	if r.Method == "GET" {
		name, ok := r.URL.Query()["name"]
		if !ok || len(name[0]) < 1 {
			log.Println("Url Param 'name' is missing")
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

// OnPublishHandler is the http handler for "/on_publish".
func (c *Controller) OnPublishHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	streamName := r.Form.Get("name")
	streamKey := r.Form.Get("key")
	p, err := c.getPublisher(streamName)
	if err != nil {
		log.Warnf("on_publish unauthorized: %s with key: %s\n", p.Name, p.Key)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if streamKey != p.Key {
		log.Warnf("on_publish unauthorized: %s with key: %s\n", p.Name, p.Key)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	log.Printf("on_publish authorized: %s with key: %s\n", p.Name, p.Key)
	w.WriteHeader(http.StatusCreated)
}

// OnPublishDoneHandler is the http handler for "/on_publish_done".
func (c *Controller) OnPublishDoneHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	streamName := r.Form.Get("name")
	streamKey := r.Form.Get("key")
	p, err := c.getPublisher(streamName)
	if err != nil {
		log.Warnf("on_publish_done unauthorized: %s with key: %s\n", p.Name, p.Key)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if streamKey != p.Key {
		log.Warnf("on_publish_done unauthorized: %s with key: %s\n", p.Name, p.Key)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	log.Printf("on_publish_done authorized: %s with key: %s\n", p.Name, p.Key)
	w.WriteHeader(http.StatusCreated)
}
