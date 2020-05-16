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
}

// ParseEnv parses configurations from environment environment variables
func (c *Config) ParseEnv() error {
	c.ServerFQDN = os.Getenv("SERVER_FQDN")
	c.PublishWebhook = os.Getenv("PUBLISH_WEBHOOK")
	_, err := strconv.ParseBool(os.Getenv("WEBHOOK_ENABLED"))
	if err == nil {
		c.WebhookEnabled = true
	}
	return nil
}
