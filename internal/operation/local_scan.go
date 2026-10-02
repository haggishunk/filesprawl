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

// LocalScanner walks the local filesystem and computes hashes
type LocalScanner struct {
	Root string
	Repo *repository.ObjectRepository
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

	info, err := os.Stat(absRoot)
	if err != nil {
		return fmt.Errorf("local root path not accessible: %w", err)
	}
	if !info.IsDir() {
		return fmt.Errorf("local root path must be a directory")
	}

	return filepath.WalkDir(absRoot, func(filePath string, d os.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("error accessing path %s: %w", filePath, err)
		}

		if d.IsDir() {
			return nil
		}

		relPath, _ := filepath.Rel(absRoot, filePath)
		return persistLocalFile(ctx, ls, filePath, relPath)
	})
}

func persistLocalFile(ctx context.Context, ls *LocalScanner, filePath string, relPath string) error {
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
		object.WithMetaSize(info.Size()))

	// Create local remote
	hostname, err := os.Hostname()
	if err != nil {
		return fmt.Errorf("failed to get hostname: %w", err)
	}
	rem := remote.NewRemote(hostname, ls.Root,
		remote.WithRemoteType("local"))

	// Persist
	return ls.Repo.PersistLocalResult(ctx, &rem, &meta, hashes)
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
