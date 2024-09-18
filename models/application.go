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

type Application struct {
	menu []ManuItem
}

func NewApplication() *Application {

	menu := []ManuItem{
		{
			Name:     "File",
			Selected: true,
		},
	}

	return &Application{
		menu: menu,
	}
}

func (application *Application) Draw(gtx layout.Context, dt float64) layout.Dimensions {
	// Layout the left side (menu)

	menuWidget := layout.Rigid(func(gtx layout.Context) layout.Dimensions {
		buttin := Button{}
		return buttin.ClickColorBox(gtx, image.Pt(200, gtx.Constraints.Max.Y), ui.Red)
	})

	// Layout the right side (content)
	contentWidget := layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
		buttin := Button{}
		return buttin.ClickColorBox(gtx, image.Pt(gtx.Constraints.Max.X, gtx.Constraints.Max.Y), ui.Blue)
	})

	// leftSide := layout.Flex{
	// 	Axis:      layout.Vertical,
	// 	Alignment: layout.Middle,
	// }.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
	// 	return layout.List{Axis: layout.Vertical}.Layout(gtx, len(application.menu), func(gtx layout.Context, i int) layout.Dimensions {
	// 		item := &application.menu[i]
	// 		btn := material.Button(theme, item.Name).Layout(gtx)
	// 		if btn.Clicked() {
	// 			// Handle menu item selection
	// 		}
	// 		return btn
	// 	})
	// })

	// // Layout the right side (content)
	// rightSide := layout.Flex{
	// 	Axis:      layout.Vertical,
	// 	Alignment: layout.Middle,
	// }.Layout(ctx, func(gtx layout.Context) layout.Dimensions {
	// 	selectedItem := application.getSelectedMenuItem()
	// 	if selectedItem != nil {
	// 		return selectedItem.Content.Layout(gtx)
	// 	}
	// 	return layout.Dimensions{}
	// })

	// Combine left and right sides
	return layout.Flex{
		// Axis:      layout.Horizontal,
		// Alignment: layout.Middle,
	}.Layout(gtx, menuWidget, contentWidget)
}

type Button struct {
	pressed bool
}

func (b *Button) ClickColorBox(gtx layout.Context, size image.Point, color color.NRGBA) layout.Dimensions {
	area := clip.Rect{Max: size}.Push(gtx.Ops)
	event.Op(gtx.Ops, b)

	for {
		ev, ok := gtx.Event(pointer.Filter{
			Target: b,
			Kinds:  pointer.Press | pointer.Release,
		})
		fmt.Println("click event", ok)
		if !ok {
			break
		}

		e, ok := ev.(pointer.Event)
		if !ok {
			continue
		}

		switch e.Kind {
		case pointer.Press:
			b.pressed = true
		case pointer.Release:
			b.pressed = false
		}

	}

	area.Pop()

	col := color
	if b.pressed {
		col = ui.Green
	}

	defer clip.Rect{Max: size}.Push(gtx.Ops).Pop()
	paint.ColorOp{Color: col}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)
	return layout.Dimensions{Size: size}
}
