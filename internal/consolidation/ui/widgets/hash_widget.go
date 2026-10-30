package widgets

import (
	"fmt"
	"github.com/metaspartan/gotui/v5"
	"github.com/metaspartan/gotui/v5/widgets"
	"github.com/haggishunk/filesprawl/internal/consolidation/browsing"
)

// HashCentricWidget displays files grouped by content hash
type HashCentricWidget struct {
	list      *BaseListWidget
	duplicates []browsing.DuplicateGroup
	indicators map[int]browsing.PresenceIndicator
	analyzer  *browsing.PresenceAnalyzer
}

// NewHashCentricWidget creates a new hash-centric browsing widget
func NewHashCentricWidget() *HashCentricWidget {
	return &HashCentricWidget{
		list:       NewBaseListWidget("📎 HASH VIEW - Content-Centric"),
		duplicates: []browsing.DuplicateGroup{},
		indicators: make(map[int]browsing.PresenceIndicator),
		analyzer:   browsing.NewPresenceAnalyzer(),
	}
}

// GetWidget returns the underlying gotui widget
func (w *HashCentricWidget) GetWidget() gotui.Drawable {
	return w.list.GetWidget()
}

// UpdateData updates the widget with duplicate groups
func (w *HashCentricWidget) UpdateData(data interface{}) error {
	groups, ok := data.([]browsing.DuplicateGroup)
	if !ok {
		return fmt.Errorf("expected []DuplicateGroup, got %T", data)
	}

	w.duplicates = groups
	w.renderList()
	return nil
}

// renderList renders the list based on duplicates
func (w *HashCentricWidget) renderList() {
	rows := []string{}

	for idx, group := range w.duplicates {
		indicator := w.analyzer.AnalyzeDuplicateGroup(group)
		w.indicators[idx] = indicator

		// Format: indicator hash... (N copies) [FLAG]
		flagStr := ""
		if group.Flag != nil {
			flagStr = fmt.Sprintf(" [%s]", group.Flag.FlagType)
		}

		row := fmt.Sprintf("%s Hash: %s... (%d copies)%s",
			indicator,
			w.truncateHash(group.Identifier, 8),
			group.Count,
			flagStr,
		)
		rows = append(rows, row)
	}

	w.list.SetRows(rows)
}

// GetSelectedGroup returns the currently selected duplicate group
func (w *HashCentricWidget) GetSelectedGroup() *browsing.DuplicateGroup {
	idx := w.list.GetSelectedIndex()
	if idx >= 0 && idx < len(w.duplicates) {
		return &w.duplicates[idx]
	}
	return nil
}

// GetSelectedIndex returns the currently selected item index
func (w *HashCentricWidget) GetSelectedIndex() int {
	return w.list.GetSelectedIndex()
}

// SetSelectedIndex sets the selected item index
func (w *HashCentricWidget) SetSelectedIndex(index int) error {
	return w.list.SetSelectedIndex(index)
}

// GetTitle returns the widget title
func (w *HashCentricWidget) GetTitle() string {
	return w.list.GetTitle()
}

// HandleKeypress handles keyboard input
func (w *HashCentricWidget) HandleKeypress(key rune) error {
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
func (w *HashCentricWidget) SetRect(x, y, width, height int) {
	w.list.SetRect(x, y, width, height)
}

// GetDuplicateCount returns total number of duplicates
func (w *HashCentricWidget) GetDuplicateCount() int {
	return len(w.duplicates)
}

// GetPresenceIndicator returns the indicator for a group index
func (w *HashCentricWidget) GetPresenceIndicator(index int) browsing.PresenceIndicator {
	if indicator, ok := w.indicators[index]; ok {
		return indicator
	}
	return browsing.PresenceGreen
}

// truncateHash truncates a hash for display
func (w *HashCentricWidget) truncateHash(hash string, length int) string {
	if len(hash) <= length {
		return hash
	}
	return hash[:length] + "..."
}

// GetGroupDetails returns formatted details for a group
func (w *HashCentricWidget) GetGroupDetails(index int) string {
	if index < 0 || index >= len(w.duplicates) {
		return "No group selected"
	}

	group := w.duplicates[index]
	detail := fmt.Sprintf("\n📎 Hash Details: %s\n", group.Identifier)
	detail += fmt.Sprintf("Total Copies: %d\n\n", group.Count)

	// Count locations
	locations := make(map[string]int)
	for _, version := range group.Versions {
		for _, loc := range version.Locations {
			locations[loc.Name]++
		}
	}

	detail += "Locations:\n"
	for locName, count := range locations {
		indicator := w.analyzer.AnalyzeFileVersion(browsing.FileVersion{
			Locations: make([]browsing.Location, count),
		})
		detail += fmt.Sprintf("  %s %s (%d copies)\n", indicator, locName, count)
	}

	return detail
}
