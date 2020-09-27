package app

import (
	"fmt"
	"net/http"
	"os"

	"github.com/bcambl/rtmpauth/config"
	"github.com/bcambl/rtmpauth/controllers"
	log "github.com/sirupsen/logrus"
	bolt "go.etcd.io/bbolt"
)

func init() {
	//log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)

	// Initialize the database
	db, err := bolt.Open("rtmpauth.db", 0700, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Initialize the datatabase publisher bucket
	db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("PublisherBucket"))
		if err != nil {
			return fmt.Errorf("error creating bucket: %s", err)
		}
		return nil
	})

}

// Run performs setup and starts the server.
func Run() {

	var config config.Config
	err := config.ParseEnv()
	if err != nil {
		log.Fatal(err)
	}

	db, err := bolt.Open("rtmpauth.db", 0700, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	c := controllers.Controller{Config: config, DB: db}

	// Root Handler
	http.HandleFunc("/", c.IndexHandler)

	// Play Handlers
	http.HandleFunc("/on_play", c.OnPlayHandler)
	http.HandleFunc("/on_play_done", c.OnPlayDoneHandler)

	// Publish Handlers
	http.HandleFunc("/on_publish", c.OnPublishHandler)
	http.HandleFunc("/on_publish_done", c.OnPublishDoneHandler)

	// API Endpoints
	http.HandleFunc("/api/publisher", c.PublisherAPIHandler)

	// Serve
	log.Info("starting rtmpauth server")
	http.ListenAndServe("127.0.0.1:9090", nil)

}
