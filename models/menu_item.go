package models

import (
	"image"

	"gioui.org/layout"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

type Content interface {
	Update(gtx layout.Context, dt float64)
	Draw(gtx layout.Context, size image.Point) layout.Dimensions
	Menu(gtx layout.Context, theme *material.Theme) layout.Dimensions
}

type MainMenuItem struct {
	Btn     *widget.Clickable
	Name    string
	Content Content
}

type ControlMenuItem struct {
	Btn  interface{}
	Name string
}
