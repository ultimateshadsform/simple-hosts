package main

import (
	"embed"
	"log/slog"
	"runtime"

	"github.com/txn2/txeh"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS
var hosts *txeh.Hosts
var filePath string

func main() {
	// Create an instance of the app structure
	app := NewApp()

	host, err := txeh.NewHostsDefault()
	if err != nil {
		slog.Error("[App] Error loading hosts file", "error", err.Error())
	}

	switch runtime.GOOS {
	case "windows":
		filePath = "C:\\Windows\\System32\\drivers\\etc\\hosts"
	default:
		filePath = "/etc/hosts"
	}

	slog.Info("[App] Host file path set to", "path", filePath)

	hosts = host

	// Create application with options
	err = wails.Run(&options.App{
		Title:  "Simple Hosts",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 25, G: 25, B: 25, A: 1},
		OnStartup:        app.startup,
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		slog.Error(err.Error())
	}
}
