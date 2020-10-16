package controllers

import (
	"net/http"

	"github.com/bcambl/rtmpauth/config"
	log "github.com/sirupsen/logrus"
	bolt "go.etcd.io/bbolt"
)

// Controller struct to provide the database to all handlers
type Controller struct {
	Config config.Config
	DB     *bolt.DB
}

// IndexHandler is the http handler for "/".
func (c *Controller) IndexHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	w.Write([]byte(`{"handler": "index"}`))

	//var err error

	// publishers, err := c.getAllPublisher()
	// if err != nil {
	// 	log.Error(err)
	// }

	streams, err := c.getStreams()
	if err != nil {
		log.Error(err)
	}

	err = c.updateLiveStatus(streams)
	if err != nil {
		log.Error(err)
	}

	// log.Printf("config token: ", c.Config.TwitchAccessToken)

	// accessToken, err := c.getCachedAccessToken()
	// if err != nil {
	// 	log.Error(err)
	// }
	// log.Printf("config token: ", accessToken)

	// err = validateAccessToken(accessToken)
	// if err != nil {
	// 	log.Error(err)
	// }

	// accessToken, err := c.TwitchAuthToken()
	// if err != nil {
	// 	log.Error(err)
	// }

	// log.Println(accessToken)

}
