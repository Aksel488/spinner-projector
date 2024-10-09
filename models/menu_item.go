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
	Menu(gtx layout.Context, theme *material.Theme)
}

type ManuItem struct {
	btn      *widget.Clickable
	Name     string
	Selected bool
	Content  Content
}
