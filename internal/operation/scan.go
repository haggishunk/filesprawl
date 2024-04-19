package operation

import (
	"context"
	"fmt"

	"github.com/haggishunk/filesprawl/internal/object"
	"github.com/haggishunk/filesprawl/internal/rclone"
	"github.com/haggishunk/filesprawl/internal/remote"
	"github.com/haggishunk/filesprawl/internal/repository"
)

// scanning:
// recursively descending through directories
// retrieving objects and persisting to a database
// persist object aspects in relationships with scan timestamp

// synthesize a list response with the targeted remote
type ScanResult struct {
	Remote *remote.Remote
	Hash   *object.Hash
	Object *object.Meta
	// Timestamp
}

// Scanner is a:
// reader and writer
// lister and persister
// set it loose on a host, remote and path
// and it will save traverse saving state
// picking up where it left off last time
type Scanner struct {
	Remote *remote.Remote
	Repo   *repository.ObjectRepository
	// writer channel
}

type ScannerOpt func(s *Scanner)

func NewScanner(options ...ScannerOpt) Scanner {
	s := Scanner{}

	for _, opt := range options {
		opt(&s)
	}

	return s
}

func WithRemote(r *remote.Remote) ScannerOpt {
	return func(s *Scanner) {
		s.Remote = r
	}
}

func WithRepo(r *repository.ObjectRepository) ScannerOpt {
	return func(s *Scanner) {
		s.Repo = r
	}
}

// Scan handles scanning of a remote
func Scan(ctx context.Context, s *Scanner, lc rclone.ListConfig) error {
	// construct initial options, eg path
	lr, err := rclone.ListJSON(ctx, lc)
	if err != nil {
		return fmt.Errorf("failed to list json: %w", err)
	}
	// when dirs only -- iterate and recursively descend
	// starting with files only
	if lc.Opt.DirsOnly {
		for _, lri := range lr.List {
			fmt.Printf("Descending into path: %s\n", lri.Path)
			newLo := rclone.NewListOption(rclone.ListOptionFilesOnly())
			newLc := rclone.NewListConfig(lc.Fs, lri.Path, &newLo)
			err := Scan(ctx, s, newLc)
			if err != nil {
				return fmt.Errorf("failed scan: %w", err)
			}
		}
	}
	// when files only -- persist
	// then rescan same path but for dirs only
	if lc.Opt.FilesOnly {
		// do files iteration and persistence
		for _, lri := range lr.List {
			fmt.Printf("Found file: %s\n", lri)
			err := s.Repo.PersistResult(ctx, lri)
			if err != nil {
				return fmt.Errorf("failed to persist: %w", err)
			}
		}
		newLo := rclone.NewListOption(rclone.ListOptionDirsOnly())
		newLc := rclone.NewListConfig(lc.Fs, lc.Remote, &newLo)
		err := Scan(ctx, s, newLc)
		if err != nil {
			return fmt.Errorf("failed scan: %w", err)
		}
	}
	return nil
}
