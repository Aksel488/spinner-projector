package layout

import (
	"image"

	"gioui.org/layout"
	"gioui.org/op"
)

type SplitVisual struct{}

func (s SplitVisual) Layout(gtx layout.Context, left, rigth layout.Widget) layout.Dimensions {
	leftSize := gtx.Constraints.Max.X / 2
	rigthSize := gtx.Constraints.Max.X - leftSize

	{
		gtx := gtx
		gtx.Constraints = layout.Exact(image.Pt(leftSize, gtx.Constraints.Max.Y))
		left(gtx)
	}

	{
		gtx := gtx
		gtx.Constraints = layout.Exact(image.Pt(rigthSize, gtx.Constraints.Max.Y))
		trans := op.Offset(image.Pt(leftSize, 0)).Push(gtx.Ops)
		rigth(gtx)
		trans.Pop()
	}

	return layout.Dimensions{Size: gtx.Constraints.Max}
}

// func exampleSplitVisual(gtx layout.Context, th *material.Theme) layout.Dimensions {
// 	return SplitVisual{}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
// 		return FillWithLabel(gtx, th, "Left", color.NRGBA{R: 255, G: 0, B: 0, A: 255})
// 	}, func(gtx layout.Context) layout.Dimensions {
// 		return FillWithLabel(gtx, th, "Right", color.NRGBA{R: 0, G: 0, B: 255, A: 255})
// 	})
// }

// func FillWithLabel(gtx layout.Context, th *material.Theme, text string, backgroundColor color.NRGBA) layout.Dimensions {
// 	ColorBox(gtx, gtx.Constraints.Max, backgroundColor)
// 	return layout.Center.Layout(gtx, material.H3(th, text).Layout)
// }

// func ColorBox(gtx layout.Context, size image.Point, backgroundColor color.NRGBA) layout.Dimensions {
// 	// Set the background color
// 	colorOp := paint.ColorOp{
// 		Color: backgroundColor,
// 	}
// 	colorOp.Add(gtx.Ops)

// 	// Draw a rectangle covering the entire area
// 	image.Rectangle{
// 		Size:   size,
// 		Color:  backgroundColor,
// 		Border: nil,
// 	}
// 	dims := rect.Layout(gtx)

// 	return dims
// }
