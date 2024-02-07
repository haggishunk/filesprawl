package database

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Database interface {
	QueryRow(ctx context.Context, query string, args ...interface{}) pgx.Row
	Exec(ctx context.Context, query string, args ...interface{}) (pgconn.CommandTag, error)
}

type PgxDatabase struct {
	pool *pgxpool.Pool
}

func NewPgxDatabase(p *pgxpool.Pool) *PgxDatabase {
	return &PgxDatabase{pool: p}
}

func (d *PgxDatabase) QueryRow(ctx context.Context, query string, args ...any) pgx.Row {
	return d.pool.QueryRow(ctx, query, args...)
}

func (d *PgxDatabase) Exec(ctx context.Context, query string, args ...interface{}) (pgconn.CommandTag, error) {
	return d.pool.Exec(ctx, query, args...)
}
