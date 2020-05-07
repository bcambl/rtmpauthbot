package controllers

import (
	"net/http"

	log "github.com/sirupsen/logrus"
)

// OnPlayHandler is the http handler for "/on_play".
func (c *Controller) OnPlayHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	streamName := r.Form.Get("name")
	log.Printf("playing %s\n", streamName)
}

// OnPlayDoneHandler is the http handler for "/on_play_done".
func (c *Controller) OnPlayDoneHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	streamName := r.Form.Get("name")
	log.Printf("playing  %s\n", streamName)
}
