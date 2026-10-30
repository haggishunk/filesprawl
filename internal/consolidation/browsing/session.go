package browsing

import (
	"context"
	"fmt"
	"github.com/haggishunk/filesprawl/internal/consolidation/flags"
)

// SessionManager manages browsing sessions
type SessionManager struct {
	userID        int
	currentView   ViewType
	duplicates    []DuplicateGroup
	flags         map[int]*flags.Flag
	flagRepo      flags.Repository
	expandedItems map[string]bool
	cursor        int
}

// NewSessionManager creates a new browsing session
func NewSessionManager(userID int, flagRepo flags.Repository) *SessionManager {
	return &SessionManager{
		userID:        userID,
		currentView:   ViewTypeHash,
		duplicates:    []DuplicateGroup{},
		flags:         make(map[int]*flags.Flag),
		flagRepo:      flagRepo,
		expandedItems: make(map[string]bool),
		cursor:        0,
	}
}

// LoadDuplicates loads a list of duplicate groups for browsing
func (sm *SessionManager) LoadDuplicates(duplicates []DuplicateGroup) {
	sm.duplicates = duplicates
}

// LoadExistingFlags loads previously flagged decisions from database
func (sm *SessionManager) LoadExistingFlags(ctx context.Context) error {
	if sm.flagRepo == nil {
		return nil
	}

	userFlags, err := sm.flagRepo.GetFlagsByUser(ctx, sm.userID)
	if err != nil {
		return fmt.Errorf("failed to load existing flags: %w", err)
	}

	for _, f := range userFlags {
		flag := f // capture copy
		sm.flags[f.FileID] = &flag
	}

	return nil
}

// FlagFile sets a flag on a file and persists it
func (sm *SessionManager) FlagFile(ctx context.Context, fileID int, flagType flags.FlagType, reason *string, archiveDest *string) error {
	if sm.flagRepo == nil {
		return fmt.Errorf("flag repository not initialized")
	}

	// Check if flag already exists
	if existingFlag, ok := sm.flags[fileID]; ok {
		// Update existing flag
		updated, err := sm.flagRepo.UpdateFlag(ctx, existingFlag.ID, flagType, reason, archiveDest)
		if err != nil {
			return fmt.Errorf("failed to update flag: %w", err)
		}
		sm.flags[fileID] = updated
	} else {
		// Create new flag
		newFlag, err := sm.flagRepo.CreateFlag(ctx, fileID, sm.userID, flagType, reason, archiveDest)
		if err != nil {
			return fmt.Errorf("failed to create flag: %w", err)
		}
		sm.flags[fileID] = newFlag
	}

	return nil
}

// UnflagFile removes a flag from a file
func (sm *SessionManager) UnflagFile(ctx context.Context, fileID int) error {
	if sm.flagRepo == nil {
		return fmt.Errorf("flag repository not initialized")
	}

	if flag, ok := sm.flags[fileID]; ok {
		if err := sm.flagRepo.DeleteFlag(ctx, flag.ID); err != nil {
			return fmt.Errorf("failed to delete flag: %w", err)
		}
		delete(sm.flags, fileID)
	}

	return nil
}

// SwitchView toggles between hash and filename views
func (sm *SessionManager) SwitchView() {
	if sm.currentView == ViewTypeHash {
		sm.currentView = ViewTypeFilename
	} else {
		sm.currentView = ViewTypeHash
	}
	sm.cursor = 0 // Reset cursor when switching views
}

// GetCurrentView returns the current view type
func (sm *SessionManager) GetCurrentView() ViewType {
	return sm.currentView
}

// ExpandItem toggles expansion of an item
func (sm *SessionManager) ExpandItem(identifier string) {
	sm.expandedItems[identifier] = !sm.expandedItems[identifier]
}

// IsExpanded checks if an item is expanded
func (sm *SessionManager) IsExpanded(identifier string) bool {
	return sm.expandedItems[identifier]
}

// GetFlags returns all flags in this session
func (sm *SessionManager) GetFlags() map[int]*flags.Flag {
	return sm.flags
}

// CountByType returns flag counts grouped by type
func (sm *SessionManager) CountByType() map[flags.FlagType]int {
	counts := make(map[flags.FlagType]int)
	for _, flag := range sm.flags {
		counts[flag.FlagType]++
	}
	return counts
}

// GetSession returns the current session state
func (sm *SessionManager) GetSession() *BrowsingSession {
	return &BrowsingSession{
		UserID:          sm.userID,
		CurrentView:     sm.currentView,
		DuplicateGroups: sm.duplicates,
		ExpandedItems:   sm.expandedItems,
		Flags:           sm.flags,
		Cursor:          sm.cursor,
	}
}

// MoveCursorUp moves cursor up in the current view
func (sm *SessionManager) MoveCursorUp() {
	if sm.cursor > 0 {
		sm.cursor--
	}
}

// MoveCursorDown moves cursor down in the current view
func (sm *SessionManager) MoveCursorDown() {
	maxCursor := len(sm.duplicates) - 1
	if maxCursor >= 0 && sm.cursor < maxCursor {
		sm.cursor++
	}
}

// GetCursor returns current cursor position
func (sm *SessionManager) GetCursor() int {
	return sm.cursor
}

// SetCursor sets cursor position
func (sm *SessionManager) SetCursor(pos int) {
	if pos >= 0 && pos < len(sm.duplicates) {
		sm.cursor = pos
	}
}
