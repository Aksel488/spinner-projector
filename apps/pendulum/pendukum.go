package pendulum

import (
	"image"
	"image/color"
	"math"
	"spinner-projector/ui"

	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op"
)

type DoublePendulumSystem struct {
	pendulum1 *pendulum
	pendulum2 *pendulum

	gravity float64
}

func NewDoublePendulumSystem() *DoublePendulumSystem {
	pendulum1 := NewPendulum(500, 400, 0, 1)
	x, y := pendulum1.getEnd()
	pendulum2 := NewPendulum(x, y, 0, 2)

	return &DoublePendulumSystem{
		pendulum1: pendulum1,
		pendulum2: pendulum2,
		gravity:   9.81,
	}
}

func (s *DoublePendulumSystem) Draw(gtx layout.Context, size image.Point) layout.Dimensions {
	s.pendulum1.Draw(gtx.Ops)
	s.pendulum2.Draw(gtx.Ops)

	return layout.Dimensions{Size: size}
}

func (s *DoublePendulumSystem) Update(gtx layout.Context, dt float64) {
	s.pendulum1.Update(dt)

	x, y := s.pendulum1.getEnd()
	s.pendulum2.px = x
	s.pendulum2.py = y
	s.pendulum2.Update(dt)

	// s.pendulum1.angleVel += (s.gravity*math.Sin(s.pendulum1.angle) - s.pendulum2.angleVel*math.Sin(s.pendulum1.angle-s.pendulum2.angle)) / s.pendulum2.angle * dt
	// s.pendulum2.angleVel += (s.gravity*math.Sin(s.pendulum2.angle) - s.pendulum1.angleVel*math.Sin(s.pendulum2.angle-s.pendulum1.angle)) / s.pendulum1.angle * dt
}

type pendulum struct {
	px, py   float64
	angle    float64
	angleVel float64
	r        int
	color    color.NRGBA
}

func NewPendulum(x, y, angle, angleVel float64) *pendulum {
	return &pendulum{
		px:       x,
		py:       y,
		angle:    angle,
		angleVel: angleVel,
		r:        100,
		color:    ui.Blue,
	}
}

func (p *pendulum) getEnd() (x, y float64) {
	x = float64(p.r)*math.Cos(p.angle) + p.px
	y = float64(p.r)*math.Sin(p.angle) + p.py
	return x, y
}

func (p *pendulum) Draw(ops *op.Ops) {
	x, y := p.getEnd()
	ui.Line(ops, f32.Pt(float32(p.px), float32(p.py)), f32.Pt(float32(x), float32(y)), 2, p.color)
}

func (p *pendulum) Update(dt float64) {
	p.angle += p.angleVel * dt

	// p.px = 200 + 100*math.Cos(p.angle)
	// p.py = 200 + 100*math.Sin(p.angle)
}
