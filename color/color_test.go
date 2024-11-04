package color

import (
	"math"
	"testing"
)

func approxEqual(a, b float64) bool {
	return math.Abs(a-b) < 0.02 // Allow a tolerance of 0.02
}

func TestHexToXY(t *testing.T) {
	tests := []struct {
		description string
		hexColor    string
		expected    XY
		expectErr   bool
	}{
		{
			description: "Valid hex color #ff0000 (Red)",
			hexColor:    "#ff0000",
			expected:    XY{X: 0.64, Y: 0.33},
			expectErr:   false,
		},
		{
			description: "Valid hex color #ff5722 (Deep Orange)",
			hexColor:    "#ff5722",
			expected:    XY{X: 0.57, Y: 0.36},
			expectErr:   false,
		},
		{
			description: "Valid hex color without hash prefix",
			hexColor:    "ff0000",
			expected:    XY{X: 0.64, Y: 0.33},
			expectErr:   false,
		},
		{
			description: "Invalid hex color with non-hex characters",
			hexColor:    "#xyz123",
			expectErr:   true,
		},
		{
			description: "Hex color with incorrect length (too short)",
			hexColor:    "#12345",
			expectErr:   true,
		},
		{
			description: "Hex color with incorrect length (too long)",
			hexColor:    "#1234567",
			expectErr:   true,
		},
		{
			description: "Hex color with invalid red component",
			hexColor:    "#gg0000",
			expectErr:   true,
		},
		{
			description: "Hex color with invalid green component",
			hexColor:    "#00gg00",
			expectErr:   true,
		},
		{
			description: "Hex color with invalid blue component",
			hexColor:    "#0000gg",
			expectErr:   true,
		},
		{
			description: "Empty hex color string",
			hexColor:    "",
			expectErr:   true,
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			result, err := HexToXY(test.hexColor)
			if test.expectErr {
				if err == nil {
					t.Errorf("Expected error for input '%s', got nil", test.hexColor)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error for input '%s': %v", test.hexColor, err)
				}
				if !approxEqual(result.X, test.expected.X) || !approxEqual(result.Y, test.expected.Y) {
					t.Errorf("HexToXY(%s) = %v; expected %v", test.hexColor, result, test.expected)
				}
			}
		})
	}
}

func TestRGBToXY(t *testing.T) {
	tests := []struct {
		description string
		r           float64
		g           float64
		b           float64
		expected    XY
	}{
		{
			description: "RGB(255, 0, 0) (Red)",
			r:           255,
			g:           0,
			b:           0,
			expected:    XY{X: 0.64, Y: 0.33},
		},
		{
			description: "RGB(0, 255, 0) (Green)",
			r:           0,
			g:           255,
			b:           0,
			expected:    XY{X: 0.3, Y: 0.6},
		},
		{
			description: "RGB(0, 0, 255) (Blue)",
			r:           0,
			g:           0,
			b:           255,
			expected:    XY{X: 0.15, Y: 0.06},
		},
		{
			description: "RGB(0, 0, 0) leading to zero total",
			r:           0,
			g:           0,
			b:           0,
			expected:    XY{X: 0, Y: 0},
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			result := RGBToXY(test.r, test.g, test.b)
			if !approxEqual(result.X, test.expected.X) || !approxEqual(result.Y, test.expected.Y) {
				t.Errorf("RGBToXY(%f, %f, %f) = %v; expected %v", test.r, test.g, test.b, result, test.expected)
			}
		})
	}
}

func TestGammaCorrect(t *testing.T) {
	tests := []struct {
		description string
		input       float64
		expected    float64
	}{
		{
			description: "Value just below threshold (0.04044)",
			input:       0.04044,
			expected:    0.04044 / 12.92,
		},
		{
			description: "Value at threshold (0.04045)",
			input:       0.04045,
			expected:    0.04045 / 12.92,
		},
		{
			description: "Value just above threshold (0.04046)",
			input:       0.04046,
			expected:    math.Pow((0.04046+0.055)/1.055, 2.4),
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			result := gammaCorrect(test.input)
			if !approxEqual(result, test.expected) {
				t.Errorf("gammaCorrect(%f) = %f; expected %f", test.input, result, test.expected)
			}
		})
	}
}
