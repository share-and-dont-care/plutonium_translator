package layouts

import (
	"gioui.org/app"
	"gioui.org/layout"
	"gioui.org/widget/material"
	"fmt"
)

type LayoutProject struct {
	window *app.Window
	projectPath string
}

func (w *LayoutProject) FrameEventHandler(theme *material.Theme, gtx layout.Context) LayoutWindow {
	layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return material.H4(theme, fmt.Sprintf("Project opened: %s", w.projectPath)).Layout(gtx)
	})
	return w
}

func (w *LayoutProject) Init(window *app.Window) {
	w.window = window
}