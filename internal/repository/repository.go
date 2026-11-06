package repository

import (
	"context"
	"fmt"

	"github.com/haggishunk/filesprawl/internal/database"
	"github.com/haggishunk/filesprawl/internal/location"
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

// FilenameGroup represents a filename that appears in more than one location/with variations.
type FilenameGroup struct {
	Filename string
	Count    int
}

// FileRecord holds per-file metadata returned by GetFilesByHash.
type FileRecord struct {
	ID           int
	Name         string
	Path         string
	MimeType     string
	Size         int64
	LocationId   int
	AbsolutePath string
	LocationName string
	LocationType string
}

type ObjectRepository struct {
	db database.Database
}

func NewObjectRepository(d database.Database) *ObjectRepository {
	return &ObjectRepository{db: d}
}

// GetDatabase returns the underlying Database interface
func (r *ObjectRepository) GetDatabase() database.Database {
	return r.db
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

	// Extract directory from path (everything before the last /)
	// BasePath is already calculated in object.NewMeta using path.Dir()
	statement := `
		INSERT INTO object_meta (object_name, object_path, object_mime_type, object_size, is_symlink, location_id, absolute_path)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id;
	`

	row := r.db.QueryRow(ctx, statement, m.Name, m.Path, m.MimeType, m.Size, m.IsSymlink, m.LocationId, m.AbsolutePath)
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
	var isSymlink bool
	var locationId int
	var absolutePath string

	statement := `
		SELECT id, object_size, is_symlink, location_id, COALESCE(absolute_path, '')
		FROM object_meta
		WHERE object_name = $1 AND object_path = $2 AND object_mime_type = $3
		  AND absolute_path IS NOT DISTINCT FROM NULLIF($4, '');
	`

	row := r.db.QueryRow(ctx, statement, m.Name, m.Path, m.MimeType, m.AbsolutePath)
	err := row.Scan(&id, &size, &isSymlink, &locationId, &absolutePath)

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
	m.IsSymlink = isSymlink
	m.LocationId = locationId
	m.AbsolutePath = absolutePath
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

// ReadRemote is deprecated and no longer used with consolidated schema
func (r *ObjectRepository) ReadRemote(ctx context.Context, rem *remote.Remote) error {
	return fmt.Errorf("ReadRemote is deprecated - use Location directly instead")
}

// WriteRemote is deprecated and no longer used with consolidated schema
func (r *ObjectRepository) WriteRemote(ctx context.Context, rem *remote.Remote) error {
	return fmt.Errorf("WriteRemote is deprecated - use Location directly instead")
}

func (r *ObjectRepository) ReadLocation(ctx context.Context, loc *location.Location) error {
	var id int

	statement := `
		SELECT id
		FROM location
		WHERE name = $1 AND type = $2;
	`

	row := r.db.QueryRow(ctx, statement, loc.Name, loc.Type)
	err := row.Scan(&id)
	if err != nil {
		if err == pgx.ErrNoRows {
			fmt.Printf("No location found in db.\n")
			return nil
		}
		return fmt.Errorf("failed to retrieve location: %w", err)
	}

	loc.Id = id
	loc.Persisted = true
	return nil
}

func (r *ObjectRepository) WriteLocation(ctx context.Context, loc *location.Location) error {
	var id int

	statement := `
		INSERT INTO location (name, type)
		VALUES ($1, $2)
		RETURNING id;
	`

	row := r.db.QueryRow(ctx, statement, loc.Name, loc.Type)
	err := row.Scan(&id)
	if err != nil {
		return fmt.Errorf("failed to persist location: %w", err)
	}

	loc.Id = id
	loc.Persisted = true
	return nil
}

func (r *ObjectRepository) EnsureLocation(ctx context.Context, loc *location.Location) error {
	if loc == nil {
		return fmt.Errorf("location is required")
	}
	if loc.Persisted || loc.Id > 0 {
		return nil
	}

	err := r.ReadLocation(ctx, loc)
	if err != nil {
		return fmt.Errorf("failed to get location from db: %w", err)
	}
	if loc.Persisted {
		return nil
	}

	err = r.WriteLocation(ctx, loc)
	if err != nil {
		return fmt.Errorf("failed to put location into db: %w", err)
	}

	return nil
}

// ReadObjectRemoteJunction is deprecated and no longer used with consolidated schema
func (r *ObjectRepository) ReadObjectRemoteJunction(ctx context.Context, orj *object.ObjectRemoteJunction) error {
	return fmt.Errorf("ReadObjectRemoteJunction is deprecated - object_remote_junction table no longer exists")
}

// WriteObjectRemoteJunction is deprecated and no longer used with consolidated schema
func (r *ObjectRepository) WriteObjectRemoteJunction(ctx context.Context, orj *object.ObjectRemoteJunction) error {
	return fmt.Errorf("WriteObjectRemoteJunction is deprecated - object_remote_junction table no longer exists")
}

func (r *ObjectRepository) PersistResult(ctx context.Context, rem *remote.Remote, lri rclone.ListResponseItem, scanRoot string) error {
	// Ensure location exists
	if rem.Location.Id <= 0 || !rem.Location.Persisted {
		err := r.EnsureLocation(ctx, &rem.Location)
		if err != nil {
			return fmt.Errorf("failed to ensure location: %w", err)
		}
	}

	// Build absolute path (remote format: remote_name:path)
	absolutePath := rem.Location.Name + ":" + lri.Path

	// get or set object meta
	// LocationId and AbsolutePath are required fields
	m := object.NewMeta(lri.Name, lri.Path, lri.MimeType,
		object.WithMetaSize(lri.Size),
		object.WithMetaScanRoot(scanRoot),
		object.WithMetaLocationId(rem.Location.Id),
		object.WithMetaAbsolutePath(absolutePath))
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

// PersistLocalResult persists a local file scan result to the database
// Note: LocationId and AbsolutePath must be set on the Meta before calling this method
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

// FindDuplicatesWithinLocation returns hashes shared by more than one file within
// a specific location. Optionally filter by hashType (empty = all types).
func (r *ObjectRepository) FindDuplicatesWithinLocation(ctx context.Context, locationId int, hashType string, limit, offset int) ([]HashGroup, error) {
	var args []interface{}
	args = append(args, locationId)

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
		WHERE om.location_id = $1
		%s
		GROUP BY oh.hash_value, oh.hash_type
		HAVING COUNT(DISTINCT om.id) > 1
		ORDER BY duplicate_count DESC
		%s;
	`, hashFilter, limitClause)

	rows, err := r.db.Query(ctx, statement, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query duplicates within location: %w", err)
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
// a specific location by looking up the location by name and type.
// This is a convenience wrapper around FindDuplicatesWithinLocation.
// Optionally filter by hashType (empty = all types).
func (r *ObjectRepository) FindDuplicatesWithinRemote(ctx context.Context, locationName, hashType string, limit, offset int) ([]HashGroup, error) {
	// Look up the location by name to get its ID
	var locationId int
	statement := `SELECT id FROM location WHERE name = $1 LIMIT 1`
	row := r.db.QueryRow(ctx, statement, locationName)
	err := row.Scan(&locationId)
	if err != nil {
		if err == pgx.ErrNoRows {
			// No location found with this name
			return []HashGroup{}, nil
		}
		return nil, fmt.Errorf("failed to find location by name %q: %w", locationName, err)
	}

	return r.FindDuplicatesWithinLocation(ctx, locationId, hashType, limit, offset)
}

// GetFilesByHash returns all files that share the given hash value and type.
// Use minSize=0 to return files of any size.
func (r *ObjectRepository) GetFilesByHash(ctx context.Context, hashValue, hashType string, minSize int64) ([]FileRecord, error) {
	statement := `
		SELECT om.id, om.object_name, om.object_path, om.object_mime_type, om.object_size,
		       om.location_id, om.absolute_path, COALESCE(l.name, ''), COALESCE(l.type, '')
		FROM object_meta om
		JOIN object_hash_junction ohj ON om.id = ohj.object_meta_id
		JOIN object_hash oh ON ohj.object_hash_id = oh.id
		LEFT JOIN location l ON om.location_id = l.id
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
		if err := rows.Scan(&f.ID, &f.Name, &f.Path, &f.MimeType, &f.Size, &f.LocationId, &f.AbsolutePath, &f.LocationName, &f.LocationType); err != nil {
			return nil, fmt.Errorf("failed to scan file record: %w", err)
		}
		records = append(records, f)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating file records: %w", err)
	}
	return records, nil
}

// HashTypePair represents a hash value and type combination.
type HashTypePair struct {
	HashValue string
	HashType  string
}

// FilesByHashMap maps a hash to its file records for batch operations.
// The key is a combination of "hashValue:hashType" for stable lookup.
type FilesByHashMap map[string][]FileRecord

// GetFilesByHashes returns all files for multiple hash values in batches to avoid stack depth limits.
// Use minSize=0 to return files of any size.
// Returns a map keyed by "hashValue:hashType" for efficient group enrichment.
// Processes hashes in batches of 1000 to avoid PostgreSQL stack depth limit exceeded errors.
func (r *ObjectRepository) GetFilesByHashes(ctx context.Context, hashes []HashTypePair, minSize int64) (FilesByHashMap, error) {
	if len(hashes) == 0 {
		return FilesByHashMap{}, nil
	}

	result := make(FilesByHashMap)

	// Process hashes in batches to avoid stack depth limit exceeded error
	// PostgreSQL has limits on the number of elements in an IN clause
	const batchSize = 1000

	for batchStart := 0; batchStart < len(hashes); batchStart += batchSize {
		batchEnd := batchStart + batchSize
		if batchEnd > len(hashes) {
			batchEnd = len(hashes)
		}

		batch := hashes[batchStart:batchEnd]

		// Build the WHERE clause for this batch: (hash_value, hash_type) IN (...)
		statement := `
			SELECT om.id, om.object_name, om.object_path, om.object_mime_type, om.object_size,
			       om.location_id, om.absolute_path, COALESCE(l.name, ''), COALESCE(l.type, ''),
			       oh.hash_value, oh.hash_type::text
			FROM object_meta om
			JOIN object_hash_junction ohj ON om.id = ohj.object_meta_id
			JOIN object_hash oh ON ohj.object_hash_id = oh.id
			LEFT JOIN location l ON om.location_id = l.id
			WHERE (oh.hash_value, oh.hash_type::text) IN (
		`

		// Build the VALUES clause dynamically for this batch
		args := []interface{}{minSize}
		for i, h := range batch {
			if i > 0 {
				statement += ", "
			}
			statement += fmt.Sprintf("($%d::text, $%d::text)", len(args)+1, len(args)+2)
			args = append(args, h.HashValue, h.HashType)
		}

		statement += `
			)
			AND om.object_size >= $1
			ORDER BY oh.hash_value, oh.hash_type;
		`

		rows, err := r.db.Query(ctx, statement, args...)
		if err != nil {
			return nil, fmt.Errorf("failed to batch query files by hashes (batch %d-%d): %w", batchStart, batchEnd, err)
		}
		defer rows.Close()

		for rows.Next() {
			var f FileRecord
			var hashValue, hashType string
			if err := rows.Scan(&f.ID, &f.Name, &f.Path, &f.MimeType, &f.Size, &f.LocationId, &f.AbsolutePath, &f.LocationName, &f.LocationType, &hashValue, &hashType); err != nil {
				return nil, fmt.Errorf("failed to scan file record in batch: %w", err)
			}
			key := hashValue + ":" + hashType
			result[key] = append(result[key], f)
		}

		if err := rows.Err(); err != nil {
			return nil, fmt.Errorf("error iterating file records in batch (batch %d-%d): %w", batchStart, batchEnd, err)
		}
	}

	return result, nil
}

// FindDuplicateFilenamesAcrossRemotes returns filenames that appear in more than one location.
// Optionally filter by minSize (0 means no size filter).
// Use limit=0 and offset=0 for no pagination.
func (r *ObjectRepository) FindDuplicateFilenamesAcrossRemotes(ctx context.Context, minSize int64, limit, offset int) ([]FilenameGroup, error) {
	var args []interface{}
	where := ""
	if minSize > 0 {
		args = append(args, minSize)
		where = fmt.Sprintf("WHERE om.object_size >= $%d", len(args))
	}

	limitClause := ""
	if limit > 0 {
		args = append(args, limit)
		args = append(args, offset)
		limitClause = fmt.Sprintf("LIMIT $%d OFFSET $%d", len(args)-1, len(args))
	}

	statement := fmt.Sprintf(`
		SELECT om.object_name, COUNT(DISTINCT om.id) AS duplicate_count
		FROM object_meta om
		%s
		GROUP BY om.object_name
		HAVING COUNT(DISTINCT om.id) > 1
		ORDER BY duplicate_count DESC
		%s;
	`, where, limitClause)

	rows, err := r.db.Query(ctx, statement, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query duplicate filenames: %w", err)
	}
	defer rows.Close()

	var groups []FilenameGroup
	for rows.Next() {
		var g FilenameGroup
		if err := rows.Scan(&g.Filename, &g.Count); err != nil {
			return nil, fmt.Errorf("failed to scan filename group: %w", err)
		}
		groups = append(groups, g)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating filename groups: %w", err)
	}

	return groups, nil
}

// FindDuplicateFilenamesWithinRemote returns filenames that appear more than once within a specific location.
// Optionally filter by minSize (0 means no size filter).
// Use limit=0 and offset=0 for no pagination.
func (r *ObjectRepository) FindDuplicateFilenamesWithinRemote(ctx context.Context, locationName string, minSize int64, limit, offset int) ([]FilenameGroup, error) {
	// Look up the location by name to get its ID
	var locationId int
	statement := `SELECT id FROM location WHERE name = $1 LIMIT 1`
	row := r.db.QueryRow(ctx, statement, locationName)
	err := row.Scan(&locationId)
	if err != nil {
		if err == pgx.ErrNoRows {
			// No location found with this name
			return []FilenameGroup{}, nil
		}
		return nil, fmt.Errorf("failed to find location %q: %w", locationName, err)
	}

	var args []interface{}
	args = append(args, locationId)
	where := fmt.Sprintf("WHERE om.location_id = $%d", len(args))
	if minSize > 0 {
		args = append(args, minSize)
		where += fmt.Sprintf(" AND om.object_size >= $%d", len(args))
	}

	limitClause := ""
	if limit > 0 {
		args = append(args, limit)
		args = append(args, offset)
		limitClause = fmt.Sprintf("LIMIT $%d OFFSET $%d", len(args)-1, len(args))
	}

	statement = fmt.Sprintf(`
		SELECT om.object_name, COUNT(DISTINCT om.id) AS duplicate_count
		FROM object_meta om
		%s
		GROUP BY om.object_name
		HAVING COUNT(DISTINCT om.id) > 1
		ORDER BY duplicate_count DESC
		%s;
	`, where, limitClause)

	rows, err := r.db.Query(ctx, statement, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query duplicate filenames within location: %w", err)
	}
	defer rows.Close()

	var groups []FilenameGroup
	for rows.Next() {
		var g FilenameGroup
		if err := rows.Scan(&g.Filename, &g.Count); err != nil {
			return nil, fmt.Errorf("failed to scan filename group: %w", err)
		}
		groups = append(groups, g)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating filename groups: %w", err)
	}

	return groups, nil
}

// GetFilesByFilename returns all files with the given name (across all locations).
// Optionally filter by minSize (0 means no size filter).
func (r *ObjectRepository) GetFilesByFilename(ctx context.Context, filename string, minSize int64) ([]FileRecord, error) {
	var args []interface{}
	args = append(args, filename)
	where := fmt.Sprintf("WHERE om.object_name = $%d", len(args))
	if minSize > 0 {
		args = append(args, minSize)
		where += fmt.Sprintf(" AND om.object_size >= $%d", len(args))
	}

	statement := fmt.Sprintf(`
		SELECT om.id, om.object_name, om.object_path, om.object_mime_type, om.object_size,
		       om.location_id, om.absolute_path, l.name, l.type
		FROM object_meta om
		JOIN location l ON om.location_id = l.id
		%s
		ORDER BY om.absolute_path;
	`, where)

	rows, err := r.db.Query(ctx, statement, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query files by filename: %w", err)
	}
	defer rows.Close()

	var records []FileRecord
	for rows.Next() {
		var f FileRecord
		if err := rows.Scan(&f.ID, &f.Name, &f.Path, &f.MimeType, &f.Size,
			&f.LocationId, &f.AbsolutePath, &f.LocationName, &f.LocationType); err != nil {
			return nil, fmt.Errorf("failed to scan file record: %w", err)
		}
		records = append(records, f)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating file records: %w", err)
	}

	return records, nil
}
