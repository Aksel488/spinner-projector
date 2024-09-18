package models

import "gioui.org/layout"

type ManuItem struct {
	Name     string
	Selected bool
	Content  layout.Widget
}
