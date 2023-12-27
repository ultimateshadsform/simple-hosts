package main

import (
	"context"
	"fmt"
	"os"
	"os/user"
	"runtime"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

func (a *App) CheckAdmin() (bool, error) {
	if runtime.GOOS == "windows" {
		fileHandle, err := os.Open("\\\\.\\PHYSICALDRIVE0")
		if err != nil {
			return false, err
		}
		defer fileHandle.Close()

		return true, nil
	} else {
		user, err := user.Current()
		if err != nil {
			return false, err
		}
		if user.Uid == "0" {
			return true, nil
		}

		return false, nil
	}
}

func (a *App) GetHosts() ([]Host, error) {
	hosts, err := readHostFile()
	if err != nil {
		return nil, err
	}

	return hosts, nil
}
