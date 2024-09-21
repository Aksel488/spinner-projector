package models

import (
	"image"
	"image/color"
	"math/rand/v2"

	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
)

// /////////// BALLS /////////////
type balls struct {
	Balls []*ball
}

func NewBalls(n int) balls {
	var ballList []*ball
	for range n {
		ballList = append(ballList, NewBall())
	}

	b := balls{
		Balls: ballList,
	}
	return b
}

func (b *balls) AddBall() {
	b.Balls = append(b.Balls, NewBall())
}

func (b *balls) RemoveBall() {
	if len(b.Balls) == 0 {
		return
	}
	b.Balls = b.Balls[:len(b.Balls)-1]
}

func (b *balls) Update(dt float64, maxX, maxY int) {
	for _, ball := range b.Balls {
		ball.Update(dt, maxX, maxY)
	}
}

func (b *balls) Draw(ops *op.Ops) {
	for _, ball := range b.Balls {
		ball.Draw(ops)
	}
}

// /////////// BALL /////////////
type ball struct {
	px, py float64
	r      int
	color  color.NRGBA
	vx, vy float64
}

func NewBall() *ball {
	return &ball{
		px:    float64(rand.IntN(500)),
		py:    float64(rand.IntN(500)),
		r:     rand.IntN(10) + 10,
		color: color.NRGBA{R: uint8(rand.IntN(255)), G: uint8(rand.IntN(255)), B: uint8(rand.IntN(255)), A: 255},
		vx:    float64(rand.IntN(300)),
		vy:    float64(rand.IntN(300)),
	}
}

func (b *ball) Update(dt float64, maxX, maxY int) {
	ay := float64(9.81)
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
