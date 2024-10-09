package pendulumsystem

import (
	"image"
	"image/color"
	"spinner-projector/ui"

	"gioui.org/layout"
	"gioui.org/widget/material"
)

type RainbowPendulumSystem struct {
	pendulums []*DoublePendulumSystem
}

func NewRainbowPendulumSystem(numPendulums int, innerAngle, outerAngle, space float64, colorSpace string) *RainbowPendulumSystem {
	var colors []color.NRGBA

	switch colorSpace {
	case "hcl":
		colors = ui.GenerateEvenHclColors(numPendulums)
	default:
		colors = ui.GenerateEvenRGBColors(numPendulums)
	}

	pends := make([]*DoublePendulumSystem, numPendulums)

	for i := range numPendulums {
		angle := (float64(i+1) * space) + outerAngle
		pends[i] = NewDoublePendulumSystem(innerAngle, angle, colors[i])
	}

	return &RainbowPendulumSystem{
		pendulums: pends,
	}
}

func (s *RainbowPendulumSystem) Menu(gtx layout.Context, theme *material.Theme) {

}

func (s *RainbowPendulumSystem) Draw(gtx layout.Context, size image.Point) layout.Dimensions {
	for _, pendulum := range s.pendulums {
		pendulum.Draw(gtx, size)
	}

	return layout.Dimensions{Size: size}
}

func (s *RainbowPendulumSystem) Update(gtx layout.Context, dt float64) {
	for _, pendulum := range s.pendulums {
		pendulum.Update(gtx, dt)
	}
}
