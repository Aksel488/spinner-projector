package pendulumsystem

import (
	"image"
	"image/color"
	"math"
	"spinner-projector/models"
	"spinner-projector/ui"

	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op"
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
			Name: "Trail",
			Btn:  &widget.Clickable{},
		},
		{
			Name: "Reset",
			Btn:  &widget.Clickable{},
		},
		{
			Name: "Gravity",
			Btn:  ui.NewLinearSlider(9.81, 1, 20),
		},
		{
			Name: "Angle 1",
			Btn:  ui.NewLinearSlider(float32(offset1), -math.Pi, math.Pi),
		},
		{
			Name: "Angle 2",
			Btn:  ui.NewLinearSlider(float32(offset2), -math.Pi, math.Pi),
		},
		{
			Name: "log test",
			Btn:  ui.NewLogSlider(0.001, 0.00001, 0.1),
		},
		{
			Name: "log test 2",
			Btn:  ui.NewLogSlider(500, 1, 10000),
		},
		{
			Name: "log test 3",
			Btn:  ui.NewLogSlider(1, 0.0001, 1000),
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

func (s *DoublePendulumSystem) Menu(gtx layout.Context, theme *material.Theme) layout.Dimensions {
	btnList := layout.List{Axis: layout.Vertical, Alignment: layout.Baseline}

	return layout.Flex{
		Axis:      layout.Vertical,
		Alignment: layout.Middle,
	}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return btnList.Layout(gtx, len(s.menu), func(gtx layout.Context, i int) layout.Dimensions {
				menuBtn := s.menu[i]
				switch menuBtn.Btn.(type) {
				case *widget.Clickable:
					clickable := menuBtn.Btn.(*widget.Clickable)
					btn := material.Button(theme, clickable, menuBtn.Name)

					if clickable.Clicked(gtx) {
						s.handleBtnClick(menuBtn.Name)
					}

					return ui.DefaultButton(gtx, &btn)
				case ui.Slider:
					slider := menuBtn.Btn.(ui.Slider)

					if menuBtn.Name == "Gravity" {
						s.g = float64(slider.Value())
					} else if menuBtn.Name == "Angle 1" {
						s.inputs.offset1 = float64(slider.Value())
						if slider.GetFloat().Dragging() {
							s.Init()
						}
					} else if menuBtn.Name == "Angle 2" {
						s.inputs.offset2 = float64(slider.Value())
						if slider.GetFloat().Dragging() {
							s.Init()
						}
					}

					return slider.Layout(gtx, theme, menuBtn.Name)
				default:
					btnText := material.Body1(theme, menuBtn.Name+" (NI)")
					btnText.Color = color.NRGBA{R: 200, A: 255}
					dims := btnText.Layout(gtx)
					return dims
				}
			})
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			// pendulum 1 settings
			return layout.Inset{
				Top: 20,
			}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{
					Axis:      layout.Vertical,
					Alignment: layout.Middle,
				}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						label := material.H6(theme, "Inner line")
						label.Color = ui.White
						return label.Layout(gtx)
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return s.p1.Menu(gtx, theme)
					}),
				)
			})
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			// pendulum 2 settings
			return layout.Inset{
				Top: 20,
			}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{
					Axis:      layout.Vertical,
					Alignment: layout.Middle,
				}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						label := material.H6(theme, "Outer line")
						label.Color = ui.White
						return label.Layout(gtx)
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return s.p2.Menu(gtx, theme)
					}),
				)
			})
		}),
	)
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
	l1 := 1.0
	l2 := 1.0

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
	menu     []models.ControlMenuItem
}

func NewPendulum(x, y, angle float64, color color.NRGBA) *pendulum {
	pendulum := &pendulum{
		px:       x,
		py:       y,
		angle:    angle,
		angleVel: 0,
		r:        40,
		m:        5,
		color:    color,
	}

	pendulum.menu = []models.ControlMenuItem{
		{
			Name: "Length",
			Btn:  ui.NewLinearSlider(40, 10, 100),
		},
		{
			Name: "Mass",
			Btn:  ui.NewLinearSlider(5, 1, 10),
		},
	}

	return pendulum
}

func (s *pendulum) Menu(gtx layout.Context, theme *material.Theme) layout.Dimensions {
	btnList := layout.List{Axis: layout.Vertical, Alignment: layout.Baseline}
	return btnList.Layout(gtx, len(s.menu), func(gtx layout.Context, i int) layout.Dimensions {
		menuBtn := s.menu[i]
		switch menuBtn.Btn.(type) {
		case ui.Slider:
			slider := menuBtn.Btn.(ui.Slider)

			if menuBtn.Name == "Length" {
				s.r = float64(slider.Value())
			} else if menuBtn.Name == "Mass" {
				s.m = float64(slider.Value())
			}

			return slider.Layout(gtx, theme, menuBtn.Name)
		default:
			btnText := material.Body1(theme, menuBtn.Name+" (NI)")
			btnText.Color = color.NRGBA{R: 200, A: 255}
			dims := btnText.Layout(gtx)
			return dims
		}
	})
}

func (p *pendulum) getEnd() (x, y float64) {
	length := p.r * p.m
	x = length*math.Cos(p.angle+math.Pi/2) + p.px
	y = length*math.Sin(p.angle+math.Pi/2) + p.py
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
