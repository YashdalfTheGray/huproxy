package types

import "github.com/YashdalfTheGray/huproxy/color"

// Config holds the environment configuration.
type Config struct {
	BridgeAddress  string
	GroupedLightID string
	HueUsername    string
	StartColorHex  string
	JumpColorHex   string
	StartColorXY   color.XY
	JumpColorXY    color.XY
	DurationMS     int
}

// Response represents the structure of responses sent to clients.
type Response struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
}

// Success creates a success response.
func Success() Response {
	return Response{Status: "okay"}
}

// Error creates an error response with a message.
func Error(message string) Response {
	return Response{Status: "broke", Message: message}
}
