package components

import (
	"fmt"

	"github.com/haggishunk/filesprawl/internal/consolidation/browsing"
	"github.com/haggishunk/filesprawl/internal/consolidation/flags"
)

// SafetyChecker validates deletion safety constraints
type SafetyChecker struct {
	flagManager *browsing.PersistentFlagManager
}

// NewSafetyChecker creates a new safety checker
func NewSafetyChecker(flagManager *browsing.PersistentFlagManager) *SafetyChecker {
	return &SafetyChecker{
		flagManager: flagManager,
	}
}

// CanDeleteFile checks if a file can be safely flagged for deletion
// Returns (allowed, reason) where reason explains why deletion is not allowed
func (sc *SafetyChecker) CanDeleteFile(fileToDelete *browsing.FileVersion, group *browsing.DuplicateGroup) (bool, string) {
	if fileToDelete == nil || fileToDelete.FileInfo == nil || group == nil {
		return false, "Invalid file or group"
	}

	// Get the (remote, filename) pair for this file
	remote := fileToDelete.FileInfo.ScanRoot
	filename := fileToDelete.FileInfo.Name

	// Count non-deleted copies of this (remote, filename) pair
	var nonDeletedCount int
	var flaggedDeleteCount int

	for _, version := range group.Versions {
		if version.FileInfo == nil {
			continue
		}

		// Only count copies with same (remote, filename) pair
		if version.FileInfo.ScanRoot != remote || version.FileInfo.Name != filename {
			continue
		}

		// Check if this copy is flagged for deletion
		isFlaggedDelete := version.Flag != nil && version.Flag.FlagType == flags.FlagTypeDelete

		// Count the target file being deleted
		if version.FileID == fileToDelete.FileID {
			if isFlaggedDelete {
				flaggedDeleteCount++
			}
			continue
		}

		// Count other copies
		if !isFlaggedDelete {
			nonDeletedCount++
		} else {
			flaggedDeleteCount++
		}
	}

	// Safety rule: Must keep at least one copy of (remote, filename) pair
	if nonDeletedCount == 0 && flaggedDeleteCount > 0 {
		return false, fmt.Sprintf(
			"Cannot delete: Would remove all %d copies of %s on %s. Must keep at least one.",
			flaggedDeleteCount+1, filename, remote,
		)
	}

	return true, ""
}

// CheckDeletionSafety performs comprehensive safety checks before deletion
// Returns true if safe to delete, false with explanation otherwise
func (sc *SafetyChecker) CheckDeletionSafety(fileToDelete *browsing.FileVersion, group *browsing.DuplicateGroup) (bool, string) {
	if fileToDelete == nil || group == nil {
		return false, "Invalid file or group data"
	}

	// Check (remote, filename) pair safety
	allowed, reason := sc.CanDeleteFile(fileToDelete, group)
	if !allowed {
		return false, reason
	}

	return true, ""
}

// GetCopyCount returns the number of copies for a (remote, filename) pair
func (sc *SafetyChecker) GetCopyCount(remote, filename string, group *browsing.DuplicateGroup) int {
	if group == nil {
		return 0
	}

	var count int
	for _, version := range group.Versions {
		if version.FileInfo == nil {
			continue
		}
		if version.FileInfo.ScanRoot == remote && version.FileInfo.Name == filename {
			count++
		}
	}
	return count
}

// GetNonDeletedCopyCount returns count of non-deleted copies for a (remote, filename) pair
func (sc *SafetyChecker) GetNonDeletedCopyCount(remote, filename string, group *browsing.DuplicateGroup) int {
	if group == nil {
		return 0
	}

	var count int
	for _, version := range group.Versions {
		if version.FileInfo == nil {
			continue
		}
		if version.FileInfo.ScanRoot != remote || version.FileInfo.Name != filename {
			continue
		}

		// Count if NOT flagged for deletion
		if version.Flag == nil || version.Flag.FlagType != flags.FlagTypeDelete {
			count++
		}
	}
	return count
}
