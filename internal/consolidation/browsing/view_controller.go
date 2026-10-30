package browsing

import (
	"fmt"
)

// ViewController manages view rendering and switching
type ViewController struct {
	session         *SessionManager
	hashRenderer    *HashViewRenderer
	filenameRenderer *FilenameViewRenderer
	presenceAnalyzer *PresenceAnalyzer
}

// NewViewController creates a new view controller
func NewViewController(session *SessionManager) *ViewController {
	return &ViewController{
		session:          session,
		hashRenderer:     NewHashViewRenderer(),
		filenameRenderer: NewFilenameViewRenderer(),
		presenceAnalyzer: NewPresenceAnalyzer(),
	}
}

// RenderCurrentView renders the current view based on active mode
func (vc *ViewController) RenderCurrentView() string {
	session := vc.session.GetSession()

	switch session.CurrentView {
	case ViewTypeHash:
		return vc.hashRenderer.Render(session, 20)
	case ViewTypeFilename:
		return vc.filenameRenderer.Render(session, 20)
	default:
		return "Unknown view type"
	}
}

// GetViewName returns human-readable name of current view
func (vc *ViewController) GetViewName() string {
	switch vc.session.GetCurrentView() {
	case ViewTypeHash:
		return "Hash-Centric (Content)"
	case ViewTypeFilename:
		return "Filename-Centric (Version)"
	default:
		return "Unknown"
	}
}

// GetViewDescription returns description of current view
func (vc *ViewController) GetViewDescription() string {
	switch vc.session.GetCurrentView() {
	case ViewTypeHash:
		return "Browse by content hash - see all identical files across locations"
	case ViewTypeFilename:
		return "Browse by filename - see all versions (different hashes) of same name"
	default:
		return "Unknown view"
	}
}

// SwitchViewAndRender switches to the other view and renders it
func (vc *ViewController) SwitchViewAndRender() string {
	currentView := vc.session.GetCurrentView()
	vc.session.SwitchView()

	output := fmt.Sprintf("\n>>> Switched from %s to %s view <<<\n",
		getViewLabel(currentView),
		getViewLabel(vc.session.GetCurrentView()),
	)

	output += fmt.Sprintf("Tip: %s\n\n", vc.GetViewDescription())
	output += vc.RenderCurrentView()

	return output
}

// HandleKeypress handles keyboard input and returns updated view
func (vc *ViewController) HandleKeypress(key rune) (string, error) {
	switch key {
	case 'k', 'K':
		// Flag current item as KEEP
		return "Flagged as KEEP\n", nil
	case 'a', 'A':
		// Flag current item as ARCHIVE
		return "Flagged as ARCHIVE\n", nil
	case 'd', 'D':
		// Flag current item as DELETE
		return "Flagged as DELETE\n", nil
	case '\t':
		// Switch views (TAB key)
		return vc.SwitchViewAndRender(), nil
	case 'Q':
		// Quit without consolidating
		return "", nil
	case '\n':
		// Expand/collapse current item
		currentItem := vc.getCurrentItem()
		if currentItem != "" {
			vc.session.ExpandItem(currentItem)
		}
		return vc.RenderCurrentView(), nil
	case 'u', 'U':
		// Move cursor up
		vc.session.MoveCursorUp()
		return vc.RenderCurrentView(), nil
	case 'd', 'd':
		// Move cursor down
		vc.session.MoveCursorDown()
		return vc.RenderCurrentView(), nil
	default:
		return "", fmt.Errorf("unknown key: %c", key)
	}
}

// getCurrentItem returns the identifier of the currently selected item
func (vc *ViewController) getCurrentItem() string {
	session := vc.session.GetSession()
	cursor := vc.session.GetCursor()

	if cursor >= 0 && cursor < len(session.DuplicateGroups) {
		return session.DuplicateGroups[cursor].Identifier
	}
	return ""
}

// getViewLabel returns a label for a view type
func getViewLabel(viewType ViewType) string {
	switch viewType {
	case ViewTypeHash:
		return "Hash"
	case ViewTypeFilename:
		return "Filename"
	default:
		return "Unknown"
	}
}

// GetStatusLine returns a status line showing current view and position
func (vc *ViewController) GetStatusLine() string {
	session := vc.session.GetSession()
	viewName := vc.GetViewName()
	cursor := vc.session.GetCursor()
	total := len(session.DuplicateGroups)
	flagged := len(session.Flags)

	return fmt.Sprintf("View: %s | Item: %d/%d | Flagged: %d | TAB=switch Q=quit",
		viewName, cursor+1, total, flagged)
}
