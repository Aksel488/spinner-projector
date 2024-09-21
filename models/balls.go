package models

import (
	"fmt"
	"image"
	"image/color"
	"math/rand/v2"
	"spinner-projector/ui"

	"gioui.org/f32"
	"gioui.org/io/event"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
)

// /////////// BALLS /////////////
type balls struct {
	Balls          []*ball
	isClicked      *bool
	clickX, clickY *int
	mouseX, mouseY *int
}

func NewBalls(n int) *balls {
	var ballList []*ball
	for range n {
		ballList = append(ballList, NewRandBall())
	}

	var isClicked bool

	b := &balls{
		Balls:     ballList,
		isClicked: &isClicked,
		clickX:    new(int),
		clickY:    new(int),
		mouseX:    new(int),
		mouseY:    new(int),
	}
	return b
}

func (b *balls) AddBall(x, y, vx, vy float64) {
	b.Balls = append(b.Balls, NewBall(x, y, vx, vy))
}

func (b *balls) RemoveBall(other *ball) {
	for i, ball := range b.Balls {
		if ball == other {
			b.Balls = append(b.Balls[:i], b.Balls[i+1:]...)
			return
		}
	}
}

func (b *balls) BallsAt(x, y int) []*ball {
	balls := []*ball{}

	for _, ball := range b.Balls {
		dist := Distance(x, y, int(ball.px), int(ball.py))
		fmt.Println(fmt.Sprintf("(%d, %d)", x, y), fmt.Sprintf("(%d, %d)", int(ball.px), int(ball.py)), "dist", dist, "r", ball.r)
		if dist < float64(ball.r) {
			balls = append(balls, ball)
		}
	}

	return balls
}

func (b *balls) Update(gtx layout.Context, dt float64) {
	// handle click
	area := clip.Rect{Max: gtx.Constraints.Max}.Push(gtx.Ops)
	event.Op(gtx.Ops, b.isClicked)

	for {
		ev, ok := gtx.Event(pointer.Filter{
			Target: b.isClicked,
			Kinds:  pointer.Press | pointer.Release | pointer.Drag,
		})
		if !ok {
			break
		}
		// fmt.Println("click event", ev, b.isClicked)

		e, ok := ev.(pointer.Event)
		if !ok {
			continue
		}

		switch e.Kind {
		case pointer.Press:
			*b.clickX = int(e.Position.X)
			*b.clickY = int(e.Position.Y)
			*b.mouseX = int(e.Position.X)
			*b.mouseY = int(e.Position.Y)

			// remove clicked balls
			clickedBalls := b.BallsAt(*b.clickX, *b.clickY)
			fmt.Println(clickedBalls)
			for _, ball := range clickedBalls {
				b.RemoveBall(ball)
			}

			if len(clickedBalls) > 0 {
				*b.isClicked = true
			}
		case pointer.Drag:
			*b.mouseX = int(e.Position.X)
			*b.mouseY = int(e.Position.Y)
		case pointer.Release:
			if *b.isClicked {
				*b.isClicked = false
				dx := (*b.clickX - *b.mouseX) * 2
				dy := (*b.clickY - *b.mouseY) * 2
				b.AddBall(float64(*b.clickX), float64(*b.clickY), float64(dx), float64(dy))
			}
		}
	}

	area.Pop()

	for _, ball := range b.Balls {
		ball.Update(dt, gtx.Constraints.Max.X, gtx.Constraints.Max.Y)
		if ball.ttl <= 0 {
			b.RemoveBall(ball)
		}
	}
}

func (b *balls) Draw(gtx layout.Context, size image.Point) layout.Dimensions {
	for _, ball := range b.Balls {
		ball.Draw(gtx.Ops)
	}

	// draw line from clickX, clickY to cursor
	if *b.isClicked {
		start := f32.Point{X: float32(*b.clickX), Y: float32(*b.clickY)}
		dx := *b.clickX - *b.mouseX
		dy := *b.clickY - *b.mouseY
		end := f32.Point{X: float32(*b.clickX + dx), Y: float32(*b.clickY + dy)}
		ui.Line(gtx.Ops, start, end, 2, ui.Red)
	}

	return layout.Dimensions{Size: size}
}

// /////////// BALL /////////////
type ball struct {
	px, py float64
	r      int
	color  color.NRGBA
	vx, vy float64
	ay     float64
	ttl    float64 // seconds to be removed
}

func NewRandBall() *ball {
	return &ball{
		px:    float64(rand.IntN(500)),
		py:    float64(rand.IntN(500)),
		r:     rand.IntN(10) + 10,
		color: color.NRGBA{R: uint8(rand.IntN(255)), G: uint8(rand.IntN(255)), B: uint8(rand.IntN(255)), A: 255},
		vx:    float64(rand.IntN(300)),
		vy:    float64(rand.IntN(300)),
		ay:    float64(rand.Float64()*5 + 5),
		ttl:   60,
	}
}

func NewBall(x, y, vx, vy float64) *ball {
	return &ball{
		px:    x,
		py:    y,
		r:     rand.IntN(10) + 10,
		color: color.NRGBA{R: uint8(rand.IntN(255)), G: uint8(rand.IntN(255)), B: uint8(rand.IntN(255)), A: 255},
		vx:    vx,
		vy:    vy,
		ay:    float64(rand.Float64()*5 + 5),
		ttl:   float64(rand.IntN(10) + 10),
	}
}

func (b *ball) Update(dt float64, maxX, maxY int) {
	// decrese ttl
	b.ttl -= dt

	// ay := float64(9.81)
	ay := float64(2.81)
	// ay := b.ay
	b.vy += ay

	dx := b.px + (b.vx * dt)
	dy := b.py + (b.vy * dt)

	b.px = dx
	b.py = dy

	if int(b.px)+b.r > maxX || int(b.px)-b.r < 0 {
		b.vx *= -1
	}
	if int(b.py)+b.r > maxY || int(b.py)-b.r < 0 {
		b.vy *= -1
	}

	if int(b.px)+b.r > maxX {
		b.px = float64(maxX - b.r)
	} else if int(b.px)-b.r < 0 {
		b.px = float64(b.r)
	}
	if int(b.py)+b.r > maxY {
		b.py = float64(maxY - b.r)
	} else if int(b.py)-b.r < 0 {
		b.py = float64(b.r)
	}
}

func (b *ball) Draw(ops *op.Ops) {
	bounds := image.Rect(int(b.px)-b.r, int(b.py)-b.r, int(b.px)+b.r, int(b.py)+b.r)
	defer clip.RRect{Rect: bounds, SE: b.r, SW: b.r, NW: b.r, NE: b.r}.Push(ops).Pop()
	paint.ColorOp{Color: b.color}.Add(ops)
	paint.PaintOp{}.Add(ops)
}
