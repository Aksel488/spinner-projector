package main

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"os"
	"spinner-projector/models"
	"time"

	"gioui.org/app"
	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/widget/material"
)

func main() {
	go func() {
		window := new(app.Window)

		err := run(window)
		if err != nil {
			log.Fatal(err)
		}

		os.Exit(0)
	}()

	app.Main()
}

func run(window *app.Window) error {
	theme := material.NewTheme()
	windowState := NewState()
	balls := models.NewBalls(50)

	var ops op.Ops
	for {
		switch e := window.Event().(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)
			dt := time.Since(windowState.frameDraw).Seconds()
			// update functions
			windowState.update()
			balls.Update(dt, gtx.Constraints.Max.X, gtx.Constraints.Max.Y)

			// draw functions
			draw(gtx.Ops)
			balls.Draw(gtx.Ops)
			windowState.draw(gtx, theme)

			// draw the window and trigger redraw
			e.Frame(gtx.Ops)
			window.Invalidate()
		}
	}
}

func welcomeText(gtx layout.Context, theme *material.Theme) {
	title := material.H1(theme, "Hi, I'm Giggles")
	maroon := color.NRGBA{R: 127, G: 0, B: 0, A: 255}
	title.Color = maroon
	title.Alignment = text.Middle
	title.Layout(gtx)
}

func draw(ops *op.Ops) {
	// white := color.NRGBA{R: 255, G: 255, B: 255, A: 255}
	purple := color.NRGBA{R: 54, G: 1, B: 64, A: 255}
	paint.ColorOp{Color: purple}.Add(ops)
	paint.PaintOp{}.Add(ops)

	// offsetRect(ops)
	// line(ops, f32.Point{X: 0, Y: 0}, f32.Point{X: 400, Y: 200}, 4, color.NRGBA{R: 0, G: 0, B: 255, A: 255})
	// strokeTriangle(ops)
	// redButtonBackground(ops)
}

func rect(ops *op.Ops) {
	defer clip.Rect{Max: image.Pt(200, 100)}.Push(ops).Pop()
	paint.ColorOp{Color: color.NRGBA{R: 255, G: 0, B: 0, A: 255}}.Add(ops)
	paint.PaintOp{}.Add(ops)
}

func offsetRect(ops *op.Ops) {
	defer op.Offset(image.Pt(100, 20)).Push(ops).Pop()
	redButtonBackground(ops)
}

func redButtonBackground(ops *op.Ops) {
	const r = 10 // roundness
	bounds := image.Rect(0, 0, 100, 100)
	defer clip.RRect{Rect: bounds, SE: r, SW: r, NW: r, NE: r}.Push(ops).Pop()
	rect(ops)
}

func line(ops *op.Ops, start f32.Point, end f32.Point, width float32, color color.NRGBA) {
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

func strokeTriangle(ops *op.Ops) {
	var path clip.Path
	path.Begin(ops)
	path.MoveTo(f32.Pt(30, 30))
	path.LineTo(f32.Pt(70, 30))
	path.LineTo(f32.Pt(50, 70))
	path.Close()

	green := color.NRGBA{R: 0, G: 255, B: 0, A: 255}

	paint.FillShape(ops, green,
		clip.Stroke{
			Path:  path.End(),
			Width: 4,
		}.Op())
}

type state struct {
	frameDraw     time.Time
	fpsTimer      time.Time
	countedFrames int
	fpsText       string
}

func NewState() state {
	return state{
		frameDraw:     time.Now(),
		fpsTimer:      time.Now(),
		countedFrames: 0,
		fpsText:       "",
	}
}

func (s *state) update() {
	secounds := time.Since(s.fpsTimer).Seconds()
	if secounds > 0.5 {
		avgFPS := float64(s.countedFrames) / secounds
		if avgFPS > 2000000 {
			avgFPS = 0
		}

		s.fpsText = fmt.Sprintf("%.0f FPS", avgFPS)
		s.fpsTimer = time.Now()
		s.countedFrames = 0
	}
	s.countedFrames++
	s.frameDraw = time.Now()
}

func (s *state) draw(gtx layout.Context, theme *material.Theme) {
	fpsText := s.fpsText
	if fpsText == "" {
		fpsText = "FPS: --.--"
	}

	fpsDisplayText := material.H5(theme, fpsText)
	fpsDisplayText.Color = color.NRGBA{G: 200, A: 255}
	// fpsDisplayText.Alignment = text.Middle
	fpsDisplayText.Alignment = text.End
	fpsDisplayText.Layout(gtx)
}
