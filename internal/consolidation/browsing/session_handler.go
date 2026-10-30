package browsing

import (
	"context"
	"fmt"
	"time"
)

// SessionHandler manages session lifecycle including save and exit
type SessionHandler struct {
	session              *SessionManager
	flagManager          *PersistentFlagManager
	summaryRenderer      *SummaryRenderer
	lastSaveTime         time.Time
	unsavedChanges       bool
	autoSaveInterval     time.Duration
}

// NewSessionHandler creates a new session handler
func NewSessionHandler(session *SessionManager, flagManager *PersistentFlagManager) *SessionHandler {
	return &SessionHandler{
		session:          session,
		flagManager:      flagManager,
		summaryRenderer:  NewSummaryRenderer(),
		lastSaveTime:     time.Now(),
		unsavedChanges:   false,
		autoSaveInterval: 30 * time.Second,
	}
}

// QuitWithoutConsolidation exits the session without executing consolidation
// All flags are already persisted in database, so we just clean up
func (sh *SessionHandler) QuitWithoutConsolidation(ctx context.Context) (string, error) {
	var output string

	// Show quit confirmation
	output += "\n╔════════════════════════════════════════════════════════════════╗\n"
	output += "║                    EXITING CONSOLIDATION TOOL                   ║\n"
	output += "╚════════════════════════════════════════════════════════════════╝\n\n"

	// Get summary of flagged items
	sessionData := sh.session.GetSession()
	flaggedCount := len(sessionData.Flags)

	if flaggedCount > 0 {
		output += fmt.Sprintf("✅ All %d flagged decisions have been saved to the database\n\n", flaggedCount)

		// Show what was saved
		keep, archive, delete := sh.summaryRenderer.countDecisions(sessionData.Flags)
		output += fmt.Sprintf("Saved flags:\n")
		output += fmt.Sprintf("  🟢 KEEP:    %d files\n", keep)
		output += fmt.Sprintf("  📦 ARCHIVE: %d files\n", archive)
		output += fmt.Sprintf("  🗑️  DELETE:  %d files\n\n", delete)

		output += "Resume flagging later by running:\n"
		output += "  filesprawl consolidate --review\n\n"
	} else {
		output += "No flags were set. Exiting without changes.\n\n"
	}

	output += "Thank you for using the consolidation tool!\n"

	return output, nil
}

// SaveProgress saves all current flags to database
func (sh *SessionHandler) SaveProgress(ctx context.Context) error {
	if sh.flagManager == nil {
		return nil
	}

	// Flags are persisted on each keystroke via flagManager
	// This is a checkpoint save that could flush other metadata
	sh.lastSaveTime = time.Now()
	sh.unsavedChanges = false

	return nil
}

// HasUnsavedChanges returns true if there are unsaved changes
func (sh *SessionHandler) HasUnsavedChanges() bool {
	return sh.unsavedChanges
}

// MarkChanged marks session as having unsaved changes
func (sh *SessionHandler) MarkChanged() {
	sh.unsavedChanges = true
}

// ShouldAutoSave returns true if auto-save interval has passed
func (sh *SessionHandler) ShouldAutoSave() bool {
	return time.Since(sh.lastSaveTime) > sh.autoSaveInterval
}

// GetSessionSummary returns a summary of the current session
func (sh *SessionHandler) GetSessionSummary() string {
	sessionData := sh.session.GetSession()

	var output string
	output += "\n═══════════════════════════════════════════════════════════════\n"
	output += "               SESSION SUMMARY\n"
	output += "═══════════════════════════════════════════════════════════════\n\n"

	// User and view info
	output += fmt.Sprintf("User ID:        %d\n", sessionData.UserID)
	output += fmt.Sprintf("Current View:   %v\n", sessionData.CurrentView)
	output += fmt.Sprintf("Total Groups:   %d\n", len(sessionData.DuplicateGroups))

	// Flags summary
	keep, archive, delete := sh.summaryRenderer.countDecisions(sessionData.Flags)
	output += fmt.Sprintf("\nFlagged Decisions:\n")
	output += fmt.Sprintf("  KEEP:    %3d\n", keep)
	output += fmt.Sprintf("  ARCHIVE: %3d\n", archive)
	output += fmt.Sprintf("  DELETE:  %3d\n", delete)
	output += fmt.Sprintf("  ─────────────\n")
	output += fmt.Sprintf("  TOTAL:   %3d\n", keep+archive+delete)

	// Expanded items
	output += fmt.Sprintf("\nExpanded Items: %d\n", len(sessionData.ExpandedItems))

	// Time info
	output += fmt.Sprintf("\nLast Save: %s\n", sh.lastSaveTime.Format("2006-01-02 15:04:05"))

	output += "\n═══════════════════════════════════════════════════════════════\n"

	return output
}

// RenderQuitPrompt shows the quit confirmation prompt
func (sh *SessionHandler) RenderQuitPrompt() string {
	sessionData := sh.session.GetSession()
	flaggedCount := len(sessionData.Flags)

	var output string
	output += "\n┌─────────────────────────────────────────────────────────────┐\n"
	output += "│                   QUIT CONSOLIDATION?                       │\n"
	output += "└─────────────────────────────────────────────────────────────┘\n\n"

	if flaggedCount > 0 {
		output += fmt.Sprintf("ℹ️  You have %d flagged decisions\n", flaggedCount)
		output += "✅ All decisions have been saved to the database\n"
		output += "📝 You can resume and continue consolidating later\n\n"
		output += "Press:\n"
		output += "  Q - Confirm quit (flags saved)\n"
		output += "  C - Continue consolidating\n"
		output += "  V - View summary before quitting\n"
	} else {
		output += "No flags have been set\n"
		output += "Press:\n"
		output += "  Q - Confirm quit\n"
		output += "  C - Continue consolidating\n"
	}

	return output
}

// GetResumeInstructions returns instructions for resuming the session
func (sh *SessionHandler) GetResumeInstructions() string {
	sessionData := sh.session.GetSession()
	flaggedCount := len(sessionData.Flags)

	var output string
	output += "\n╔════════════════════════════════════════════════════════════════╗\n"
	output += "║                      HOW TO RESUME                              ║\n"
	output += "╚════════════════════════════════════════════════════════════════╝\n\n"

	output += "Your consolidation session has been saved with:\n"
	output += fmt.Sprintf("  • %d flagged files\n", flaggedCount)
	output += "  • All decisions preserved in the database\n"
	output += "  • Multi-user review capability enabled\n\n"

	output += "To resume consolidating later:\n"
	output += "  $ filesprawl consolidate --review [--user <username>]\n\n"

	output += "To review decisions made by others:\n"
	output += "  $ filesprawl consolidate --review-all\n\n"

	output += "To execute all approved consolidations:\n"
	output += "  $ filesprawl consolidate --execute\n\n"

	return output
}
