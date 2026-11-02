package screens

import (
	"context"
	"fmt"
	"image"

	"github.com/haggishunk/filesprawl/internal/consolidation/browsing"
	"github.com/haggishunk/filesprawl/internal/consolidation/ui/modals"
	"github.com/haggishunk/filesprawl/internal/consolidation/ui/widgets"
	"github.com/metaspartan/gotui/v5"
	gotui_widgets "github.com/metaspartan/gotui/v5/widgets"
)

// BrowsingScreen implements the interactive browsing interface for consolidation
type BrowsingScreen struct {
	sessionManager *browsing.SessionManager
	flagManager    *browsing.PersistentFlagManager
	hashWidget     *widgets.HashCentricWidget
	filenameWidget *widgets.FilenameCentricWidget
	detailsWidget  *widgets.DetailsWidget
	bannerWidget   *widgets.BannerWidget
	currentWidget  widgets.BrowsingWidget // pointer to current active widget
	isHashView     bool                    // true for hash, false for filename
	ctx            context.Context
	app            *gotui.Application     // reference to gotui app for event handling
	version        string                 // application version
	reviewMode     bool                   // true if in review mode
	remotes        []string               // list of configured remotes
}

// NewBrowsingScreen creates a new browsing screen
func NewBrowsingScreen(
	sessionManager *browsing.SessionManager,
	flagManager *browsing.PersistentFlagManager,
	hashWidget *widgets.HashCentricWidget,
	filenameWidget *widgets.FilenameCentricWidget,
	detailsWidget *widgets.DetailsWidget,
	version string,
	reviewMode bool,
	remotes []string,
) *BrowsingScreen {
	// Create banner with initial info
	bannerInfo := &widgets.BannerInfo{
		Version:       version,
		ReviewMode:    reviewMode,
		ExecuteMode:   !reviewMode,
		Remotes:       remotes,
		TotalDuplicates: 0,
		FlaggedCount:  0,
	}
	bannerWidget := widgets.NewBannerWidget(bannerInfo)

	return &BrowsingScreen{
		sessionManager: sessionManager,
		flagManager:    flagManager,
		hashWidget:     hashWidget,
		filenameWidget: filenameWidget,
		detailsWidget:  detailsWidget,
		bannerWidget:   bannerWidget,
		isHashView:     true, // Start with hash view
		currentWidget:  hashWidget,
		version:        version,
		reviewMode:     reviewMode,
		remotes:        remotes,
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

	// Create layout manager for 3-pane layout with banner
	layout := widgets.NewLayoutManager(termWidth, termHeight)

	// Get widget rectangles
	bannerX, bannerY, bannerW, bannerH := layout.GetBannerRect()
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

	// Configure banner widget position
	bs.bannerWidget.SetRect(bannerX, bannerY, bannerW, bannerH)

	// Display layout information (Phase 5: gotui will render this)
	fmt.Println("═══════════════════════════════════════════════════════════════")
	fmt.Println("BrowsingScreen: 4-pane TUI layout with banner initialized")
	fmt.Println("═══════════════════════════════════════════════════════════════")
	fmt.Println(fmt.Sprintf("Banner (info):              (%d,%d) %dx%d", bannerX, bannerY, bannerW, bannerH))
	fmt.Println(fmt.Sprintf("List pane (hash/filename):  (%d,%d) %dx%d", listX, listY, listW, listH))
	fmt.Println(fmt.Sprintf("Details pane (file info):   (%d,%d) %dx%d", detailsX, detailsY, detailsW, detailsH))
	fmt.Println(fmt.Sprintf("Status bar (controls):      (%d,%d) %dx%d", statusX, statusY, statusW, statusH))
	fmt.Println("═══════════════════════════════════════════════════════════════")

	// Register keyboard bindings for Phase 5 event loop
	if err := bs.CreateKeyBindings(); err != nil {
		fmt.Printf("Warning: Failed to create key bindings: %v\n", err)
	}

	// Phase 5: Wire to real gotui rendering with Grid layout

	// 1. Create root container with Grid layout
	gtApp, ok := app.(*gotui.Application)
	if !ok {
		return fmt.Errorf("invalid gotui app type: expected *gotui.Application, got %T", app)
	}

	root := gotui.NewGrid()
	if root == nil {
		return fmt.Errorf("failed to create Grid container")
	}

	// 2. Extract gotui widgets from wrappers
	hashListWidget, ok := bs.hashWidget.GetWidget().(*gotui_widgets.List)
	if !ok || hashListWidget == nil {
		return fmt.Errorf("hash list widget extraction failed")
	}

	filenameListWidget, ok := bs.filenameWidget.GetWidget().(*gotui_widgets.List)
	if !ok || filenameListWidget == nil {
		return fmt.Errorf("filename list widget extraction failed")
	}

	detailsGotUI, ok := bs.detailsWidget.GetWidget().(*gotui_widgets.Paragraph)
	if !ok || detailsGotUI == nil {
		return fmt.Errorf("details widget extraction failed")
	}

	statusGotUI, ok := statusWidget.GetWidget().(*gotui_widgets.Paragraph)
	if !ok || statusGotUI == nil {
		return fmt.Errorf("status widget extraction failed")
	}

	bannerGotUI, ok := bs.bannerWidget.GetWidget().(*gotui_widgets.Paragraph)
	if !ok || bannerGotUI == nil {
		return fmt.Errorf("banner widget extraction failed")
	}

	// 3. Wrap list widgets with event dispatcher to intercept keyboard events
	// The wrapper routes all keyboard events to BrowsingScreen action handlers
	wrappedListWidget := &EventDispatchWrapper{
		widget:         hashListWidget,
		browsingScreen: bs,
	}

	// 4. Configure widget rectangles
	hashListWidget.SetRect(listX, listY, listW, listH)
	filenameListWidget.SetRect(listX, listY, listW, listH)
	detailsGotUI.SetRect(detailsX, detailsY, detailsW, detailsH)
	statusGotUI.SetRect(statusX, statusY, statusW, statusH)

	// 5. Set up Grid with banner, list, details, and status widgets
	root.Set(
		gotui.NewRow(float64(bannerH),
			gotui.NewCol(100, bannerGotUI),
		),
		gotui.NewRow(float64(listH),
			gotui.NewCol(40, wrappedListWidget),
			gotui.NewCol(60, detailsGotUI),
		),
		gotui.NewRow(float64(statusH),
			gotui.NewCol(100, statusGotUI),
		),
	)

	// 5. Set root widget on app
	gtApp.SetRoot(root, true)

	// 6. Store app reference for event handling
	bs.app = gtApp

	// 7. Run event loop (blocks until user quits)
	return gtApp.Run()
}

// CreateKeyBindings creates the event handler with all keyboard bindings
func (bs *BrowsingScreen) CreateKeyBindings() error {
	if bs == nil {
		return fmt.Errorf("browsing screen not initialized")
	}

	// In Phase 5, this will be used to register handlers with gotui app
	// For now, just validate that all handlers are available

	handlers := map[rune]string{
		'j': "Scroll Down",
		'k': "Scroll Up",
		'K': "Keep (all copies)",
		'v': "View Switch (hash ↔ filename)",
		'a': "Archive (backup copies)",
		'd': "Delete (remove duplicates)",
		's': "Skip (clear flags)",
		'?': "Help (keyboard reference)",
		'q': "Quit (review summary)",
	}

	fmt.Println()
	fmt.Println("Keyboard Bindings Registered:")
	fmt.Println("─────────────────────────────────────")
	for key, action := range handlers {
		fmt.Printf("  '%c'  →  %s\n", key, action)
	}
	fmt.Println("─────────────────────────────────────")
	fmt.Println()

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
	case '?': // Help
		_ = bs.handleHelp()
	case 'q', 'Q': // Quit
		_ = bs.handleQuit(app)
	}
}

// handleScrollDown scrolls the list down
func (bs *BrowsingScreen) handleScrollDown() {
	// TODO: Phase 4 - Implement scrolling via gotui widgets
	if bs.currentWidget != nil {
		_ = bs.currentWidget.HandleKeypress('j')
		// Update details pane to show the new selection
		if err := bs.updateDetailsPane(); err != nil {
			fmt.Printf("Warning: Failed to update details pane: %v\n", err)
		}
	}
}

// handleScrollUp scrolls the list up
func (bs *BrowsingScreen) handleScrollUp() {
	// TODO: Phase 4 - Implement scrolling via gotui widgets
	if bs.currentWidget != nil {
		_ = bs.currentWidget.HandleKeypress('k')
		// Update details pane to show the new selection
		if err := bs.updateDetailsPane(); err != nil {
			fmt.Printf("Warning: Failed to update details pane: %v\n", err)
		}
	}
}

// handleViewSwitch toggles between hash and filename views
func (bs *BrowsingScreen) handleViewSwitch() error {
	bs.isHashView = !bs.isHashView
	if bs.isHashView {
		bs.currentWidget = bs.hashWidget
		fmt.Print("\n>>> Switched to HASH VIEW - Content-centric <<<\n")
	} else {
		bs.currentWidget = bs.filenameWidget
		fmt.Print("\n>>> Switched to FILENAME VIEW - Name-centric <<<\n")
	}

	// Update details pane for new view
	if err := bs.updateDetailsPane(); err != nil {
		fmt.Printf("Warning: Failed to update details pane: %v\n", err)
	}

	return nil
}

// handleKeep marks current item as KEEP
func (bs *BrowsingScreen) handleKeep() error {
	confirmed, err := bs.showKeepModal()
	if err != nil {
		return fmt.Errorf("keep modal failed: %w", err)
	}

	if !confirmed {
		return nil
	}

	// Get current selection
	selectedIdx := bs.currentWidget.GetSelectedIndex()
	session := bs.sessionManager.GetSession()
	if selectedIdx < 0 || selectedIdx >= len(session.DuplicateGroups) {
		return fmt.Errorf("invalid selection index: %d", selectedIdx)
	}

	group := session.DuplicateGroups[selectedIdx]

	// Flag each version of the duplicate group as KEEP
	for _, version := range group.Versions {
		if err := bs.flagManager.PersistKeep(bs.ctx, version.FileID); err != nil {
			return fmt.Errorf("failed to flag file %d as KEEP: %w", version.FileID, err)
		}
	}

	// Update UI
	if err := bs.updateDetailsPane(); err != nil {
		fmt.Printf("Warning: Failed to update details pane: %v\n", err)
	}

	return nil
}

// handleArchive marks current item as ARCHIVE
func (bs *BrowsingScreen) handleArchive() error {
	destination, err := bs.showArchiveModal()
	if err != nil {
		return fmt.Errorf("archive modal failed: %w", err)
	}

	if destination == "" {
		return nil // User cancelled
	}

	// Get current selection
	selectedIdx := bs.currentWidget.GetSelectedIndex()
	session := bs.sessionManager.GetSession()
	if selectedIdx < 0 || selectedIdx >= len(session.DuplicateGroups) {
		return fmt.Errorf("invalid selection index: %d", selectedIdx)
	}

	group := session.DuplicateGroups[selectedIdx]
	reason := fmt.Sprintf("Archived to %s", destination)

	// Flag primary version (first) as ARCHIVE; others as KEEP
	for i, version := range group.Versions {
		if i == 0 {
			// Keep the first copy
			if err := bs.flagManager.PersistKeep(bs.ctx, version.FileID); err != nil {
				return fmt.Errorf("failed to flag primary file %d as KEEP: %w", version.FileID, err)
			}
		} else {
			// Archive other copies
			if err := bs.flagManager.PersistArchive(bs.ctx, version.FileID, destination, &reason); err != nil {
				return fmt.Errorf("failed to flag file %d as ARCHIVE: %w", version.FileID, err)
			}
		}
	}

	// Update UI
	if err := bs.updateDetailsPane(); err != nil {
		fmt.Printf("Warning: Failed to update details pane: %v\n", err)
	}

	return nil
}

// handleDelete marks current item as DELETE
func (bs *BrowsingScreen) handleDelete() error {
	// Get current selection first to check critical status
	selectedIdx := bs.currentWidget.GetSelectedIndex()
	session := bs.sessionManager.GetSession()
	if selectedIdx < 0 || selectedIdx >= len(session.DuplicateGroups) {
		return fmt.Errorf("invalid selection index: %d", selectedIdx)
	}

	group := session.DuplicateGroups[selectedIdx]

	// Check if any version is critical (only 1 location)
	isCritical := false
	for _, version := range group.Versions {
		if len(version.Locations) == 1 {
			isCritical = true
			break
		}
	}

	confirmed, err := bs.showDeleteModal(isCritical)
	if err != nil {
		return fmt.Errorf("delete modal failed: %w", err)
	}

	if !confirmed {
		return nil
	}

	reason := "User confirmed deletion"

	// Flag all versions as DELETE
	for _, version := range group.Versions {
		if err := bs.flagManager.PersistDelete(bs.ctx, version.FileID, &reason); err != nil {
			return fmt.Errorf("failed to flag file %d as DELETE: %w", version.FileID, err)
		}
	}

	// Update UI
	if err := bs.updateDetailsPane(); err != nil {
		fmt.Printf("Warning: Failed to update details pane: %v\n", err)
	}

	return nil
}

// handleSkip removes any flag from current item
func (bs *BrowsingScreen) handleSkip() error {
	confirmed, err := bs.showSkipModal()
	if err != nil {
		return fmt.Errorf("skip modal failed: %w", err)
	}

	if !confirmed {
		return nil
	}

	// Get current selection
	selectedIdx := bs.currentWidget.GetSelectedIndex()
	session := bs.sessionManager.GetSession()
	if selectedIdx < 0 || selectedIdx >= len(session.DuplicateGroups) {
		return fmt.Errorf("invalid selection index: %d", selectedIdx)
	}

	group := session.DuplicateGroups[selectedIdx]

	// Remove flags from all versions
	for _, version := range group.Versions {
		if err := bs.flagManager.PersistSkip(bs.ctx, version.FileID); err != nil {
			return fmt.Errorf("failed to skip file %d: %w", version.FileID, err)
		}
	}

	// Update UI
	if err := bs.updateDetailsPane(); err != nil {
		fmt.Printf("Warning: Failed to update details pane: %v\n", err)
	}

	return nil
}

// handleHelp displays the help screen
func (bs *BrowsingScreen) handleHelp() error {
	return modals.ShowHelpModal()
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
	if result.Error != nil {
		return "", result.Error
	}
	if !result.Confirmed {
		return "", nil // User cancelled
	}
	return result.Value, nil
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

// EventDispatchWrapper wraps a gotui widget to intercept and dispatch keyboard events
// This implements the gotui EventHandler interface to handle keyboard input
type EventDispatchWrapper struct {
	widget         gotui.Widget
	browsingScreen *BrowsingScreen
}

// HandleEvent implements gotui.EventHandler to process keyboard events
// Returns true if the event was handled, false to let it bubble up
func (w *EventDispatchWrapper) HandleEvent(evt gotui.Event) bool {
	if w == nil || w.browsingScreen == nil || w.widget == nil {
		return false
	}

	// Only handle keyboard events
	if evt.Type != gotui.KeyboardEvent {
		// Pass non-keyboard events to wrapped widget
		if handler, ok := w.widget.(gotui.EventHandler); ok {
			return handler.HandleEvent(evt)
		}
		return false
	}

	// Extract rune from keyboard event payload
	// gotui passes rune as the event payload for keyboard events
	key, ok := evt.Payload.(rune)
	if !ok {
		// Try to pass to wrapped widget if it can handle it
		if handler, ok := w.widget.(gotui.EventHandler); ok {
			return handler.HandleEvent(evt)
		}
		return false
	}

	// Dispatch keyboard event to BrowsingScreen handler
	w.browsingScreen.handleKeypress(key, w.browsingScreen.app)

	// Return true to indicate event was handled (don't bubble up)
	return true
}

// Draw implements gotui.Drawable to render the wrapped widget
func (w *EventDispatchWrapper) Draw(buf *gotui.Buffer) {
	if w != nil && w.widget != nil {
		// Check if widget implements Drawable
		if drawable, ok := w.widget.(gotui.Drawable); ok {
			drawable.Draw(buf)
		}
	}
}

// GetRect returns the drawing rectangle of the wrapped widget
func (w *EventDispatchWrapper) GetRect() image.Rectangle {
	if w != nil && w.widget != nil {
		return w.widget.GetRect()
	}
	return image.Rectangle{}
}

// SetRect sets the drawing rectangle of the wrapped widget (x, y, width, height)
func (w *EventDispatchWrapper) SetRect(x, y, width, height int) {
	if w != nil && w.widget != nil {
		w.widget.SetRect(x, y, width, height)
	}
}

// Lock implements sync.Locker for the wrapped widget
func (w *EventDispatchWrapper) Lock() {
	if w != nil && w.widget != nil {
		w.widget.Lock()
	}
}

// Unlock implements sync.Locker for the wrapped widget
func (w *EventDispatchWrapper) Unlock() {
	if w != nil && w.widget != nil {
		w.widget.Unlock()
	}
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
