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

var logger *slog.Logger

func initLogger() *slog.Logger {
	logFile, err := os.OpenFile("info.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, utils.FileMode)
	if err != nil {
		slog.Error("failed to open log file", "err", err)
		os.Exit(1)
	}
	handler := slog.NewTextHandler(logFile, &slog.HandlerOptions{
		Level:     slog.LevelInfo,
		AddSource: true,
	})
	return slog.New(handler)
}

func main() {
	cfg, err := config.Load("config.json")
	if err != nil {
		slog.Error("failed to load config", "err", err)
		os.Exit(1)
	}
	logger = initLogger()

	srv := server.New(cfg)
	if err := file.InitFiles(cfg); err != nil {
		logger.Error("failed to init files", "err", err)
		os.Exit(1)
	}
	if err := file.StartWriting(cfg); err != nil {
		logger.Error("failed to start writer", "err", err)
		os.Exit(1)
	}
	srv.Addr = fmt.Sprintf("%s:%d", cfg.IP, cfg.Port)
	slog.Info("server starting", "address", srv.Addr)
	if err := srv.ListenAndServe(); err != nil {
		logger.Error("server stopped with error", "err", err)
		os.Exit(1)
	}
}
