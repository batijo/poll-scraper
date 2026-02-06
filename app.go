package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/batijo/poll-scraper/config"
	"github.com/batijo/poll-scraper/models"
	"github.com/batijo/poll-scraper/server"
	"github.com/batijo/poll-scraper/utils"
	"github.com/batijo/poll-scraper/utils/file"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type App struct {
	ctx context.Context
	cfg *config.Config
	srv *server.Server
}

func NewApp() *App {
	return &App{}
}

func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx

	cfg, err := config.Load("config.json")
	if err != nil {
		runtime.LogFatal(ctx, fmt.Sprintf("failed to load config: %v", err))
		return
	}
	a.cfg = cfg

	if err := a.initLogger(cfg.Debug); err != nil {
		runtime.LogFatal(ctx, fmt.Sprintf("failed to init logger: %v", err))
		return
	}

	srv := server.New(cfg)
	srv.Addr = fmt.Sprintf("%s:%d", cfg.IP, cfg.Port)
	a.srv = srv

	go func() {
		slog.Info("server starting", "address", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("server stopped", "err", err)
		}
	}()

	if err := file.InitFiles(cfg); err != nil {
		slog.Error("failed to init files", "err", err)
	}

	if err := file.StartWriting(cfg); err != nil {
		slog.Error("failed to start writer", "err", err)
	}
}

func (a *App) Shutdown(ctx context.Context) {
	slog.Info("application shutting down")
}

func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, from Go!", name)
}

func (a *App) GetConfig() *config.Config {
	return a.cfg
}

func (a *App) UpdateConfig(cfg config.Config) error {
	if err := cfg.Save("config.json"); err != nil {
		return err
	}

	a.cfg = &cfg
	return nil
}

func (a *App) EmitScraperData(data []models.Data) {
	payload := map[string]interface{}{
		"data":      data,
		"timestamp": time.Now().Unix(),
	}
	runtime.EventsEmit(a.ctx, "polled:data", payload)
}

func (a *App) EmitScraperState(state string) {
	runtime.EventsEmit(a.ctx, "polled:state", state)
}

func (a *App) EmitScraperError(message string) {
	payload := map[string]interface{}{
		"message":   message,
		"timestamp": time.Now().Unix(),
	}
	runtime.EventsEmit(a.ctx, "polled:error", payload)
}

func (a *App) initLogger(debug bool) error {
	logFile, err := os.OpenFile("error.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, utils.FileMode)
	if err != nil {
		return fmt.Errorf("failed to open log file: %w", err)
	}

	errorLog := slog.New(slog.NewTextHandler(logFile, &slog.HandlerOptions{
		Level:     slog.LevelError,
		AddSource: true,
	}))
	slog.SetDefault(errorLog)

	if debug {
		slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		})))
	}

	return nil
}
