package widgets

import (
	"fmt"

	"github.com/haggishunk/filesprawl/internal/consolidation/browsing"
)

// FilenameCentricWidget displays files grouped by filename
// NOTE: Phase 5 TUI integration pending. Stub implementation for now.
type FilenameCentricWidget struct {
	list       *BaseListWidget
	duplicates []browsing.DuplicateGroup
	indicators map[int]browsing.PresenceIndicator
	analyzer   *browsing.PresenceAnalyzer
}

// NewFilenameCentricWidget creates a new filename-centric browsing widget
func NewFilenameCentricWidget() *FilenameCentricWidget {
	return &FilenameCentricWidget{
		list:       NewBaseListWidget("📋 FILENAME VIEW - Version-Centric"),
		duplicates: []browsing.DuplicateGroup{},
		indicators: make(map[int]browsing.PresenceIndicator),
		analyzer:   browsing.NewPresenceAnalyzer(),
	}
}

// GetWidget returns the underlying widget
func (w *FilenameCentricWidget) GetWidget() interface{} {
	return w.list.GetWidget()
}

// UpdateData updates the widget with duplicate groups
func (w *FilenameCentricWidget) UpdateData(data interface{}) error {
	groups, ok := data.([]browsing.DuplicateGroup)
	if !ok {
		return fmt.Errorf("expected []DuplicateGroup, got %T", data)
	}

	w.duplicates = groups
	w.renderList()
	return nil
}

// renderList renders the list based on duplicates
func (w *FilenameCentricWidget) renderList() {
	rows := []string{}

	for idx, group := range w.duplicates {
		indicator := w.analyzer.AnalyzeDuplicateGroup(group)
		w.indicators[idx] = indicator

		// Format: indicator filename (N versions) [FLAG]
		versionCount := len(group.Versions)
		versionWord := "version"
		if versionCount > 1 {
			versionWord = "versions"
		}

		flagStr := ""
		if group.Flag != nil {
			flagStr = fmt.Sprintf(" [%s]", group.Flag.FlagType)
		}

		row := fmt.Sprintf("%s %s (%d %s)%s",
			indicator,
			group.Identifier,
			versionCount,
			versionWord,
			flagStr,
		)
		rows = append(rows, row)
	}

	w.list.SetRows(rows)
}

// GetSelectedGroup returns the currently selected duplicate group
func (w *FilenameCentricWidget) GetSelectedGroup() *browsing.DuplicateGroup {
	idx := w.list.GetSelectedIndex()
	if idx >= 0 && idx < len(w.duplicates) {
		return &w.duplicates[idx]
	}
	return nil
}

// GetSelectedIndex returns the currently selected item index
func (w *FilenameCentricWidget) GetSelectedIndex() int {
	return w.list.GetSelectedIndex()
}

// SetSelectedIndex sets the selected item index
func (w *FilenameCentricWidget) SetSelectedIndex(index int) error {
	return w.list.SetSelectedIndex(index)
}

// GetTitle returns the widget title
func (w *FilenameCentricWidget) GetTitle() string {
	return w.list.GetTitle()
}

// HandleKeypress handles keyboard input
func (w *FilenameCentricWidget) HandleKeypress(key rune) error {
	switch key {
	case 'j', 'k': // vim-style navigation
		if key == 'j' {
			w.list.ScrollDown()
		} else {
			w.list.ScrollUp()
		}
	}
	return nil
}

// SetRect sets the drawing rectangle
func (w *FilenameCentricWidget) SetRect(x, y, width, height int) {
	w.list.SetRect(x, y, width, height)
}

// GetDuplicateCount returns total number of duplicate filenames
func (w *FilenameCentricWidget) GetDuplicateCount() int {
	return len(w.duplicates)
}

// GetPresenceIndicator returns the indicator for a group index
func (w *FilenameCentricWidget) GetPresenceIndicator(index int) browsing.PresenceIndicator {
	if indicator, ok := w.indicators[index]; ok {
		return indicator
	}
	return browsing.PresenceGreen
}

// GetGroupDetails returns formatted details for a filename group
func (w *FilenameCentricWidget) GetGroupDetails(index int) string {
	if index < 0 || index >= len(w.duplicates) {
		return "No group selected"
	}

	group := w.duplicates[index]
	detail := fmt.Sprintf("\n📋 Filename: %s\n", group.Identifier)
	detail += fmt.Sprintf("Different Versions: %d\n\n", len(group.Versions))

	// Show each version
	for vIdx, version := range group.Versions {
		indicator := w.analyzer.AnalyzeFileVersion(version)
		detail += fmt.Sprintf("Version %d: %s (%d location%s)\n",
			vIdx+1,
			w.truncateHash(version.Hash.Hash, 12),
			len(version.Locations),
			w.plural(len(version.Locations)),
		)

		// Show locations for this version
		for _, loc := range version.Locations {
			detail += fmt.Sprintf("  %s %s\n", indicator, loc.Name)
		}
	}

	return detail
}

// truncateHash truncates a hash for display
func (w *FilenameCentricWidget) truncateHash(hash string, length int) string {
	if len(hash) <= length {
		return hash
	}
	return hash[:length] + "..."
}

// plural returns "s" if count != 1
func (w *FilenameCentricWidget) plural(count int) string {
	if count == 1 {
		return ""
	}
	return "s"
}
