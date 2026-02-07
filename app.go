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
	"github.com/batijo/poll-scraper/scraper"
	"github.com/batijo/poll-scraper/server"
	"github.com/batijo/poll-scraper/utils"
	"github.com/batijo/poll-scraper/utils/file"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type App struct {
	ctx          context.Context
	cfg          *config.Config
	srv          *server.Server
	stopWriter   context.CancelFunc
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

	stopWriter, err := file.StartWriting(cfg, a)
	if err != nil {
		slog.Error("failed to start writer", "err", err)
	}
	a.stopWriter = stopWriter
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

	oldCfg := a.cfg
	a.cfg = &cfg

	// Restart the writer goroutine with new config
	if a.stopWriter != nil {
		a.stopWriter()
	}
	stopWriter, err := file.StartWriting(a.cfg, a)
	if err != nil {
		slog.Error("failed to restart writer", "err", err)
		return err
	}
	a.stopWriter = stopWriter

	// Restart HTTP server so handler picks up the new config pointer
	if a.srv != nil {
		if err := a.srv.Close(); err != nil {
			slog.Error("failed to stop server", "err", err)
		}
	}
	srv := server.New(a.cfg)
	srv.Addr = fmt.Sprintf("%s:%d", cfg.IP, cfg.Port)
	a.srv = srv
	go func() {
		slog.Info("server restarting", "address", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("server stopped", "err", err)
		}
	}()

	// Reinit output files if paths or toggles changed
	if oldCfg.WriteToCSV != cfg.WriteToCSV || oldCfg.CSVPath != cfg.CSVPath ||
		oldCfg.WriteToTXT != cfg.WriteToTXT || oldCfg.TXTPath != cfg.TXTPath {
		if err := file.InitFiles(a.cfg); err != nil {
			slog.Error("failed to reinit files", "err", err)
		}
	}

	// Reinit logger if debug mode changed
	if oldCfg.Debug != cfg.Debug {
		if err := a.initLogger(cfg.Debug); err != nil {
			slog.Error("failed to reinit logger", "err", err)
		}
	}

	return nil
}

func (a *App) EmitScraperData(data []models.Data) {
	payload := map[string]interface{}{
		"data":      data,
		"timestamp": time.Now().Format(time.RFC3339),
	}
	runtime.EventsEmit(a.ctx, "polled:data", payload)
}

func (a *App) EmitScraperState(state string) {
	runtime.EventsEmit(a.ctx, "polled:state", state)
}

func (a *App) PreviewURL(url string) []models.Data {
	return scraper.ScrapeURL(url, a.cfg.WithEq)
}

func (a *App) EmitURLStatus(statuses []models.URLStatus) {
	runtime.EventsEmit(a.ctx, "polled:url-status", statuses)
}

func (a *App) EmitLog(entry utils.LogEntry) {
	runtime.EventsEmit(a.ctx, "polled:log", entry)
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

	var inner slog.Handler
	if debug {
		inner = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		})
	} else {
		inner = slog.NewTextHandler(logFile, &slog.HandlerOptions{
			Level:     slog.LevelError,
			AddSource: true,
		})
	}

	// Wrap with frontend handler to forward logs to the UI
	handler := utils.NewFrontendHandler(inner, a)
	slog.SetDefault(slog.New(handler))

	return nil
}
