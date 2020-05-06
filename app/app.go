package app

import (
	"net/http"

	"github.com/bcambl/rtmpauth/controllers"
	log "github.com/sirupsen/logrus"
)

// Run performs setup and starts the server.
func Run() {
	c := controllers.Controller{}
	// Load handlers
	http.HandleFunc("/", c.IndexHandler)

	// Play Handlers
	http.HandleFunc("/on_play", c.OnPlayHandler)
	http.HandleFunc("/on_play_done", c.OnPlayDoneHandler)

	// Publish Handlers
	http.HandleFunc("/on_publish", c.OnPublishHandler)
	http.HandleFunc("/on_publish_done", c.OnPublishDoneHandler)

	// Serve
	log.Info("starting rtmpauth server")
	http.ListenAndServe("127.0.0.1:9090", nil)

}
