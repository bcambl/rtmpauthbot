package controllers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	bolt "go.etcd.io/bbolt"
	"golang.org/x/oauth2/clientcredentials"
	"golang.org/x/oauth2/twitch"
)

const (
	defaultClientID     = "abcd1234"
	defaultClientSecret = "abcd1234"
)

// TwitchStreamsResponse to marshal the json response from /helix/streams/
type TwitchStreamsResponse struct {
	Data []StreamData `json:"data"`
}

// StreamData to marshal the inner data of the TwitchStreamsResponse
type StreamData struct {
	ID          string `json:"id"`
	UserID      string `json:"user_id"`
	UserName    string `json:"user_name"`
	GameID      string `json:"game_id"`
	Type        string `json:"type"`
	Title       string `json:"title"`
	ViewerCount int    `json:"viewer_count"`
	StartedAt   string `json:"started_at"`
}

// retrieve cached twitch access token from database and set in the
// Config struct. This is only called when the token is not set in Config
func (c *Controller) getCachedAccessToken() (string, error) {
	var tokenBytes []byte
	c.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("ConfigBucket"))
		tokenBytes = b.Get([]byte("twitchAccessToken"))
		return nil
	})
	if len(tokenBytes) < 1 {
		return "", errors.New("cached twitch access token not found in db")
	}
	return string(tokenBytes), nil
}

// update the cached access token record in the database
func (c *Controller) updateCachedAccessToken(accessToken string) error {
	var err error
	if accessToken == "" {
		return errors.New("updateCachedAccessToken: no token provided")
	}
	c.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("ConfigBucket"))
		err = b.Put([]byte("twitchAccessToken"), []byte(accessToken))
		return err
	})
	return nil
}

func validateAccessToken(accessToken string) error {
	if accessToken == "" {
		err := errors.New("token validation fail - not set")
		return err
	}
	r, err := http.NewRequest("GET", "https://id.twitch.tv/oauth2/validate", nil)
	if err != nil {
		log.Error(err)
	}
	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("Authorization", "OAuth "+accessToken)

	resp, err := http.DefaultClient.Do(r)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return errors.New("token validation response status code != 200")
	}

	return nil
}

func (c *Controller) getNewAuthToken() error {
	var oauth2Config *clientcredentials.Config

	oauth2Config = &clientcredentials.Config{
		ClientID:     c.Config.TwitchClientID,
		ClientSecret: c.Config.TwitchClientSecret,
		TokenURL:     twitch.Endpoint.TokenURL,
	}

	token, err := oauth2Config.Token(context.Background())
	if err != nil {
		return err
	}

	log.Debug("New Access Token: ", token.AccessToken)
	err = c.updateCachedAccessToken(token.AccessToken)
	if err != nil {
		return err
	}
	return nil

}

func (c *Controller) validateClientCredentials() error {
	if c.Config.TwitchClientID == defaultClientID || c.Config.TwitchClientID == "" {
		err := errors.New("Default twitch client id value detected. Skipping twitch call")
		return err
	}
	if c.Config.TwitchClientSecret == defaultClientSecret || c.Config.TwitchClientSecret == "" {
		err := errors.New("Default twitch client secret value detected. Skipping twitch call")
		return err
	}
	return nil
}

//twitchAuthToken handles the lifecycle of the twitch access token
func (c *Controller) twitchAuthToken() (string, error) {
	var token string
	var err error

	token, err = c.getCachedAccessToken()
	if err != nil {
		log.Debug(err)
	}

	err = validateAccessToken(token)
	if err != nil {
		err = c.getNewAuthToken()
		if err != nil {
			return "", err
		}
	}

	token, err = c.getCachedAccessToken()
	if err != nil {
		return "", err
	}

	return token, nil
}

func streamQueryURL(publishers []Publisher) (string, error) {
	var userQuery string
	for i := range publishers {
		if publishers[i].TwitchStream == "" {
			continue
		}
		if userQuery != "" {
			userQuery = userQuery + "&"
		}
		userQuery = userQuery + fmt.Sprintf("user_login=%s", publishers[i].TwitchStream)
	}

	if userQuery == "" {
		err := errors.New("no streams to query")
		return "", err
	}

	//log.Debug("stream userQuery: ", userQuery)
	return "https://api.twitch.tv/helix/streams/?" + userQuery, nil
}

func (c *Controller) getStreams() ([]StreamData, error) {

	var (
		err         error
		streamQuery string
	)

	err = c.validateClientCredentials()
	if err != nil {
		return nil, err
	}

	accessToken, err := c.twitchAuthToken()
	if err != nil {
		return nil, err
	}

	publishers, err := c.getAllPublisher()
	if err != nil {
		return nil, err
	}

	streamQuery, err = streamQueryURL(publishers)
	if err != nil {
		return nil, err
	}

	r, err := http.NewRequest("GET", streamQuery, nil)
	if err != nil {
		log.Error(err)
	}
	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("client-id", c.Config.TwitchClientID)
	r.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := http.DefaultClient.Do(r)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	streamResponse := TwitchStreamsResponse{}
	err = json.Unmarshal(body, &streamResponse)
	if err != nil {
		return nil, err
	}

	if len(streamResponse.Data) == 0 {
		log.Debug("no twitch streams currently live")
	}
	for i := range streamResponse.Data {
		log.Debug("Live Now:", streamResponse.Data[i].UserName)
	}

	return streamResponse.Data, nil
}

func (c *Controller) updateLiveStatus(streams []StreamData) error {

	var live bool
	publishers, err := c.getAllPublisher()
	if err != nil {
		return err
	}

	// mark previous live streams -> offline
	for i := range publishers {
		live = false
		p := &publishers[i]
		if p.IsTwitchLive() {
			for x := range streams {
				s := streams[x]
				if strings.ToLower(s.UserName) == strings.ToLower(p.TwitchStream) {
					live = true
				}
			}
			if !live {
				c.setTwitchLive(p, "")
				notification := fmt.Sprintf("%s is no longer live on twitch", p.Name)
				c.setTwitchNotification(p, notification)
			}
		}
	}

	// mark live twitch streams -> online
	for x := range streams {
		s := streams[x]
		for i := range publishers {
			p := &publishers[i]
			if p.TwitchStream == "" {
				continue
			}
			if strings.ToLower(s.UserName) == strings.ToLower(p.TwitchStream) {
				if !p.IsTwitchLive() {
					c.setTwitchLive(p, s.Type)
					streamLink := fmt.Sprintf("http://twitch.tv/%s", p.TwitchStream)
					notification := fmt.Sprintf("%s is live on twitch - %s - %s", p.Name, s.Title, streamLink)
					c.setTwitchNotification(p, notification)
				}
			}
		}
	}

	return nil
}

func (c *Controller) processNotifications() error {

	publishers, err := c.getAllPublisher()
	if err != nil {
		return err
	}

	for i := range publishers {
		p := publishers[i]
		if p.TwitchLive == "" && p.TwitchNotification == "" {
			continue
		}
		if p.TwitchLive != "" && p.TwitchNotification == "" {
			continue
		}
		log.Debug("notification: ", p.TwitchNotification)
		if c.Config.DiscordEnabled {
			log.Debug("sending discord notification: %s", p.TwitchNotification)
			err := c.callWebhook(p.TwitchNotification)
			if err != nil {
				return err
			}
		}
		if p.TwitchNotification != "" {
			log.Debugf("resetting notification for %s (%s)", p.Name, p.TwitchStream)
			err = c.setTwitchNotification(&p, "")
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (c *Controller) twitchMain() {
	streams, err := c.getStreams()
	if err != nil {
		log.Debug(err)
		return
	}

	err = c.updateLiveStatus(streams)
	if err != nil {
		log.Error(err)
		return
	}

	err = c.processNotifications()
	if err != nil {
		log.Error(err)
		return
	}
}

// TwitchScheduler launches the twitch stream query & notification background processes
func (c *Controller) TwitchScheduler(ctx context.Context, pollRate time.Duration) {
	ticker := time.NewTicker(pollRate)
	go func() {
		for {
			select {
			case <-ticker.C:
				c.twitchMain()
			case <-ctx.Done():
				ticker.Stop()
				return
			}
		}
	}()
}
