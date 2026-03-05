package repository

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/haggishunk/filesprawl/internal/database"
	"github.com/jackc/pgx/v5/pgxpool"
)

// openTestDB returns a real *pgxpool.Pool for integration tests.
// The test is skipped when DATABASE_URL is not set.
func openTestDB(t *testing.T) *pgxpool.Pool {
	t.Helper()
	url := os.Getenv("DATABASE_URL")
	if url == "" {
		t.Skip("DATABASE_URL not set; skipping integration test")
	}
	pool, err := pgxpool.New(context.Background(), url)
	if err != nil {
		t.Fatalf("failed to connect to database: %v", err)
	}
	t.Cleanup(pool.Close)
	return pool
}

// seedDuplicates inserts two object_meta rows that share the same hash.
// Returns a cleanup function that removes the inserted rows.
func seedDuplicates(t *testing.T, repo *ObjectRepository) func() {
	t.Helper()
	ctx := context.Background()

	// Two files with the same md5 hash
	hashVal := "integ_test_hash_" + time.Now().Format("150405.000")

	hash := writeRaw(t, ctx, repo, hashVal, "md5")
	meta1 := writeMetaRaw(t, ctx, repo, "dup_a.txt", "/integ/dup_a.txt", "text/plain")
	meta2 := writeMetaRaw(t, ctx, repo, "dup_b.txt", "/integ/dup_b.txt", "text/plain")
	j1 := writeJunctionRaw(t, ctx, repo, meta1, hash)
	j2 := writeJunctionRaw(t, ctx, repo, meta2, hash)

	return func() {
		ctx := context.Background()
		repo.db.Exec(ctx, `DELETE FROM object_hash_junction WHERE id = $1 OR id = $2`, j1, j2)
		repo.db.Exec(ctx, `DELETE FROM object_meta WHERE id = $1 OR id = $2`, meta1, meta2)
		repo.db.Exec(ctx, `DELETE FROM object_hash WHERE id = $1`, hash)
	}
}

func writeRaw(t *testing.T, ctx context.Context, repo *ObjectRepository, hashVal, hashType string) int {
	t.Helper()
	var id int
	err := repo.db.QueryRow(ctx,
		`INSERT INTO object_hash (hash_value, hash_type) VALUES ($1, $2::hash_type_enum) RETURNING id`,
		hashVal, hashType,
	).Scan(&id)
	if err != nil {
		t.Fatalf("seed hash: %v", err)
	}
	return id
}

func writeMetaRaw(t *testing.T, ctx context.Context, repo *ObjectRepository, name, path, mime string) int {
	t.Helper()
	var id int
	err := repo.db.QueryRow(ctx,
		`INSERT INTO object_meta (object_name, object_path, object_mime_type, object_size) VALUES ($1, $2, $3, 512) RETURNING id`,
		name, path, mime,
	).Scan(&id)
	if err != nil {
		t.Fatalf("seed meta: %v", err)
	}
	return id
}

func writeJunctionRaw(t *testing.T, ctx context.Context, repo *ObjectRepository, metaID, hashID int) int {
	t.Helper()
	var id int
	err := repo.db.QueryRow(ctx,
		`INSERT INTO object_hash_junction (scan_time, object_meta_id, object_hash_id) VALUES (NOW(), $1, $2) RETURNING id`,
		metaID, hashID,
	).Scan(&id)
	if err != nil {
		t.Fatalf("seed junction: %v", err)
	}
	return id
}

func TestIntegration_FindDuplicatesAcrossRemotes(t *testing.T) {
	pool := openTestDB(t)
	repo := NewObjectRepository(database.NewPgxDatabase(pool))
	cleanup := seedDuplicates(t, repo)
	defer cleanup()

	groups, err := repo.FindDuplicatesAcrossRemotes(context.Background(), "md5", 0, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(groups) == 0 {
		t.Fatal("expected at least one duplicate group, got none")
	}
	// Verify at least one group has count >= 2
	found := false
	for _, g := range groups {
		if g.Count >= 2 {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("no group with count >= 2 found; groups: %+v", groups)
	}
}

func TestIntegration_FindDuplicatesAcrossRemotes_HashTypeFilter(t *testing.T) {
	pool := openTestDB(t)
	repo := NewObjectRepository(database.NewPgxDatabase(pool))
	cleanup := seedDuplicates(t, repo)
	defer cleanup()

	// sha1 filter should return no duplicates (we only seeded md5)
	groups, err := repo.FindDuplicatesAcrossRemotes(context.Background(), "sha1", 0, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, g := range groups {
		if g.HashType != "sha1" {
			t.Errorf("expected only sha1 groups, got %s", g.HashType)
		}
	}
}

func TestIntegration_GetFilesByHash(t *testing.T) {
	pool := openTestDB(t)
	repo := NewObjectRepository(database.NewPgxDatabase(pool))

	hashVal := "integ_files_hash_" + time.Now().Format("150405.000")
	ctx := context.Background()

	hash := writeRaw(t, ctx, repo, hashVal, "md5")
	meta1 := writeMetaRaw(t, ctx, repo, "file_a.txt", "/integ/file_a.txt", "text/plain")
	meta2 := writeMetaRaw(t, ctx, repo, "file_b.txt", "/integ/file_b.txt", "text/plain")
	j1 := writeJunctionRaw(t, ctx, repo, meta1, hash)
	j2 := writeJunctionRaw(t, ctx, repo, meta2, hash)
	defer func() {
		repo.db.Exec(ctx, `DELETE FROM object_hash_junction WHERE id = $1 OR id = $2`, j1, j2)
		repo.db.Exec(ctx, `DELETE FROM object_meta WHERE id = $1 OR id = $2`, meta1, meta2)
		repo.db.Exec(ctx, `DELETE FROM object_hash WHERE id = $1`, hash)
	}()

	records, err := repo.GetFilesByHash(ctx, hashVal, "md5", 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(records) != 2 {
		t.Errorf("expected 2 records, got %d", len(records))
	}
}

func TestIntegration_GetFilesByHash_MinSizeFilter(t *testing.T) {
	pool := openTestDB(t)
	repo := NewObjectRepository(database.NewPgxDatabase(pool))

	hashVal := "integ_size_hash_" + time.Now().Format("150405.000")
	ctx := context.Background()

	hash := writeRaw(t, ctx, repo, hashVal, "md5")
	// seed meta with size=512 (set in writeMetaRaw)
	meta := writeMetaRaw(t, ctx, repo, "small.txt", "/integ/small.txt", "text/plain")
	j := writeJunctionRaw(t, ctx, repo, meta, hash)
	defer func() {
		repo.db.Exec(ctx, `DELETE FROM object_hash_junction WHERE id = $1`, j)
		repo.db.Exec(ctx, `DELETE FROM object_meta WHERE id = $1`, meta)
		repo.db.Exec(ctx, `DELETE FROM object_hash WHERE id = $1`, hash)
	}()

	// minSize larger than 512 should exclude the record
	records, err := repo.GetFilesByHash(ctx, hashVal, "md5", 1024)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(records) != 0 {
		t.Errorf("expected 0 records above minSize, got %d", len(records))
	}
}

