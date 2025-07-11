package layouts

import (
	"os"

	"gioui.org/app"
	"gioui.org/layout"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"gioui.org/x/explorer"
)

type LayoutMain struct {
	createButton    widget.Clickable
	openButton      widget.Clickable
	explorer        *explorer.Explorer
	window          *app.Window
	openProjectChan chan string
}

func (w *LayoutMain) Init(window *app.Window) {
	w.window = window
	w.explorer = explorer.NewExplorer(window)
	w.createButton = widget.Clickable{}
	w.openButton = widget.Clickable{}
	w.openProjectChan = make(chan string, 1)
}

func (w *LayoutMain) FrameEventHandler(theme *material.Theme, gtx layout.Context) LayoutWindow {
	var layoutResult LayoutWindow = w

	select {
	case path := <-w.openProjectChan:
		if path != "" {
			layoutProject := &LayoutProject{}
			layoutProject.Init(w.window)
			layoutProject.projectPath = path
			return layoutProject
		}
	default:
	}

	layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{
			Axis:    layout.Vertical,
			Spacing: layout.SpaceEvenly,
		}.Layout(gtx,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				btn := material.Button(theme, &w.createButton, "Create New Project")
				if w.createButton.Clicked(gtx) {
					layoutCreate := &LayoutCreateProject{}
					layoutCreate.Init(w.window)
					layoutResult = layoutCreate
				}
				return btn.Layout(gtx)
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				btn := material.Button(theme, &w.openButton, "Open Project")
				if w.openButton.Clicked(gtx) {
					w.onOpen()
				}
				return btn.Layout(gtx)
			}),
		)
	})
	return layoutResult
}

func (w *LayoutMain) onOpen() {
	go func() {
		reader, err := w.explorer.ChooseFile(".langproj")
		if err == nil && reader != nil {
			defer reader.Close()
			if f, ok := reader.(*os.File); ok {
				w.openProjectChan <- f.Name()
			}
		}
	}()
}
