package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Config struct {
	// URL is the connection string (required).
	// Format: postgres://user:password@host:port/database?sslmode=disable
	URL string

	// MaxConns is the maximum number of connections in the pool.
	// Default: 25
	MaxConns int32

	// MinConns is the minimum number of connections to keep open.
	// Default: 5
	MinConns int32

	// MaxConnLifetime is how long a connection can be reused.
	// Default: 1 hour
	MaxConnLifetime time.Duration

	// MaxConnIdleTime is how long a connection can be idle before closing.
	// Default: 30 minutes
	MaxConnIdleTime time.Duration

	// HealthCheckPeriod is how often to check connection health.
	// Default: 1 minute
	HealthCheckPeriod time.Duration
}

// DefaultConfig returns production-ready defaults.
func DefaultConfig(url string) Config {
	return Config{
		URL:               url,
		MaxConns:          25,
		MinConns:          5,
		MaxConnLifetime:   time.Hour,
		MaxConnIdleTime:   30 * time.Minute,
		HealthCheckPeriod: time.Minute,
	}
}

// Open opens a PostgreSQL connection pool with default settings.
//
// Example:
//
//	pool, err := postgres.Open(ctx, "postgres://user:pass@localhost:5432/mydb")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	defer pool.Close()
//
//	var name string
//	err = pool.QueryRow(ctx, "SELECT name FROM users WHERE id = $1", 1).Scan(&name)
func Open(ctx context.Context, url string) (*pgxpool.Pool, error) {
	return OpenWithConfig(ctx, DefaultConfig(url))
}

func OpenWithConfig(ctx context.Context, cfg Config) (*pgxpool.Pool, error) {
	poolCfg, err := pgxpool.ParseConfig(cfg.URL)
	if err != nil {
		return nil, fmt.Errorf("invalid postgres url: %w", err)
	}

	if cfg.MaxConns > 0 {
		poolCfg.MaxConns = cfg.MaxConns
	}
	if cfg.MinConns > 0 {
		poolCfg.MinConns = cfg.MinConns
	}
	if cfg.MaxConnLifetime > 0 {
		poolCfg.MaxConnLifetime = cfg.MaxConnLifetime
	}
	if cfg.MaxConnIdleTime > 0 {
		poolCfg.MaxConnIdleTime = cfg.MaxConnIdleTime
	}
	if cfg.HealthCheckPeriod > 0 {
		poolCfg.HealthCheckPeriod = cfg.HealthCheckPeriod
	}
	pool, err := pgxpool.NewWithConfig(ctx, poolCfg)
	if err != nil {
		return nil, fmt.Errorf("create postgres pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("ping postgres: %w", err)
	}
	return pool, nil
}
