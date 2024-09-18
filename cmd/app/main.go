package main

import (
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/kshuta/fetchSRE/healthCheck"
)

func getLogLevel() slog.Level {
	switch os.Getenv("LOG_LEVEL") {
	case "DEBUG":
		return slog.LevelDebug
	case "INFO":
		return slog.LevelInfo
	case "WARN":
		return slog.LevelWarn
	case "ERROR":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

func main() {
	args := os.Args
	if len(args) != 2 {
		fmt.Printf("usage: %s {endpoint file}\n", args[0])
		return
	}

	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: getLogLevel(),
	}))

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	err := healthCheck.Run(args[1], sigChan, logger)
	if err != nil {
		logger.Error("running health check", "err", err)
		return
	}
}
