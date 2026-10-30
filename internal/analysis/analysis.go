package analysis

import (
	"context"
	"fmt"
	"strings"

	"github.com/haggishunk/filesprawl/internal/repository"
)

// FilterOptions controls which duplicates are returned.
type FilterOptions struct {
	// HashType restricts detection to a specific hash type (md5, sha1, sha256, dropbox).
	// Empty string means all hash types.
	HashType string
	// MinSize excludes files smaller than this many bytes. 0 means no size filter.
	MinSize int64
	// Limit is the maximum number of duplicate groups to return. 0 means no limit.
	Limit int
	// Offset skips this many duplicate groups for pagination.
	Offset int
}

// FileInfo holds per-file detail within a duplicate group.
type FileInfo struct {
	ID         int
	Name       string
	Path       string
	MimeType   string
	Size       int64
	RemoteName string
	Hostname   string
	ScanRoot   string
}

// DuplicateGroup is a set of files that share the same content hash.
type DuplicateGroup struct {
	HashValue string
	HashType  string
	Count     int
	Files     []FileInfo
}

// duplicateRepo is the subset of ObjectRepository used by DuplicateDetector.
type duplicateRepo interface {
	FindDuplicatesAcrossRemotes(ctx context.Context, hashType string, limit, offset int) ([]repository.HashGroup, error)
	FindDuplicatesWithinRemote(ctx context.Context, remoteName, hashType string, limit, offset int) ([]repository.HashGroup, error)
	GetFilesByHash(ctx context.Context, hashValue, hashType string, minSize int64) ([]repository.FileRecord, error)
	GetFilesByHashes(ctx context.Context, hashes []repository.HashTypePair, minSize int64) (repository.FilesByHashMap, error)
}

// DuplicateDetector finds and groups duplicate files using the repository.
type DuplicateDetector struct {
	repo duplicateRepo
}

// NewDuplicateDetector creates a DuplicateDetector backed by the given repository.
func NewDuplicateDetector(r duplicateRepo) *DuplicateDetector {
	return &DuplicateDetector{repo: r}
}

// FindAcrossRemotes returns duplicate groups across all remotes.
func (d *DuplicateDetector) FindAcrossRemotes(ctx context.Context, opts FilterOptions) ([]DuplicateGroup, error) {
	groups, err := d.repo.FindDuplicatesAcrossRemotes(ctx, opts.HashType, opts.Limit, opts.Offset)
	if err != nil {
		return nil, fmt.Errorf("failed to find duplicates across remotes: %w", err)
	}
	return d.enrich(ctx, groups, opts.MinSize)
}

// FindWithinRemote returns duplicate groups within a single named remote.
func (d *DuplicateDetector) FindWithinRemote(ctx context.Context, remoteName string, opts FilterOptions) ([]DuplicateGroup, error) {
	if remoteName == "" {
		return nil, fmt.Errorf("remoteName must not be empty")
	}
	groups, err := d.repo.FindDuplicatesWithinRemote(ctx, remoteName, opts.HashType, opts.Limit, opts.Offset)
	if err != nil {
		return nil, fmt.Errorf("failed to find duplicates within remote %q: %w", remoteName, err)
	}
	return d.enrich(ctx, groups, opts.MinSize)
}

// enrich fetches file details for all HashGroups in a single batch query and assembles DuplicateGroups.
func (d *DuplicateDetector) enrich(ctx context.Context, groups []repository.HashGroup, minSize int64) ([]DuplicateGroup, error) {
	if len(groups) == 0 {
		return []DuplicateGroup{}, nil
	}

	// Convert groups to HashTypePair for batch query
	hashes := make([]repository.HashTypePair, len(groups))
	for i, g := range groups {
		hashes[i] = repository.HashTypePair{
			HashValue: g.HashValue,
			HashType:  g.HashType,
		}
	}

	// Fetch all files for all hashes in a single query
	filesByHash, err := d.repo.GetFilesByHashes(ctx, hashes, minSize)
	if err != nil {
		return nil, fmt.Errorf("failed to get files for hashes: %w", err)
	}

	// Build result from batch query data
	result := make([]DuplicateGroup, 0, len(groups))
	for _, g := range groups {
		key := g.HashValue + ":" + g.HashType
		records := filesByHash[key]

		// After applying minSize filter, skip groups that no longer have duplicates.
		if len(records) < 2 {
			continue
		}

		files := make([]FileInfo, len(records))
		for i, rec := range records {
			files[i] = FileInfo{
				ID:         rec.ID,
				Name:       rec.Name,
				Path:       rec.Path,
				MimeType:   rec.MimeType,
				Size:       rec.Size,
				RemoteName: rec.RemoteName,
				Hostname:   rec.Hostname,
				ScanRoot:   rec.ScanRoot,
			}
		}
		result = append(result, DuplicateGroup{
			HashValue: g.HashValue,
			HashType:  g.HashType,
			Count:     len(records),
			Files:     files,
		})
	}
	return result, nil
}

// FormatReport returns a human-readable report of duplicate groups.
func FormatReport(groups []DuplicateGroup) string {
	if len(groups) == 0 {
		return "No duplicates found.\n"
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "Found %d duplicate group(s):\n\n", len(groups))
	for i, g := range groups {
		fmt.Fprintf(&sb, "Group %d — %s (%s), %d copies\n", i+1, g.HashValue[:min(len(g.HashValue), 12)]+"...", g.HashType, g.Count)
		for _, f := range g.Files {
			remote := f.RemoteName
			if f.Hostname != "" {
				remote = f.Hostname + "/" + f.RemoteName
			}
			root := ""
			if f.ScanRoot != "" {
				root = "  root: " + f.ScanRoot
			}
			fmt.Fprintf(&sb, "  • %s  [%s]  %d bytes  remote: %s%s\n", f.Name, f.Path, f.Size, remote, root)
		}
		sb.WriteString("\n")
	}
	return sb.String()
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
