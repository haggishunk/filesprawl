package repository

import (
	"context"
	"fmt"

	"github.com/haggishunk/filesprawl/internal/database"
	"github.com/haggishunk/filesprawl/internal/object"
	"github.com/haggishunk/filesprawl/internal/rclone"
	"github.com/jackc/pgx/v5"
)

// TODO: check out pgx.RowToStructByName to match structs and table columns

type ObjectRepository struct {
	db database.Database
}

func NewObjectRepository(d database.Database) *ObjectRepository {
	return &ObjectRepository{db: d}
}

func (r *ObjectRepository) WriteHash(ctx context.Context, h *object.Hash) error {
	// write to database if not exists
	// return new or existing id
	var id int

	statement := `
		INSERT INTO object_hash (hash_value, hash_type)
		VALUES ($1, $2)
		RETURNING id;
	`

	row := r.db.QueryRow(ctx, statement, h.Hash, h.Type)
	err := row.Scan(&id)

	if err != nil {
		return fmt.Errorf("failed to persist hash: %w", err)
	}

	// modify the hash ref
	h.Id = id
	h.Persisted = true

	return nil
}

func (r *ObjectRepository) ReadHash(ctx context.Context, h *object.Hash) error {
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
		if err == pgx.ErrNoRows {
			// hash is left alone and no error
			fmt.Printf("No hash found in db.\n")
			return nil
		}
		// hash is left alone but we throw an error
		return fmt.Errorf("failed to retrieve hash: %w", err)
	}

	// modify the hash ref
	h.Id = id
	h.Persisted = true

	return nil
}

func (r *ObjectRepository) WriteMeta(ctx context.Context, m *object.Meta) error {
	// write to database if not exists
	// return new or existing id
	var id int

	statement := `
		INSERT INTO object_meta (object_name, object_path, object_mime_type)
		VALUES ($1, $2, $3)
		RETURNING id;
	`

	row := r.db.QueryRow(ctx, statement, m.Name, m.Path, m.MimeType)
	err := row.Scan(&id)

	if err != nil {
		return fmt.Errorf("failed to persist object meta: %w", err)
	}

	m.Id = id
	m.Persisted = true
	return nil
}

func (r *ObjectRepository) ReadMeta(ctx context.Context, m *object.Meta) error {
	// write to database if not exists
	// return new or existing id
	var id int

	statement := `
		SELECT id
		FROM object_meta
		WHERE object_name = $1 AND object_path = $2 AND object_mime_type = $3;
	`

	row := r.db.QueryRow(ctx, statement, m.Name, m.Path, m.MimeType)
	err := row.Scan(&id)

	if err != nil {
		if err == pgx.ErrNoRows {
			// hash is left alone and no error
			fmt.Printf("No meta found in db.\n")
			return nil
		}
		// hash is left alone but we throw an error
		return fmt.Errorf("failed to retrieve meta: %w", err)
	}

	m.Id = id
	m.Persisted = true
	return nil
}

func (r *ObjectRepository) ReadMetaHashJunction(ctx context.Context, mhj *object.MetaHashJunction) error {
	// write to database if not exists
	// return new or existing id
	var id int

	statement := `
		SELECT id
		FROM object_hash_junction
		WHERE scan_time = $1 and object_meta_id = $2 AND object_hash_id = $3;
	`

	row := r.db.QueryRow(ctx, statement, mhj.ScanTime, mhj.MetaId, mhj.HashId)
	err := row.Scan(&id)

	if err != nil {
		if err == pgx.ErrNoRows {
			// hash is left alone and no error
			fmt.Printf("No meta-hash junction found in db.\n")
			return nil
		}
		// hash is left alone but we throw an error
		return fmt.Errorf("failed to retrieve meta-hash junction: %w", err)
	}

	mhj.Id = id
	mhj.Persisted = true
	return nil
}

func (r *ObjectRepository) WriteMetaHashJunction(ctx context.Context, mhj *object.MetaHashJunction) error {
	// write to database if not exists
	// return new or existing id
	var id int

	statement := `
		INSERT INTO object_hash_junction (scan_time, object_meta_id, object_hash_id)
		VALUES ($1, $2, $3)
		RETURNING id;
	`

	row := r.db.QueryRow(ctx, statement, mhj.ScanTime, mhj.MetaId, mhj.HashId)
	err := row.Scan(&id)

	if err != nil {
		return fmt.Errorf("failed to persist object hash junction: %w", err)
	}

	mhj.Id = id
	mhj.Persisted = true
	return nil
}

func (r *ObjectRepository) PersistResult(ctx context.Context, lri rclone.ListResponseItem) error {
	// get or set object meta
	m := object.NewMeta(lri.Name, lri.Path, lri.MimeType)
	fmt.Printf("Found object meta: %s\n", m)
	err := r.ReadMeta(ctx, &m)
	if err != nil {
		return fmt.Errorf("failed to get object meta from db: %w", err)
	}
	if m.Persisted {
		fmt.Printf("Read object meta id: %d\n", m.Id)
	} else {
		fmt.Printf("Persisting meta...  ")
		err = r.WriteMeta(ctx, &m)
		if err != nil {
			return fmt.Errorf("failed to put object into db: %w", err)
		}
		fmt.Printf("as %d\n", m.Id)
	}

	// get or set hash
	for hType, hVal := range lri.Hashes {
		h := object.NewHash(hVal, hType)
		fmt.Printf("Found hash: %s\n", h)

		err := r.ReadHash(ctx, &h)
		if err != nil {
			return fmt.Errorf("failed to get hash from db: %w", err)
		}
		if h.Persisted {
			fmt.Printf("Read hash id: %d\n", h.Id)
		} else {
			fmt.Printf("Persisting hash...  ")
			err = r.WriteHash(ctx, &h)
			if err != nil {
				return fmt.Errorf("failed to put hash into db: %w", err)
			}
			fmt.Printf("as %d\n", h.Id)
		}

		// get or set object hash junction
		mhj := object.NewMetaHashJunction(m.Id, h.Id)
		fmt.Printf("Found object meta-hash junction: %s\n", mhj)
		err = r.ReadMetaHashJunction(ctx, &mhj)
		if err != nil {
			return fmt.Errorf("failed to get object meta hash junction from db: %w", err)
		}
		if mhj.Persisted {
			fmt.Printf("Read object meta hash junction id: %d\n", mhj.Id)
		} else {
			fmt.Printf("Persisting meta-hash junction...  ")
			err = r.WriteMetaHashJunction(ctx, &mhj)
			if err != nil {
				return fmt.Errorf("failed to set object meta-hash junction in db: %w", err)
			}
			fmt.Printf("as id: %d\n", mhj.Id)
		}
	}
	return nil
}
