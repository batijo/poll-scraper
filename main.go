package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/batijo/poll-scraper/config"
	"github.com/batijo/poll-scraper/server"
	"github.com/batijo/poll-scraper/utils"
	"github.com/batijo/poll-scraper/utils/file"
)

var errorLog *slog.Logger

func initLogger(debug bool) {
	logFile, err := os.OpenFile("error.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, utils.FileMode)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open log file: %v\n", err)
		os.Exit(1)
	}
	errorLog = slog.New(slog.NewTextHandler(logFile, &slog.HandlerOptions{
		Level:     slog.LevelError,
		AddSource: true,
	}))

	if debug {
		slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		})))
	}
}

func main() {
	cfg, err := config.Load("config.json")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load config: %v\n", err)
		os.Exit(1)
	}
	initLogger(cfg.Debug)

	srv := server.New(cfg)
	if err := file.InitFiles(cfg); err != nil {
		errorLog.Error("failed to init files", "err", err)
		os.Exit(1)
	}
	if err := file.StartWriting(cfg); err != nil {
		errorLog.Error("failed to start writer", "err", err)
		os.Exit(1)
	}
	srv.Addr = fmt.Sprintf("%s:%d", cfg.IP, cfg.Port)
	slog.Info("server starting", "address", srv.Addr)
	if err := srv.ListenAndServe(); err != nil {
		errorLog.Error("server stopped with error", "err", err)
		os.Exit(1)
	}
}
