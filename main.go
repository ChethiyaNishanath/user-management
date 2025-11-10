package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
	"user-management/internal/app"
	"user-management/internal/config"
	"user-management/internal/db"

	"github.com/go-chi/chi/v5"
)

func main() {

	config := config.Load()

	slog.SetLogLoggerLevel(slog.LevelInfo)
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	dbConn := db.Connect(config.DBDsn)
	defer dbConn.Close()

	newApp := app.NewApp(dbConn)
	r := chi.NewRouter()
	newApp.RegisterRoutes(r)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", config.Port),
		Handler: r,
	}

	go func() {
		slog.Info(fmt.Sprintf("Server starting on port %s", config.Port))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("Server failed", "error", err)
		}
	}()

	<-ctx.Done()
	slog.Info("Shutting down gracefully...")

	stop()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), time.Duration(func() int {
		v, _ := strconv.Atoi(config.ShutdownTimeout)
		return v
	}())*time.Second)

	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		slog.Error("Forced server shutdown", "error", err)
	}

	if err := dbConn.Close(); err != nil {
		slog.Error("Error closing DB", "error", err)
	}

	slog.Info("Shutdown complete")
}
