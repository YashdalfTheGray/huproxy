package handlers

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"io"
	"net/http"

	"github.com/YashdalfTheGray/huproxy/types"

	"github.com/sirupsen/logrus"
)

type Handler struct {
	Config *types.Config
	Log    *logrus.Logger
}

// NewHandler creates a new Handler with the given Config and Logger.
func NewHandler(config *types.Config, log *logrus.Logger) *Handler {
	return &Handler{
		Config: config,
		Log:    log,
	}
}

func (h *Handler) PingHandler(w http.ResponseWriter, r *http.Request) {
	h.Log.Infof("Received /ping request from %s", r.RemoteAddr)
	var response types.Response

	if h.Config.BridgeAddress == "" || h.Config.GroupedLightID == "" || h.Config.HueUsername == "" {
		response = types.Error("")
		h.Log.Warn("Missing one or more environment variables.")
	} else {
		response = types.Success()
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) PageHandler(w http.ResponseWriter, r *http.Request) {
	h.Log.Infof("Received /page request from %s", r.RemoteAddr)

	if h.Config.BridgeAddress == "" || h.Config.GroupedLightID == "" || h.Config.HueUsername == "" {
		response := types.Error("")
		h.Log.Warn("Environment variables are not properly set.")
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	url := "https://" + h.Config.BridgeAddress + "/clip/v2/resource/grouped_light/" + h.Config.GroupedLightID

	body := map[string]interface{}{
		"signaling": map[string]interface{}{
			"signal":   "alternating",
			"duration": h.Config.DurationMS,
			"colors": []map[string]interface{}{
				{"xy": h.Config.StartColorXY},
				{"xy": h.Config.JumpColorXY},
			},
		},
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		response := types.Error("")
		h.Log.Error("Failed to marshal JSON body.")
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	// custom http client that ignores certificate verification
	// this is some shit that the hue api imposes on us
	httpClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	req, err := http.NewRequest("PUT", url, bytes.NewReader(jsonBody))
	if err != nil {
		response := types.Error("")
		h.Log.Error("Error creating HTTP request.")
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	req.Header.Add("hue-application-key", h.Config.HueUsername)
	req.Header.Set("Content-Type", "application/json")

	resp, err := httpClient.Do(req)
	if err != nil {
		response := types.Error("")
		h.Log.Error("Error sending HTTP request.")
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}
	defer resp.Body.Close()

	_, err = io.ReadAll(resp.Body)
	if err != nil {
		response := types.Error("")
		h.Log.Error("Error reading response body.")
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	var response types.Response

	if resp.StatusCode == 200 {
		response = types.Success()
		h.Log.Info("Successfully sent command to Hue Bridge.")
	} else {
		response = types.Error("")
		h.Log.Warnf("Received non-200 status code from Hue Bridge: %d", resp.StatusCode)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
