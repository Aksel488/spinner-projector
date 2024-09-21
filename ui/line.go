package ui

import (
	"image/color"

	"gioui.org/f32"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
)

func Line(ops *op.Ops, start f32.Point, end f32.Point, width float32, color color.NRGBA) {
	var path clip.Path
	path.Begin(ops)
	path.MoveTo(start)
	path.LineTo(end)
	path.Close()
	line := clip.Stroke{
		Path:  path.End(),
		Width: width,
	}
	paint.FillShape(ops, color, line.Op())
}
