package operation

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/haggishunk/filesprawl/internal/object"
	"github.com/haggishunk/filesprawl/internal/remote"
)

func TestLocalScannerDefaults(t *testing.T) {
	ls := NewLocalScanner()
	if ls.Root != "" {
		t.Fatalf("expected empty root, got %s", ls.Root)
	}
	if ls.Repo != nil {
		t.Fatalf("expected nil repo")
	}
}

func TestLocalScannerWithOptions(t *testing.T) {
	ls := NewLocalScanner(
		WithLocalRoot("/tmp/test"),
	)
	if ls.Root != "/tmp/test" {
		t.Fatalf("expected root /tmp/test, got %s", ls.Root)
	}
}

func TestComputeHashesValidFile(t *testing.T) {
	// Create a temporary file
	tmpFile, err := os.CreateTemp("", "test_hash_*")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	content := []byte("test content for hashing")
	if _, err := tmpFile.Write(content); err != nil {
		t.Fatalf("failed to write test content: %v", err)
	}

	// Seek back to beginning
	tmpFile.Seek(0, 0)

	hashes := make(map[string]string)
	if err := computeHashes(tmpFile, hashes); err != nil {
		t.Fatalf("computeHashes failed: %v", err)
	}

	tmpFile.Close()

	// Verify we have hashes
	expectedHashes := []string{"md5", "sha1", "sha256"}
	for _, hashType := range expectedHashes {
		if _, ok := hashes[hashType]; !ok {
			t.Fatalf("missing hash type: %s", hashType)
		}
	}
}

func TestScanLocalInvalidPath(t *testing.T) {
	ls := NewLocalScanner(WithLocalRoot(""))
	ctx := context.Background()

	err := ScanLocal(ctx, &ls)
	if err == nil {
		t.Fatalf("expected error for empty root, got nil")
	}
}

func TestScanLocalNonExistentPath(t *testing.T) {
	ls := NewLocalScanner(WithLocalRoot("/nonexistent/path/that/does/not/exist"))
	ctx := context.Background()

	err := ScanLocal(ctx, &ls)
	if err == nil {
		t.Fatalf("expected error for nonexistent path, got nil")
	}
}

func TestScanLocalFilePath(t *testing.T) {
	// Create a temp file
	tmpFile, err := os.CreateTemp("", "test_*")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	tmpFile.Close()

	// Try to scan a file (not directory)
	ls := NewLocalScanner(WithLocalRoot(tmpFile.Name()))
	ctx := context.Background()

	err = ScanLocal(ctx, &ls)
	if err == nil {
		t.Fatalf("expected error when scanning a file instead of directory")
	}
}

// recordedLocalPersist captures one PersistLocalResult call
type recordedLocalPersist struct {
	rem    remote.Remote
	meta   object.Meta
	hashes map[string]string
}

// fakeLocalPersister records persisted entries instead of writing to a database
type fakeLocalPersister struct {
	calls []recordedLocalPersist
}

func (f *fakeLocalPersister) PersistLocalResult(_ context.Context, rem *remote.Remote, m *object.Meta, hashes map[string]string) error {
	f.calls = append(f.calls, recordedLocalPersist{rem: *rem, meta: *m, hashes: hashes})
	return nil
}

func newFakeLocalScanner(root string) (*LocalScanner, *fakeLocalPersister) {
	fake := &fakeLocalPersister{}
	ls := NewLocalScanner(WithLocalRoot(root))
	ls.Repo = fake
	return &ls, fake
}

func resolvedRoot(t *testing.T, root string) string {
	t.Helper()
	if resolved, err := filepath.EvalSymlinks(root); err == nil {
		return resolved
	}
	return root
}

func persistsByPath(fake *fakeLocalPersister) map[string]recordedLocalPersist {
	byPath := make(map[string]recordedLocalPersist, len(fake.calls))
	for _, c := range fake.calls {
		byPath[c.meta.Path] = c
	}
	return byPath
}

func TestScanLocalSkipsDirectorySymlink(t *testing.T) {
	root := t.TempDir()

	subDir := filepath.Join(root, "subdir")
	if err := os.Mkdir(subDir, 0755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(subDir, "nested.txt"), []byte("nested"), 0644); err != nil {
		t.Fatalf("write nested: %v", err)
	}
	// Symlink to the directory: must not crash the scan or be persisted
	if err := os.Symlink(subDir, filepath.Join(root, "dirlink")); err != nil {
		t.Fatalf("symlink dirlink: %v", err)
	}

	ls, fake := newFakeLocalScanner(root)
	if err := ScanLocal(context.Background(), ls); err != nil {
		t.Fatalf("ScanLocal failed: %v", err)
	}

	byPath := persistsByPath(fake)
	if len(fake.calls) != 1 {
		t.Fatalf("expected only subdir/nested.txt to be persisted, got %d entries", len(fake.calls))
	}
	if _, ok := byPath[filepath.Join("subdir", "nested.txt")]; !ok {
		t.Fatalf("expected subdir/nested.txt to be persisted, got %v", byPath)
	}
	if _, ok := byPath["dirlink"]; ok {
		t.Fatalf("directory symlink must not be persisted")
	}
}

func TestScanLocalMarksFileSymlink(t *testing.T) {
	root := t.TempDir()

	if err := os.WriteFile(filepath.Join(root, "real.txt"), []byte("content"), 0644); err != nil {
		t.Fatalf("write real: %v", err)
	}
	if err := os.Symlink(filepath.Join(root, "real.txt"), filepath.Join(root, "link.txt")); err != nil {
		t.Fatalf("symlink link.txt: %v", err)
	}

	ls, fake := newFakeLocalScanner(root)
	if err := ScanLocal(context.Background(), ls); err != nil {
		t.Fatalf("ScanLocal failed: %v", err)
	}

	byPath := persistsByPath(fake)
	if len(fake.calls) != 2 {
		t.Fatalf("expected 2 persisted entries, got %d", len(fake.calls))
	}

	link, ok := byPath["link.txt"]
	if !ok {
		t.Fatalf("expected link.txt to be persisted")
	}
	if !link.meta.IsSymlink {
		t.Fatalf("expected link.txt to be marked as symlink")
	}
	if len(link.hashes) != 0 {
		t.Fatalf("symlink must not be hashed, got %v", link.hashes)
	}

	real, ok := byPath["real.txt"]
	if !ok {
		t.Fatalf("expected real.txt to be persisted")
	}
	if real.meta.IsSymlink {
		t.Fatalf("real.txt must not be marked as symlink")
	}
	for _, ht := range []string{"md5", "sha1", "sha256"} {
		if _, ok := real.hashes[ht]; !ok {
			t.Fatalf("expected %s hash for real.txt", ht)
		}
	}
}

func TestScanLocalRecordsBrokenSymlink(t *testing.T) {
	root := t.TempDir()

	if err := os.Symlink(filepath.Join(root, "missing.txt"), filepath.Join(root, "broken.txt")); err != nil {
		t.Fatalf("symlink broken.txt: %v", err)
	}

	ls, fake := newFakeLocalScanner(root)
	if err := ScanLocal(context.Background(), ls); err != nil {
		t.Fatalf("ScanLocal failed: %v", err)
	}

	byPath := persistsByPath(fake)
	link, ok := byPath["broken.txt"]
	if !ok {
		t.Fatalf("expected broken.txt to be persisted")
	}
	if !link.meta.IsSymlink {
		t.Fatalf("expected broken.txt to be marked as symlink")
	}
	if len(link.hashes) != 0 {
		t.Fatalf("broken symlink must not be hashed, got %v", link.hashes)
	}
}

func TestScanLocalScanRootTracking(t *testing.T) {
	root := t.TempDir()

	if err := os.WriteFile(filepath.Join(root, "file.txt"), []byte("data"), 0644); err != nil {
		t.Fatalf("write file: %v", err)
	}

	ls, fake := newFakeLocalScanner(root)
	if err := ScanLocal(context.Background(), ls); err != nil {
		t.Fatalf("ScanLocal failed: %v", err)
	}

	if len(fake.calls) != 1 {
		t.Fatalf("expected 1 persisted entry, got %d", len(fake.calls))
	}

	want := resolvedRoot(t, root)
	call := fake.calls[0]
	if call.meta.ScanRoot != want {
		t.Fatalf("expected scan root %q, got %q", want, call.meta.ScanRoot)
	}
	if call.meta.Path != "file.txt" {
		t.Fatalf("expected path relative to scan root, got %q", call.meta.Path)
	}
	if call.rem.Name != want || !call.rem.Location.IsLocal() {
		t.Fatalf("expected local remote named %q, got %q", want, call.rem.Name)
	}
}
