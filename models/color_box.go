package models

import (
	"fmt"
	"image"
	"image/color"
	"spinner-projector/ui"

	"gioui.org/io/event"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
)

type ColorBox struct {
	color   color.NRGBA
	clicked *bool
	counter int
}

func NewColorBox(color color.NRGBA) *ColorBox {
	var cliked bool
	return &ColorBox{color: color, clicked: &cliked}
}

func (cb *ColorBox) Draw(gtx layout.Context, size image.Point) layout.Dimensions {
	area := clip.Rect{Max: size}.Push(gtx.Ops)
	event.Op(gtx.Ops, cb.clicked)

	for {
		ev, ok := gtx.Event(pointer.Filter{
			Target: cb.clicked,
			Kinds:  pointer.Press | pointer.Release,
		})
		if !ok {
			break
		}
		fmt.Println("click event", ev, "counter", cb.counter, "size", size)

		e, ok := ev.(pointer.Event)
		if !ok {
			continue
		}

		switch e.Kind {
		case pointer.Press:
			cb.counter++
			*cb.clicked = true
		case pointer.Release:
			*cb.clicked = false
		}
	}

	col := cb.color
	if *cb.clicked {
		col = ui.Green
	}
	area.Pop()

	defer clip.Rect{Max: size}.Push(gtx.Ops).Pop()
	paint.ColorOp{Color: col}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)

	return layout.Dimensions{Size: size}
}

func (cb *ColorBox) Update(gtx layout.Context, dt float64) {}
