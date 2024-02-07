package repository

import (
	"context"
	"fmt"

	"github.com/haggishunk/filesprawl/internal/database"
	"github.com/haggishunk/filesprawl/internal/object"
)

// TODO: check out pgx.RowToStructByName to match structs and table columns

type ObjectRepository struct {
	db database.Database
}

func NewObjectRepository(d database.Database) *ObjectRepository {
	return &ObjectRepository{db: d}
}

func (r *ObjectRepository) SetOrGetHash(ctx context.Context, h object.Hash) (int, error) {
	// write to database if not exists
	// return new or existing id
	var id int

	statement := `
		INSERT INTO object_hash (hash_value, hash_type)
		VALUES ($1, $2)
		ON CONFLICT (hash_value) DO UPDATE
		SET hash_value = EXCLUDED.hash_value
		RETURNING id;
	`

	row := r.db.QueryRow(ctx, statement, h.Hash, h.Type)
	err := row.Scan(&id)

	if err != nil {
		return 0, fmt.Errorf("failed to persist hash: %w", err)
	}

	return id, nil
}

func (r *ObjectRepository) SetHash(ctx context.Context, h object.Hash) (int, error) {
	// write to database if not exists
	// return new or existing id
	var id int

	statement := `
		INSERT INTO object_hash (hash_value, hash_type)
		VALUES ($1, $2)

		WHERE hash_value = $1 AND hash_type = $2;
	`

	row := r.db.QueryRow(ctx, statement, h.Hash, h.Type)
	err := row.Scan(&id)

	if err != nil {
		return 0, fmt.Errorf("failed to persist hash: %w", err)
	}

	return id, nil
}

func (r *ObjectRepository) GetHash(ctx context.Context, h object.Hash) (int, error) {
	// write to database if not exists
	// return new or existing id
	var id int

	statement := `
		SELECT id
		FROM object_hash
		WHERE hash_value = $1 AND hash_type = $2;
	`

	row := r.db.QueryRow(ctx, statement, h.Hash, h.Type)
	err := row.Scan(&id)

	if err != nil {
		return 0, fmt.Errorf("failed to persist hash: %w", err)
	}

	return id, nil
}
