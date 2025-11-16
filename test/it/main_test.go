package it

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"user-management/internal/app"
	"user-management/internal/db"

	"github.com/go-chi/chi/v5"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

var TestRouter http.Handler
var TestDB *sql.DB

func TestMain(m *testing.M) {
	_, filename, _, _ := runtime.Caller(0)
	projectRoot := filepath.Join(filepath.Dir(filename), "..")
	testdataDir := filepath.Join(projectRoot, "testdata")

	ctx := context.Background()

	dbName := "usermanagementdb"
	dbUser := "postgres"
	dbPassword := "Test12344"

	postgresContainer, err := postgres.Run(ctx,
		"postgres:16-alpine",
		postgres.WithInitScripts(filepath.Join(testdataDir, "init-user-db.sh")),
		postgres.WithConfigFile(filepath.Join(testdataDir, "test-container-postgres.conf")),
		postgres.WithDatabase(dbName),
		postgres.WithUsername(dbUser),
		postgres.WithPassword(dbPassword),
		postgres.BasicWaitStrategies(),
	)

	defer func() {
		if err := testcontainers.TerminateContainer(postgresContainer); err != nil {
			slog.Error(fmt.Sprintf("failed to terminate container: %s", err))
		}
	}()
	if err != nil {
		slog.Error(fmt.Sprintf("failed to start container: %s", err))
		return
	}

	port, _ := postgresContainer.MappedPort(ctx, "5432/tcp")
	host, _ := postgresContainer.Host(ctx)

	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		dbUser, dbPassword, host, port.Port(), dbName,
	)
	TestDB := db.Connect(dsn)

	defer TestDB.Close()

	newApp := app.NewApp(TestDB, &ctx)

	r := chi.NewRouter()
	newApp.RegisterRoutes(r)

	TestRouter = r

	os.Exit(m.Run())
}
