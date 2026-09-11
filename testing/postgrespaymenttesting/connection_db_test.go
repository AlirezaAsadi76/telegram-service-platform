package postgrespaymenttesting

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func newTestPool(t *testing.T) *pgxpool.Pool {
	t.Helper()
	
	host := envOrDefault("TEST_POSTGRES_HOST", "localhost")
	port := envOrDefault("TEST_POSTGRES_PORT", "5433")
	user := envOrDefault("TEST_POSTGRES_USER", "postgres")
	database := envOrDefault(
		"TEST_POSTGRES_DATABASE",
		"telegram_service_platform_test",
	)
	password := envOrDefault("TEST_POSTGRES_PASSWORD", "postgres")

	if password == "" {
		t.Fatal("TEST_POSTGRES_PASSWORD is required")
	}

	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		user,
		password,
		host,
		port,
		database,
	)

	ctx, cancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)
	defer cancel()

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		t.Fatalf("create postgres pool: %v", err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		t.Fatalf("ping postgres: %v", err)
	}

	t.Cleanup(func() {
		pool.Close()
	})

	return pool
}

func envOrDefault(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	return value
}
