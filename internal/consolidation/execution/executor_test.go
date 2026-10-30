package execution

import (
	"context"
	"testing"
	"github.com/haggishunk/filesprawl/internal/consolidation/flags"
)

func TestNewExecutor(t *testing.T) {
	executor := NewExecutor(nil, 1)

	if executor.userID != 1 {
		t.Errorf("expected userID 1, got %d", executor.userID)
	}

	if executor.dryRun {
		t.Errorf("expected dry-run to be false initially")
	}
}

func TestSetDryRun(t *testing.T) {
	executor := NewExecutor(nil, 1)

	executor.SetDryRun(true)
	if !executor.dryRun {
		t.Errorf("expected dry-run to be true after setting")
	}

	executor.SetDryRun(false)
	if executor.dryRun {
		t.Errorf("expected dry-run to be false after unsetting")
	}
}

func TestExecutionResultCreation(t *testing.T) {
	result := &ExecutionResult{
		KeepCount:    10,
		ArchiveCount: 5,
		DeleteCount:  3,
		SuccessCount: 18,
		FailureCount: 0,
		Errors:       []string{},
	}

	if result.KeepCount != 10 {
		t.Errorf("expected KeepCount 10, got %d", result.KeepCount)
	}

	if result.ArchiveCount != 5 {
		t.Errorf("expected ArchiveCount 5, got %d", result.ArchiveCount)
	}

	if result.DeleteCount != 3 {
		t.Errorf("expected DeleteCount 3, got %d", result.DeleteCount)
	}

	if result.SuccessCount != 18 {
		t.Errorf("expected SuccessCount 18, got %d", result.SuccessCount)
	}

	if result.FailureCount != 0 {
		t.Errorf("expected FailureCount 0, got %d", result.FailureCount)
	}

	if len(result.Errors) != 0 {
		t.Errorf("expected no errors, got %d", len(result.Errors))
	}
}

func TestGetExecutionSummary(t *testing.T) {
	executor := NewExecutor(nil, 1)

	result := &ExecutionResult{
		KeepCount:    10,
		ArchiveCount: 5,
		DeleteCount:  3,
		SuccessCount: 18,
		FailureCount: 0,
		Errors:       []string{},
	}

	summary := executor.GetExecutionSummary(result)

	if summary == "" {
		t.Errorf("expected non-empty summary")
	}

	// Check that summary contains key numbers
	if !contains(summary, "10") {
		t.Errorf("expected keep count in summary")
	}

	if !contains(summary, "5") {
		t.Errorf("expected archive count in summary")
	}

	if !contains(summary, "3") {
		t.Errorf("expected delete count in summary")
	}
}

func TestGetExecutionSummaryWithErrors(t *testing.T) {
	executor := NewExecutor(nil, 1)

	result := &ExecutionResult{
		KeepCount:    10,
		ArchiveCount: 5,
		DeleteCount:  3,
		SuccessCount: 17,
		FailureCount: 1,
		Errors:       []string{"File 100: access denied"},
	}

	summary := executor.GetExecutionSummary(result)

	if !contains(summary, "access denied") {
		t.Errorf("expected error in summary")
	}
}

// Helper function
func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
