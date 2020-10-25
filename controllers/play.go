package controllers

import (
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"
)

// OnPlayHandler is the http handler for "/on_play".
func (c *Controller) OnPlayHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	streamName := r.Form.Get("name")
	p, err := c.getPublisher(streamName)
	if err != nil {
		log.Warnf("on_play: stream not found: %s\n", streamName)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	log.Printf("on_play: %s\n", p.Name)

	if c.Config.DiscordEnabled {
		content := fmt.Sprintf(":chart_with_upwards_trend: %s gained a viewer.", streamName)
		err := c.callWebhook(content)
		if err != nil {
			log.Error(err)
		}
	}

	w.WriteHeader(http.StatusCreated)
}

// OnPlayDoneHandler is the http handler for "/on_play_done".
func (c *Controller) OnPlayDoneHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	streamName := r.Form.Get("name")
	p, err := c.getPublisher(streamName)
	if err != nil {
		log.Warnf("on_play_done: stream not found: %s\n", streamName)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	log.Printf("on_play_done: %s\n", p.Name)

	if c.Config.DiscordEnabled {
		content := fmt.Sprintf(":chart_with_downwards_trend: %s lost a viewer.", streamName)
		err := c.callWebhook(content)
		if err != nil {
			log.Error(err)
		}
	}

	w.WriteHeader(http.StatusCreated)
}
