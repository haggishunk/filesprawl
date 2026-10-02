package operation

import (
	"context"
	"os"
	"testing"
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
