package analysis

import (
	"context"
	"errors"
	"testing"

	"github.com/haggishunk/filesprawl/internal/repository"
)

// mockRepo implements duplicateRepo for testing.
type mockRepo struct {
	hashGroups  []repository.HashGroup
	fileRecords []repository.FileRecord
	filesMap    repository.FilesByHashMap
	acrossErr   error
	withinErr   error
	filesErr    error
	batchErr    error
}

func (m *mockRepo) FindDuplicatesAcrossRemotes(_ context.Context, _ string, _, _ int) ([]repository.HashGroup, error) {
	return m.hashGroups, m.acrossErr
}

func (m *mockRepo) FindDuplicatesWithinRemote(_ context.Context, _, _ string, _, _ int) ([]repository.HashGroup, error) {
	return m.hashGroups, m.withinErr
}

func (m *mockRepo) GetFilesByHash(_ context.Context, _, _ string, _ int64) ([]repository.FileRecord, error) {
	return m.fileRecords, m.filesErr
}

func (m *mockRepo) GetFilesByHashes(_ context.Context, _ []repository.HashTypePair, _ int64) (repository.FilesByHashMap, error) {
	return m.filesMap, m.batchErr
}

var twoFiles = []repository.FileRecord{
	{ID: 1, Name: "a.txt", Path: "/docs/a.txt", Size: 1024, RemoteName: "dbox:", Hostname: "host1"},
	{ID: 2, Name: "a.txt", Path: "/backup/a.txt", Size: 1024, RemoteName: "s3:", Hostname: "host2"},
}

func TestFindAcrossRemotes_ReturnsDuplicateGroups(t *testing.T) {
	filesMap := repository.FilesByHashMap{
		"abc123:md5": twoFiles,
	}
	repo := &mockRepo{
		hashGroups: []repository.HashGroup{{HashValue: "abc123", HashType: "md5", Count: 2}},
		filesMap:   filesMap,
	}
	det := NewDuplicateDetector(repo)
	groups, err := det.FindAcrossRemotes(context.Background(), FilterOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(groups) != 1 {
		t.Fatalf("expected 1 group, got %d", len(groups))
	}
	if groups[0].HashValue != "abc123" {
		t.Errorf("expected hash abc123, got %s", groups[0].HashValue)
	}
	if len(groups[0].Files) != 2 {
		t.Errorf("expected 2 files, got %d", len(groups[0].Files))
	}
}

func TestFindAcrossRemotes_RepoError(t *testing.T) {
	repo := &mockRepo{acrossErr: errors.New("db down")}
	det := NewDuplicateDetector(repo)
	_, err := det.FindAcrossRemotes(context.Background(), FilterOptions{})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestFindWithinRemote_EmptyNameErrors(t *testing.T) {
	det := NewDuplicateDetector(&mockRepo{})
	_, err := det.FindWithinRemote(context.Background(), "", FilterOptions{})
	if err == nil {
		t.Fatal("expected error for empty remote name, got nil")
	}
}

func TestFindWithinRemote_ReturnsDuplicates(t *testing.T) {
	filesMap := repository.FilesByHashMap{
		"def456:sha256": twoFiles,
	}
	repo := &mockRepo{
		hashGroups: []repository.HashGroup{{HashValue: "def456", HashType: "sha256", Count: 2}},
		filesMap:   filesMap,
	}
	det := NewDuplicateDetector(repo)
	groups, err := det.FindWithinRemote(context.Background(), "dbox:", FilterOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(groups) != 1 {
		t.Fatalf("expected 1 group, got %d", len(groups))
	}
}

func TestFindAcrossRemotes_MinSizeFiltersSmallFiles(t *testing.T) {
	// minSize filter is applied by GetFilesByHash at the DB layer.
	// If fewer than 2 records come back, the group is dropped.
	repo := &mockRepo{
		hashGroups:  []repository.HashGroup{{HashValue: "abc123", HashType: "md5", Count: 2}},
		fileRecords: twoFiles[:1], // only 1 file returned (below min size)
	}
	det := NewDuplicateDetector(repo)
	groups, err := det.FindAcrossRemotes(context.Background(), FilterOptions{MinSize: 2048})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(groups) != 0 {
		t.Errorf("expected 0 groups after size filter, got %d", len(groups))
	}
}

func TestFindAcrossRemotes_NoDuplicates(t *testing.T) {
	repo := &mockRepo{hashGroups: nil}
	det := NewDuplicateDetector(repo)
	groups, err := det.FindAcrossRemotes(context.Background(), FilterOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(groups) != 0 {
		t.Errorf("expected 0 groups, got %d", len(groups))
	}
}

func TestFormatReport_Empty(t *testing.T) {
	out := FormatReport(nil)
	if out != "No duplicates found.\n" {
		t.Errorf("unexpected output: %q", out)
	}
}

func TestFormatReport_WithGroups(t *testing.T) {
	groups := []DuplicateGroup{
		{
			HashValue: "abc123def456",
			HashType:  "md5",
			Count:     2,
			Files: []FileInfo{
				{Name: "a.txt", Path: "/docs/a.txt", Size: 1024, RemoteName: "dbox:", Hostname: "host1"},
				{Name: "a.txt", Path: "/backup/a.txt", Size: 1024, RemoteName: "s3:", Hostname: "host2"},
			},
		},
	}
	out := FormatReport(groups)
	if out == "No duplicates found.\n" {
		t.Error("expected non-empty report")
	}
	if len(out) == 0 {
		t.Error("report should not be empty")
	}
}

// TestFindAcrossRemotes_UsesBatchQuery verifies that FindAcrossRemotes uses GetFilesByHashes (batch query)
// instead of calling GetFilesByHash multiple times (N+1 problem).
func TestFindAcrossRemotes_UsesBatchQuery(t *testing.T) {
	// Create a batch-ready map with multiple hash groups
	filesMap := repository.FilesByHashMap{
		"abc123:md5": {
			{ID: 1, Name: "file1.txt", Path: "/path/file1.txt", Size: 1024, RemoteName: "dbox:", Hostname: "host1"},
			{ID: 2, Name: "file1.txt", Path: "/backup/file1.txt", Size: 1024, RemoteName: "s3:", Hostname: "host2"},
		},
		"def456:sha256": {
			{ID: 3, Name: "file2.pdf", Path: "/docs/file2.pdf", Size: 2048, RemoteName: "gdrive:", Hostname: "host1"},
			{ID: 4, Name: "file2.pdf", Path: "/archive/file2.pdf", Size: 2048, RemoteName: "s3:", Hostname: "host2"},
		},
	}

	repo := &mockRepo{
		hashGroups: []repository.HashGroup{
			{HashValue: "abc123", HashType: "md5", Count: 2},
			{HashValue: "def456", HashType: "sha256", Count: 2},
		},
		filesMap: filesMap,
	}

	det := NewDuplicateDetector(repo)
	groups, err := det.FindAcrossRemotes(context.Background(), FilterOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(groups) != 2 {
		t.Fatalf("expected 2 groups, got %d", len(groups))
	}

	// Verify first group
	if groups[0].HashValue != "abc123" || groups[0].HashType != "md5" || len(groups[0].Files) != 2 {
		t.Errorf("group 0: expected (abc123, md5, 2 files), got (%s, %s, %d files)",
			groups[0].HashValue, groups[0].HashType, len(groups[0].Files))
	}

	// Verify second group
	if groups[1].HashValue != "def456" || groups[1].HashType != "sha256" || len(groups[1].Files) != 2 {
		t.Errorf("group 1: expected (def456, sha256, 2 files), got (%s, %s, %d files)",
			groups[1].HashValue, groups[1].HashType, len(groups[1].Files))
	}
}

// TestFindAcrossRemotes_BatchQueryWithMinSize tests that batch query respects size filtering.
func TestFindAcrossRemotes_BatchQueryWithMinSize(t *testing.T) {
	// Simulate: 2 hash groups, but one ends up with only 1 file after min-size filter
	filesMap := repository.FilesByHashMap{
		"abc123:md5": {
			{ID: 1, Name: "large.bin", Path: "/large.bin", Size: 5000, RemoteName: "dbox:", Hostname: "host1"},
			{ID: 2, Name: "large.bin", Path: "/backup/large.bin", Size: 5000, RemoteName: "s3:", Hostname: "host2"},
		},
		"small123:md5": {
			{ID: 3, Name: "tiny.txt", Path: "/tiny.txt", Size: 100, RemoteName: "dbox:", Hostname: "host1"},
			// Second file filtered out by minSize, so only 1 remains
		},
	}

	repo := &mockRepo{
		hashGroups: []repository.HashGroup{
			{HashValue: "abc123", HashType: "md5", Count: 2},
			{HashValue: "small123", HashType: "md5", Count: 2},
		},
		filesMap: filesMap,
	}

	det := NewDuplicateDetector(repo)
	groups, err := det.FindAcrossRemotes(context.Background(), FilterOptions{MinSize: 1000})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Only the large file group should remain (small group dropped because < 2 files)
	if len(groups) != 1 {
		t.Fatalf("expected 1 group after size filter, got %d", len(groups))
	}
	if groups[0].HashValue != "abc123" {
		t.Errorf("expected remaining group to be abc123, got %s", groups[0].HashValue)
	}
}

// TestFindAcrossRemotes_BatchQueryBatchError tests error handling in batch query.
func TestFindAcrossRemotes_BatchQueryBatchError(t *testing.T) {
	repo := &mockRepo{
		hashGroups: []repository.HashGroup{
			{HashValue: "abc123", HashType: "md5", Count: 2},
		},
		batchErr: errors.New("batch query failed"),
	}

	det := NewDuplicateDetector(repo)
	_, err := det.FindAcrossRemotes(context.Background(), FilterOptions{})
	if err == nil {
		t.Fatal("expected error from batch query, got nil")
	}
}

// TestFindWithinRemote_UsesBatchQuery verifies batch query works for within-remote duplicates.
func TestFindWithinRemote_UsesBatchQuery(t *testing.T) {
	filesMap := repository.FilesByHashMap{
		"xyz789:sha1": {
			{ID: 1, Name: "doc.docx", Path: "/doc.docx", Size: 50000, RemoteName: "dbox:"},
			{ID: 2, Name: "doc.docx", Path: "/archives/doc.docx", Size: 50000, RemoteName: "dbox:"},
		},
	}

	repo := &mockRepo{
		hashGroups: []repository.HashGroup{
			{HashValue: "xyz789", HashType: "sha1", Count: 2},
		},
		filesMap: filesMap,
	}

	det := NewDuplicateDetector(repo)
	groups, err := det.FindWithinRemote(context.Background(), "dbox:", FilterOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(groups) != 1 {
		t.Fatalf("expected 1 group, got %d", len(groups))
	}
	if len(groups[0].Files) != 2 {
		t.Errorf("expected 2 files in group, got %d", len(groups[0].Files))
	}
}

// TestEnrich_EmptyGroups tests batch query with no groups.
func TestEnrich_EmptyGroups(t *testing.T) {
	repo := &mockRepo{
		hashGroups: []repository.HashGroup{},
		filesMap:   repository.FilesByHashMap{},
	}

	det := NewDuplicateDetector(repo)
	groups, err := det.FindAcrossRemotes(context.Background(), FilterOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(groups) != 0 {
		t.Errorf("expected 0 groups, got %d", len(groups))
	}
}

// TestGetLargestDuplicateGroups_ReturnsTopN tests that largest groups are returned in descending order.
func TestGetLargestDuplicateGroups_ReturnsTopN(t *testing.T) {
	groups := []DuplicateGroup{
		{HashValue: "hash1", Count: 5, Files: []FileInfo{{Name: "a"}}},
		{HashValue: "hash2", Count: 20, Files: []FileInfo{{Name: "b"}}},
		{HashValue: "hash3", Count: 10, Files: []FileInfo{{Name: "c"}}},
		{HashValue: "hash4", Count: 8, Files: []FileInfo{{Name: "d"}}},
		{HashValue: "hash5", Count: 15, Files: []FileInfo{{Name: "e"}}},
	}

	largest := GetLargestDuplicateGroups(groups, 3)

	if len(largest) != 3 {
		t.Fatalf("expected 3 groups, got %d", len(largest))
	}
	if largest[0].Count != 20 {
		t.Errorf("expected first group count 20, got %d", largest[0].Count)
	}
	if largest[1].Count != 15 {
		t.Errorf("expected second group count 15, got %d", largest[1].Count)
	}
	if largest[2].Count != 10 {
		t.Errorf("expected third group count 10, got %d", largest[2].Count)
	}
}

// TestGetLargestDuplicateGroups_RequestMoreThanAvailable returns all groups.
func TestGetLargestDuplicateGroups_RequestMoreThanAvailable(t *testing.T) {
	groups := []DuplicateGroup{
		{HashValue: "hash1", Count: 5},
		{HashValue: "hash2", Count: 10},
	}

	largest := GetLargestDuplicateGroups(groups, 100)

	if len(largest) != 2 {
		t.Fatalf("expected 2 groups (all available), got %d", len(largest))
	}
}

// TestGetLargestDuplicateGroups_ZeroOrNegative returns empty slice.
func TestGetLargestDuplicateGroups_ZeroOrNegative(t *testing.T) {
	groups := []DuplicateGroup{
		{HashValue: "hash1", Count: 5},
		{HashValue: "hash2", Count: 10},
	}

	if largest := GetLargestDuplicateGroups(groups, 0); len(largest) != 0 {
		t.Errorf("expected 0 groups for n=0, got %d", len(largest))
	}

	if largest := GetLargestDuplicateGroups(groups, -5); len(largest) != 0 {
		t.Errorf("expected 0 groups for n=-5, got %d", len(largest))
	}
}

// TestGetLargestDuplicateGroups_EmptyInput returns empty slice.
func TestGetLargestDuplicateGroups_EmptyInput(t *testing.T) {
	var groups []DuplicateGroup
	largest := GetLargestDuplicateGroups(groups, 10)
	if len(largest) != 0 {
		t.Errorf("expected 0 groups for empty input, got %d", len(largest))
	}
}

// TestGetLargestDuplicateGroups_PreservesOriginal verifies that the original slice is not modified.
func TestGetLargestDuplicateGroups_PreservesOriginal(t *testing.T) {
	groups := []DuplicateGroup{
		{HashValue: "hash1", Count: 5},
		{HashValue: "hash2", Count: 20},
		{HashValue: "hash3", Count: 10},
	}
	original := make([]DuplicateGroup, len(groups))
	copy(original, groups)

	GetLargestDuplicateGroups(groups, 2)

	// Verify original was not modified
	if groups[0].Count != original[0].Count || groups[1].Count != original[1].Count {
		t.Error("original slice was modified by GetLargestDuplicateGroups")
	}
}

