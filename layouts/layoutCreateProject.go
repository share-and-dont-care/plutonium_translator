package layouts

import (
	"os"
	"path/filepath"

	"gioui.org/app"
	"gioui.org/layout"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"gioui.org/x/explorer"
)

type LayoutCreateProject struct {
	dbPath           string
	savePath         string
	dbPathEdit       widget.Editor
	savePathEdit     widget.Editor
	selectDBButton   widget.Clickable
	selectSaveButton widget.Clickable
	createBtn        widget.Clickable
	explorer         *explorer.Explorer
	selectDBCh       chan string
	selectSaveCh     chan string
	window           *app.Window
}

func (w *LayoutCreateProject) Init(window *app.Window) {
	w.window = window
	w.explorer = explorer.NewExplorer(window)
	w.selectDBCh = make(chan string, 1)
	w.selectSaveCh = make(chan string, 1)
}

func (w *LayoutCreateProject) FrameEventHandler(theme *material.Theme, gtx layout.Context) LayoutWindow {
	var layoutResult LayoutWindow = w

	select {
	case path := <-w.selectDBCh:
		w.dbPath = filepath.Dir(path)
		w.dbPathEdit.SetText(w.dbPath)
	case path := <-w.selectSaveCh:
		if path != "" {
			w.savePath = path
			w.savePathEdit.SetText(w.savePath)
		}
	default:
	}

	layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{
			Axis:    layout.Vertical,
			Spacing: layout.SpaceEvenly,
		}.Layout(gtx,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
					layout.Flexed(0.7, func(gtx layout.Context) layout.Dimensions {
						return material.Editor(theme, &w.dbPathEdit, "Select plutonium data directory (select any file in data directory)").Layout(gtx)
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						btn := material.Button(theme, &w.selectDBButton, "...")
						if w.selectDBButton.Clicked(gtx) {
							w.onSelectDB()
						}
						return btn.Layout(gtx)
					}),
				)
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
					layout.Flexed(0.7, func(gtx layout.Context) layout.Dimensions {
						return material.Editor(theme, &w.savePathEdit, "Save new project to ...").Layout(gtx)
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						btn := material.Button(theme, &w.selectSaveButton, "...")
						if w.selectSaveButton.Clicked(gtx) {
							w.onSelectSave()
						}
						return btn.Layout(gtx)
					}),
				)
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				btn := material.Button(theme, &w.createBtn, "Create Project")
				if w.createBtn.Clicked(gtx) {
					layoutProject := &LayoutProject{}
					layoutProject.Init(w.window)
					layoutProject.projectPath = w.savePath
					layoutResult = layoutProject
				}
				return btn.Layout(gtx)
			}),
		)
	})
	return layoutResult
}

func (w *LayoutCreateProject) onSelectDB() {
	go func() {
		reader, err := w.explorer.ChooseFile()
		if err == nil && reader != nil {
			defer reader.Close()
			if f, ok := reader.(*os.File); ok {
				w.selectDBCh <- f.Name()
			}
		}
	}()
}

func (w *LayoutCreateProject) onSelectSave() {
	go func() {
		reader, err := w.explorer.ChooseFile()
		if err == nil && reader != nil {
			defer reader.Close()
			if f, ok := reader.(*os.File); ok {
				w.selectSaveCh <- f.Name()
			}
		}
	}()
}
