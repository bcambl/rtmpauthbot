package controllers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/bcambl/rtmpauth/models"
	log "github.com/sirupsen/logrus"
)

const defaultWebhookURL = "https://discordapp.com/api/webhooks/1234567890/abcdefghijklmnopqrstuvwxyz1234567890"

func (c *Controller) callWebhook(message string) error {

	webhookURL := c.Config.PublishWebhook
	if webhookURL == defaultWebhookURL {
		err := errors.New("Default webhook value detected. Skipping webhook call")
		return err
	}
	body := models.WebhookPost{}
	body.Content = message

	b, err := json.Marshal(body)
	if err != nil {
		return err
	}

	contentType := "application/json"

	resp, err := http.Post(c.Config.PublishWebhook, contentType, bytes.NewBuffer(b))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	log.Infof("message posted to webhook: ", message)

	return nil
}
