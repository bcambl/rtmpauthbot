package controllers

import (
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"
	bolt "go.etcd.io/bbolt"
)

// Publisher struct contains streamer names and stream keys
type Publisher struct {
	Name string
	Key  string
}

func (c *Controller) getPublisher(name string) {
	c.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("PublisherBucket"))
		v := b.Get([]byte(name))
		log.Printf("The key for %s is: %s\n", name, v)
		return nil
	})
}

func (c *Controller) updatePublisher(p Publisher) {
	//log.Debug(p)
	c.DB.Update(func(tx *bolt.Tx) error {
		log.Debug(p)
		b := tx.Bucket([]byte("PublisherBucket"))
		err := b.Put([]byte(p.Name), []byte(p.Key))
		return err
	})
}

func (c *Controller) deletePublisher(name string) {
	c.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("PublisherBucket"))
		b.Delete([]byte(name))
		log.Printf("publisher deleted: %s\n", name)
		return nil
	})
}

// PublisherhHandler adds a publisher to the database
func (c *Controller) PublisherhHandler(w http.ResponseWriter, r *http.Request) {

	var p Publisher

	if r.Method == "GET" {
		name, err := r.URL.Query()["name"]
		if !err || len(name[0]) < 1 {
			log.Println("Url Param 'name' is missing")
			return
		}
		// URL.Query() returns a []string
		n := name[0]
		c.getPublisher(n)
	}
	if r.Method == "POST" {
		err := json.NewDecoder(r.Body).Decode(&p)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		c.updatePublisher(p)
	}
	if r.Method == "DELETE" {
		err := json.NewDecoder(r.Body).Decode(&p)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		c.deletePublisher(p.Name)
	}

}

// GetPublisher adds a publisher to the database
func (c *Controller) GetPublisher(w http.ResponseWriter, r *http.Request) {

}

// OnPublishHandler is the http handler for "/on_publish".
func (c *Controller) OnPublishHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	streamName := r.Form.Get("name")
	streamKey := r.Form.Get("key")
	log.Printf("publishing %s with key: %s\n", streamName, streamKey)
}

// OnPublishDoneHandler is the http handler for "/on_publish_done".
func (c *Controller) OnPublishDoneHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	streamName := r.Form.Get("name")
	streamKey := r.Form.Get("key")
	log.Printf("publishing %s with key: %s\n", streamName, streamKey)
}
