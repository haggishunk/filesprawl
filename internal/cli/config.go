// Package cli holds CLI-specific helpers used by the filesprawl command-line
// binary. These helpers are intentionally kept in an internal package so they
// remain outside the public library API exposed by the root filesprawl module.
package cli

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/haggishunk/filesprawl/internal/database"
	"github.com/haggishunk/filesprawl/internal/repository"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// OpenRepository reads DATABASE_URL from the environment, opens a pgx pool,
// and returns an ObjectRepository together with a cleanup function the caller
// must invoke when the repository is no longer needed.
func OpenRepository(ctx context.Context) (*repository.ObjectRepository, func(), error) {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		return nil, nil, fmt.Errorf("DATABASE_URL is required")
	}

	config, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		return nil, nil, fmt.Errorf("error parsing database url: %w", err)
	}
	config.AfterConnect = func(_ context.Context, conn *pgx.Conn) error {
		log.Printf("Connected to database with pid %d", conn.PgConn().PID())
		return nil
	}

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to connect: %w", err)
	}

	cleanup := func() {
		pool.Close()
	}

	return repository.NewObjectRepository(database.NewPgxDatabase(pool)), cleanup, nil
}
