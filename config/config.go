package config

import (
	"os"
	"strconv"

	"github.com/YashdalfTheGray/huproxy/color"
	"github.com/YashdalfTheGray/huproxy/types"

	"github.com/sirupsen/logrus"
)

// LoadConfig reads the environment variables, sets defaults, validates,
// and returns a Config object.
func LoadConfig(log *logrus.Logger) (*types.Config, error) {
	config := &types.Config{
		BridgeAddress:          os.Getenv("HUE_BRIDGE_ADDRESS"),
		ErrorDiscordWebhookUrl: os.Getenv("ERROR_DISCORD_WEBHOOK_URL"),
		GroupedLightID:         os.Getenv("GROUPED_LIGHT_ID"),
		HueUsername:            os.Getenv("HUE_USERNAME"),
		StartColorHex:          os.Getenv("START_COLOR"),
		JumpColorHex:           os.Getenv("JUMP_COLOR"),
	}

	if config.ErrorDiscordWebhookUrl == "" {
		log.Warn("No Discord webhook URL set for error notifications")
	}

	if config.StartColorHex == "" {
		config.StartColorHex = "#ff5722"
	}
	if config.JumpColorHex == "" {
		config.JumpColorHex = "#ff0000"
	}

	startColorXY, err := color.HexToXY(config.StartColorHex)
	if err != nil {
		log.Warn("Invalid START_COLOR value, using default.")
		config.StartColorHex = "#ff5722"
		startColorXY, _ = color.HexToXY(config.StartColorHex)
	}
	config.StartColorXY = startColorXY

	jumpColorXY, err := color.HexToXY(config.JumpColorHex)
	if err != nil {
		log.Warn("Invalid JUMP_COLOR value, using default.")
		config.JumpColorHex = "#ff0000"
		jumpColorXY, _ = color.HexToXY(config.JumpColorHex)
	}
	config.JumpColorXY = jumpColorXY

	durationStr := os.Getenv("DURATION_SECONDS")
	if durationStr == "" {
		durationStr = "15"
	}
	durationSeconds, err := strconv.Atoi(durationStr)
	if err != nil || durationSeconds <= 0 {
		log.Warn("Invalid DURATION_SECONDS value, using default of 15 seconds.")
		durationSeconds = 15
	}
	config.DurationMS = durationSeconds * 1000

	return config, nil
}
