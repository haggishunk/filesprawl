package database

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Defines what functional scope and signatures a Database should have
type Database interface {
	QueryRow(ctx context.Context, query string, args ...interface{}) pgx.Row
	Query(ctx context.Context, query string, args ...interface{}) (pgx.Rows, error)
	Exec(ctx context.Context, query string, args ...interface{}) (pgconn.CommandTag, error)
}

// PgxDatabase is the root functional object implementing Database interface for postgres (pgx)
type PgxDatabase struct {
	pool *pgxpool.Pool
}

// NewPgxDatabase is a factory function receiving a pgx pool and returning
// a PgxDatabase struct with the pgx pool assigned
func NewPgxDatabase(p *pgxpool.Pool) *PgxDatabase {
	return &PgxDatabase{pool: p}
}

// QueryRow implements the Database interface for the function of the same name
func (d *PgxDatabase) QueryRow(ctx context.Context, query string, args ...any) pgx.Row {
	return d.pool.QueryRow(ctx, query, args...)
}

// Query implements the Database interface for multi-row queries
func (d *PgxDatabase) Query(ctx context.Context, query string, args ...any) (pgx.Rows, error) {
	return d.pool.Query(ctx, query, args...)
}

// Exec implements the Database interface for the function of the same name
func (d *PgxDatabase) Exec(ctx context.Context, query string, args ...interface{}) (pgconn.CommandTag, error) {
	return d.pool.Exec(ctx, query, args...)
}
