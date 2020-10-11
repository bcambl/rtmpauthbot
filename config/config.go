package config

import (
	"os"
	"strconv"
)

// Config contains config vars parsed from the environment
type Config struct {
	ServerIP              string
	ServerPort            string
	ServerFQDN            string
	TwitchClientID        string
	TwitchClientSecret    string
	DiscordWebhook        string
	DiscordWebhookEnabled bool
}

// ParseEnv parses configurations from environment environment variables
func (c *Config) ParseEnv() error {
	c.ServerIP = os.Getenv("SERVER_IP")
	c.ServerPort = os.Getenv("SERVER_PORT")
	c.ServerFQDN = os.Getenv("SERVER_FQDN")
	c.TwitchClientID = os.Getenv("TWITCH_CLIENT_ID")
	c.TwitchClientSecret = os.Getenv("TWITCH_CLIENT_SECRET")
	c.DiscordWebhook = os.Getenv("DISCORD_WEBHOOK")
	_, err := strconv.ParseBool(os.Getenv("DISCORD_WEBHOOK_ENABLED"))
	if err == nil {
		c.DiscordWebhookEnabled = true
	}
	return nil
}
