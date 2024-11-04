package color

import (
	"fmt"
	"math"
	"strconv"
)

// XY represents the CIE xy color space values
type XY struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

// HexToXY converts a hex color code to CIE xy values
func HexToXY(hexColor string) (XY, error) {
	if len(hexColor) == 0 {
		return XY{}, fmt.Errorf("empty hex color string")
	}

	if hexColor[0] == '#' {
		hexColor = hexColor[1:]
	}

	if len(hexColor) != 6 {
		return XY{}, fmt.Errorf("invalid hex color: %s", hexColor)
	}

	r, err := strconv.ParseUint(hexColor[0:2], 16, 8)
	if err != nil {
		return XY{}, err
	}
	g, err := strconv.ParseUint(hexColor[2:4], 16, 8)
	if err != nil {
		return XY{}, err
	}
	b, err := strconv.ParseUint(hexColor[4:6], 16, 8)
	if err != nil {
		return XY{}, err
	}

	return RGBToXY(float64(r), float64(g), float64(b)), nil
}

// RGBToXY converts RGB values to CIE xy values using the sRGB color space and D65 white point
func RGBToXY(r, g, b float64) XY {
	r /= 255
	g /= 255
	b /= 255

	r = gammaCorrect(r)
	g = gammaCorrect(g)
	b = gammaCorrect(b)

	X := r*0.4124 + g*0.3576 + b*0.1805
	Y := r*0.2126 + g*0.7152 + b*0.0722
	Z := r*0.0193 + g*0.1192 + b*0.9505

	total := X + Y + Z
	if total == 0 {
		return XY{X: 0, Y: 0}
	}

	x := X / total
	y := Y / total

	x = math.Round(x*100) / 100
	y = math.Round(y*100) / 100

	return XY{X: x, Y: y}
}

// gammaCorrect applies gamma correction to a color value
func gammaCorrect(color float64) float64 {
	if color > 0.04045 {
		return math.Pow((color+0.055)/(1.055), 2.4)
	}
	return color / 12.92
}
