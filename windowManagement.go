package main

import (
	"gioui.org/app"
	"gioui.org/op"
	"gioui.org/widget/material"

	"example.com/main/layouts"
)

func run(window *app.Window) error {
	theme := material.NewTheme()
	var ops op.Ops

	var currentWindow layouts.LayoutWindow = &layouts.LayoutMain{}
	currentWindow.Init(window)

	for {
		e := window.Event()
		// exp.ListenEvents(e)
		switch e := e.(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)
			currentWindow = currentWindow.FrameEventHandler(theme, gtx)
			e.Frame(gtx.Ops)
		}
	}
}
