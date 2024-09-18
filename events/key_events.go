package events

import (
	"fmt"
	"image"

	"gioui.org/io/event"
	"gioui.org/io/key"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op/clip"
)

var tag = new(bool)

// Click event
func PointerEvent(gtx layout.Context) {
	defer clip.Rect{Max: image.Pt(gtx.Constraints.Max.X, gtx.Constraints.Max.Y)}.Push(gtx.Ops).Pop()
	event.Op(gtx.Ops, tag)

	for {
		ev, ok := gtx.Event(pointer.Filter{
			Target: tag,
			Kinds:  pointer.Press | pointer.Release,
		})
		if !ok {
			break
		}

		fmt.Println("click event", ev, "tag", *tag)
	}
}

// this not working
func KeyEvent(gtx layout.Context) {
	defer clip.Rect{Max: image.Pt(gtx.Constraints.Max.X, gtx.Constraints.Max.Y)}.Push(gtx.Ops).Pop()
	event.Op(gtx.Ops, tag)

	for {
		ev, ok := gtx.Event(key.Filter{
			Focus: tag,
			Name:  key.NameSpace,
		})
		if !ok {
			fmt.Println("no key event")
			break
		}

		fmt.Println("key event", ev, "tag", *tag)
	}
}
