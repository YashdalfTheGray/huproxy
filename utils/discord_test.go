package utils

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/YashdalfTheGray/huproxy/types"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestDiscordNotifier_SendErrorNotification(t *testing.T) {
	log := logrus.New()
	log.SetOutput(bytes.NewBuffer(nil))

	cfg := &types.Config{
		ErrorDiscordWebhookUrl: "something",
	}
	notifier := NewDiscordNotifier(cfg, log)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

		var payload map[string]string
		err := json.NewDecoder(r.Body).Decode(&payload)
		assert.NoError(t, err)

		content, ok := payload["content"]
		assert.True(t, ok, "content key should be present in the payload")
		assert.Regexp(t, `\[(INFO|WARN|ERROR)\]`, content, "message should contain one of the log levels: INFO, WARN, ERROR")
		assert.Contains(t, content, time.Now().Format(time.RFC3339), "message should contain the date")
		assert.Contains(t, content, "test message", "message should contain the provided message")

		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	// Replace the webhook URL with the test server URL
	cfg.ErrorDiscordWebhookUrl = server.URL

	err := notifier.SendErrorNotification("test message")
	assert.NoError(t, err)
}
