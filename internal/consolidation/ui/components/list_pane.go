package components

import (
	"fmt"
	"strings"

	"github.com/haggishunk/filesprawl/internal/consolidation/browsing"
)

// ListPane renders the left panel showing duplicate groups (hash view)
type ListPane struct {
	groups        []browsing.DuplicateGroup
	selectedIndex int
	width         int
	height        int
	scrollOffset  int // For scrolling when list is longer than height
}

// NewListPane creates a new ListPane component
func NewListPane(groups []browsing.DuplicateGroup, width, height int) *ListPane {
	return &ListPane{
		groups:        groups,
		selectedIndex: 0,
		width:         width,
		height:        height,
		scrollOffset:  0,
	}
}

// SetDimensions updates the pane dimensions
func (lp *ListPane) SetDimensions(width, height int) {
	lp.width = width
	lp.height = height
}

// SelectionUp moves selection up
func (lp *ListPane) SelectionUp() {
	if lp.selectedIndex > 0 {
		lp.selectedIndex--
	}
	// Adjust scroll offset to keep selected item visible
	if lp.selectedIndex < lp.scrollOffset {
		lp.scrollOffset = lp.selectedIndex
	}
}

// SelectionDown moves selection down
func (lp *ListPane) SelectionDown() {
	if lp.selectedIndex < len(lp.groups)-1 {
		lp.selectedIndex++
	}
	// Adjust scroll offset to keep selected item visible
	maxScroll := lp.selectedIndex - (lp.height - 2)
	if maxScroll > lp.scrollOffset {
		lp.scrollOffset = maxScroll
	}
}

// GetSelectedIndex returns the currently selected group index
func (lp *ListPane) GetSelectedIndex() int {
	return lp.selectedIndex
}

// GetSelectedGroup returns the selected DuplicateGroup
func (lp *ListPane) GetSelectedGroup() *browsing.DuplicateGroup {
	if lp.selectedIndex >= 0 && lp.selectedIndex < len(lp.groups) {
		return &lp.groups[lp.selectedIndex]
	}
	return nil
}

// Render returns the rendered view of the list pane
func (lp *ListPane) Render() string {
	var buf strings.Builder

	// Calculate available space
	contentHeight := lp.height - 2
	if contentHeight < 1 {
		contentHeight = 1
	}

	// Render header
	buf.WriteString("┌─ HASH LIST ─────────────────────┐\n")

	// Calculate which groups to show
	endIdx := lp.scrollOffset + contentHeight
	if endIdx > len(lp.groups) {
		endIdx = len(lp.groups)
	}

	// Render visible groups
	for i := lp.scrollOffset; i < endIdx; i++ {
		lp.renderGroupLine(&buf, lp.groups[i], i == lp.selectedIndex)
	}

	// Fill remaining space
	for i := endIdx - lp.scrollOffset; i < contentHeight; i++ {
		buf.WriteString("│                                 │\n")
	}

	return buf.String()
}

// renderGroupLine renders a single group line
func (lp *ListPane) renderGroupLine(buf *strings.Builder, group browsing.DuplicateGroup, isSelected bool) {
	indicator := getPresenceIndicator(group)
	hashDisplay := truncateHash(group.Identifier, 6)

	// Count unique locations
	locationSet := make(map[string]bool)
	for _, v := range group.Versions {
		for _, loc := range v.Locations {
			locationSet[loc.Name] = true
		}
	}
	locCount := len(locationSet)

	// Build the line
	var prefix string
	if isSelected {
		prefix = "▶"
	} else {
		prefix = " "
	}

	content := fmt.Sprintf("%s %s %s (%d copies, %d loc)",
		prefix, indicator, hashDisplay, group.Count, locCount)

	// Pad to width and truncate if needed
	lineWidth := lp.width - 2 // Account for borders
	if len(content) > lineWidth {
		content = content[:lineWidth]
	}
	content = fmt.Sprintf("%-*s", lineWidth, content)

	fmt.Fprintf(buf, "│%s│\n", content)
}

// getPresenceIndicator returns the presence indicator for a group
func getPresenceIndicator(group browsing.DuplicateGroup) string {
	// Get unique locations
	locationSet := make(map[string]bool)
	for _, v := range group.Versions {
		for _, loc := range v.Locations {
			locationSet[loc.Name] = true
		}
	}

	locCount := len(locationSet)

	// Determine indicator based on location count
	if locCount >= 3 {
		return "🟢" // Safe - multiple locations
	} else if locCount == 2 {
		return "🟡" // Caution - limited replication
	}
	return "🔴" // Critical - only one location
}

// truncateHash returns a truncated hash string
func truncateHash(hash string, length int) string {
	if len(hash) > length {
		return hash[:length] + "…"
	}
	return hash
}
