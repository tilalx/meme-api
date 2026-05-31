package server

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/getsentry/sentry-go"
)

var httpServer *http.Server

// Init : Initialize the routes and server
func Init() {
	r := NewRouter()

	httpServer = &http.Server{
		Addr:         ":8080",
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		sentry.CaptureException(err)
		slog.Error("server error", "error", err)
	}
}

// Shutdown gracefully stops the HTTP server with the given context.
func Shutdown(ctx context.Context) error {
	if httpServer == nil {
		return nil
	}
	return httpServer.Shutdown(ctx)
}
