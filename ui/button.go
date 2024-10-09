package ui

import (
	"image"
	"image/color"

	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

func DefaultButton(gtx layout.Context, btn *material.ButtonStyle) layout.Dimensions {
	btn.Inset = layout.Inset{
		Top:    10,
		Bottom: 10,
		Left:   40,
		Right:  40,
	}

	border := widget.Border{
		Color:        color.NRGBA{R: 0, G: 0, B: 0, A: 255},
		CornerRadius: unit.Dp(0),
		Width:        unit.Dp(1),
	}

	// Center the button, apply the border and layout the button
	return layout.Flex{
		Axis:      layout.Horizontal,
		Spacing:   layout.SpaceAround,
		Alignment: layout.Middle,
	}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return border.Layout(gtx, func(gtx layout.Context) layout.Dimensions {

				btnLayout := btn.Layout(gtx)
				btnLayout.Size = image.Pt(150, 50)
				return btnLayout
			})
		}),
	)
}
