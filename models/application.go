package models

import (
	"fmt"
	"image"
	"image/color"
	"spinner-projector/apps/balls"
	"spinner-projector/apps/pendulumsystem"
	"spinner-projector/ui"

	"gioui.org/io/key"
	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

type Application struct {
	menu          []ManuItem
	menuState     string
	escapePressed bool
	theme         *material.Theme
	splitVisual   SplitVisual
}

func NewApplication() *Application {
	menu := []ManuItem{
		{
			Name:    "pendulum",
			Content: pendulumsystem.NewDoublePendulumSystem(1, -1.5, ui.Red),
			btn:     &widget.Clickable{},
		},
		{
			Name:    "spider",
			Content: pendulumsystem.NewMultiPendulumSystem(9),
			btn:     &widget.Clickable{},
		},
		{
			Name:    "Rainbow rgb",
			Content: pendulumsystem.NewRainbowPendulumSystem(500, 0, 0.1, 0.0001, "rgb"),
			btn:     &widget.Clickable{},
		},
		{
			Name:    "Rainbow slow",
			Content: pendulumsystem.NewRainbowPendulumSystem(700, 1, -1, 0.00000001, "rgb"),
			btn:     &widget.Clickable{},
		},
		{
			Name:    "Rainbow hcl",
			Content: pendulumsystem.NewRainbowPendulumSystem(500, 0, 0, 0.01, "hcl"),
			btn:     &widget.Clickable{},
		},
		{
			Name:    "Rainbow slow low",
			Content: pendulumsystem.NewRainbowPendulumSystem(500, 1, -1.5, 0.0000001, "rgb"),
			btn:     &widget.Clickable{},
		},
		{
			Name:    "ball flinger",
			Content: balls.NewBalls(10),
			btn:     &widget.Clickable{},
		},
		{
			Name:    "Red",
			Content: NewColorBox(ui.Red),
			btn:     &widget.Clickable{},
		},
		{
			Name:    "Purple",
			Content: NewColorBox(ui.Purple),
			btn:     &widget.Clickable{},
		},
		{
			Name:    "blue",
			Content: balls.NewBalls(20),
			btn:     &widget.Clickable{},
		},
		{
			Name:    "example 40",
			Content: balls.NewBalls(40),
			btn:     &widget.Clickable{},
		},
		{
			Name:    "example 2",
			Content: balls.NewBalls(2),
			btn:     &widget.Clickable{},
		},
	}

	menu[0].Selected = true

	return &Application{
		menu:        menu,
		menuState:   "main",
		theme:       material.NewTheme(),
		splitVisual: SplitVisual{},
	}
}

func menuWidget(gtx layout.Context) layout.Dimensions {
	// Center the box widget inside both horizontal and vertical layout
	return layout.Flex{
		Axis:      layout.Vertical,
		Alignment: layout.Middle, // Center vertically
	}.Layout(gtx,
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			// box := layout.Dimensions{Size: gtx.Constraints.Max}
			return layout.Flex{
				Axis:      layout.Horizontal,
				Alignment: layout.Middle, // Center horizontally
			}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					// align := layout.Alignment(layout.Middle)
					// return align.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					// 	return boxWidget(gtx)
					// })
					return boxWidget(gtx)
				}),
			)
		}),
	)
}

func boxWidget(gtx layout.Context) layout.Dimensions {
	// Create a vertical list for the boxes

	colors := []color.NRGBA{
		{R: 255, G: 0, B: 0, A: 255},
		{R: 0, G: 255, B: 0, A: 255},
		{R: 0, G: 0, B: 255, A: 255},
	}

	boxList := layout.List{Axis: layout.Vertical}
	return boxList.Layout(gtx, len(colors), func(gtx layout.Context, i int) layout.Dimensions {
		boxSize := unit.Dp(100)

		defer clip.Rect{Max: image.Pt(100, 100)}.Push(gtx.Ops).Pop()
		paint.ColorOp{Color: colors[i]}.Add(gtx.Ops)
		paint.PaintOp{}.Add(gtx.Ops)

		return layout.Dimensions{
			Size: image.Pt(int(boxSize), int(boxSize)),
		}
	})
}

func (application *Application) Draw(gtx layout.Context, dt float64) layout.Dimensions {
	// menuWidget := layout.Rigid(func(gtx layout.Context) layout.Dimensions {
	// 	return application.menu[0].Content.Draw(gtx, image.Pt(200, gtx.Constraints.Max.Y))
	// })

	// menuWidget := layout.Rigid(func(gtx layout.Context) layout.Dimensions {
	// 	return layout.Flex{
	// 		Axis:      layout.Vertical,
	// 		Alignment: layout.Middle, // Center vertically
	// 	}.Layout(gtx,
	// 		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
	// 			return boxWidget(gtx)
	// 		}),
	// 	)
	// })

	// handle key events

	for {
		event, ok := gtx.Event(
			key.Filter{
				Name: key.NameEscape,
			},
		)
		if !ok {
			break
		}
		ev, ok := event.(key.Event)
		if !ok {
			continue
		}

		if !application.escapePressed && ev.State == key.Press {
			if application.menuState == "main" {
				application.menuState = "controls"
			} else if application.menuState == "controls" {
				application.menuState = "hidden"
			} else if application.menuState == "hidden" {
				application.menuState = "main"
			}
		}

		application.escapePressed = ev.State == key.Press
	}

	menuWidget := layout.Rigid(func(gtx layout.Context) layout.Dimensions {
		size := image.Pt(200, gtx.Constraints.Max.Y)
		btnHolder := layout.Dimensions{Size: size}

		defer clip.Rect{Max: size}.Push(gtx.Ops).Pop()
		paint.ColorOp{Color: ui.Purple}.Add(gtx.Ops)
		paint.PaintOp{}.Add(gtx.Ops)

		if application.menuState == "main" {
			btnList := layout.List{Axis: layout.Vertical, Alignment: layout.Baseline}
			btnList.Layout(gtx, len(application.menu), func(gtx layout.Context, i int) layout.Dimensions {
				menuItem := &application.menu[i]
				btn := material.Button(application.theme, menuItem.btn, menuItem.Name)

				if menuItem.btn.Pressed() && !menuItem.Selected {
					fmt.Println("cliked happen", menuItem.btn)
					menuItem.Selected = true
					for j := range application.menu {
						if i != j {
							application.menu[j].Selected = false
						}
					}
				}

				return ui.DefaultButton(gtx, &btn)
			})
		} else if application.menuState == "controls" {
			var selected Content
			for i := range application.menu {
				if application.menu[i].Selected {
					selected = application.menu[i].Content
				}
			}
			selected.Menu(gtx, application.theme)
		}

		return btnHolder
	})

	contentWidget := layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
		var content Content
		for _, appl := range application.menu {
			if appl.Selected {
				content = appl.Content
			}
		}
		content.Update(gtx, dt)
		return content.Draw(gtx, image.Pt(gtx.Constraints.Max.X, gtx.Constraints.Max.Y))
	})

	// Combine left and right sides
	if application.menuState == "hidden" {
		return layout.Flex{}.Layout(gtx, contentWidget)
	} else {
		return layout.Flex{}.Layout(gtx, menuWidget, contentWidget)
	}
}

// return application.splitVisual.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
// 	return FillWithLabel(gtx, application.theme, "left", ui.Red, application.left)
// }, func(gtx layout.Context) layout.Dimensions {
// 	return FillWithLabel(gtx, application.theme, "right", ui.Blue, application.rigth)
// })

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
