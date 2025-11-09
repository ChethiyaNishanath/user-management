package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"user-management/internal/app"
	"user-management/internal/db"

	"github.com/go-chi/chi/v5"
)

func main() {
	slog.SetLogLoggerLevel(slog.LevelInfo)
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	dbConn := db.Connect()
	defer dbConn.Close()

	newApp := app.NewApp(dbConn)
	r := chi.NewRouter()
	newApp.RegisterRoutes(r)

	server := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	go func() {
		slog.Info("Server running on http://localhost:8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("Server failed", "error", err)
		}
	}()

	<-ctx.Done()
	slog.Info("Shutting down gracefully...")

	stop()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		slog.Error("Forced server shutdown", "error", err)
	}

	if err := dbConn.Close(); err != nil {
		slog.Error("Error closing DB", "error", err)
	}

	slog.Info("Shutdown complete")
}
