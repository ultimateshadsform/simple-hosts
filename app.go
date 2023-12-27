package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/user"
	"runtime"

	"github.com/txn2/txeh"
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
	slog.Info("[App] Starting up...")
	a.ctx = ctx
	slog.Info("[App] Startup complete")
}

func (a *App) CheckAdmin() (bool, error) {
	slog.Info("Checking if user is admin")
	if runtime.GOOS == "windows" {
		fileHandle, err := os.Open("\\\\.\\PHYSICALDRIVE0")
		if err != nil {
			slog.Error(fmt.Sprintf("[App]: [Windows] Error checking if user is admin: %s", err.Error()))
			return false, err
		}
		defer fileHandle.Close()

		slog.Info("[App]: [Windows] User is admin")
		return true, nil
	} else {
		user, err := user.Current()
		if err != nil {
			slog.Error(fmt.Sprintf("[App]: [Unix] Error checking if user is admin: %s", err.Error()))
			return false, err
		}
		if user.Uid == "0" {
			slog.Info("[App]: [Unix] User is admin")
			return true, nil
		}

		return false, nil
	}
}

func (a *App) GetHosts() ([]Host, error) {
	slog.Info("[App]: [GetHosts] Getting hosts file")

	lines, err := txeh.ParseHosts(filePath)
	if err != nil {
		slog.Error("[App]: [GetHosts] Error parsing hosts file", "error", err.Error())
		return nil, err
	}

	var hosts []Host

	for _, line := range lines {
		if line.LineType == txeh.EMPTY || line.LineType == txeh.COMMENT {
			continue
		}

		hosts = append(hosts, Host{
			IP:       line.Address,
			Hostname: line.Hostnames[0],
			Comment:  line.Comment,
		})
	}

	slog.Info("[App]: [GetHosts] Returning hosts", "hosts", hosts)

	return hosts, nil
}

func (a *App) UpdateHost(host Host) error {
	slog.Info("[App]: [UpdateHost] Updating host", "host", host)

	hosts.AddHost(host.IP, host.Hostname)

	err := hosts.Save()
	if err != nil {
		slog.Error("[App]: [UpdateHost] Error saving hosts file", "error", err.Error())
		return err
	}

	slog.Info("[App]: [UpdateHost] Host updated successfully")

	return nil
}
