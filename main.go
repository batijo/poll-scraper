package main

import (
	"log/slog"
	"os"

	"github.com/batijo/poll-scraper/server"
	"github.com/batijo/poll-scraper/utils/file"
	"github.com/joho/godotenv"
)

var logger *slog.Logger

func initLogger() *slog.Logger {
	logFile, err := os.OpenFile("info.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
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

func init() {
	if err := godotenv.Load(); err != nil {
		slog.Warn(".env file not found")
	}
	logger = initLogger()
}

func main() {
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
		os.Setenv("PS_IP", "localhost")
	}
	addr := os.Getenv("PS_IP") + ":" + os.Getenv("PS_PORT")
	slog.Info("server starting", "address", addr)
	if err := srv.App.Listen(addr); err != nil {
		logger.Error("server stopped with error", "err", err)
		os.Exit(1)
	}
}
