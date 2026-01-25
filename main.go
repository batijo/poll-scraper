package main

import (
	"log/slog"
	"os"

	"github.com/joho/godotenv"

	"github.com/batijo/poll-scraper/server"
	"github.com/batijo/poll-scraper/utils/file"
)

const fileMode = 0o600

var logger *slog.Logger

func initLogger() *slog.Logger {
	logFile, err := os.OpenFile("info.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, fileMode)
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
	if err := godotenv.Load(); err != nil {
		slog.Warn(".env file not found")
	}
	logger = initLogger()

	srv := server.New()
	if err := file.InitFiles(); err != nil {
		logger.Error("failed to init files", "err", err)
		os.Exit(1)
	}
	if err := file.StartWriting(); err != nil {
		logger.Error("failed to start writer", "err", err)
		os.Exit(1)
	}
	if os.Getenv("PS_IP") == "" {
		if err := os.Setenv("PS_IP", "localhost"); err != nil {
			logger.Error("failed to set PS_IP", "err", err)
		}
	}
	srv.Addr = os.Getenv("PS_IP") + ":" + os.Getenv("PS_PORT")
	slog.Info("server starting", "address", srv.Addr)
	if err := srv.ListenAndServe(); err != nil {
		logger.Error("server stopped with error", "err", err)
		os.Exit(1)
	}
}
