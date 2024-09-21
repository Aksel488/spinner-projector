package models

import (
	"image"
	"spinner-projector/ui"

	"gioui.org/layout"
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
			Name:     "File 2",
			Selected: false,
			Content:  NewColorBox(ui.Blue),
			btn:      &widget.Clickable{},
		},
		{
			Name:     "File 3",
			Selected: true,
			Content:  NewColorBox(ui.Black),
			btn:      &widget.Clickable{},
		},
	}

	return &Application{
		menu:        menu,
		theme:       material.NewTheme(),
		splitVisual: SplitVisual{},
	}
}

func (application *Application) Draw(gtx layout.Context, dt float64) layout.Dimensions {
	menuWidget := layout.Rigid(func(gtx layout.Context) layout.Dimensions {
		return application.menu[0].Content.Draw(gtx, image.Pt(200, gtx.Constraints.Max.Y))
	})

	// menuWidget := layout.Rigid(func(gtx layout.Context) layout.Dimensions {
	// 	leftSide := layout.Flex{
	// 		Axis:      layout.Vertical,
	// 		Alignment: layout.Middle,
	// 	}.Layout(gtx, func(gtx layout.Context) layout.FlexChild {

	// 		return layout.Flexed(gtx, func(gtx layout.Context) layout.Dimensions {
	// 			asd := layout.List{Axis: layout.Vertical}
	// 			return asd.Layout(gtx, len(application.menu), func(gtx layout.Context, i int) layout.Dimensions {
	// 				item := &application.menu[i]
	// 				btn := material.Button(application.theme, item.btn, item.Name).Layout(gtx)
	// 				if item.btn.Pressed() {
	// 					fmt.Println("cliked happen sad", item.btn)
	// 					item.Selected = true
	// 					for j := range application.menu {
	// 						if i != j {
	// 							application.menu[j].Selected = false
	// 						}
	// 					}
	// 				}
	// 				return btn
	// 			})
	// 		})
	// 		return alks
	// 	})
	// 	leftSide.Size = image.Pt(200, gtx.Constraints.Max.Y)
	// 	return leftSide
	// })

	contentWidget := layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
		var content Content
		for _, appl := range application.menu {
			if appl.Selected {
				content = appl.Content
			}
		}
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
