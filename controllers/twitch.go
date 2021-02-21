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

// TwitchGamesResponse to marshal the json response from /helix/games/
type TwitchGamesResponse struct {
	Data []GameData `json:"data"`
}

// GameData to marshal the inner data of
type GameData struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	BoxArtURL string `json:"box_art_url"`
}

// retrieve cached twitch access token from database and set in the
// Config struct. This is only called when the token is not set in Config
func (c *Controller) getCachedAccessToken() (string, error) {
	var tokenBytes []byte
	var err error
	tokenBytes, err = c.getBucketValue("ConfigBucket", "twitchAccessToken")
	if err != nil {
		return "", err
	}
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
	err = c.setBucketValue("ConfigBucket", "twitchAccessToken", accessToken)
	if err != nil {
		return err
	}
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

func (c *Controller) getGame(gameID string) (GameData, error) {

	var (
		err        error
		gamesQuery string
		g          GameData
	)

	err = c.validateClientCredentials()
	if err != nil {
		return g, err
	}

	accessToken, err := c.twitchAuthToken()
	if err != nil {
		return g, err
	}

	gamesQuery = fmt.Sprintf("https://api.twitch.tv/helix/games?id=%s", gameID)

	r, err := http.NewRequest("GET", gamesQuery, nil)
	if err != nil {
		log.Error(err)
	}
	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("client-id", c.Config.TwitchClientID)
	r.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := http.DefaultClient.Do(r)
	if err != nil {
		return g, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return g, err
	}

	gamesResponse := TwitchGamesResponse{}
	err = json.Unmarshal(body, &gamesResponse)
	if err != nil {
		return g, err
	}

	if len(gamesResponse.Data) != 1 {
		err = fmt.Errorf("game query for '%s' did not return exactly 1 result", gameID)
		return g, err
	}

	return gamesResponse.Data[0], nil
}

func (c *Controller) updateLiveStatus(streams []StreamData) error {

	var live bool
	publishers, err := c.getAllPublisher()
	if err != nil {
		return err
	}

	// mark previous live streams -> offline and notify of stream info change
	for i := range publishers {
		live = false
		p := &publishers[i]
		if p.IsTwitchLive() {
			for x := range streams {
				s := streams[x]
				if strings.ToLower(s.UserName) == strings.ToLower(p.TwitchStream) {
					live = true
					// retieve stream info
					g, err := c.getGame(s.GameID)
					if err != nil {
						return err
					}
					streamInfo := fmt.Sprintf("title: %s\ngame: %s", s.Title, g.Name)
					if p.StreamInfo != "" && p.StreamInfo != streamInfo {
						// streamer changed their stream info, set notification
						notification := fmt.Sprintf("%s switched it up!\n%s", p.Name, p.StreamInfo)
						c.setBucketValue("TwitchNotificationBucket", p.Name, notification)
					}
					p.StreamInfo = streamInfo
				}
			}
			if !live {
				c.setBucketValue("TwitchLiveBucket", p.Name, "")
				c.setBucketValue("StreamInfoBucket", p.Name, "")
				notification := fmt.Sprintf(":checkered_flag: %s finished streaming on twitch", p.Name)
				c.setBucketValue("TwitchNotificationBucket", p.Name, notification)
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
					c.setBucketValue("TwitchLiveBucket", p.Name, s.Type)
					streamLink := fmt.Sprintf("https://twitch.tv/%s", p.TwitchStream)
					notification := fmt.Sprintf(":movie_camera: %s started streaming on twitch!"+
						"\n%s\nwatch now: `%s`", p.Name, p.StreamInfo, streamLink)
					c.setBucketValue("TwitchNotificationBucket", p.Name, notification)
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
			log.Debug("sending discord notification: ", p.TwitchNotification)
			err := c.callWebhook(p.TwitchNotification)
			if err != nil {
				return err
			}
		}
		if p.TwitchNotification != "" {
			log.Debugf("resetting notification for %s (%s)", p.Name, p.TwitchStream)
			err = c.setBucketValue("TwitchNotificationBucket", p.Name, "")
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
