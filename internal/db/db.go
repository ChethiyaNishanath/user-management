package db

import (
	"context"
	"database/sql"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func Connect(connStr string) *sql.DB {
	const (
		maxRetries = 5
		retryDelay = 3 * time.Second
	)

	var pool *sql.DB
	var err error

	for attempt := 1; attempt <= maxRetries; attempt++ {
		pool, err = sql.Open("pgx", connStr)
		if err != nil {
			slog.Error("Failed to open DB connection", "attempt", attempt, "error", err)
			time.Sleep(retryDelay)
			continue
		}

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		if err = pool.PingContext(ctx); err == nil {
			slog.Info("Connected to PostgreSQL", "attempt", attempt)
			break
		}

		slog.Warn("Database not ready, retrying...", "attempt", attempt, "error", err)
		time.Sleep(retryDelay)
	}

	if err != nil {
		slog.Error("Unable to connect to PostgreSQL after retries", "error", err)
		panic(err)
	}

	pool.SetConnMaxLifetime(time.Hour)
	pool.SetMaxIdleConns(5)
	pool.SetMaxOpenConns(10)

	go handleShutdown(pool)

	return pool
}

func handleShutdown(pool *sql.DB) {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	<-signalChan
	slog.Info("Shutting down database connection...")
	if err := pool.Close(); err != nil {
		slog.Error("Error closing DB connection", "error", err)
	} else {
		slog.Info("Database connection closed.")
	}
}
