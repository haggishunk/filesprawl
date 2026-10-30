package screens

import (
	"context"

	"github.com/haggishunk/filesprawl/internal/consolidation/flags"
)

// SessionLoaderScreen handles loading previous sessions with pending flags
// NOTE: This is a stub implementation. Full gotui integration is in progress.
type SessionLoaderScreen struct {
	flagRepo   flags.Repository
	userID     int
	hasPending bool
	flagCount  int
}

// NewSessionLoaderScreen creates a new session loader screen
func NewSessionLoaderScreen(flagRepo flags.Repository, userID int) *SessionLoaderScreen {
	return &SessionLoaderScreen{
		flagRepo: flagRepo,
		userID:   userID,
	}
}

// Show displays the session loader screen (stub implementation)
func (sls *SessionLoaderScreen) Show(ctx context.Context) (bool, error) {
	// Check for pending flags
	userFlags, err := sls.flagRepo.GetFlagsByUser(ctx, sls.userID)
	if err != nil {
		return true, nil // Continue with new session on error
	}

	hasPending := len(userFlags) > 0
	sls.hasPending = hasPending
	sls.flagCount = len(userFlags)

	if !hasPending {
		// No pending flags, start fresh
		return true, nil
	}

	// TODO: Implement full gotui-based session loader screen
	// For now, automatically resume session with pending flags
	return true, nil
}

// HasPending returns true if there are pending flags
func (sls *SessionLoaderScreen) HasPending() bool {
	return sls.hasPending
}

// GetPendingCount returns the number of pending flags
func (sls *SessionLoaderScreen) GetPendingCount() int {
	return sls.flagCount
}
