package app

import (
	"fmt"
	"net/http"

	"github.com/bcambl/rtmpauth/controllers"
	log "github.com/sirupsen/logrus"
	bolt "go.etcd.io/bbolt"
)

// Run performs setup and starts the server.
func Run() {

	db, err := bolt.Open("rtmpauth.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("PublisherBucket"))
		if err != nil {
			return fmt.Errorf("error creating bucket: %s", err)
		}
		return nil
	})

	c := controllers.Controller{DB: db}
	// Load handlers
	http.HandleFunc("/", c.IndexHandler)

	// Play Handlers
	http.HandleFunc("/on_play", c.OnPlayHandler)
	http.HandleFunc("/on_play_done", c.OnPlayDoneHandler)

	// Publish Handlers
	http.HandleFunc("/publisher", c.PublisherhHandler)
	http.HandleFunc("/on_publish", c.OnPublishHandler)
	http.HandleFunc("/on_publish_done", c.OnPublishDoneHandler)

	// Serve
	log.Info("starting rtmpauth server")
	http.ListenAndServe("127.0.0.1:9090", nil)

}
