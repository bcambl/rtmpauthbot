package controllers

import (
	"net/http"

	"github.com/bcambl/rtmpauthbot/config"
	bolt "go.etcd.io/bbolt"
)

// Controller struct to provide the database to all handlers
type Controller struct {
	Config *config.Config
	DB     *bolt.DB
}

// IndexHandler is the http handler for "/".
func (c *Controller) IndexHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	w.Write([]byte(`{"handler": "index"}`))

}

func (c *Controller) setBucketValue(bucket, key, value string) error {
	c.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		err := b.Put([]byte(key), []byte(value))
		return err
	})
	return nil
}

func (c *Controller) getBucketValue(bucket, key string) ([]byte, error) {
	var result []byte
	err := c.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		result = b.Get([]byte(key))
		return nil
	})
	return result, err
}
