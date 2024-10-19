package pendulumsystem

import (
	"fmt"
	"image"
	"image/color"
	"math"
	"spinner-projector/models"
	"spinner-projector/ui"

	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

type DoublePendulumSystem struct {
	p1        *pendulum
	p2        *pendulum
	drawTrail bool
	trail     [][]int
	g         float64
	menu      []models.ControlMenuItem
	inputs    inputParams
}

type inputParams struct {
	offset1, offset2 float64
	color            color.NRGBA
}

func NewDoublePendulumSystem(offset1, offset2 float64, color color.NRGBA) *DoublePendulumSystem {
	inputs := inputParams{offset1, offset2, color}

	pendulumSystem := &DoublePendulumSystem{}
	pendulumSystem.inputs = inputs
	pendulumSystem.Init()
	pendulumSystem.drawTrail = false
	pendulumSystem.g = 9.81

	pendulumSystem.menu = []models.ControlMenuItem{
		{
			Btn:  &widget.Clickable{},
			Name: "Trail",
		},
		{
			Btn:  &widget.Clickable{},
			Name: "Reset",
		},
		{
			Btn:  &widget.Float{},
			Name: "Gravity",
		},
	}

	return pendulumSystem
}

func (s *DoublePendulumSystem) handleBtnClick(btnName string) {
	switch btnName {
	case "Trail":
		s.drawTrail = !s.drawTrail
		s.trail = nil
	case "Reset":
		s.Init()
	}
}

func (s *DoublePendulumSystem) Menu(gtx layout.Context, theme *material.Theme) {
	btnList := layout.List{Axis: layout.Vertical, Alignment: layout.Baseline}
	btnList.Layout(gtx, len(s.menu), func(gtx layout.Context, i int) layout.Dimensions {
		menuBtn := s.menu[i]
		switch menuBtn.Btn.(type) {
		case *widget.Clickable:
			clickable := menuBtn.Btn.(*widget.Clickable)
			btn := material.Button(theme, clickable, menuBtn.Name)

			if clickable.Clicked(gtx) {
				s.handleBtnClick(menuBtn.Name)
			}

			return ui.DefaultButton(gtx, &btn)
		case *widget.Float:
			slider := menuBtn.Btn.(*widget.Float)

			if menuBtn.Name == "Gravity" {
				s.g = float64(slider.Value) * 9.81
			}

			// sliderStyle := material.Slider(theme, slider)
			// return ui.DefaultSlider(gtx, &sliderStyle)

			return layout.Flex{
				Axis:      layout.Vertical,
				Alignment: layout.Middle,
			}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return layout.Inset{
						Top:   unit.Dp(20),
						Left:  unit.Dp(20),
						Right: unit.Dp(20),
					}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						gtx.Constraints.Max.X = 200
						// Slider Layout
						sliderStyle := material.Slider(theme, slider).Layout(gtx)
						// sliderStyle.Baseline = 20
						return sliderStyle
					})
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					// Show current slider value
					label := material.H5(theme, "Value: "+fmt.Sprintf("%.2f", slider.Value*9.81))
					label.Color = ui.White
					return label.Layout(gtx)
				}),
			)

		default:
			btnText := material.H5(theme, menuBtn.Name)
			btnText.Color = color.NRGBA{G: 200, A: 255}
			// fpsDisplayText.Alignment = text.Middle
			dims := btnText.Layout(gtx)
			dims.Size = image.Pt(150, 50)
			return dims
		}
	})
}

func (s *DoublePendulumSystem) Init() {
	pendulum1 := NewPendulum(500, 400, -math.Pi+s.inputs.offset1, s.inputs.color)
	x, y := pendulum1.getEnd()
	pendulum2 := NewPendulum(x, y, -math.Pi+s.inputs.offset2, s.inputs.color)

	s.p1 = pendulum1
	s.p2 = pendulum2
}

func (s *DoublePendulumSystem) Draw(gtx layout.Context, size image.Point) layout.Dimensions {
	s.p1.Draw(gtx.Ops)
	s.p2.Draw(gtx.Ops)

	if s.drawTrail {
		for i, p := range s.trail {
			if i+2 < len(s.trail) {
				np := s.trail[i+1]
				c := color.NRGBA{R: 0, G: 0, B: 0, A: uint8(i)}
				ui.Line(gtx.Ops, f32.Pt(float32(p[0]), float32(p[1])), f32.Pt(float32(np[0]), float32(np[1])), 1, c)
			}
		}
	}

	return layout.Dimensions{Size: size}
}

func (s *DoublePendulumSystem) Update(gtx layout.Context, dt float64) {
	s.p1.px = float64(gtx.Constraints.Max.X / 2)
	s.p1.py = float64(gtx.Constraints.Max.Y / 2)

	a1 := s.p1.angle
	a2 := s.p2.angle
	av1 := s.p1.angleVel
	av2 := s.p2.angleVel
	m1 := s.p1.m
	m2 := s.p2.m
	l1 := 10.0
	l2 := 10.0

	mass1 := -s.g * (2*m1 + m2) * math.Sin(a1)
	mass2 := -m2 * s.g * math.Sin(a1-2*a2)

	interaction := -2 * math.Sin(a1-a2) * m2 * math.Cos(math.Pow(av2, 2)*l2+math.Pow(av1, 2)*l1*math.Cos(a1-a2))
	normalization := l1 * (2*m1 + m2 - m2*math.Cos(2*a1-2*a2))
	angle1Ddot := (mass1 + mass2 + interaction) / normalization

	system := 2 * math.Sin(a1-a2) * (math.Pow(av1, 2)*l1*(m1+m2) + s.g*(m1+m2)*math.Cos(a1) + math.Pow(av2, 2)*l2*m2*math.Cos(a1-a2))
	normalization = l1 * (2*m1 + m2 - m2*math.Cos(2*a1-2*a2))
	angle2Ddot := system / normalization

	s.p1.angleVel += angle1Ddot * dt
	s.p1.angle += s.p1.angleVel * dt
	s.p2.angleVel += angle2Ddot * dt
	s.p2.angle += s.p2.angleVel * dt

	// fmt.Println("angle1", s.p1.angle, "angle2", s.p2.angle)

	x, y := s.p1.getEnd()
	s.p2.px = x
	s.p2.py = y

	if s.drawTrail {
		endX, endY := s.p2.getEnd()
		point := []int{int(endX), int(endY)}
		s.trail = append(s.trail, point)
		if len(s.trail) > 200 {
			s.trail = s.trail[1:]
		}
	}
}

type pendulum struct {
	px, py   float64
	angle    float64
	angleVel float64
	r        float64
	m        float64
	color    color.NRGBA
}

func NewPendulum(x, y, angle float64, color color.NRGBA) *pendulum {
	return &pendulum{
		px:       x,
		py:       y,
		angle:    angle,
		angleVel: 0,
		r:        200,
		m:        1,
		color:    color,
	}
}

func (p *pendulum) getEnd() (x, y float64) {
	x = p.r*math.Cos(p.angle+math.Pi/2) + p.px
	y = p.r*math.Sin(p.angle+math.Pi/2) + p.py
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
