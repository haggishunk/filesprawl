package operations

import (
	"context"
	"testing"
)

func TestNewRcloneExecutor(t *testing.T) {
	executor := NewRcloneExecutor()
	if executor == nil {
		t.Errorf("expected non-nil executor")
	}
	if executor.IsDryRun() {
		t.Errorf("expected dry-run to be false initially")
	}
}

func TestSetDryRun(t *testing.T) {
	executor := NewRcloneExecutor()

	executor.SetDryRun(true)
	if !executor.IsDryRun() {
		t.Errorf("expected dry-run to be true after setting")
	}

	executor.SetDryRun(false)
	if executor.IsDryRun() {
		t.Errorf("expected dry-run to be false after unsetting")
	}
}

func TestArchiveDryRun(t *testing.T) {
	executor := NewRcloneExecutor()
	executor.SetDryRun(true)

	file := ConsolidationFile{
		ID:     1,
		Path:   "documents/report.pdf",
		Remote: "media",
		Hash:   "abc123",
	}

	result := executor.Archive(context.Background(), file, "archive:/consolidated/2024")

	if !result.Success {
		t.Errorf("expected dry-run archive to succeed")
	}
	if result.Error != "" {
		t.Errorf("expected no error in dry-run, got: %s", result.Error)
	}
	if result.Operation != OperationArchive {
		t.Errorf("expected operation type ARCHIVE, got %v", result.Operation)
	}
}

func TestDeleteDryRun(t *testing.T) {
	executor := NewRcloneExecutor()
	executor.SetDryRun(true)

	file := ConsolidationFile{
		ID:     2,
		Path:   "documents/old.txt",
		Remote: "local",
		Hash:   "def456",
	}

	result := executor.Delete(context.Background(), file)

	if !result.Success {
		t.Errorf("expected dry-run delete to succeed")
	}
	if result.Error != "" {
		t.Errorf("expected no error in dry-run, got: %s", result.Error)
	}
	if result.Operation != OperationDelete {
		t.Errorf("expected operation type DELETE, got %v", result.Operation)
	}
}

func TestMoveDryRun(t *testing.T) {
	executor := NewRcloneExecutor()
	executor.SetDryRun(true)

	file := ConsolidationFile{
		ID:     3,
		Path:   "documents/temp.txt",
		Remote: "local",
		Hash:   "ghi789",
	}

	result := executor.Move(context.Background(), file, "backup:/archive/temp.txt")

	if !result.Success {
		t.Errorf("expected dry-run move to succeed")
	}
	if result.Error != "" {
		t.Errorf("expected no error in dry-run, got: %s", result.Error)
	}
	if result.Operation != OperationMove {
		t.Errorf("expected operation type MOVE, got %v", result.Operation)
	}
}

func TestArchiveInvalidDestination(t *testing.T) {
	executor := NewRcloneExecutor()

	file := ConsolidationFile{
		ID:     1,
		Path:   "documents/report.pdf",
		Remote: "media",
	}

	// Invalid destination format (no colon)
	result := executor.Archive(context.Background(), file, "invalid_destination")

	if result.Success {
		t.Errorf("expected archive with invalid destination to fail")
	}
	if result.Error == "" {
		t.Errorf("expected error message for invalid destination")
	}
}

func TestExtractFileName(t *testing.T) {
	tests := []struct {
		path     string
		expected string
	}{
		{"documents/report.pdf", "report.pdf"},
		{"folder/subfolder/file.txt", "file.txt"},
		{"file.txt", "file.txt"},
		{"path/with/trailing/slash/", "slash"},
		{"/absolute/path/to/file.md", "file.md"},
	}

	for _, tc := range tests {
		result := extractFileName(tc.path)
		if result != tc.expected {
			t.Errorf("extractFileName(%s): expected %s, got %s", tc.path, tc.expected, result)
		}
	}
}

func TestOperationResult(t *testing.T) {
	result := &OperationResult{
		FileID:    42,
		Operation: OperationArchive,
		Success:   true,
		Details:   "Archived file to destination",
	}

	if result.FileID != 42 {
		t.Errorf("expected FileID 42, got %d", result.FileID)
	}
	if result.Operation != OperationArchive {
		t.Errorf("expected operation ARCHIVE, got %v", result.Operation)
	}
	if !result.Success {
		t.Errorf("expected Success to be true")
	}
}
