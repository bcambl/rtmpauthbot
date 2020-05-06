package controllers

import (
	"net/http"
)

// OnPublishHandler is the http handler for "/on_publish".
func (c *Controller) OnPublishHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	w.Write([]byte(`{"handler": "on_publish"}`))
}

// OnPublishDoneHandler is the http handler for "/on_publish_done".
func (c *Controller) OnPublishDoneHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	w.Write([]byte(`{"handler": "on_publish_done"}`))
}
