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

// TestGetFilesByHashes tests the batch query method for multiple hashes.
// This integration test verifies that GetFilesByHashes returns all files for
// multiple hashes in a single query, addressing the N+1 performance issue.
func TestGetFilesByHashes(t *testing.T) {
	pool := openTestDB(t)
	repo := NewObjectRepository(database.NewPgxDatabase(pool))
	ctx := context.Background()

	// Create 3 hashes with 2 files each = 6 total files
	hash1 := writeRaw(t, ctx, repo, "batch_hash_1", "md5")
	hash2 := writeRaw(t, ctx, repo, "batch_hash_2", "sha256")
	hash3 := writeRaw(t, ctx, repo, "batch_hash_3", "sha1")

	meta1a := writeMetaRaw(t, ctx, repo, "file1a.txt", "/path/file1a.txt", "text/plain")
	meta1b := writeMetaRaw(t, ctx, repo, "file1b.txt", "/path/file1b.txt", "text/plain")
	meta2a := writeMetaRaw(t, ctx, repo, "file2a.txt", "/path/file2a.txt", "text/plain")
	meta2b := writeMetaRaw(t, ctx, repo, "file2b.txt", "/path/file2b.txt", "text/plain")
	meta3a := writeMetaRaw(t, ctx, repo, "file3a.txt", "/path/file3a.txt", "text/plain")
	meta3b := writeMetaRaw(t, ctx, repo, "file3b.txt", "/path/file3b.txt", "text/plain")

	writeJunctionRaw(t, ctx, repo, meta1a, hash1)
	writeJunctionRaw(t, ctx, repo, meta1b, hash1)
	writeJunctionRaw(t, ctx, repo, meta2a, hash2)
	writeJunctionRaw(t, ctx, repo, meta2b, hash2)
	writeJunctionRaw(t, ctx, repo, meta3a, hash3)
	writeJunctionRaw(t, ctx, repo, meta3b, hash3)

	// Batch query for all 3 hashes
	hashes := []HashTypePair{
		{HashValue: "batch_hash_1", HashType: "md5"},
		{HashValue: "batch_hash_2", HashType: "sha256"},
		{HashValue: "batch_hash_3", HashType: "sha1"},
	}

	result, err := repo.GetFilesByHashes(ctx, hashes, 0)
	if err != nil {
		t.Fatalf("GetFilesByHashes failed: %v", err)
	}

	// Verify we got results for all 3 hashes
	if len(result) != 3 {
		t.Fatalf("expected 3 hash keys, got %d", len(result))
	}

	// Verify each hash has 2 files
	for _, key := range []string{"batch_hash_1:md5", "batch_hash_2:sha256", "batch_hash_3:sha1"} {
		if len(result[key]) != 2 {
			t.Errorf("hash %q: expected 2 files, got %d", key, len(result[key]))
		}
	}

	// Verify file names are correct
	files1 := result["batch_hash_1:md5"]
	if len(files1) == 2 {
		names := make(map[string]bool)
		for _, f := range files1 {
			names[f.Name] = true
		}
		if !names["file1a.txt"] || !names["file1b.txt"] {
			t.Errorf("hash1 files: expected file1a.txt and file1b.txt, got %v", names)
		}
	}

	// Cleanup
	repo.db.Exec(ctx, `DELETE FROM object_hash_junction WHERE object_hash_id IN ($1, $2, $3)`, hash1, hash2, hash3)
	repo.db.Exec(ctx, `DELETE FROM object_meta WHERE id IN ($1, $2, $3, $4, $5, $6)`,
		meta1a, meta1b, meta2a, meta2b, meta3a, meta3b)
	repo.db.Exec(ctx, `DELETE FROM object_hash WHERE id IN ($1, $2, $3)`, hash1, hash2, hash3)
}

// TestGetFilesByHashes_WithMinSize tests batch query respects size filter.
func TestGetFilesByHashes_WithMinSize(t *testing.T) {
	pool := openTestDB(t)
	repo := NewObjectRepository(database.NewPgxDatabase(pool))
	ctx := context.Background()

	hash := writeRaw(t, ctx, repo, "size_test_hash", "md5")
	// File with size 5000 (passes minSize=1000)
	meta1 := writeMetaWithSize(t, ctx, repo, "large.bin", "/path/large.bin", "application/octet-stream", 5000)
	// File with size 500 (fails minSize=1000)
	meta2 := writeMetaWithSize(t, ctx, repo, "small.bin", "/path/small.bin", "application/octet-stream", 500)

	writeJunctionRaw(t, ctx, repo, meta1, hash)
	writeJunctionRaw(t, ctx, repo, meta2, hash)

	// Query with minSize=1000 should only return large.bin
	hashes := []HashTypePair{{HashValue: "size_test_hash", HashType: "md5"}}
	result, err := repo.GetFilesByHashes(ctx, hashes, 1000)
	if err != nil {
		t.Fatalf("GetFilesByHashes failed: %v", err)
	}

	files := result["size_test_hash:md5"]
	if len(files) != 1 {
		t.Errorf("expected 1 file after size filter, got %d", len(files))
	}
	if files[0].Name != "large.bin" {
		t.Errorf("expected large.bin, got %s", files[0].Name)
	}

	// Cleanup
	repo.db.Exec(ctx, `DELETE FROM object_hash_junction WHERE object_hash_id = $1`, hash)
	repo.db.Exec(ctx, `DELETE FROM object_meta WHERE id IN ($1, $2)`, meta1, meta2)
	repo.db.Exec(ctx, `DELETE FROM object_hash WHERE id = $1`, hash)
}

// TestGetFilesByHashes_Empty tests batch query with empty hash list.
func TestGetFilesByHashes_Empty(t *testing.T) {
	pool := openTestDB(t)
	repo := NewObjectRepository(database.NewPgxDatabase(pool))
	ctx := context.Background()

	result, err := repo.GetFilesByHashes(ctx, []HashTypePair{}, 0)
	if err != nil {
		t.Fatalf("GetFilesByHashes failed: %v", err)
	}
	if len(result) != 0 {
		t.Errorf("expected 0 results for empty hash list, got %d", len(result))
	}
}

// Helper to write metadata with specific size
func writeMetaWithSize(t *testing.T, ctx context.Context, repo *ObjectRepository, name, path, mime string, size int64) int {
	t.Helper()
	m := object.NewMeta(name, path, mime)
	m.Size = size
	if err := repo.WriteMeta(ctx, &m); err != nil {
		t.Fatalf("WriteMeta failed: %v", err)
	}
	return m.Id
}

// TestFindDuplicatesWithinRemote_WithHostnamePrefix tests Issue #7 fix:
// FindDuplicatesWithinRemote should accept both "remote_name" and "hostname/remote_name" formats.
func TestFindDuplicatesWithinRemote_WithHostnamePrefix(t *testing.T) {
	pool := openTestDB(t)
	repo := NewObjectRepository(database.NewPgxDatabase(pool))
	ctx := context.Background()

	// Create a remote with hostname and name
	host := "testhost"
	remName := "testremote"
	rem := remote.NewRemote(host, remName)
	remID, err := repo.WriteRemote(ctx, rem)
	if err != nil {
		t.Fatalf("WriteRemote failed: %v", err)
	}

	// Create duplicate files and associate with remote
	hash := writeRaw(t, ctx, repo, "dup_hash_with_host", "md5")
	meta1 := writeMetaRaw(t, ctx, repo, "file1.txt", "/path/file1.txt", "text/plain")
	meta2 := writeMetaRaw(t, ctx, repo, "file2.txt", "/path/file2.txt", "text/plain")
	writeJunctionRaw(t, ctx, repo, meta1, hash)
	writeJunctionRaw(t, ctx, repo, meta2, hash)

	// Link metadata to remote
	if _, err := repo.db.Exec(ctx, `INSERT INTO object_remote_junction (object_meta_id, remote_id) VALUES ($1, $2), ($3, $2)`, meta1, remID, meta2); err != nil {
		t.Fatalf("failed to link metadata to remote: %v", err)
	}

	// Test 1: Query with "hostname/remote_name" format (as shown in report)
	groups1, err := repo.FindDuplicatesWithinRemote(ctx, host+"/"+remName, "", 0, 0)
	if err != nil {
		t.Fatalf("FindDuplicatesWithinRemote with hostname/name failed: %v", err)
	}
	if len(groups1) == 0 {
		t.Errorf("expected duplicates with hostname/name format, got none")
	}

	// Test 2: Query with just "remote_name" format (traditional)
	groups2, err := repo.FindDuplicatesWithinRemote(ctx, remName, "", 0, 0)
	if err != nil {
		t.Fatalf("FindDuplicatesWithinRemote with name only failed: %v", err)
	}
	if len(groups2) == 0 {
		t.Errorf("expected duplicates with name-only format, got none")
	}

	// Both queries should return the same results
	if len(groups1) != len(groups2) {
		t.Errorf("different results for same remote: hostname/name got %d, name-only got %d", len(groups1), len(groups2))
	}

	// Cleanup
	repo.db.Exec(ctx, `DELETE FROM object_remote_junction WHERE remote_id = $1`, remID)
	repo.db.Exec(ctx, `DELETE FROM object_hash_junction WHERE object_hash_id = $1`, hash)
	repo.db.Exec(ctx, `DELETE FROM object_meta WHERE id IN ($1, $2)`, meta1, meta2)
	repo.db.Exec(ctx, `DELETE FROM object_hash WHERE id = $1`, hash)
	repo.db.Exec(ctx, `DELETE FROM remote WHERE id = $1`, remID)
}

// TestFindDuplicatesWithinRemote_HostnameInNameWithSlash tests edge case:
// If a remote name contains "/" it should be parsed correctly.
func TestFindDuplicatesWithinRemote_HostnamePrefixedPath(t *testing.T) {
	pool := openTestDB(t)
	repo := NewObjectRepository(database.NewPgxDatabase(pool))
	ctx := context.Background()

	// Create a remote representing a local scan path with hostname
	host := "archlinux"
	remName := "/tmp/fs-test.xyz/light" // Local scan path
	rem := remote.NewRemote(host, remName)
	remID, err := repo.WriteRemote(ctx, rem)
	if err != nil {
		t.Fatalf("WriteRemote failed: %v", err)
	}

	// Create duplicates
	hash := writeRaw(t, ctx, repo, "path_hash", "sha256")
	meta1 := writeMetaRaw(t, ctx, repo, "doc.txt", "/doc.txt", "text/plain")
	meta2 := writeMetaRaw(t, ctx, repo, "doc.txt", "/backup/doc.txt", "text/plain")
	writeJunctionRaw(t, ctx, repo, meta1, hash)
	writeJunctionRaw(t, ctx, repo, meta2, hash)

	if _, err := repo.db.Exec(ctx, `INSERT INTO object_remote_junction (object_meta_id, remote_id) VALUES ($1, $2), ($3, $2)`, meta1, remID, meta2); err != nil {
		t.Fatalf("failed to link metadata: %v", err)
	}

	// Query with full "hostname/remote_name" format (as it appears in report)
	groups, err := repo.FindDuplicatesWithinRemote(ctx, host+"/"+remName, "", 0, 0)
	if err != nil {
		t.Fatalf("FindDuplicatesWithinRemote failed: %v", err)
	}
	if len(groups) != 1 {
		t.Errorf("expected 1 duplicate group, got %d", len(groups))
	}

	// Cleanup
	repo.db.Exec(ctx, `DELETE FROM object_remote_junction WHERE remote_id = $1`, remID)
	repo.db.Exec(ctx, `DELETE FROM object_hash_junction WHERE object_hash_id = $1`, hash)
	repo.db.Exec(ctx, `DELETE FROM object_meta WHERE id IN ($1, $2)`, meta1, meta2)
	repo.db.Exec(ctx, `DELETE FROM object_hash WHERE id = $1`, hash)
	repo.db.Exec(ctx, `DELETE FROM remote WHERE id = $1`, remID)
}
