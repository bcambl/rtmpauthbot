package controllers

import (
	"net/http"
)

// OnPlayHandler is the http handler for "/on_play".
func (c *Controller) OnPlayHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	w.Write([]byte(`{"handler": "on_play"}`))
}

// OnPlayDoneHandler is the http handler for "/on_play_done".
func (c *Controller) OnPlayDoneHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	w.Write([]byte(`{"handler": "on_play_done"}`))
}
