package controllers

import "net/http"

type Controller struct {
}

// IndexHandler is the http handler for "/".
func (c *Controller) IndexHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	w.Write([]byte(`{"handler": "index"}`))
}
