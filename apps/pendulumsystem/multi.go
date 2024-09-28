package pendulumsystem

import (
	"image"
	"math"
	"spinner-projector/ui"

	"gioui.org/layout"
)

type MultiPendulumSystem struct {
	pendulums []*DoublePendulumSystem
}

func NewMultiPendulumSystem(numPendulums int) *MultiPendulumSystem {
	pends := make([]*DoublePendulumSystem, numPendulums)
	for i := range numPendulums {
		angle := float64(i)*2*math.Pi/float64(numPendulums) + math.Pi/float64(numPendulums)
		pends[i] = NewDoublePendulumSystem(0, angle, ui.Blue)
	}

	return &MultiPendulumSystem{
		pendulums: pends,
	}
}

func (s *MultiPendulumSystem) Draw(gtx layout.Context, size image.Point) layout.Dimensions {
	for _, pendulum := range s.pendulums {
		pendulum.Draw(gtx, size)
	}

	return layout.Dimensions{Size: size}
}

func (s *MultiPendulumSystem) Update(gtx layout.Context, dt float64) {
	for _, pendulum := range s.pendulums {
		pendulum.Update(gtx, dt)
	}
}
