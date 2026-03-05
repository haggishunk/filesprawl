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
	acrossErr   error
	withinErr   error
	filesErr    error
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

var twoFiles = []repository.FileRecord{
	{ID: 1, Name: "a.txt", Path: "/docs/a.txt", Size: 1024, RemoteName: "dbox:", Hostname: "host1"},
	{ID: 2, Name: "a.txt", Path: "/backup/a.txt", Size: 1024, RemoteName: "s3:", Hostname: "host2"},
}

func TestFindAcrossRemotes_ReturnsDuplicateGroups(t *testing.T) {
	repo := &mockRepo{
		hashGroups:  []repository.HashGroup{{HashValue: "abc123", HashType: "md5", Count: 2}},
		fileRecords: twoFiles,
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
	repo := &mockRepo{
		hashGroups:  []repository.HashGroup{{HashValue: "def456", HashType: "sha256", Count: 2}},
		fileRecords: twoFiles,
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

