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

	_ "user-management/docs"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/go-chi/httprate"
)

// @title User Management API
// @version 1.0
// @description This is User Management server.
// @termsOfService http://swagger.io/terms/

// @contact.name Chethiya Viharagama
// @contact.url http://www.swagger.io/support
// @contact.email chethiya.viharagama@yaalalabs.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @BasePath /
func main() {

	config := config.Load()

	slog.SetLogLoggerLevel(slog.LevelDebug)
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	dbConn := db.Connect(config.DBDsn)
	defer dbConn.Close()

	newApp := app.NewApp(dbConn, &ctx)
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Use(middleware.Timeout(60 * time.Second))

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	r.Use(httprate.LimitByIP(100, 1*time.Minute))

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
