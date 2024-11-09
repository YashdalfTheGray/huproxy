// config/config_test.go
package config

import (
	"math"
	"os"
	"strings"
	"testing"

	"github.com/YashdalfTheGray/huproxy/color"
	"github.com/YashdalfTheGray/huproxy/types"

	"github.com/sirupsen/logrus"
)

type logWriter struct {
	logs *[]string
}

func (w *logWriter) Write(p []byte) (n int, err error) {
	*w.logs = append(*w.logs, string(p))
	return len(p), nil
}

func TestLoadConfig(t *testing.T) {
	tests := []struct {
		description      string
		envVars          map[string]string
		expectedConfig   *types.Config
		expectedWarnings []string
	}{
		{
			description: "All environment variables set correctly",
			envVars: map[string]string{
				"HUE_BRIDGE_ADDRESS":        "192.168.1.2",
				"GROUPED_LIGHT_ID":          "group1",
				"HUE_USERNAME":              "user123",
				"START_COLOR":               "#00ff00",
				"JUMP_COLOR":                "#0000ff",
				"DURATION_SECONDS":          "20",
				"ERROR_DISCORD_WEBHOOK_URL": "http://localhost:49152",
			},
			expectedConfig: &types.Config{
				BridgeAddress:          "192.168.1.2",
				GroupedLightID:         "group1",
				HueUsername:            "user123",
				StartColorHex:          "#00ff00",
				JumpColorHex:           "#0000ff",
				StartColorXY:           color.XY{X: 0.3, Y: 0.6},
				JumpColorXY:            color.XY{X: 0.15, Y: 0.06},
				DurationMS:             20000,
				ErrorDiscordWebhookUrl: "http://localhost:49152",
			},
		},
		{
			description: "Missing optional environment variables, using defaults",
			envVars: map[string]string{
				"HUE_BRIDGE_ADDRESS":        "192.168.1.2",
				"GROUPED_LIGHT_ID":          "group1",
				"HUE_USERNAME":              "user123",
				"ERROR_DISCORD_WEBHOOK_URL": "http://localhost:49152",
			},
			expectedConfig: &types.Config{
				BridgeAddress:          "192.168.1.2",
				GroupedLightID:         "group1",
				HueUsername:            "user123",
				StartColorHex:          "#ff5722",
				JumpColorHex:           "#ff0000",
				StartColorXY:           color.XY{X: 0.57, Y: 0.36},
				JumpColorXY:            color.XY{X: 0.64, Y: 0.33},
				DurationMS:             15000,
				ErrorDiscordWebhookUrl: "http://localhost:49152",
			},
		},
		{
			description: "Invalid START_COLOR, should use default and log warning",
			envVars: map[string]string{
				"HUE_BRIDGE_ADDRESS":        "192.168.1.2",
				"GROUPED_LIGHT_ID":          "group1",
				"HUE_USERNAME":              "user123",
				"START_COLOR":               "invalid",
				"JUMP_COLOR":                "#0000ff",
				"DURATION_SECONDS":          "15",
				"ERROR_DISCORD_WEBHOOK_URL": "http://localhost:49152",
			},
			expectedConfig: &types.Config{
				BridgeAddress:          "192.168.1.2",
				GroupedLightID:         "group1",
				HueUsername:            "user123",
				StartColorHex:          "#ff5722",
				JumpColorHex:           "#0000ff",
				StartColorXY:           color.XY{X: 0.57, Y: 0.36},
				JumpColorXY:            color.XY{X: 0.15, Y: 0.06},
				DurationMS:             15000,
				ErrorDiscordWebhookUrl: "http://localhost:49152",
			},
			expectedWarnings: []string{
				"Invalid START_COLOR value, using default.",
			},
		},
		{
			description: "Invalid JUMP_COLOR, should use default and log warning",
			envVars: map[string]string{
				"HUE_BRIDGE_ADDRESS":        "192.168.1.2",
				"GROUPED_LIGHT_ID":          "group1",
				"HUE_USERNAME":              "user123",
				"START_COLOR":               "#00ff00",
				"JUMP_COLOR":                "invalid",
				"DURATION_SECONDS":          "15",
				"ERROR_DISCORD_WEBHOOK_URL": "http://localhost:49152",
			},
			expectedConfig: &types.Config{
				BridgeAddress:          "192.168.1.2",
				GroupedLightID:         "group1",
				HueUsername:            "user123",
				StartColorHex:          "#00ff00",
				JumpColorHex:           "#ff0000",
				StartColorXY:           color.XY{X: 0.3, Y: 0.6},
				JumpColorXY:            color.XY{X: 0.64, Y: 0.33},
				DurationMS:             15000,
				ErrorDiscordWebhookUrl: "http://localhost:49152",
			},
			expectedWarnings: []string{
				"Invalid JUMP_COLOR value, using default.",
			},
		},
		{
			description: "Invalid DURATION_SECONDS, should use default and log warning",
			envVars: map[string]string{
				"HUE_BRIDGE_ADDRESS":        "192.168.1.2",
				"GROUPED_LIGHT_ID":          "group1",
				"HUE_USERNAME":              "user123",
				"DURATION_SECONDS":          "invalid",
				"ERROR_DISCORD_WEBHOOK_URL": "http://localhost:49152",
			},
			expectedConfig: &types.Config{
				BridgeAddress:          "192.168.1.2",
				GroupedLightID:         "group1",
				HueUsername:            "user123",
				StartColorHex:          "#ff5722",
				JumpColorHex:           "#ff0000",
				StartColorXY:           color.XY{X: 0.57, Y: 0.36},
				JumpColorXY:            color.XY{X: 0.64, Y: 0.33},
				DurationMS:             15000,
				ErrorDiscordWebhookUrl: "http://localhost:49152",
			},
			expectedWarnings: []string{
				"Invalid DURATION_SECONDS value, using default of 15 seconds.",
			},
		},
		{
			description: "Negative DURATION_SECONDS, should use default and log warning",
			envVars: map[string]string{
				"HUE_BRIDGE_ADDRESS":        "192.168.1.2",
				"GROUPED_LIGHT_ID":          "group1",
				"HUE_USERNAME":              "user123",
				"DURATION_SECONDS":          "-10",
				"ERROR_DISCORD_WEBHOOK_URL": "http://localhost:49152",
			},
			expectedConfig: &types.Config{
				BridgeAddress:          "192.168.1.2",
				GroupedLightID:         "group1",
				HueUsername:            "user123",
				StartColorHex:          "#ff5722",
				JumpColorHex:           "#ff0000",
				StartColorXY:           color.XY{X: 0.57, Y: 0.36},
				JumpColorXY:            color.XY{X: 0.64, Y: 0.33},
				DurationMS:             15000,
				ErrorDiscordWebhookUrl: "http://localhost:49152",
			},
			expectedWarnings: []string{
				"Invalid DURATION_SECONDS value, using default of 15 seconds.",
			},
		},
		{
			description: "Missing discord notification URL, logs out warning",
			envVars: map[string]string{
				"HUE_BRIDGE_ADDRESS": "192.168.1.2",
				"GROUPED_LIGHT_ID":   "group1",
				"HUE_USERNAME":       "user123",
				"DURATION_SECONDS":   "15",
			},
			expectedConfig: &types.Config{
				BridgeAddress:          "192.168.1.2",
				GroupedLightID:         "group1",
				HueUsername:            "user123",
				StartColorHex:          "#ff5722",
				JumpColorHex:           "#ff0000",
				StartColorXY:           color.XY{X: 0.57, Y: 0.36},
				JumpColorXY:            color.XY{X: 0.64, Y: 0.33},
				DurationMS:             15000,
				ErrorDiscordWebhookUrl: "",
			},
			expectedWarnings: []string{
				"No Discord webhook URL set for error notifications",
			},
		},
		{
			description: "All required environment variables missing",
			envVars:     map[string]string{},
			expectedConfig: &types.Config{
				BridgeAddress:          "",
				GroupedLightID:         "",
				HueUsername:            "",
				StartColorHex:          "#ff5722",
				JumpColorHex:           "#ff0000",
				StartColorXY:           color.XY{X: 0.57, Y: 0.36},
				JumpColorXY:            color.XY{X: 0.64, Y: 0.33},
				DurationMS:             15000,
				ErrorDiscordWebhookUrl: "",
			},
			expectedWarnings: []string{
				"No Discord webhook URL set for error notifications",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			for key, value := range test.envVars {
				os.Setenv(key, value)
			}

			var logOutput []string
			log := logrus.New()
			log.SetLevel(logrus.WarnLevel)
			log.SetFormatter(&logrus.TextFormatter{
				DisableTimestamp: true,
			})
			log.SetOutput(&logWriter{logs: &logOutput})

			cfg, err := LoadConfig(log)
			if err != nil {
				t.Fatalf("LoadConfig returned an unexpected error: %v", err)
			}

			if cfg.BridgeAddress != test.expectedConfig.BridgeAddress {
				t.Errorf("Expected BridgeAddress '%s', got '%s'", test.expectedConfig.BridgeAddress, cfg.BridgeAddress)
			}
			if cfg.GroupedLightID != test.expectedConfig.GroupedLightID {
				t.Errorf("Expected GroupedLightID '%s', got '%s'", test.expectedConfig.GroupedLightID, cfg.GroupedLightID)
			}
			if cfg.HueUsername != test.expectedConfig.HueUsername {
				t.Errorf("Expected HueUsername '%s', got '%s'", test.expectedConfig.HueUsername, cfg.HueUsername)
			}
			if cfg.StartColorHex != test.expectedConfig.StartColorHex {
				t.Errorf("Expected StartColorHex '%s', got '%s'", test.expectedConfig.StartColorHex, cfg.StartColorHex)
			}
			if cfg.JumpColorHex != test.expectedConfig.JumpColorHex {
				t.Errorf("Expected JumpColorHex '%s', got '%s'", test.expectedConfig.JumpColorHex, cfg.JumpColorHex)
			}
			if !approxEqual(cfg.StartColorXY.X, test.expectedConfig.StartColorXY.X) || !approxEqual(cfg.StartColorXY.Y, test.expectedConfig.StartColorXY.Y) {
				t.Errorf("Expected StartColorXY %v, got %v", test.expectedConfig.StartColorXY, cfg.StartColorXY)
			}
			if !approxEqual(cfg.JumpColorXY.X, test.expectedConfig.JumpColorXY.X) || !approxEqual(cfg.JumpColorXY.Y, test.expectedConfig.JumpColorXY.Y) {
				t.Errorf("Expected JumpColorXY %v, got %v", test.expectedConfig.JumpColorXY, cfg.JumpColorXY)
			}
			if cfg.DurationMS != test.expectedConfig.DurationMS {
				t.Errorf("Expected DurationMS %d, got %d", test.expectedConfig.DurationMS, cfg.DurationMS)
			}

			if len(test.expectedWarnings) != len(logOutput) {
				t.Errorf("Expected %d warning(s), got %d", len(test.expectedWarnings), len(logOutput))
			} else {
				for i, expectedWarning := range test.expectedWarnings {
					if !contains(logOutput[i], expectedWarning) {
						t.Errorf("Expected warning '%s', got '%s'", expectedWarning, logOutput[i])
					}
				}
			}

			for key := range test.envVars {
				os.Unsetenv(key)
			}
		})
	}
}

func approxEqual(a, b float64) bool {
	return math.Abs(a-b) < 0.02
}

func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}
