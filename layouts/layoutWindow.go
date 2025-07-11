package layouts

import (
	"gioui.org/layout"
	"gioui.org/widget/material"
	"gioui.org/app"
)

type LayoutWindow interface {
	Init(window *app.Window)
	FrameEventHandler(theme *material.Theme, gtx layout.Context) LayoutWindow
}
