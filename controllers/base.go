package controllers

import (
	"net/http"

	bolt "go.etcd.io/bbolt"
)

// Controller struct to provide the database to all handlers
type Controller struct {
	DB *bolt.DB
}

// IndexHandler is the http handler for "/".
func (c *Controller) IndexHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	w.Write([]byte(`{"handler": "index"}`))
}
