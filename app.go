package main

import (
	"context"
	"fmt"
	"log/slog"
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
	slog.Info("[App]: [GetHosts] Getting hosts")
	hosts, err := readHostFile()
	if err != nil {
		return nil, err
	}

	slog.Info("[App]: [GetHosts] Returning hosts", "hosts", fmt.Sprintf("%+v", hosts))

	return hosts, nil
}

func (a *App) UpdateHost(host Host) error {
	slog.Info("[App]: [UpdateHost] Updating host", "host", fmt.Sprintf("%+v", host))
	hosts, err := readHostFile()
	if err != nil {
		return err
	}

	// Check if host exists and update
	// If host does not exist, add it
	for i, h := range hosts {
		if h.IP == host.IP {
			hosts[i] = host
			break
		}
	}

	// If host does not exist, add it using if check
	if !containsHost(hosts, host) {
		slog.Warn("[App]: [UpdateHost] Host does not exist, adding it")
		hosts = append(hosts, host)
	}

	slog.Info("[App]: [UpdateHost] Attempting to write hosts to file")

	// Write hosts to file
	err = writeHostFile(hosts)
	if err != nil {
		slog.Error("[App]: [UpdateHost] Error writing hosts to file")
		return err
	}

	slog.Info("[App]: [UpdateHost] Written hosts to file")

	return nil
}
