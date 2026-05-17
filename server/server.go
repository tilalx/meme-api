package server

import (
	"context"
	"log"
	"net/http"

	"github.com/getsentry/sentry-go"
)

var httpServer *http.Server

// Init : Initialize the routes and server
func Init() {
	r := NewRouter()

	httpServer = &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		sentry.CaptureException(err)
		log.Println("Server error:", err)
	}
}

// Shutdown gracefully stops the HTTP server with the given context.
func Shutdown(ctx context.Context) error {
	if httpServer == nil {
		return nil
	}
	return httpServer.Shutdown(ctx)
}
