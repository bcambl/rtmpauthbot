package config

import (
	"os"
	"strconv"
	"time"
)

// Config contains config vars parsed from the environment
type Config struct {
	ServerIP              string
	ServerPort            string
	ServerFQDN            string
	TwitchEnabled         bool
	TwitchClientID        string
	TwitchClientSecret    string
	DiscordWebhook        string
	DiscordWebhookEnabled bool
	TwitchPollRate        time.Duration
}

// ParseEnv parses configurations from environment environment variables
func (c *Config) ParseEnv() error {
	var (
		err         error
		pollRateSec int64
	)
	c.ServerIP = os.Getenv("SERVER_IP")
	c.ServerPort = os.Getenv("SERVER_PORT")
	c.ServerFQDN = os.Getenv("SERVER_FQDN")
	c.TwitchClientID = os.Getenv("TWITCH_CLIENT_ID")
	c.TwitchClientSecret = os.Getenv("TWITCH_CLIENT_SECRET")
	c.DiscordWebhook = os.Getenv("DISCORD_WEBHOOK")
	c.DiscordWebhookEnabled, err = strconv.ParseBool(os.Getenv("DISCORD_WEBHOOK_ENABLED"))
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
