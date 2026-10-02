package repository

import (
	"context"
	"fmt"

	"github.com/haggishunk/filesprawl/internal/database"
	"github.com/haggishunk/filesprawl/internal/object"
	"github.com/haggishunk/filesprawl/internal/rclone"
	"github.com/haggishunk/filesprawl/internal/remote"
	"github.com/jackc/pgx/v5"
)

// TODO: check out pgx.RowToStructByName to match structs and table columns

// HashGroup represents a hash that appears on more than one file (a duplicate).
type HashGroup struct {
	HashValue string
	HashType  string
	Count     int
}

// FileRecord holds per-file metadata returned by GetFilesByHash.
type FileRecord struct {
	ID         int
	Name       string
	Path       string
	MimeType   string
	Size       int64
	RemoteName string
	Hostname   string
}

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
		INSERT INTO object_meta (object_name, object_path, object_mime_type, object_size)
		VALUES ($1, $2, $3, $4)
		RETURNING id;
	`

	row := r.db.QueryRow(ctx, statement, m.Name, m.Path, m.MimeType, m.Size)
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
	var size int64

	statement := `
		SELECT id, object_size
		FROM object_meta
		WHERE object_name = $1 AND object_path = $2 AND object_mime_type = $3;
	`

	row := r.db.QueryRow(ctx, statement, m.Name, m.Path, m.MimeType)
	err := row.Scan(&id, &size)

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
	m.Size = size
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

func (r *ObjectRepository) ReadRemote(ctx context.Context, rem *remote.Remote) error {
	var id int
	var remoteType string

	statement := `
		SELECT id, remote_type
		FROM remote
		WHERE remote_name = $1 AND hostname = $2 AND remote_type = $3;
	`

	row := r.db.QueryRow(ctx, statement, rem.Name, rem.Hostname, rem.Type)
	err := row.Scan(&id, &remoteType)
	if err != nil {
		if err == pgx.ErrNoRows {
			fmt.Printf("No remote found in db.\n")
			return nil
		}
		return fmt.Errorf("failed to retrieve remote: %w", err)
	}

	rem.Id = id
	rem.Type = remoteType
	rem.Persisted = true
	return nil
}

func (r *ObjectRepository) WriteRemote(ctx context.Context, rem *remote.Remote) error {
	var id int

	statement := `
		INSERT INTO remote (remote_name, remote_type, hostname)
		VALUES ($1, $2, $3)
		RETURNING id;
	`

	row := r.db.QueryRow(ctx, statement, rem.Name, rem.Type, rem.Hostname)
	err := row.Scan(&id)
	if err != nil {
		return fmt.Errorf("failed to persist remote: %w", err)
	}

	rem.Id = id
	rem.Persisted = true
	return nil
}

func (r *ObjectRepository) ReadObjectRemoteJunction(ctx context.Context, orj *object.ObjectRemoteJunction) error {
	var id int

	statement := `
		SELECT id
		FROM object_remote_junction
		WHERE object_meta_id = $1 AND remote_id = $2;
	`

	row := r.db.QueryRow(ctx, statement, orj.MetaId, orj.RemoteId)
	err := row.Scan(&id)
	if err != nil {
		if err == pgx.ErrNoRows {
			fmt.Printf("No object-remote junction found in db.\n")
			return nil
		}
		return fmt.Errorf("failed to retrieve object remote junction: %w", err)
	}

	orj.Id = id
	orj.Persisted = true
	return nil
}

func (r *ObjectRepository) WriteObjectRemoteJunction(ctx context.Context, orj *object.ObjectRemoteJunction) error {
	var id int

	statement := `
		INSERT INTO object_remote_junction (object_meta_id, remote_id)
		VALUES ($1, $2)
		RETURNING id;
	`

	row := r.db.QueryRow(ctx, statement, orj.MetaId, orj.RemoteId)
	err := row.Scan(&id)
	if err != nil {
		return fmt.Errorf("failed to persist object remote junction: %w", err)
	}

	orj.Id = id
	orj.Persisted = true
	return nil
}

func (r *ObjectRepository) ensureRemote(ctx context.Context, rem *remote.Remote) error {
	if rem == nil {
		return fmt.Errorf("remote is required")
	}
	if rem.Persisted {
		return nil
	}

	err := r.ReadRemote(ctx, rem)
	if err != nil {
		return fmt.Errorf("failed to get remote from db: %w", err)
	}
	if rem.Persisted {
		return nil
	}

	err = r.WriteRemote(ctx, rem)
	if err != nil {
		return fmt.Errorf("failed to put remote into db: %w", err)
	}

	return nil
}

func (r *ObjectRepository) PersistResult(ctx context.Context, rem *remote.Remote, lri rclone.ListResponseItem) error {
	// get or set object meta
	m := object.NewMeta(lri.Name, lri.Path, lri.MimeType, object.WithMetaSize(lri.Size))
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

	err = r.ensureRemote(ctx, rem)
	if err != nil {
		return err
	}

	orj := object.NewObjectRemoteJunction(m.Id, rem.Id)
	fmt.Printf("Found object-remote junction: %+v\n", orj)
	err = r.ReadObjectRemoteJunction(ctx, &orj)
	if err != nil {
		return fmt.Errorf("failed to get object remote junction from db: %w", err)
	}
	if orj.Persisted {
		fmt.Printf("Read object remote junction id: %d\n", orj.Id)
	} else {
		fmt.Printf("Persisting object-remote junction...  ")
		err = r.WriteObjectRemoteJunction(ctx, &orj)
		if err != nil {
			return fmt.Errorf("failed to set object remote junction in db: %w", err)
		}
		fmt.Printf("as id: %d\n", orj.Id)
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

// PersistLocalResult persists a local file scan result to the database
func (r *ObjectRepository) PersistLocalResult(ctx context.Context, rem *remote.Remote, m *object.Meta, hashes map[string]string) error {
	// Persist metadata
	fmt.Printf("Found local object meta: %s\n", m)
	err := r.ReadMeta(ctx, m)
	if err != nil {
		return fmt.Errorf("failed to get object meta from db: %w", err)
	}
	if m.Persisted {
		fmt.Printf("Read object meta id: %d\n", m.Id)
	} else {
		fmt.Printf("Persisting local meta...  ")
		err = r.WriteMeta(ctx, m)
		if err != nil {
			return fmt.Errorf("failed to persist local object meta: %w", err)
		}
		fmt.Printf("as %d\n", m.Id)
	}

	// Ensure remote (local type)
	err = r.ensureRemote(ctx, rem)
	if err != nil {
		return err
	}

	// Persist object-remote junction
	orj := object.NewObjectRemoteJunction(m.Id, rem.Id)
	fmt.Printf("Found local object-remote junction: %+v\n", orj)
	err = r.ReadObjectRemoteJunction(ctx, &orj)
	if err != nil {
		return fmt.Errorf("failed to get object remote junction from db: %w", err)
	}
	if orj.Persisted {
		fmt.Printf("Read object remote junction id: %d\n", orj.Id)
	} else {
		fmt.Printf("Persisting local object-remote junction...  ")
		err = r.WriteObjectRemoteJunction(ctx, &orj)
		if err != nil {
			return fmt.Errorf("failed to persist local object remote junction: %w", err)
		}
		fmt.Printf("as id: %d\n", orj.Id)
	}

	// Persist hashes
	for hType, hVal := range hashes {
		h := object.NewHash(hVal, hType)
		fmt.Printf("Found local hash: %s\n", h)

		err := r.ReadHash(ctx, &h)
		if err != nil {
			return fmt.Errorf("failed to get hash from db: %w", err)
		}
		if h.Persisted {
			fmt.Printf("Read hash id: %d\n", h.Id)
		} else {
			fmt.Printf("Persisting local hash...  ")
			err = r.WriteHash(ctx, &h)
			if err != nil {
				return fmt.Errorf("failed to persist local hash: %w", err)
			}
			fmt.Printf("as %d\n", h.Id)
		}

		// Persist meta-hash junction
		mhj := object.NewMetaHashJunction(m.Id, h.Id)
		fmt.Printf("Found local meta-hash junction: %s\n", mhj)
		err = r.ReadMetaHashJunction(ctx, &mhj)
		if err != nil {
			return fmt.Errorf("failed to get meta-hash junction from db: %w", err)
		}
		if mhj.Persisted {
			fmt.Printf("Read meta-hash junction id: %d\n", mhj.Id)
		} else {
			fmt.Printf("Persisting local meta-hash junction...  ")
			err = r.WriteMetaHashJunction(ctx, &mhj)
			if err != nil {
				return fmt.Errorf("failed to persist local meta-hash junction: %w", err)
			}
			fmt.Printf("as id: %d\n", mhj.Id)
		}
	}
	return nil
}

// FindDuplicatesAcrossRemotes returns hashes that are shared by more than one
// file across all remotes. Optionally filter by hashType (empty = all types).
// Use limit=0 and offset=0 for no pagination.
func (r *ObjectRepository) FindDuplicatesAcrossRemotes(ctx context.Context, hashType string, limit, offset int) ([]HashGroup, error) {
	var args []interface{}
	where := ""
	if hashType != "" {
		args = append(args, hashType)
		where = "WHERE oh.hash_type = $1::hash_type_enum"
	}

	limitClause := ""
	if limit > 0 {
		args = append(args, limit)
		args = append(args, offset)
		limitClause = fmt.Sprintf("LIMIT $%d OFFSET $%d", len(args)-1, len(args))
	}

	statement := fmt.Sprintf(`
		SELECT oh.hash_value, oh.hash_type::text, COUNT(DISTINCT om.id) AS duplicate_count
		FROM object_hash oh
		JOIN object_hash_junction ohj ON oh.id = ohj.object_hash_id
		JOIN object_meta om ON ohj.object_meta_id = om.id
		%s
		GROUP BY oh.hash_value, oh.hash_type
		HAVING COUNT(DISTINCT om.id) > 1
		ORDER BY duplicate_count DESC
		%s;
	`, where, limitClause)

	rows, err := r.db.Query(ctx, statement, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query duplicates across remotes: %w", err)
	}
	defer rows.Close()

	var groups []HashGroup
	for rows.Next() {
		var g HashGroup
		if err := rows.Scan(&g.HashValue, &g.HashType, &g.Count); err != nil {
			return nil, fmt.Errorf("failed to scan duplicate row: %w", err)
		}
		groups = append(groups, g)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating duplicate rows: %w", err)
	}
	return groups, nil
}

// FindDuplicatesWithinRemote returns hashes shared by more than one file within
// a specific remote. Optionally filter by hashType (empty = all types).
func (r *ObjectRepository) FindDuplicatesWithinRemote(ctx context.Context, remoteName, hashType string, limit, offset int) ([]HashGroup, error) {
	args := []interface{}{remoteName}
	hashFilter := ""
	if hashType != "" {
		args = append(args, hashType)
		hashFilter = fmt.Sprintf("AND oh.hash_type = $%d::hash_type_enum", len(args))
	}

	limitClause := ""
	if limit > 0 {
		args = append(args, limit)
		args = append(args, offset)
		limitClause = fmt.Sprintf("LIMIT $%d OFFSET $%d", len(args)-1, len(args))
	}

	statement := fmt.Sprintf(`
		SELECT oh.hash_value, oh.hash_type::text, COUNT(DISTINCT om.id) AS duplicate_count
		FROM object_hash oh
		JOIN object_hash_junction ohj ON oh.id = ohj.object_hash_id
		JOIN object_meta om ON ohj.object_meta_id = om.id
		JOIN object_remote_junction orj ON om.id = orj.object_meta_id
		JOIN remote r ON orj.remote_id = r.id
		WHERE r.remote_name = $1
		%s
		GROUP BY oh.hash_value, oh.hash_type
		HAVING COUNT(DISTINCT om.id) > 1
		ORDER BY duplicate_count DESC
		%s;
	`, hashFilter, limitClause)

	rows, err := r.db.Query(ctx, statement, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query duplicates within remote: %w", err)
	}
	defer rows.Close()

	var groups []HashGroup
	for rows.Next() {
		var g HashGroup
		if err := rows.Scan(&g.HashValue, &g.HashType, &g.Count); err != nil {
			return nil, fmt.Errorf("failed to scan duplicate row: %w", err)
		}
		groups = append(groups, g)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating duplicate rows: %w", err)
	}
	return groups, nil
}

// GetFilesByHash returns all files that share the given hash value and type.
// Use minSize=0 to return files of any size.
func (r *ObjectRepository) GetFilesByHash(ctx context.Context, hashValue, hashType string, minSize int64) ([]FileRecord, error) {
	statement := `
		SELECT om.id, om.object_name, om.object_path, om.object_mime_type, om.object_size,
		       COALESCE(rem.remote_name, ''), COALESCE(rem.hostname, '')
		FROM object_meta om
		JOIN object_hash_junction ohj ON om.id = ohj.object_meta_id
		JOIN object_hash oh ON ohj.object_hash_id = oh.id
		LEFT JOIN object_remote_junction orj ON om.id = orj.object_meta_id
		LEFT JOIN remote rem ON orj.remote_id = rem.id
		WHERE oh.hash_value = $1
		  AND oh.hash_type = $2::hash_type_enum
		  AND om.object_size >= $3;
	`

	rows, err := r.db.Query(ctx, statement, hashValue, hashType, minSize)
	if err != nil {
		return nil, fmt.Errorf("failed to query files by hash: %w", err)
	}
	defer rows.Close()

	var records []FileRecord
	for rows.Next() {
		var f FileRecord
		if err := rows.Scan(&f.ID, &f.Name, &f.Path, &f.MimeType, &f.Size, &f.RemoteName, &f.Hostname); err != nil {
			return nil, fmt.Errorf("failed to scan file record: %w", err)
		}
		records = append(records, f)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating file records: %w", err)
	}
	return records, nil
}
