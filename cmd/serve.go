package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
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
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Application startup",
	Run: func(cmd *cobra.Command, args []string) {
		run()
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
	serveCmd.Flags().Int("server.port", 8080, "Port to run the server on")
	serveCmd.Flags().Int("server.shutdownTimeout", 10, "Server shutdown timeout")
	serveCmd.Flags().String("database.dsn", "postgres://postgres:Test1234@localhost:5432/usermanagementdb?sslmode=disable", "DB DSN")
	serveCmd.Flags().String("logging.level", "INFO", "DB DSN")
}

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
func run() {

	cfg, err := loadConfig()

	if err != nil {
		slog.Error("Invalid config", "error", err)
		os.Exit(1)
	}

	if b, err := json.MarshalIndent(cfg, "", "  "); err == nil {
		slog.Debug("CONFIG LOADED\n", "values", string(b))
	}

	slog.SetLogLoggerLevel(parseLogLevel(cfg.Logging.Level))
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	dbConn := db.Connect(cfg.Database.Dsn)
	defer dbConn.Close()

	newApp := app.NewApp(dbConn)
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

	port := cfg.Server.Port

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: r,
	}

	go func() {
		slog.Info(fmt.Sprintf("Server starting on port %s", port))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("Server failed", "error", err)
		}
	}()

	<-ctx.Done()
	slog.Info("Shutting down gracefully...")

	stop()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), time.Duration(func() int {
		v, _ := strconv.Atoi(cfg.Server.ShutdownTimeout)
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

func loadConfig() (*config.Config, error) {
	var cfg config.Config

	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}
	return &cfg, nil
}

func parseLogLevel(level string) slog.Level {
	switch strings.ToLower(level) {
	case "debug":
		return slog.LevelDebug
	case "warn", "warning":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
