package config

import (
	"os"
	"strconv"
)

// Config contains config vars parsed from the environment
type Config struct {
	ServerFQDN     string
	PublishWebhook string
	WebhookEnabled bool
	ServerIP       string
	ServerPort     string
}

// ParseEnv parses configurations from environment environment variables
func (c *Config) ParseEnv() error {
	c.ServerIP = os.Getenv("SERVER_IP")
	c.ServerPort = os.Getenv("SERVER_PORT")
	c.ServerFQDN = os.Getenv("SERVER_FQDN")
	c.PublishWebhook = os.Getenv("PUBLISH_WEBHOOK")
	_, err := strconv.ParseBool(os.Getenv("WEBHOOK_ENABLED"))
	if err == nil {
		c.WebhookEnabled = true
	}
	return nil
}
