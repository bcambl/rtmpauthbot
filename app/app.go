package app

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/bcambl/rtmpauth/config"
	"github.com/bcambl/rtmpauth/controllers"
	log "github.com/sirupsen/logrus"
	bolt "go.etcd.io/bbolt"
)

func init() {
	debugFlag := flag.Bool("debug", false, "enable debug logging")
	flag.Parse()

	logLevel := log.InfoLevel
	if *debugFlag {
		logLevel = log.DebugLevel
	}
	log.SetOutput(os.Stdout)
	log.SetLevel(logLevel)

	// Initialize the database
	db, err := bolt.Open(config.DatabasePath(), 0700, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	bucketList := []string{
		"ConfigBucket",             // General configuration & caching
		"PublisherBucket",          // Local publishers -> rtmp stream keys
		"RTMPLiveBucket",           // Local publishers -> rtmp live stream status
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

	var conf config.Config
	err := conf.ParseEnv()
	if err != nil {
		log.Fatal(err)
	}

	db, err := bolt.Open(config.DatabasePath(), 0700, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	c := controllers.Controller{Config: &conf, DB: db}

	// Start Twitch polling scheduler if integration is enabled
	if c.Config.TwitchEnabled {
		log.Infof("twitch integration enabled")
		log.Infof("starting twitch scheduler (poll rate: %s)", c.Config.TwitchPollRate.String())
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		c.TwitchScheduler(ctx, c.Config.TwitchPollRate)
	} else {
		log.Infof("twitch integration disabled")
	}

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
	if conf.AuthServerIP == "" {
		conf.AuthServerIP = "127.0.0.1"
	}
	if conf.AuthServerPort == "" {
		conf.AuthServerPort = "9090"
	}
	listenAddress := fmt.Sprintf("%s:%s", conf.AuthServerIP, conf.AuthServerPort)

	// Serve
	log.Infof("starting rtmpauth server on %s", listenAddress)
	err = http.ListenAndServe(listenAddress, nil)
	if err != nil {
		log.Fatal(err)
	}
}
