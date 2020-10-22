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
	TwitchStream       string `json:"twitch_stream"`
	TwitchLive         string `json:"twitch_live"`
	TwitchNotification string `json:"twitch_notification"`
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

func (c *Controller) setTwitchLive(p *Publisher, status string) error {
	c.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("TwitchLiveBucket"))
		err := b.Put([]byte(p.Name), []byte(status))
		return err
	})
	return nil
}

func (c *Controller) setTwitchNotification(p *Publisher, notification string) error {
	c.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("TwitchNotificationBucket"))
		err := b.Put([]byte(p.Name), []byte(notification))
		return err
	})
	return nil
}

func (c *Controller) getAllPublisher() ([]Publisher, error) {
	var stream, live, notification []byte
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
		c.DB.View(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte("TwitchStreamBucket"))
			stream = b.Get([]byte(p.Name))
			return nil
		})
		c.DB.View(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte("TwitchLiveBucket"))
			live = b.Get([]byte(p.Name))
			return nil
		})
		c.DB.View(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte("TwitchNotificationBucket"))
			notification = b.Get([]byte(p.Name))
			return nil
		})
		p.TwitchStream = string(stream)
		p.TwitchLive = string(live)
		p.TwitchNotification = string(notification)
	}

	return publishers, nil
}

func (c *Controller) getPublisher(name string) (Publisher, error) {
	var key, stream, live, notification []byte
	p := Publisher{}
	c.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("PublisherBucket"))
		key = b.Get([]byte(name))
		return nil
	})
	if len(key) < 1 {
		return p, errors.New("publisher not found")
	}

	c.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("TwitchStreamBucket"))
		stream = b.Get([]byte(name))
		return nil
	})
	c.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("TwitchLiveBucket"))
		live = b.Get([]byte(name))
		return nil
	})
	c.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("TwitchNotificationBucket"))
		notification = b.Get([]byte(name))
		return nil
	})

	p.Name = name
	p.Key = string(key)
	p.TwitchLive = string(live)
	p.TwitchStream = string(stream)
	p.TwitchNotification = string(notification)

	return p, nil
}

func (c *Controller) updatePublisher(p Publisher) error {
	var err error
	c.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("PublisherBucket"))
		err = b.Put([]byte(p.Name), []byte(p.Key))
		return err
	})

	if p.TwitchStream != "" {
		c.DB.Update(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte("TwitchStreamBucket"))
			err = b.Put([]byte(p.Name), []byte(p.TwitchStream))
			return err
		})
	}

	return nil
}

func (c *Controller) deletePublisher(name string) error {
	log.Debug("deleting ", name)
	buckets := []string{
		"PublisherBucket",
		"TwitchStreamBucket",
		"TwitchLiveBucket",
		"TwitchNotificationBucket",
	}
	for i := range buckets {
		c.DB.Update(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte(buckets[i]))
			err := b.Delete([]byte(name))
			log.Debug(err)
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
		log.Warnf("on_publish unauthorized: %s\n", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if streamKey != p.Key {
		log.Warnf("on_publish unauthorized: %s with 'key': %s\n", p.Name, streamKey)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	log.Printf("on_publish authorized: %s with key: %s\n", p.Name, p.Key)

	serverFQDN := c.Config.ServerFQDN

	if c.Config.DiscordWebhookEnabled && (serverFQDN != "") {
		content := fmt.Sprintf(":movie_camera: %s started streaming. vlc: `rtmp://%s/stream/%s`", streamName, serverFQDN, streamName)
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
		log.Warnf("on_publish_done unauthorized: %s with key: %s\n", p.Name, p.Key)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if streamKey != p.Key {
		log.Warnf("on_publish_done unauthorized: %s with key: %s\n", p.Name, p.Key)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	log.Printf("on_publish_done authorized: %s with key: %s\n", p.Name, p.Key)

	if c.Config.DiscordWebhookEnabled {
		content := fmt.Sprintf(":black_medium_small_square: %s stopped streaming.", streamName)
		err := c.callWebhook(content)
		if err != nil {
			log.Error(err)
		}
	}

	w.WriteHeader(http.StatusCreated)
}
