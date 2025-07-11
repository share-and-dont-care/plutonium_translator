package main

import (
	"log"
	"os"

	"gioui.org/app"
)

type ProjectOpenedState struct {
	ProjectPath string
}

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
