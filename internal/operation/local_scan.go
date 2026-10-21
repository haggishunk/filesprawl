package operation

import (
	"context"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"fmt"
	"hash"
	"io"
	"os"
	"path/filepath"

	"github.com/haggishunk/filesprawl/internal/object"
	"github.com/haggishunk/filesprawl/internal/remote"
	"github.com/haggishunk/filesprawl/internal/repository"
)

// localPersister is the subset of repository behavior that local scanning
// relies on. It allows tests to substitute a fake.
type localPersister interface {
	PersistLocalResult(ctx context.Context, rem *remote.Remote, m *object.Meta, hashes map[string]string) error
}

// LocalScanner walks the local filesystem and computes hashes
type LocalScanner struct {
	Root string
	Repo localPersister
}

type LocalScannerOpt func(s *LocalScanner)

func NewLocalScanner(options ...LocalScannerOpt) LocalScanner {
	s := LocalScanner{}
	for _, opt := range options {
		opt(&s)
	}
	return s
}

func WithLocalRoot(root string) LocalScannerOpt {
	return func(s *LocalScanner) {
		s.Root = root
	}
}

func WithLocalRepo(r *repository.ObjectRepository) LocalScannerOpt {
	return func(s *LocalScanner) {
		s.Repo = r
	}
}

// ScanLocal walks a local filesystem root and persists file hashes
func ScanLocal(ctx context.Context, ls *LocalScanner) error {
	if ls.Root == "" {
		return fmt.Errorf("local root path is required")
	}

	absRoot, err := filepath.Abs(ls.Root)
	if err != nil {
		return fmt.Errorf("failed to resolve local root path: %w", err)
	}
	// Resolve symlinks in the root itself so scan provenance is stable
	if resolved, err := filepath.EvalSymlinks(absRoot); err == nil {
		absRoot = resolved
	}

	info, err := os.Stat(absRoot)
	if err != nil {
		return fmt.Errorf("local root path not accessible: %w", err)
	}
	if !info.IsDir() {
		return fmt.Errorf("local root path must be a directory")
	}

	if ls.Repo == nil {
		return fmt.Errorf("local repo is required")
	}

	return filepath.WalkDir(absRoot, func(filePath string, d os.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("error accessing path %s: %w", filePath, err)
		}

		if d.IsDir() {
			return nil
		}

		relPath, err := filepath.Rel(absRoot, filePath)
		if err != nil {
			return fmt.Errorf("failed to determine path of %s relative to %s: %w", filePath, absRoot, err)
		}

		// WalkDir reports symlinks without following them
		if d.Type()&os.ModeSymlink != 0 {
			return persistLocalSymlink(ctx, ls, absRoot, filePath, relPath)
		}

		return persistLocalFile(ctx, ls, absRoot, filePath, relPath)
	})
}

func persistLocalFile(ctx context.Context, ls *LocalScanner, scanRoot string, filePath string, relPath string) error {
	info, err := os.Stat(filePath)
	if err != nil {
		return fmt.Errorf("failed to stat file %s: %w", filePath, err)
	}

	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file %s: %w", filePath, err)
	}
	defer file.Close()

	// Compute hashes
	hashes := make(map[string]string)
	if err := computeHashes(file, hashes); err != nil {
		return fmt.Errorf("failed to hash file %s: %w", filePath, err)
	}

	// Create metadata object
	meta := object.NewMeta(info.Name(), relPath, "application/octet-stream",
		object.WithMetaSize(info.Size()),
		object.WithMetaScanRoot(scanRoot))

	// Create local remote
	rem, err := localRemote(scanRoot)
	if err != nil {
		return err
	}

	// Persist
	return ls.Repo.PersistLocalResult(ctx, &rem, &meta, hashes)
}

// persistLocalSymlink records a symlink without hashing its contents, so the
// link can never be reported as a duplicate of its target. Directory symlinks
// are skipped entirely (WalkDir does not descend into them); broken symlinks
// are recorded with a warning.
func persistLocalSymlink(ctx context.Context, ls *LocalScanner, scanRoot string, filePath string, relPath string) error {
	target, err := os.Stat(filePath)
	switch {
	case err == nil && target.IsDir():
		// Directory symlink: skip it
		return nil
	case err != nil:
		fmt.Printf("Warning: symlink target not resolvable for %s: %v\n", filePath, err)
	}

	info, err := os.Lstat(filePath)
	if err != nil {
		return fmt.Errorf("failed to lstat symlink %s: %w", filePath, err)
	}

	meta := object.NewMeta(info.Name(), relPath, "application/octet-stream",
		object.WithMetaSize(info.Size()),
		object.WithMetaSymlink(true),
		object.WithMetaScanRoot(scanRoot))

	rem, err := localRemote(scanRoot)
	if err != nil {
		return err
	}

	// Symlinks are persisted without hashes so they never match file content
	return ls.Repo.PersistLocalResult(ctx, &rem, &meta, nil)
}

// localRemote builds the persistence identity of a local scan root
func localRemote(scanRoot string) (remote.Remote, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return remote.Remote{}, fmt.Errorf("failed to get hostname: %w", err)
	}
	return remote.NewRemote(hostname, scanRoot, remote.WithRemoteType("local")), nil
}

func computeHashes(file *os.File, hashes map[string]string) error {
	hashers := map[string]hash.Hash{
		"md5":    md5.New(),
		"sha1":   sha1.New(),
		"sha256": sha256.New(),
	}

	// Collect all hashers as writers
	var writers []io.Writer
	for _, h := range hashers {
		writers = append(writers, h)
	}

	// Create multi-writer to write to all hashers at once
	mw := io.MultiWriter(writers...)

	// Copy file content to all hashers
	if _, err := io.Copy(mw, file); err != nil {
		return err
	}

	// Get hash values
	for typ, h := range hashers {
		hashes[typ] = fmt.Sprintf("%x", h.Sum(nil))
	}
	return nil
}
