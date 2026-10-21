package repository

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/haggishunk/filesprawl/internal/database"
	"github.com/haggishunk/filesprawl/internal/object"
	"github.com/haggishunk/filesprawl/internal/rclone"
	"github.com/haggishunk/filesprawl/internal/remote"
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

func TestIntegration_PersistResult_PersistsRemoteAssociations(t *testing.T) {
	pool := openTestDB(t)
	repo := NewObjectRepository(database.NewPgxDatabase(pool))
	ctx := context.Background()

	stamp := time.Now().Format("150405.000")
	remoteName := "integ_remote_" + stamp + ":"
	rem := remote.NewRemote("integ-host", remoteName)

	item1 := rclone.ListResponseItem{
		Name:     "a.txt",
		Path:     "/integ/" + stamp + "/a.txt",
		MimeType: "text/plain",
		Size:     512,
		Hashes:   map[string]string{"md5": "integ-remote-hash-a-" + stamp},
	}
	item2 := rclone.ListResponseItem{
		Name:     "b.txt",
		Path:     "/integ/" + stamp + "/b.txt",
		MimeType: "text/plain",
		Size:     512,
		Hashes:   map[string]string{"md5": "integ-remote-hash-b-" + stamp},
	}

	defer func() {
		repo.db.Exec(ctx, `DELETE FROM object_remote_junction WHERE remote_id = $1`, rem.Id)
		repo.db.Exec(ctx, `DELETE FROM object_hash_junction WHERE object_meta_id IN (SELECT id FROM object_meta WHERE object_path = $1 OR object_path = $2)`, item1.Path, item2.Path)
		repo.db.Exec(ctx, `DELETE FROM object_meta WHERE object_path = $1 OR object_path = $2`, item1.Path, item2.Path)
		repo.db.Exec(ctx, `DELETE FROM object_hash WHERE hash_value = $1 OR hash_value = $2`, item1.Hashes["md5"], item2.Hashes["md5"])
		repo.db.Exec(ctx, `DELETE FROM remote WHERE id = $1`, rem.Id)
	}()

	if err := repo.PersistResult(ctx, &rem, item1, "/integ/"+stamp); err != nil {
		t.Fatalf("PersistResult item1: %v", err)
	}
	if err := repo.PersistResult(ctx, &rem, item2, "/integ/"+stamp); err != nil {
		t.Fatalf("PersistResult item2: %v", err)
	}

	var remoteCount int
	err := repo.db.QueryRow(ctx,
		`SELECT COUNT(*) FROM remote WHERE remote_name = $1 AND hostname = $2 AND remote_type = $3`,
		rem.Name, rem.Hostname, rem.Type,
	).Scan(&remoteCount)
	if err != nil {
		t.Fatalf("count remote rows: %v", err)
	}
	if remoteCount != 1 {
		t.Fatalf("expected 1 remote row, got %d", remoteCount)
	}

	var linkCount int
	err = repo.db.QueryRow(ctx,
		`SELECT COUNT(*) FROM object_remote_junction WHERE remote_id = $1`,
		rem.Id,
	).Scan(&linkCount)
	if err != nil {
		t.Fatalf("count remote links: %v", err)
	}
	if linkCount != 2 {
		t.Fatalf("expected 2 object_remote_junction rows, got %d", linkCount)
	}
}

func TestIntegration_PersistLocalResult_ScanRootDeduplication(t *testing.T) {
	pool := openTestDB(t)
	repo := NewObjectRepository(database.NewPgxDatabase(pool))
	ctx := context.Background()

	stamp := time.Now().Format("150405.000")
	rootA := "/integ/photos-" + stamp
	rootB := "/integ/photos-alt-" + stamp
	remA := remote.NewRemote("integ-host", rootA, remote.WithRemoteType("local"))
	remB := remote.NewRemote("integ-host", rootB, remote.WithRemoteType("local"))

	// Same file name and content scanned under two different roots
	hashVal := "integ-crossroot-hash-" + stamp
	hashes := map[string]string{"md5": hashVal}
	mA := object.NewMeta("a.jpg", "a.jpg", "image/jpeg", object.WithMetaSize(512), object.WithMetaScanRoot(rootA))
	mB := object.NewMeta("a.jpg", "a.jpg", "image/jpeg", object.WithMetaSize(512), object.WithMetaScanRoot(rootB))

	defer func() {
		repo.db.Exec(ctx, `DELETE FROM object_hash_junction WHERE object_meta_id IN ($1, $2)`, mA.Id, mB.Id)
		repo.db.Exec(ctx, `DELETE FROM object_remote_junction WHERE object_meta_id IN ($1, $2)`, mA.Id, mB.Id)
		repo.db.Exec(ctx, `DELETE FROM object_meta WHERE id IN ($1, $2)`, mA.Id, mB.Id)
		repo.db.Exec(ctx, `DELETE FROM object_hash WHERE hash_value = $1`, hashVal)
		repo.db.Exec(ctx, `DELETE FROM remote WHERE id IN ($1, $2)`, remA.Id, remB.Id)
	}()

	if err := repo.PersistLocalResult(ctx, &remA, &mA, hashes); err != nil {
		t.Fatalf("persist root A: %v", err)
	}
	if err := repo.PersistLocalResult(ctx, &remB, &mB, hashes); err != nil {
		t.Fatalf("persist root B: %v", err)
	}

	// Distinct records per scan root (issue #3: previously collapsed into one)
	if mA.Id == mB.Id {
		t.Fatalf("expected distinct meta records per scan root, both got id %d", mA.Id)
	}

	// The shared hash is now detected as a duplicate group
	groups, err := repo.FindDuplicatesAcrossRemotes(ctx, "md5", 0, 0)
	if err != nil {
		t.Fatalf("find duplicates: %v", err)
	}
	found := false
	for _, g := range groups {
		if g.HashValue == hashVal && g.Count == 2 {
			found = true
		}
	}
	if !found {
		t.Fatalf("expected cross-root duplicate group for hash %s", hashVal)
	}

	// GetFilesByHash reports scan provenance
	records, err := repo.GetFilesByHash(ctx, hashVal, "md5", 0)
	if err != nil {
		t.Fatalf("get files by hash: %v", err)
	}
	roots := map[string]bool{}
	for _, r := range records {
		roots[r.ScanRoot] = true
	}
	if !roots[rootA] || !roots[rootB] {
		t.Fatalf("expected scan roots %s and %s in records, got %v", rootA, rootB, roots)
	}

	// Rescanning the same root reuses the existing record (idempotent)
	mA2 := object.NewMeta("a.jpg", "a.jpg", "image/jpeg", object.WithMetaSize(512), object.WithMetaScanRoot(rootA))
	if err := repo.PersistLocalResult(ctx, &remA, &mA2, hashes); err != nil {
		t.Fatalf("rescan root A: %v", err)
	}
	if mA2.Id != mA.Id {
		t.Fatalf("expected rescan of same root to reuse meta id %d, got %d", mA.Id, mA2.Id)
	}
}

func TestIntegration_PersistLocalResult_SymlinkMarked(t *testing.T) {
	pool := openTestDB(t)
	repo := NewObjectRepository(database.NewPgxDatabase(pool))
	ctx := context.Background()

	stamp := time.Now().Format("150405.000")
	root := "/integ/links-" + stamp
	rem := remote.NewRemote("integ-host", root, remote.WithRemoteType("local"))

	m := object.NewMeta("link.txt", "link.txt", "application/octet-stream",
		object.WithMetaSize(8), object.WithMetaSymlink(true), object.WithMetaScanRoot(root))

	defer func() {
		repo.db.Exec(ctx, `DELETE FROM object_remote_junction WHERE object_meta_id = $1`, m.Id)
		repo.db.Exec(ctx, `DELETE FROM object_meta WHERE id = $1`, m.Id)
		repo.db.Exec(ctx, `DELETE FROM remote WHERE id = $1`, rem.Id)
	}()

	// Symlinks are persisted without hashes
	if err := repo.PersistLocalResult(ctx, &rem, &m, nil); err != nil {
		t.Fatalf("persist symlink: %v", err)
	}

	var isSymlink bool
	var scanRoot *string
	var hashLinks int
	err := repo.db.QueryRow(ctx,
		`SELECT is_symlink, scan_root,
		        (SELECT COUNT(*) FROM object_hash_junction WHERE object_meta_id = $1)
		 FROM object_meta WHERE id = $1`, m.Id,
	).Scan(&isSymlink, &scanRoot, &hashLinks)
	if err != nil {
		t.Fatalf("read back symlink row: %v", err)
	}
	if !isSymlink {
		t.Fatalf("expected is_symlink=true")
	}
	if scanRoot == nil || *scanRoot != root {
		t.Fatalf("expected scan_root %q, got %v", root, scanRoot)
	}
	if hashLinks != 0 {
		t.Fatalf("expected no hash junctions for symlink, got %d", hashLinks)
	}
}
