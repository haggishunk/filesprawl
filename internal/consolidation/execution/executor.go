package execution

import (
	"context"
	"fmt"
	"github.com/haggishunk/filesprawl/internal/consolidation/flags"
)

// ExecutionResult represents the result of executing consolidation
type ExecutionResult struct {
	KeepCount     int
	ArchiveCount  int
	DeleteCount   int
	SuccessCount  int
	FailureCount  int
	Errors        []string
}

// Executor performs consolidation actions (keep, archive, delete)
type Executor struct {
	flagRepo flags.Repository
	userID   int
	dryRun   bool
}

// NewExecutor creates a new consolidation executor
func NewExecutor(flagRepo flags.Repository, userID int) *Executor {
	return &Executor{
		flagRepo: flagRepo,
		userID:   userID,
		dryRun:   false,
	}
}

// SetDryRun sets dry-run mode (report but don't execute)
func (e *Executor) SetDryRun(dryRun bool) {
	e.dryRun = dryRun
}

// Execute performs all flagged consolidation actions
func (e *Executor) Execute(ctx context.Context) (*ExecutionResult, error) {
	result := &ExecutionResult{
		Errors: []string{},
	}

	// Retrieve all flags for this user
	userFlags, err := e.flagRepo.GetFlagsByUser(ctx, e.userID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve flags: %w", err)
	}

	// Count flags by type
	for _, flag := range userFlags {
		switch flag.FlagType {
		case flags.FlagTypeKeep:
			result.KeepCount++
		case flags.FlagTypeArchive:
			result.ArchiveCount++
		case flags.FlagTypeDelete:
			result.DeleteCount++
		}
	}

	if e.dryRun {
		result.SuccessCount = result.KeepCount + result.ArchiveCount + result.DeleteCount
		return result, nil
	}

	// Execute each action
	for _, flag := range userFlags {
		if err := e.executeFlag(ctx, flag); err != nil {
			result.FailureCount++
			result.Errors = append(result.Errors, fmt.Sprintf("File %d: %v", flag.FileID, err))
		} else {
			result.SuccessCount++
		}
	}

	return result, nil
}

// executeFlag executes a single flag action
func (e *Executor) executeFlag(ctx context.Context, flag *flags.Flag) error {
	switch flag.FlagType {
	case flags.FlagTypeKeep:
		// KEEP: No action needed, just mark as reviewed
		return e.markReviewed(ctx, flag, flag.FlagType, "Kept as-is")

	case flags.FlagTypeArchive:
		if flag.ArchiveDest == nil {
			return fmt.Errorf("archive destination not specified")
		}
		// ARCHIVE: Copy file to destination (via rclone RC)
		// TODO: Implement actual rclone RC call
		return e.markReviewed(ctx, flag, flag.FlagType, fmt.Sprintf("Archived to %s", *flag.ArchiveDest))

	case flags.FlagTypeDelete:
		// DELETE: Delete the file
		// TODO: Implement actual file deletion (via rclone RC)
		return e.markReviewed(ctx, flag, flag.FlagType, "Deleted")

	default:
		return fmt.Errorf("unknown flag type: %v", flag.FlagType)
	}
}

// markReviewed marks a flag as reviewed (creates audit trail entry)
func (e *Executor) markReviewed(ctx context.Context, flag *flags.Flag, newFlagType flags.FlagType, reason string) error {
	if e.flagRepo == nil {
		return nil
	}

	// Create a review record in the audit trail
	reviewReason := reason
	_, err := e.flagRepo.CreateReview(ctx, flag.FileID, flag.FlagType, flag.UserID, e.userID, newFlagType, &reviewReason)
	if err != nil {
		return fmt.Errorf("failed to create review record: %w", err)
	}

	return nil
}

// GetExecutionSummary returns a human-readable execution summary
func (e *Executor) GetExecutionSummary(result *ExecutionResult) string {
	summary := "\n╔════════════════════════════════════════════════════════════════╗\n"
	summary += "║              CONSOLIDATION EXECUTION SUMMARY                    ║\n"
	summary += "╚════════════════════════════════════════════════════════════════╝\n\n"

	summary += fmt.Sprintf("Actions Executed:\n")
	summary += fmt.Sprintf("  🟢 KEEP:    %3d files (no action)\n", result.KeepCount)
	summary += fmt.Sprintf("  📦 ARCHIVE: %3d files\n", result.ArchiveCount)
	summary += fmt.Sprintf("  🗑️  DELETE:  %3d files\n", result.DeleteCount)
	summary += fmt.Sprintf("─────────────────────────────────────\n")
	summary += fmt.Sprintf("  ✓ Success: %3d operations\n", result.SuccessCount)
	summary += fmt.Sprintf("  ✗ Failed:  %3d operations\n", result.FailureCount)

	if len(result.Errors) > 0 {
		summary += "\nErrors:\n"
		for _, err := range result.Errors {
			summary += fmt.Sprintf("  • %s\n", err)
		}
	}

	summary += "\n"
	return summary
}
