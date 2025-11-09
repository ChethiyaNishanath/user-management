package db

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

var pool *sql.DB

const (
	HOST     = "localhost"
	PORT     = 5432
	DB       = "usermanagementdb"
	SSL_MODE = "disable"
)

func Connect() *sql.DB {
	connStr := "postgres://postgres:Test1234@localhost:5432/usermanagementdb?sslmode=disable"

	pool, err := sql.Open("pgx", connStr)
	if err != nil {
		panic(fmt.Sprintf("Error opening DB: %v", err))
	}

	pool.SetConnMaxLifetime(0)
	pool.SetMaxIdleConns(3)
	pool.SetMaxOpenConns(3)

	ctx, stop := context.WithCancel(context.Background())
	defer stop()

	appSignal := make(chan os.Signal, 3)
	signal.Notify(appSignal, os.Interrupt)

	go func() {
		<-appSignal
		stop()
		pool.Close()
	}()

	if err := pool.PingContext(ctx); err != nil {
		slog.Error("unable to connect to database", "error", err)
	}

	slog.Info("Connected to PostgreSQL")
	return pool
}

func Ping(ctx context.Context) {
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	if err := pool.PingContext(ctx); err != nil {
		slog.Error("unable to connect to database", "error", err)
	}
}
