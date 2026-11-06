package components

import (
	"github.com/haggishunk/filesprawl/internal/consolidation/browsing"
)

// FilenameTablePane renders the lower right panel showing files grouped by filename
// It reuses TablePane rendering logic since both work with DuplicateGroup
type FilenameTablePane struct {
	tablePane *TablePane // Reuse table pane rendering
}

// NewFilenameTablePane creates a new FilenameTablePane component
func NewFilenameTablePane(width, height int, flagManager *browsing.PersistentFlagManager) *FilenameTablePane {
	return &FilenameTablePane{
		tablePane: NewTablePane(width, height, flagManager),
	}
}

// SetGroup sets the filename group to display
// The group parameter should be a DuplicateGroup with Identifier=filename
func (ftp *FilenameTablePane) SetGroup(group *browsing.DuplicateGroup) {
	if ftp.tablePane != nil {
		ftp.tablePane.SetGroup(group)
	}
}

// SetDimensions updates the pane dimensions
func (ftp *FilenameTablePane) SetDimensions(width, height int) {
	if ftp.tablePane != nil {
		ftp.tablePane.SetDimensions(width, height)
	}
}

// SetFocus sets whether this pane has input focus
func (ftp *FilenameTablePane) SetFocus(focused bool) {
	if ftp.tablePane != nil {
		ftp.tablePane.SetFocus(focused)
	}
}

// ToggleWideView toggles the filepath column
func (ftp *FilenameTablePane) ToggleWideView() {
	if ftp.tablePane != nil {
		ftp.tablePane.ToggleWideView()
	}
}

// GetSelectedFile returns the currently selected file version
func (ftp *FilenameTablePane) GetSelectedFile() *browsing.FileVersion {
	if ftp.tablePane != nil {
		return ftp.tablePane.GetSelectedFile()
	}
	return nil
}

// SelectionUp moves selection up
func (ftp *FilenameTablePane) SelectionUp() {
	if ftp.tablePane != nil {
		ftp.tablePane.SelectionUp()
	}
}

// SelectionDown moves selection down
func (ftp *FilenameTablePane) SelectionDown() {
	if ftp.tablePane != nil {
		ftp.tablePane.SelectionDown()
	}
}

// GetSelectedVersion returns the selected FileVersion
func (ftp *FilenameTablePane) GetSelectedVersion() *browsing.FileVersion {
	if ftp.tablePane != nil {
		return ftp.tablePane.GetSelectedVersion()
	}
	return nil
}

// Render returns the rendered view of the filename table pane
func (ftp *FilenameTablePane) Render() string {
	if ftp.tablePane == nil {
		return "(filename table pane not initialized)\n"
	}
	// The table pane renders files grouped by filename
	// Rendering is identical to hash table pane, just with different data
	return ftp.tablePane.Render()
}
