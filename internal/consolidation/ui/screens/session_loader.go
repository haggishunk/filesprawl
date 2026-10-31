package screens

import (
	"context"
	"fmt"
	"strings"

	"github.com/haggishunk/filesprawl/internal/consolidation/flags"
)

// SessionLoaderScreen handles loading previous sessions with pending flags
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

// Show displays the session loader screen with pending flags decision
func (sls *SessionLoaderScreen) Show(ctx context.Context) (bool, error) {
	// Get pending flags count
	pendingCount, err := sls.getPendingFlagsCount(ctx)
	if err != nil {
		// On error, start fresh
		return true, nil
	}

	sls.flagCount = pendingCount
	sls.hasPending = pendingCount > 0

	if !sls.hasPending {
		// No pending flags, start fresh
		return true, nil
	}

	// Show decision modal
	return sls.showDecisionModal(pendingCount)
}

// getPendingFlagsCount retrieves the count of pending flags for the user
func (sls *SessionLoaderScreen) getPendingFlagsCount(ctx context.Context) (int, error) {
	if sls.flagRepo == nil {
		return 0, fmt.Errorf("flag repository not initialized")
	}

	userFlags, err := sls.flagRepo.GetFlagsByUser(ctx, sls.userID)
	if err != nil {
		return 0, fmt.Errorf("failed to get user flags: %w", err)
	}

	return len(userFlags), nil
}

// showDecisionModal displays the modal asking whether to resume or start fresh
func (sls *SessionLoaderScreen) showDecisionModal(count int) (bool, error) {
	// TODO: Phase 4 - Implement full gotui modal for session decision
	// For now, use console-based input as placeholder

	fmt.Printf("You have %d pending flags from a previous session.\n", count)
	fmt.Print("Resume session? [y/n]: ")

	var response string
	_, err := fmt.Scanln(&response)
	if err != nil {
		// Default to resume on error
		return true, nil
	}

	response = strings.ToLower(strings.TrimSpace(response))
	if response == "y" || response == "yes" {
		return true, nil
	}

	return false, nil
}

// HasPending returns true if there are pending flags
func (sls *SessionLoaderScreen) HasPending() bool {
	return sls.hasPending
}

// GetPendingCount returns the number of pending flags
func (sls *SessionLoaderScreen) GetPendingCount() int {
	return sls.flagCount
}
