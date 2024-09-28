package ui

import (
	"image/color"

	"github.com/lucasb-eyer/go-colorful"
)

var (
	Background = color.NRGBA{R: 192, G: 192, B: 192, A: 255}
	Red        = color.NRGBA{R: 192, G: 64, B: 64, A: 255}
	Green      = color.NRGBA{R: 64, G: 192, B: 64, A: 255}
	Blue       = color.NRGBA{R: 64, G: 64, B: 192, A: 255}
	White      = color.NRGBA{R: 255, G: 255, B: 255, A: 255}
	Black      = color.NRGBA{R: 0, G: 0, B: 0, A: 255}
	Maroon     = color.NRGBA{R: 127, G: 0, B: 0, A: 255}
	Purple     = color.NRGBA{R: 54, G: 1, B: 64, A: 255}
)

func GenerateEvenHclColors(n int) []color.NRGBA {
	colors := make([]color.NRGBA, n)
	for i := range n {
		hue := float64(i) / float64(n) * 360
		c := colorful.Hsl(hue, 1, 0.5) // Fixed chroma and lightness
		r, g, b := c.RGB255()
		colors[i] = color.NRGBA{R: r, G: g, B: b, A: 150}
	}
	return colors
}

func interpolate(start, end, t float64) float64 {
	return start + (end-start)*t
}

func interpolateSmoth(start, end, t float64) float64 {
	return start + (end-start)*(t*t*(3-2*t)) // Smoothstep function
}

func GenerateEvenRGBColors(n int) []color.NRGBA {
	colors := make([]color.NRGBA, n)
	keyColors := [][3]float64{
		{1, 0, 0}, // Red
		{1, 1, 0}, // Yellow
		{0, 1, 0}, // Green
		{0, 1, 1}, // Cyan
		{0, 0, 1}, // Blue
		{1, 0, 1}, // Magenta
	}

	stepsPerTransition := n / len(keyColors)
	for i := range n {
		keyIndex := int(i / stepsPerTransition)
		if keyIndex >= len(keyColors) {
			keyIndex = len(keyColors) - 1
		}

		startColor := keyColors[keyIndex]
		endColor := keyColors[(keyIndex+1)%len(keyColors)]

		t := float64((i % stepsPerTransition)) / float64(stepsPerTransition)

		r := interpolate(startColor[0], endColor[0], t)
		g := interpolate(startColor[1], endColor[1], t)
		b := interpolate(startColor[2], endColor[2], t)
		// r := interpolateSmoth(startColor[0], endColor[0], math.Mod(t*float64(len(keyColors)), 1))
		// g := interpolateSmoth(startColor[1], endColor[1], math.Mod(t*float64(len(keyColors)), 1))
		// b := interpolateSmoth(startColor[2], endColor[2], math.Mod(t*float64(len(keyColors)), 1))

		// fmt.Printf("{%.2f, %.2f %.2f}, %v -> %v\n", r, g, b, startColor, endColor)

		colors[i] = color.NRGBA{R: uint8(r * 255), G: uint8(g * 255), B: uint8(b * 255), A: 150}
	}
	return colors
}
