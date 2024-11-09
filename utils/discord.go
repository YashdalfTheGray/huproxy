package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/YashdalfTheGray/huproxy/types"
	"github.com/sirupsen/logrus"
)

type DiscordNotifier struct {
	Config *types.Config
	Log    *logrus.Logger
}

// NewDiscordNotifier creates a new DiscordNotifier with the given Config and Logger.
func NewDiscordNotifier(config *types.Config, log *logrus.Logger) *DiscordNotifier {
	return &DiscordNotifier{
		Config: config,
		Log:    log,
	}
}

func (d *DiscordNotifier) SendErrorNotification(message string) error {
	return d.sendNotification(d.Config.ErrorDiscordWebhookUrl, message, types.LogLevelError)
}

// SendNotification sends a message to the configured Discord webhook URL.
func (d *DiscordNotifier) sendNotification(webhookURL string, message string, level types.LogLevel) error {
	if webhookURL == "" {
		d.Log.Warn("No Discord webhook URL provided, skipping notification")
		return nil
	}

	decoratedMessage := fmt.Sprintf("%s `[%s]` %s", time.Now().Format(time.RFC3339), level.String(), message)

	payload := map[string]string{"content": decoratedMessage}
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		d.Log.Error("Failed to marshal JSON payload: ", err)
		return err
	}

	req, err := http.NewRequest("POST", webhookURL, bytes.NewBuffer(jsonPayload))
	if err != nil {
		d.Log.Error("Failed to create new HTTP request: ", err)
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		d.Log.Error("Failed to send HTTP request: ", err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK {
		d.Log.Warnf("Failed to send notification, status code: %d", resp.StatusCode)
	}

	return nil
}
