package screens

import (
	"context"
	"fmt"

	"github.com/haggishunk/filesprawl/internal/consolidation/browsing"
	"github.com/haggishunk/filesprawl/internal/consolidation/ui/modals"
	"github.com/haggishunk/filesprawl/internal/consolidation/ui/widgets"
)

// BrowsingScreen implements the interactive browsing interface for consolidation
type BrowsingScreen struct {
	sessionManager *browsing.SessionManager
	flagManager    *browsing.PersistentFlagManager
	hashWidget     *widgets.HashCentricWidget
	filenameWidget *widgets.FilenameCentricWidget
	detailsWidget  *widgets.DetailsWidget
	currentWidget  widgets.BrowsingWidget // pointer to current active widget
	isHashView     bool                    // true for hash, false for filename
	ctx            context.Context
}

// NewBrowsingScreen creates a new browsing screen
func NewBrowsingScreen(
	sessionManager *browsing.SessionManager,
	flagManager *browsing.PersistentFlagManager,
	hashWidget *widgets.HashCentricWidget,
	filenameWidget *widgets.FilenameCentricWidget,
	detailsWidget *widgets.DetailsWidget,
) *BrowsingScreen {
	return &BrowsingScreen{
		sessionManager: sessionManager,
		flagManager:    flagManager,
		hashWidget:     hashWidget,
		filenameWidget: filenameWidget,
		detailsWidget:  detailsWidget,
		isHashView:     true, // Start with hash view
		currentWidget:  hashWidget,
	}
}

// Show starts the browsing UI event loop with real gotui rendering
func (bs *BrowsingScreen) Show(ctx context.Context, app interface{}) error {
	if bs == nil || ctx == nil {
		return fmt.Errorf("invalid browsing screen or context")
	}

	if app == nil {
		return fmt.Errorf("gotui app not provided")
	}

	bs.ctx = ctx

	// Set initial widget (hash view)
	bs.currentWidget = bs.hashWidget

	// For Phase 5 gotui integration, determine if it's a real gotui app
	// For now, use hardcoded terminal size or default
	termWidth, termHeight := 80, 24

	// TODO: Phase 5 - Get terminal size from actual gotui app:
	// termWidth, termHeight := app.GetTerminalSize()

	if termWidth < 80 || termHeight < 24 {
		return fmt.Errorf("terminal too small: %dx%d (minimum 80x24)", termWidth, termHeight)
	}

	// Create layout manager for 3-pane layout
	layout := widgets.NewLayoutManager(termWidth, termHeight)

	// Get widget rectangles
	listX, listY, listW, listH := layout.GetListPaneRect()
	detailsX, detailsY, detailsW, detailsH := layout.GetDetailsPaneRect()
	statusX, statusY, statusW, statusH := layout.GetStatusBarRect()

	// Configure list widget positions
	bs.hashWidget.SetRect(listX, listY, listW, listH)
	bs.filenameWidget.SetRect(listX, listY, listW, listH)

	// Configure details widget position
	bs.detailsWidget.SetRect(detailsX, detailsY, detailsW, detailsH)

	// Configure status widget
	statusWidget := widgets.NewStatusWidget()
	statusSession := bs.sessionManager.GetSession()
	if statusSession != nil {
		statusWidget.UpdateStatus(statusSession)
	}
	statusWidget.SetRect(statusX, statusY, statusW, statusH)

	// Display layout information (Phase 5: gotui will render this)
	fmt.Println("═══════════════════════════════════════════════════════════════")
	fmt.Println("BrowsingScreen: 3-pane TUI layout initialized")
	fmt.Println("═══════════════════════════════════════════════════════════════")
	fmt.Println(fmt.Sprintf("List pane (hash/filename):  (%d,%d) %dx%d", listX, listY, listW, listH))
	fmt.Println(fmt.Sprintf("Details pane (file info):   (%d,%d) %dx%d", detailsX, detailsY, detailsW, detailsH))
	fmt.Println(fmt.Sprintf("Status bar (controls):      (%d,%d) %dx%d", statusX, statusY, statusW, statusH))
	fmt.Println("═══════════════════════════════════════════════════════════════")

	// TODO: Phase 5 - Wire to real gotui rendering:
	// This is where we would:
	// 1. Create root container with Grid layout
	// 2. Add list, details, status widgets to panes
	// 3. Set keyboard handler: app.SetKeyHandler(bs.onKeyPress)
	// 4. Run event loop: return app.Run()

	return nil
}

// handleKeypress handles keyboard input for browsing
func (bs *BrowsingScreen) handleKeypress(key rune, app interface{}) {
	if bs == nil || bs.ctx == nil {
		return
	}

	switch key {
	case 'j', 'J': // Scroll down
		bs.handleScrollDown()
	case 'k', 'K': // Scroll up or keep
		if key == 'K' { // Capital K is keep
			_ = bs.handleKeep()
		} else { // lowercase k is scroll up
			bs.handleScrollUp()
		}
	case 'v', 'V': // View switch
		_ = bs.handleViewSwitch()
	case 'a', 'A': // Archive
		_ = bs.handleArchive()
	case 'd', 'D': // Delete
		_ = bs.handleDelete()
	case 's', 'S': // Skip
		_ = bs.handleSkip()
	case 'q', 'Q': // Quit
		_ = bs.handleQuit(app)
	}
}

// handleScrollDown scrolls the list down
func (bs *BrowsingScreen) handleScrollDown() {
	// TODO: Phase 4 - Implement scrolling via gotui widgets
	if bs.currentWidget != nil {
		_ = bs.currentWidget.HandleKeypress('j')
	}
}

// handleScrollUp scrolls the list up
func (bs *BrowsingScreen) handleScrollUp() {
	// TODO: Phase 4 - Implement scrolling via gotui widgets
	if bs.currentWidget != nil {
		_ = bs.currentWidget.HandleKeypress('k')
	}
}

// handleViewSwitch toggles between hash and filename views
func (bs *BrowsingScreen) handleViewSwitch() error {
	bs.isHashView = !bs.isHashView
	if bs.isHashView {
		bs.currentWidget = bs.hashWidget
	} else {
		bs.currentWidget = bs.filenameWidget
	}
	return nil
}

// handleKeep marks current item as KEEP
func (bs *BrowsingScreen) handleKeep() error {
	if _, err := bs.showKeepModal(); err != nil {
		return fmt.Errorf("keep modal failed: %w", err)
	}
	return nil
}

// handleArchive marks current item as ARCHIVE
func (bs *BrowsingScreen) handleArchive() error {
	if _, err := bs.showArchiveModal(); err != nil {
		return fmt.Errorf("archive modal failed: %w", err)
	}
	return nil
}

// handleDelete marks current item as DELETE
func (bs *BrowsingScreen) handleDelete() error {
	if _, err := bs.showDeleteModal(false); err != nil {
		return fmt.Errorf("delete modal failed: %w", err)
	}
	return nil
}

// handleSkip marks current item as SKIP
func (bs *BrowsingScreen) handleSkip() error {
	if _, err := bs.showSkipModal(); err != nil {
		return fmt.Errorf("skip modal failed: %w", err)
	}
	return nil
}

// handleQuit saves and exits
func (bs *BrowsingScreen) handleQuit(app interface{}) error {
	if _, err := bs.showQuitModal(app); err != nil {
		return fmt.Errorf("quit modal failed: %w", err)
	}
	return nil
}

// showKeepModal displays a confirmation modal for keep action
func (bs *BrowsingScreen) showKeepModal() (bool, error) {
	modal := modals.NewKeepModal()
	result := modal.Execute()
	return result.Confirmed, result.Error
}

// showArchiveModal displays a modal to get archive destination
func (bs *BrowsingScreen) showArchiveModal() (string, error) {
	modal := modals.NewArchiveModal()
	result := modal.Execute()
	return result.Value, result.Error
}

// showDeleteModal displays a delete confirmation modal
func (bs *BrowsingScreen) showDeleteModal(isCritical bool) (bool, error) {
	modal := modals.NewDeleteModal(isCritical)
	result := modal.Execute()
	return result.Confirmed, result.Error
}

// showSkipModal displays a skip confirmation modal
func (bs *BrowsingScreen) showSkipModal() (bool, error) {
	modal := modals.NewSkipModal()
	result := modal.Execute()
	return result.Confirmed, result.Error
}

// showQuitModal displays a quit confirmation modal
func (bs *BrowsingScreen) showQuitModal(app interface{}) (bool, error) {
	if app == nil {
		return false, fmt.Errorf("app not initialized")
	}
	// TODO: Phase 4 - Stop app when available
	return true, nil
}

// updateDetailsPane refreshes the details pane with current selection
func (bs *BrowsingScreen) updateDetailsPane() error {
	if bs.detailsWidget == nil || bs.sessionManager == nil {
		return nil
	}

	// Get current selection
	currentIdx := bs.sessionManager.GetCursor()
	if currentIdx < 0 || currentIdx >= len(bs.sessionManager.GetSession().DuplicateGroups) {
		return nil
	}

	group := bs.sessionManager.GetSession().DuplicateGroups[currentIdx]
	return bs.detailsWidget.ShowGroup(&group)
}

// updateStatusBar updates the status bar with flag counts
func (bs *BrowsingScreen) updateStatusBar() {
	if bs.flagManager == nil {
		return
	}

	viewMode := "HASH"
	if !bs.isHashView {
		viewMode = "FILENAME"
	}

	// TODO: Phase 4 - Get flag counts from flagManager
	keep := 0
	archive := 0
	delete := 0

	statusText := fmt.Sprintf("View: %s | 🟢 Keep: %d | 📦 Archive: %d | 🗑️  Delete: %d | [q]uit [v]iew [k]eep [a]rchive [d]elete [s]kip",
		viewMode, keep, archive, delete)

	fmt.Println(statusText)
}

// updateListPane refreshes the list if needed
func (bs *BrowsingScreen) updateListPane() error {
	// This would trigger re-rendering of the list widget
	// For now, the widgets handle this internally
	return nil
}

// persistFlag saves a flag decision to the database
func (bs *BrowsingScreen) persistFlag(fileID int, flagType string, data map[string]interface{}) error {
	if bs.flagManager == nil || bs.ctx == nil {
		return fmt.Errorf("flag manager or context not initialized")
	}

	// Implementation depends on flag type
	// This is handled by the flag manager's PersistKeep/Archive/Delete methods
	return nil
}
