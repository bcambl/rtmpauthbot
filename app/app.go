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

	bucketList := []string{
		"ConfigBucket",             // General configuration & caching
		"PublisherBucket",          // Local publishers -> rtmp stream keys
		"TwitchStreamBucket",       // Local publishers -> twitch stream names
		"TwitchLiveBucket",         // Local publishers -> twitch live stream status
		"TwitchNotificationBucket", // Local publishers -> twitch notification state
	}

	for b := range bucketList {
		log.Debug("db: ensuring bucket exists: ", bucketList[b])
		db.Update(func(tx *bolt.Tx) error {
			_, err := tx.CreateBucketIfNotExists([]byte(bucketList[b]))
			if err != nil {
				return fmt.Errorf("error creating bucket: %s", err)
			}
			return nil
		})
	}
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

	// if the listen address env variables are not set, set to sane default
	if config.ServerIP == "" {
		config.ServerIP = "127.0.0.1"
	}
	if config.ServerPort == "" {
		config.ServerPort = "9090"
	}
	listenAddress := fmt.Sprintf("%s:%s", config.ServerIP, config.ServerPort)

	// Serve
	log.Infof("starting rtmpauth server on %s", listenAddress)
	http.ListenAndServe(listenAddress, nil)

}
