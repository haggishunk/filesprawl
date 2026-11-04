package widgets

import (
	"fmt"

	"github.com/haggishunk/filesprawl/internal/consolidation/browsing"
	"github.com/metaspartan/gotui/v5/widgets"
)

// StatusWidget displays real-time status and instructions
type StatusWidget struct {
	paragraph           *widgets.Paragraph
	x, y, width, height int
}

// NewStatusWidget creates a status bar widget
func NewStatusWidget() *StatusWidget {
	return &StatusWidget{
		paragraph: widgets.NewParagraph(),
	}
}

// UpdateStatus updates the status display with current session state
func (sw *StatusWidget) UpdateStatus(session *browsing.BrowsingSession) error {
	if sw == nil || sw.paragraph == nil || session == nil {
		return fmt.Errorf("status widget or session not initialized")
	}

	// Count flags by type
	keepCount := 0
	archiveCount := 0
	deleteCount := 0

	if session.Flags != nil {
		for _, flag := range session.Flags {
			if flag == nil {
				continue
			}
			switch flag.FlagType {
			case "KEEP":
				keepCount++
			case "ARCHIVE":
				archiveCount++
			case "DELETE":
				deleteCount++
			}
		}
	}

	// Format status message
	view := string(session.CurrentView)
	status := fmt.Sprintf("View: [%s] | Navigation: [j/k] Switch: [v] | Actions: [K]eep [a]rchive [d]elete [s]kip [q]uit | Flags: 🟢Keep=%d 📦Archive=%d 🗑️Delete=%d",
		view, keepCount, archiveCount, deleteCount)

	sw.paragraph.Text = status
	sw.paragraph.Title = "Status Bar"

	return nil
}

// GetWidget returns the underlying gotui widget
func (sw *StatusWidget) GetWidget() interface{} {
	if sw == nil {
		return nil
	}
	return sw.paragraph
}

// SetRect sets the position and size
func (sw *StatusWidget) SetRect(x, y, width, height int) error {
	if sw == nil {
		return fmt.Errorf("status widget not initialized")
	}
	sw.x, sw.y, sw.width, sw.height = x, y, width, height
	if sw.paragraph != nil {
		sw.paragraph.SetRect(x, y, width, height)
	}
	return nil
}

// GetRect returns current position and size
func (sw *StatusWidget) GetRect() (int, int, int, int) {
	if sw == nil {
		return 0, 0, 0, 0
	}
	return sw.x, sw.y, sw.width, sw.height
}

// Clear clears the status display
func (sw *StatusWidget) Clear() error {
	if sw == nil || sw.paragraph == nil {
		return fmt.Errorf("status widget not initialized")
	}
	sw.paragraph.Text = ""
	return nil
}
