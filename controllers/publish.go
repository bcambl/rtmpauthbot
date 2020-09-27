package controllers

import (
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

// perform basic validations on a publisher record
func (p *Publisher) isValid() error {
	var err error
	if len(p.Name) < 1 {
		err = errors.New("missing parameter: name")
		return err
	}
	if len(p.Key) < 1 {
		err = errors.New("missing parameter: key")
		return err
	}
	return nil
}

func (c *Controller) getAllPublisher() ([]Publisher, error) {
	publishers := []Publisher{}
	c.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("PublisherBucket"))
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			var p Publisher
			p.Name = string(k)
			p.Key = string(v)
			publishers = append(publishers, p)
		}
		return nil
	})

	return publishers, nil
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

// OnPublishHandler is the http handler for "/on_publish".
func (c *Controller) OnPublishHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	streamName := r.Form.Get("name")
	streamKey := r.Form.Get("key")
	p, err := c.getPublisher(streamName)
	if err != nil {
		log.Warnf("on_publish unauthorized: %s\n", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if streamKey != p.Key {
		log.Warnf("on_publish unauthorized: %s with 'key': %s\n", p.Name, streamKey)
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
