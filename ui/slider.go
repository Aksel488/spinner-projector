package ui

import (
	"fmt"
	"math"
	"spinner-projector/util"

	"gioui.org/layout"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

type Slider interface {
	Layout(gtx layout.Context, theme *material.Theme, valueName string) layout.Dimensions
	GetFloat() *widget.Float
	SetValue(v float32)
	Value() float32
}

// Slider wraps widget.Float with a linear min/max range
type LinearSlider struct {
	Float *widget.Float
	Min   float32
	Max   float32
}

// creates a new Slider with a given range
func NewLinearSlider(value, min, max float32) *LinearSlider {
	slider := &LinearSlider{
		Float: &widget.Float{},
		Min:   min,
		Max:   max,
	}

	slider.SetValue(value)

	return slider
}

func (s *LinearSlider) GetFloat() *widget.Float {
	return s.Float
}

// Convert the 0–1 value to the min-max range
func (s *LinearSlider) Value() float32 {
	return s.Min + s.Float.Value*(s.Max-s.Min)
}

// converts the value to the 0-1 range
func (s *LinearSlider) SetValue(v float32) {
	if v < s.Min {
		v = s.Min
	} else if v > s.Max {
		v = s.Max
	}
	s.Float.Value = (v - s.Min) / (s.Max - s.Min)
}

func (s *LinearSlider) Layout(gtx layout.Context, theme *material.Theme, valueName string) layout.Dimensions {
	return layout.Flex{
		Axis:      layout.Vertical,
		Alignment: layout.Middle,
	}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			// Show current slider value
			gtx.Constraints.Max.X = gtx.Constraints.Max.X - 50
			label := material.Body1(theme, fmt.Sprintf("%s: %.2f", valueName, s.Value()))
			label.Color = White
			return label.Layout(gtx)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			gtx.Constraints.Max.X = gtx.Constraints.Max.X - 50
			return layout.Flex{
				Axis:      layout.Horizontal,
				Alignment: layout.Middle,
			}.Layout(gtx,
				layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
					sliderStyle := material.Slider(theme, s.Float).Layout(gtx)
					return sliderStyle
				}),
			)
		}),
	)
}

// LogSlider wraps widget.Float with a logarithmic min/max range
type LogSlider struct {
	Float *widget.Float
	Min   float32
	Max   float32
}

// NewLogSlider creates a new LogSlider with a given logarithmic range
func NewLogSlider(value, min, max float32) *LogSlider {
	logSlider := &LogSlider{
		Float: &widget.Float{},
		Min:   min,
		Max:   max,
	}

	logSlider.SetValue(value)

	return logSlider
}

func (s *LogSlider) GetFloat() *widget.Float {
	return s.Float
}

// Value returns the slider value converted to the logarithmic custom range
// Formula: exp(log(min) + slider_value * (log(max) - log(min)))
func (s *LogSlider) Value() float32 {
	logMin := math.Log(float64(s.Min))
	logMax := math.Log(float64(s.Max))

	// Linear slider value (0-1) is transformed into logarithmic space
	logVal := logMin + float64(s.Float.Value)*(logMax-logMin)

	return float32(math.Exp(logVal)) // Convert back using exp
}

// SetValue sets the slider value using a value from the logarithmic custom range
func (s *LogSlider) SetValue(v float32) {
	if v < s.Min {
		v = s.Min
	} else if v > s.Max {
		v = s.Max
	}

	// Convert the logarithmic value to the 0-1 range
	logMin := math.Log(float64(s.Min))
	logMax := math.Log(float64(s.Max))
	logVal := math.Log(float64(v))

	s.Float.Value = float32((logVal - logMin) / (logMax - logMin))
}

func (s *LogSlider) Layout(gtx layout.Context, theme *material.Theme, valueName string) layout.Dimensions {
	return layout.Flex{
		Axis:      layout.Vertical,
		Alignment: layout.Middle,
	}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			// Show current slider value
			gtx.Constraints.Max.X = gtx.Constraints.Max.X - 50

			minDecimals := util.CountDecimals(s.Min)
			leftPadding := util.GetLeftPadding(s.Max)
			formatLogValue := util.FormatLogValue(s.Value(), minDecimals, leftPadding)

			label := material.Body1(theme, fmt.Sprintf("%s: %s", valueName, formatLogValue))
			label.Color = White
			return label.Layout(gtx)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			gtx.Constraints.Max.X = gtx.Constraints.Max.X - 50
			return layout.Flex{
				Axis:      layout.Horizontal,
				Alignment: layout.Middle,
			}.Layout(gtx,
				layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
					sliderStyle := material.Slider(theme, s.Float).Layout(gtx)
					return sliderStyle
				}),
			)
		}),
	)
}
