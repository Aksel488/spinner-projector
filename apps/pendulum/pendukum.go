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
	p1 *pendulum
	p2 *pendulum
	g  float64
}

func NewDoublePendulumSystem() *DoublePendulumSystem {
	pendulum1 := NewPendulum(500, 400, -0.000001, 0)
	x, y := pendulum1.getEnd()
	pendulum2 := NewPendulum(x, y, 0.01, 0)

	return &DoublePendulumSystem{
		p1: pendulum1,
		p2: pendulum2,
		g:  9.81,
	}
}

func (s *DoublePendulumSystem) Draw(gtx layout.Context, size image.Point) layout.Dimensions {
	s.p1.Draw(gtx.Ops)
	s.p2.Draw(gtx.Ops)

	return layout.Dimensions{Size: size}
}

func (s *DoublePendulumSystem) Update(gtx layout.Context, dt float64) {
	a1 := s.p1.angle
	a2 := s.p2.angle
	av1 := s.p1.angleVel
	av2 := s.p2.angleVel
	m1 := s.p1.m
	m2 := s.p2.m

	mass1 := -s.g * (2*m1 + m2) * math.Sin(a1)
	mass2 := -m2 * s.g * math.Sin(a1-2*a2)

	interaction := -2 * math.Sin(a1-a2) * m2 * math.Cos(math.Pow(av2, 2)*s.p2.r+math.Pow(av1, 2)*s.p1.r*math.Cos(a1-a2))
	normalization := s.p1.r * (2*m1 + m2 - m2*math.Cos(2*a1-2*a2))
	angle1Dot := (mass1 + mass2 + interaction) / normalization

	system := 2 * math.Sin(a1-a2) * (math.Pow(av1, 2)*s.p1.r*(m1+m2) + s.g*(m1+m2)*math.Cos(a1) + math.Pow(av2, 2)*s.p2.r*m2*math.Cos(a1-a2))
	normalization = s.p1.r * (2*m1 + m2 - m2*math.Cos(2*a1-2*a2))
	angle2Dot := system / normalization

	// s.pendulum1.Update(dt)

	// x, y := s.pendulum1.getEnd()
	// s.pendulum2.px = x
	// s.pendulum2.py = y
	// s.pendulum2.Update(dt)

	// s.pendulum1.angleVel += (s.gravity*math.Sin(s.pendulum1.angle) - s.pendulum2.angleVel*math.Sin(s.pendulum1.angle-s.pendulum2.angle)) / s.pendulum2.angle * dt
	// s.pendulum2.angleVel += (s.gravity*math.Sin(s.pendulum2.angle) - s.pendulum1.angleVel*math.Sin(s.pendulum2.angle-s.pendulum1.angle)) / s.pendulum1.angle * dt
}

type pendulum struct {
	px, py   float64
	angle    float64
	angleVel float64
	r        float64
	m        float64
	color    color.NRGBA
}

func NewPendulum(x, y, angle, angleVel float64) *pendulum {
	return &pendulum{
		px:       x,
		py:       y,
		angle:    angle,
		angleVel: angleVel,
		r:        200,
		m:        1,
		color:    ui.Blue,
	}
}

func (p *pendulum) getEnd() (x, y float64) {
	x = float64(p.r)*math.Cos(p.angle-math.Pi/2) + p.px
	y = float64(p.r)*math.Sin(p.angle-math.Pi/2) + p.py
	return x, y
}

func (p *pendulum) Draw(ops *op.Ops) {
	x, y := p.getEnd()
	ui.Line(ops, f32.Pt(float32(p.px), float32(p.py)), f32.Pt(float32(x), float32(y)), 2, p.color)
}

func (p *pendulum) Update(dt float64) {
	acc := math.Sin(p.angle)
	p.angleVel += acc * dt
	p.angle += p.angleVel * dt

	// p.px = 200 + 100*math.Cos(p.angle)
	// p.py = 200 + 100*math.Sin(p.angle)
}
