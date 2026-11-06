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
	// ScanRoot is the remote path the scan was started at
	ScanRoot string
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

// WithScanRoot sets the remote path the scan starts from
func WithScanRoot(root string) ScannerOpt {
	return func(s *Scanner) {
		s.ScanRoot = root
	}
}

// Scan handles scanning of a remote
func Scan(ctx context.Context, s *Scanner, lc rclone.ListConfig) error {
	// On first call, ensure location exists for this remote
	if s.Remote.Location.Id == 0 {
		if err := s.Repo.EnsureLocation(ctx, &s.Remote.Location); err != nil {
			return fmt.Errorf("failed to ensure location: %w", err)
		}
	}

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

			// Create Meta with LocationId set
			// For remotes, object_path is the path within the remote
			meta := object.NewMeta(lri.Name, lri.Path, lri.MimeType,
				object.WithMetaSize(lri.Size),
				object.WithMetaScanRoot(s.ScanRoot),
				object.WithMetaLocationId(s.Remote.LocationId))

			// Persist meta and hashes
			err := s.persistRcloneFile(ctx, &meta, lri.Hashes)
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

// persistRcloneFile persists a single rclone file with its hashes
func (s *Scanner) persistRcloneFile(ctx context.Context, m *object.Meta, hashes map[string]string) error {
	// Check if meta already persisted
	err := s.Repo.ReadMeta(ctx, m)
	if err != nil {
		return fmt.Errorf("failed to get object meta from db: %w", err)
	}
	if !m.Persisted {
		err = s.Repo.WriteMeta(ctx, m)
		if err != nil {
			return fmt.Errorf("failed to put object into db: %w", err)
		}
	}

	// Persist hashes
	for hType, hVal := range hashes {
		h := object.NewHash(hVal, hType)

		err := s.Repo.ReadHash(ctx, &h)
		if err != nil {
			return fmt.Errorf("failed to get hash from db: %w", err)
		}
		if !h.Persisted {
			err = s.Repo.WriteHash(ctx, &h)
			if err != nil {
				return fmt.Errorf("failed to put hash into db: %w", err)
			}
		}

		// Link meta to hash
		mhj := object.NewMetaHashJunction(m.Id, h.Id)
		err = s.Repo.ReadMetaHashJunction(ctx, &mhj)
		if err != nil {
			return fmt.Errorf("failed to get object meta hash junction from db: %w", err)
		}
		if !mhj.Persisted {
			err = s.Repo.WriteMetaHashJunction(ctx, &mhj)
			if err != nil {
				return fmt.Errorf("failed to set object meta-hash junction in db: %w", err)
			}
		}
	}

	return nil
}
