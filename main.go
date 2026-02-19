package main

import (
	"embed"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/build
var assets embed.FS

const (
	windowWidth    = 920
	windowHeight   = 900
	windowMinWidth = 920
	minHeight      = 700
	bgR            = 27
	bgG            = 38
	bgB            = 54
)

func main() {
	app := NewApp()

	err := wails.Run(&options.App{
		Title:     "Poll Scraper v1.0.0",
		Width:     windowWidth,
		Height:    windowHeight,
		MinWidth:  windowMinWidth,
		MinHeight: minHeight,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: bgR, G: bgG, B: bgB, A: 1},
		OnStartup:        app.Startup,
		OnShutdown:       app.Shutdown,
		Bind: []interface{}{
			app,
		},
	})
	if err != nil {
		println("Error:", err.Error())
	}
}
