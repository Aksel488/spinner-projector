package models

import (
	"fmt"
	"image"
	"image/color"
	"spinner-projector/ui"

	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

type Application struct {
	menu        []ManuItem
	theme       *material.Theme
	splitVisual SplitVisual
	left, right bool
}

func NewApplication() *Application {
	menu := []ManuItem{
		{
			Name:     "File 1",
			Selected: true,
			Content:  NewColorBox(ui.Red),
			btn:      &widget.Clickable{},
		},
		{
			Name:     "blue",
			Selected: false,
			Content:  NewBalls(20),
			btn:      &widget.Clickable{},
		},
		{
			Name:     "example",
			Selected: false,
			Content:  NewBalls(40),
			btn:      &widget.Clickable{},
		},
		{
			Name:     "example",
			Selected: false,
			Content:  NewBalls(2),
			btn:      &widget.Clickable{},
		},
	}

	return &Application{
		menu:        menu,
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

	menuWidget := layout.Rigid(func(gtx layout.Context) layout.Dimensions {
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
		})

		btnHolder := layout.Dimensions{Size: image.Pt(200, gtx.Constraints.Max.Y)}
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
	return layout.Flex{}.Layout(gtx, menuWidget, contentWidget)
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
