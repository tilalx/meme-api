package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"Meme_Api/config"
	"Meme_Api/libraries/reddit"
	"Meme_Api/libraries/redis"
	"Meme_Api/server"

	"github.com/getsentry/sentry-go"
)

// @title           Meme API
// @version         1.0.0
// @description     A simple, free REST API serving random memes scraped from Reddit.
// @host            meme-api.aelx.de
// @BasePath        /
// @schemes         https
func main() {
	// JSON structured logging
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))

	// Validate required env vars before doing anything else
	if err := config.Validate(); err != nil {
		slog.Error("startup validation failed", "error", err)
		os.Exit(1)
	}

	// Initialize Libraries
	reddit.Init()
	redis.Init()

	// Initialize Sentry (optional — app continues without it if DSN is unset or invalid)
	if err := sentry.Init(sentry.ClientOptions{
		Dsn: os.Getenv("SENTRY_DSN"),
	}); err != nil {
		slog.Warn("Sentry init warning", "error", err)
	}
	defer sentry.Flush(2 * time.Second)

	// Listen for OS shutdown signals for graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Start server in a goroutine
	go server.Init()

	// Block until a shutdown signal is received
	<-quit
	slog.Info("shutdown signal received, shutting down gracefully")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		slog.Error("server forced to shutdown", "error", err)
	}

	slog.Info("server exited")
}
