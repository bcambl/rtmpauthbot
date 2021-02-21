package controllers

import (
	"errors"
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"
	bolt "go.etcd.io/bbolt"
)

// Publisher struct contains rtmp stream name, stream key, twitch channel name
type Publisher struct {
	Name               string `json:"name"`
	Key                string `json:"key"`
	RTMPLive           string `json:"rtmp_live"`
	TwitchStream       string `json:"twitch_stream"`
	TwitchLive         string `json:"twitch_live"`
	TwitchNotification string `json:"-"`
	StreamInfo         string `json:"-"`
}

// IsValid perform basic validations on a publisher record
func (p *Publisher) IsValid() error {
	var err error
	if len(p.Name) < 1 {
		err = errors.New("missing parameter: name")
		return err
	}
	if len(p.Key) < 1 {
		err = errors.New("missing parameter: key")
		return err
	}
	return nil
}

// IsTwitchLive returns a boolean based on string value of TwitchLive field
func (p *Publisher) IsTwitchLive() bool {
	if p.TwitchLive != "" {
		return true
	}
	return false
}

// FetchPublisher populates the publisher struct from the database
func (c *Controller) FetchPublisher(p *Publisher) error {
	var b []byte
	var err error
	b, err = c.getBucketValue("RTMPLiveBucket", p.Name)
	if err != nil {
		return err
	}
	p.RTMPLive = string(b)
	b, err = c.getBucketValue("TwitchStreamBucket", p.Name)
	if err != nil {
		return err
	}
	p.TwitchStream = string(b)
	b, err = c.getBucketValue("TwitchLiveBucket", p.Name)
	if err != nil {
		return err
	}
	p.TwitchLive = string(b)
	b, err = c.getBucketValue("TwitchNotificationBucket", p.Name)
	if err != nil {
		return err
	}
	p.TwitchNotification = string(b)
	b, err = c.getBucketValue("SteamInfoBucket", p.Name)
	if err != nil {
		return err
	}
	p.StreamInfo = string(b)

	return nil
}

func (c *Controller) getAllPublisher() ([]Publisher, error) {
	var err error
	publishers := []Publisher{}
	c.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("PublisherBucket"))
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			var p Publisher
			p.Name = string(k)
			p.Key = string(v)
			publishers = append(publishers, p)
		}
		return nil
	})

	for i := range publishers {
		p := &publishers[i]
		err = c.FetchPublisher(p)
		if err != nil {
			return nil, err
		}
	}

	return publishers, nil
}

func (c *Controller) getPublisher(name string) (Publisher, error) {
	var keyBytes []byte
	var err error

	p := Publisher{}

	keyBytes, err = c.getBucketValue("PublisherBucket", name)
	if err != nil {
		return p, err
	}
	p.Key = string(keyBytes)

	if len(p.Key) < 1 {
		return p, errors.New("publisher not found")
	}

	err = c.FetchPublisher(&p)
	if err != nil {
		return p, err
	}

	return p, nil
}

func (c *Controller) updatePublisher(p Publisher) error {
	var err error
	c.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("PublisherBucket"))
		err = b.Put([]byte(p.Name), []byte(p.Key))
		return err
	})

	// debug only. live status is managed internally
	// c.DB.Update(func(tx *bolt.Tx) error {
	// 	b := tx.Bucket([]byte("RTMPLiveBucket"))
	// 	err = b.Put([]byte(p.Name), []byte(p.LocalLive))
	// 	return err
	// })

	if p.TwitchStream != "" {
		// only update the stream if a value is provided
		c.DB.Update(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte("TwitchStreamBucket"))
			err = b.Put([]byte(p.Name), []byte(p.TwitchStream))
			return err
		})
	}

	// debug only. live status is managed internally
	// c.DB.Update(func(tx *bolt.Tx) error {
	// 	b := tx.Bucket([]byte("TwitchLiveBucket"))
	// 	err = b.Put([]byte(p.Name), []byte(p.TwitchLive))
	// 	return err
	// })

	return nil
}

func (c *Controller) deletePublisher(name string) error {
	log.Debug("deleting ", name)
	buckets := []string{
		"PublisherBucket",
		"RTMPLiveBucket",
		"TwitchStreamBucket",
		"TwitchLiveBucket",
		"TwitchNotificationBucket",
	}
	for i := range buckets {
		c.DB.Update(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte(buckets[i]))
			err := b.Delete([]byte(name))
			return err
		})
	}
	return nil
}

// OnPublishHandler is the http handler for "/on_publish".
func (c *Controller) OnPublishHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	streamName := r.Form.Get("name")
	streamKey := r.Form.Get("key")
	p, err := c.getPublisher(streamName)
	if err != nil {
		log.Warnf("on_publish unauthorized: %s", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if streamKey != p.Key {
		log.Warnf("on_publish unauthorized: %s with 'key': %s", p.Name, streamKey)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	log.Printf("on_publish authorized: %s", p.Name)

	serverFQDN := c.Config.RTMPServerFQDN
	serverPort := c.Config.RTMPServerPort

	err = c.setBucketValue("RTMPLiveBucket", p.Name, "live")
	if err != nil {
		log.Error("error enabling local live status")
	}

	if c.Config.DiscordEnabled && (serverFQDN != "") {
		content := fmt.Sprintf(":movie_camera: %s started a private stream!\nwatch now: `rtmp://%s:%s/stream/%s`", streamName, serverFQDN, serverPort, streamName)
		err := c.callWebhook(content)
		if err != nil {
			log.Error(err)
		}
	}

	w.WriteHeader(http.StatusCreated)
}

// OnPublishDoneHandler is the http handler for "/on_publish_done".
func (c *Controller) OnPublishDoneHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	streamName := r.Form.Get("name")
	streamKey := r.Form.Get("key")
	p, err := c.getPublisher(streamName)
	if err != nil {
		log.Warnf("on_publish_done unauthorized: %s", p.Name)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if streamKey != p.Key {
		log.Warnf("on_publish_done unauthorized: %s with key: %s", p.Name, p.Key)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	log.Printf("on_publish_done authorized: %s", p.Name)

	err = c.setBucketValue("RTMPLiveBucket", p.Name, "")
	if err != nil {
		log.Error("error disabling local live status")
	}

	if c.Config.DiscordEnabled {
		content := fmt.Sprintf(":checkered_flag:  %s finished streaming.", streamName)
		err := c.callWebhook(content)
		if err != nil {
			log.Error(err)
		}
	}

	w.WriteHeader(http.StatusCreated)
}
