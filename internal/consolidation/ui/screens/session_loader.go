package screens

import (
	"context"
	"fmt"
	"github.com/metaspartan/gotui/v5"
	"github.com/metaspartan/gotui/v5/widgets"
	"github.com/haggishunk/filesprawl/internal/consolidation/flags"
)

// SessionLoaderScreen handles loading previous sessions with pending flags
type SessionLoaderScreen struct {
	app       *gotui.Application
	flagRepo  flags.Repository
	userID    int
	hasPending bool
	flagCount int
}

// NewSessionLoaderScreen creates a new session loader screen
func NewSessionLoaderScreen(app *gotui.Application, flagRepo flags.Repository, userID int) *SessionLoaderScreen {
	return &SessionLoaderScreen{
		app:      app,
		flagRepo: flagRepo,
		userID:   userID,
	}
}

// Show displays the session loader screen
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

	// Display pending flags and ask to resume
	title := widgets.NewParagraph()
	title.Title = "📋 PENDING FLAGS FROM PREVIOUS SESSION"
	title.Text = fmt.Sprintf("You have %d flagged decisions from a previous session.\n\nChoose an action below:", sls.flagCount)
	title.SetRect(0, 0, gotui.MaxX(), 7)
	title.TitleAlignment = gotui.AlignCenter
	title.BorderStyle = gotui.NewStyle().Foreground(gotui.ColorCyan)

	// Count by type
	keepCount := 0
	archiveCount := 0
	deleteCount := 0
	for _, flag := range userFlags {
		switch flag.FlagType {
		case flags.FlagTypeKeep:
			keepCount++
		case flags.FlagTypeArchive:
			archiveCount++
		case flags.FlagTypeDelete:
			deleteCount++
		}
	}

	// Show summary
	summary := widgets.NewParagraph()
	summary.Title = "FLAGGED DECISIONS"
	summary.Text = fmt.Sprintf("🟢 KEEP:    %3d files\n📦 ARCHIVE: %3d files\n🗑️  DELETE:  %3d files\n", keepCount, archiveCount, deleteCount)
	summary.SetRect(5, 8, 60, 14)
	summary.BorderStyle = gotui.NewStyle().Foreground(gotui.ColorYellow)

	// Options
	options := widgets.NewList()
	options.Title = "RESUME OPTIONS"
	options.Rows = []string{
		"[R] Resume Session (load pending flags)",
		"[N] Start New Session (clear pending flags)",
		"[Q] Quit",
	}
	options.SelectedRow = 0
	options.SetRect(5, 15, 60, 20)

	// Info
	info := widgets.NewParagraph()
	info.Title = "INSTRUCTIONS"
	info.Text = "[R] Resume / [N] New / [Q] Quit\n[Enter] Select option"
	info.SetRect(5, 21, 60, 25)
	info.BorderStyle = gotui.NewStyle().Foreground(gotui.ColorGreen)

	// Event loop
	uiEvents := gotui.PollEvents()
	for {
		e := <-uiEvents
		switch e.ID {
		case "q", "<C-c>":
			return false, fmt.Errorf("user cancelled session")
		case "r", "R", "<Enter>":
			if options.SelectedRow == 0 {
				return true, nil // Resume session, load pending flags
			} else if options.SelectedRow == 1 {
				return false, nil // Start new session
			} else {
				return false, fmt.Errorf("user quit")
			}
		case "n", "N":
			return false, nil // Start new session
		case "ArrowUp":
			if options.SelectedRow > 0 {
				options.SelectedRow--
			}
		case "ArrowDown":
			if options.SelectedRow < len(options.Rows)-1 {
				options.SelectedRow++
			}
		}
	}
}

// HasPending returns true if there are pending flags
func (sls *SessionLoaderScreen) HasPending() bool {
	return sls.hasPending
}

// GetPendingCount returns the number of pending flags
func (sls *SessionLoaderScreen) GetPendingCount() int {
	return sls.flagCount
}
