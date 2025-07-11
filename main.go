package main

import (
	// "fmt"
	"image/color"
	"log"
	"os"

	"gioui.org/app"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/text"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"gioui.org/x/explorer"
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
	var ops op.Ops

	// Persistent state for editor and button
	var editor widget.Editor
	var button widget.Clickable

	// Channel to receive selected file path
	pathCh := make(chan string, 1)

	// Create explorer instance for this window
	exp := explorer.NewExplorer(window)

	for {
		e := window.Event()
		exp.ListenEvents(e) // Listen for every event
		switch e := e.(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)

			// Check for file path from goroutine
			select {
			case path := <-pathCh:
				editor.SetText(path)
			default:
			}

			// Layout: vertical stack
			layout.Flex{
				Axis:    layout.Vertical,
				Spacing: layout.SpaceStart,
			}.Layout(gtx,
				// Title label
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					title := material.H1(theme, "Hello, Gio")
					maroon := color.NRGBA{R: 127, G: 0, B: 0, A: 255}
					title.Color = maroon
					title.Alignment = text.Middle
					return title.Layout(gtx)
				}),
				// Text field
				layout.Rigid(
					func(gtx layout.Context) layout.Dimensions {
						return material.Editor(theme, &editor, "Enter text...").Layout(gtx)
					},
				),
				// Button
				layout.Rigid(
					func(gtx layout.Context) layout.Dimensions {
						if button.Clicked(gtx) {
							go func() {
								reader, err := exp.ChooseFile()
								if err == nil && reader != nil {
									defer reader.Close()
									if f, ok := reader.(*os.File); ok {
										pathCh <- f.Name()
									}
								}
							}()
						}
						return material.Button(theme, &button, "Submit").Layout(gtx)
					},
				),
			)

			e.Frame(gtx.Ops)
		}
	}
}
