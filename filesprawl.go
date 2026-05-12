// Package filesprawl provides file indexing and duplicate detection across
// multiple remote storage locations.
//
// The package re-exports the core types and constructors from the internal
// analysis and repository packages so that external Go programs can import
// github.com/haggishunk/filesprawl directly without depending on internal
// packages.
//
// Example:
//
//	import "github.com/haggishunk/filesprawl"
//
//	repo := filesprawl.NewObjectRepository(db)
//	detector := filesprawl.NewDuplicateDetector(repo)
//	groups, err := detector.FindAcrossRemotes(ctx, filesprawl.FilterOptions{
//		HashType: "sha256",
//		MinSize:  1 << 20,
//	})
package filesprawl

import (
	"github.com/haggishunk/filesprawl/internal/analysis"
	"github.com/haggishunk/filesprawl/internal/database"
	"github.com/haggishunk/filesprawl/internal/repository"
)

// DuplicateDetector finds duplicate files using content hashes recorded in the
// repository.
type DuplicateDetector = analysis.DuplicateDetector

// FilterOptions controls which duplicate groups are returned by the detector.
type FilterOptions = analysis.FilterOptions

// DuplicateGroup is a set of files that share the same content hash.
type DuplicateGroup = analysis.DuplicateGroup

// FileInfo holds per-file detail within a DuplicateGroup.
type FileInfo = analysis.FileInfo

// ObjectRepository persists indexed file metadata and exposes the queries used
// by DuplicateDetector.
type ObjectRepository = repository.ObjectRepository

// Database is the persistence interface satisfied by concrete database
// implementations (for example, the pgx-backed implementation).
type Database = database.Database

// NewDuplicateDetector creates a DuplicateDetector backed by the supplied
// repository.
func NewDuplicateDetector(repo *ObjectRepository) *DuplicateDetector {
	return analysis.NewDuplicateDetector(repo)
}

// NewObjectRepository constructs an ObjectRepository wrapping the given
// Database implementation.
func NewObjectRepository(db Database) *ObjectRepository {
	return repository.NewObjectRepository(db)
}

// FormatDuplicateReport returns a human-readable report of the supplied
// duplicate groups.
func FormatDuplicateReport(groups []DuplicateGroup) string {
	return analysis.FormatReport(groups)
}
