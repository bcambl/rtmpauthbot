package controllers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"

	log "github.com/sirupsen/logrus"
)

const defaultWebhookURL = "https://discordapp.com/api/webhooks/1234567890/abcdefghijklmnopqrstuvwxyz1234567890"

// DiscordWebhook is used to marshal the data sent to the discord webhook
type DiscordWebhook struct {
	Content string `json:"content"`
}

func (c *Controller) callWebhook(message string) error {

	webhookURL := c.Config.DiscordWebhook
	if webhookURL == defaultWebhookURL {
		err := errors.New("Default webhook value detected. Skipping webhook call")
		return err
	}
	body := DiscordWebhook{}
	body.Content = message

	b, err := json.Marshal(body)
	if err != nil {
		return err
	}

	contentType := "application/json"

	resp, err := http.Post(c.Config.DiscordWebhook, contentType, bytes.NewBuffer(b))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	log.Info("message posted to webhook: ", message)

	return nil
}
