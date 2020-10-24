package config

import (
	"os"
	"strconv"
	"time"
)

// Config contains config vars parsed from the environment
type Config struct {
	AuthServerIP       string
	AuthServerPort     string
	RTMPServerFQDN     string
	RTMPServerPort     string
	TwitchEnabled      bool
	TwitchClientID     string
	TwitchClientSecret string
	DiscordWebhook     string
	DiscordEnabled     bool
	TwitchPollRate     time.Duration
}

// ParseEnv parses configurations from environment environment variables
func (c *Config) ParseEnv() error {
	var (
		err         error
		pollRateSec int64
	)
	c.AuthServerIP = os.Getenv("AUTH_SERVER_IP")
	c.AuthServerPort = os.Getenv("AUTH_SERVER_PORT")
	c.RTMPServerFQDN = os.Getenv("RTMP_SERVER_FQDN")
	c.RTMPServerPort = os.Getenv("RTMP_SERVER_PORT")
	c.TwitchClientID = os.Getenv("TWITCH_CLIENT_ID")
	c.TwitchClientSecret = os.Getenv("TWITCH_CLIENT_SECRET")
	c.DiscordWebhook = os.Getenv("DISCORD_WEBHOOK")
	c.DiscordEnabled, err = strconv.ParseBool(os.Getenv("DISCORD_ENABLED"))
	if err != nil {
		return err
	}
	c.TwitchEnabled, err = strconv.ParseBool(os.Getenv("TWITCH_ENABLED"))
	if err != nil {
		return err
	}
	pollRateSec, err = strconv.ParseInt(os.Getenv("TWITCH_POLL_RATE"), 0, 0)
	if err != nil {
		// Default poll rate to 60sec (far below allowed rate limits)
		pollRateSec = 60
	}
	// ensure a sane minimum twitch poll rate
	if pollRateSec < 5 {
		pollRateSec = 5
	}
	c.TwitchPollRate = (time.Duration(pollRateSec) * time.Second)

	return nil
}
