package components

import (
	"fmt"
	"strings"

	"github.com/haggishunk/filesprawl/internal/consolidation/browsing"
	"github.com/haggishunk/filesprawl/internal/consolidation/flags"
)

// TablePane renders the right panel showing files in selected group
type TablePane struct {
	group         *browsing.DuplicateGroup
	selectedIndex int
	width         int
	height        int
	scrollOffset  int
	wideView      bool // Show filepath column
	focused       bool // Whether this pane has input focus
	flagManager   *browsing.PersistentFlagManager
}

// NewTablePane creates a new TablePane component
func NewTablePane(width, height int, flagManager *browsing.PersistentFlagManager) *TablePane {
	return &TablePane{
		group:         nil,
		selectedIndex: 0,
		width:         width,
		height:        height,
		scrollOffset:  0,
		wideView:      false,
		flagManager:   flagManager,
	}
}

// SetGroup sets the group to display
func (tp *TablePane) SetGroup(group *browsing.DuplicateGroup) {
	tp.group = group
	tp.selectedIndex = 0
	tp.scrollOffset = 0
}

// SetDimensions updates the pane dimensions
func (tp *TablePane) SetDimensions(width, height int) {
	tp.width = width
	tp.height = height
}

// SetFocus sets whether this pane has input focus
func (tp *TablePane) SetFocus(focused bool) {
	tp.focused = focused
}

// ToggleWideView toggles the filepath column
func (tp *TablePane) ToggleWideView() {
	tp.wideView = !tp.wideView
}

// GetSelectedFile returns the currently selected file version
func (tp *TablePane) GetSelectedFile() *browsing.FileVersion {
	if tp.group == nil || tp.selectedIndex < 0 || tp.selectedIndex >= len(tp.group.Versions) {
		return nil
	}
	return &tp.group.Versions[tp.selectedIndex]
}

// SelectionUp moves selection up
func (tp *TablePane) SelectionUp() {
	if tp.group == nil {
		return
	}
	if tp.selectedIndex > 0 {
		tp.selectedIndex--
	}
	if tp.selectedIndex < tp.scrollOffset {
		tp.scrollOffset = tp.selectedIndex
	}
}

// SelectionDown moves selection down
func (tp *TablePane) SelectionDown() {
	if tp.group == nil {
		return
	}
	if tp.selectedIndex < len(tp.group.Versions)-1 {
		tp.selectedIndex++
	}
	maxScroll := tp.selectedIndex - (tp.height - 3)
	if maxScroll > tp.scrollOffset {
		tp.scrollOffset = maxScroll
	}
}

// GetSelectedVersion returns the selected FileVersion
func (tp *TablePane) GetSelectedVersion() *browsing.FileVersion {
	if tp.group == nil || tp.selectedIndex < 0 || tp.selectedIndex >= len(tp.group.Versions) {
		return nil
	}
	return &tp.group.Versions[tp.selectedIndex]
}

// Render returns the rendered view of the table pane
func (tp *TablePane) Render() string {
	if tp.group == nil {
		return tp.renderEmpty()
	}

	var buf strings.Builder

	// Header with focus indicator
	headerChar := "─"
	if tp.focused {
		headerChar = "═"
	}
	header := fmt.Sprintf("┌%s FILES %s┐\n", strings.Repeat(headerChar, 1), strings.Repeat(headerChar, (tp.width-14)))
	buf.WriteString(header)

	// Column headers
	if tp.wideView {
		buf.WriteString("│ remote | filename | filepath | size │\n")
	} else {
		buf.WriteString("│ remote | filename | size | status │\n")
	}
	separatorWidth := tp.width - 2
	if separatorWidth < 5 {
		separatorWidth = 5
	}
	fmt.Fprintf(&buf, "├%s┤\n", strings.Repeat("─", separatorWidth-2))

	// Calculate available space
	contentHeight := tp.height - 5
	if contentHeight < 1 {
		contentHeight = 1
	}

	// Render visible files
	endIdx := tp.scrollOffset + contentHeight
	if endIdx > len(tp.group.Versions) {
		endIdx = len(tp.group.Versions)
	}

	for i := tp.scrollOffset; i < endIdx; i++ {
		tp.renderFileLine(&buf, tp.group.Versions[i], i == tp.selectedIndex)
	}

	// Fill remaining space
	for i := endIdx - tp.scrollOffset; i < contentHeight; i++ {
		buf.WriteString("│                                  │\n")
	}

	return buf.String()
}

// renderEmpty renders empty state
func (tp *TablePane) renderEmpty() string {
	return "┌─ FILES IN GROUP ─────────────────┐\n" +
		"│  (select a group on the left)     │\n" +
		"└──────────────────────────────────┘\n"
}

// renderFileLine renders a single file row
func (tp *TablePane) renderFileLine(buf *strings.Builder, version browsing.FileVersion, isSelected bool) {
	prefix := " "
	if isSelected {
		prefix = "▶"
	}

	// Get flag status if available
	flagStatus := "  "
	if version.Flag != nil {
		switch version.Flag.FlagType {
		case flags.FlagTypeKeep:
			flagStatus = "✓K"
		case flags.FlagTypeArchive:
			flagStatus = "▶A"
		case flags.FlagTypeDelete:
			flagStatus = "✗D"
		}
	}

	// Build location identifier: location name (hostname for local, remote name for rclone)
	remote := "?"
	if version.LocationName != "" {
		remote = version.LocationName
		if len(remote) > 10 {
			remote = remote[:10]
		}
	}

	filename := ""
	if version.FileInfo != nil {
		filename = version.FileInfo.Name
		if len(filename) > 12 {
			filename = filename[:12]
		}
	}

	// Build filepath display
	// AbsolutePath contains the full path ready to display
	// Local: /full/filesystem/path
	// Remote: remote_name:path/within/remote
	// Note: Path now stores the absolute path for local files
	filepath := ""
	if version.FileInfo != nil {
		filepath = version.FileInfo.Path
	}

	// Truncate for display
	maxFilepathLen := (tp.width * 60) / 100
	if maxFilepathLen < 20 {
		maxFilepathLen = 20
	}
	if len(filepath) > maxFilepathLen {
		filepath = filepath[:maxFilepathLen]
	}

	size := formatSize(version.FileInfo.Size)

	// Format based on wide view setting
	var content string
	if tp.wideView {
		// Wide view: remote | filename | filepath | size (status always shown)
		content = fmt.Sprintf("%s %s | %s | %s | %s",
			prefix, remote, filename, filepath, size)
	} else {
		// Normal view: remote | filename | size | status
		content = fmt.Sprintf("%s %s | %s | %s | %s",
			prefix, remote, filename, size, flagStatus)
	}

	lineWidth := tp.width - 2
	if len(content) > lineWidth {
		content = content[:lineWidth]
	}
	content = fmt.Sprintf("%-*s", lineWidth, content)

	fmt.Fprintf(buf, "│%s│\n", content)
}

// formatSize formats file size in human-readable format
func formatSize(bytes int64) string {
	units := []string{"B", "KB", "MB", "GB"}
	size := float64(bytes)
	unit := 0

	for size >= 1024 && unit < len(units)-1 {
		size /= 1024
		unit++
	}

	if unit == 0 {
		return fmt.Sprintf("%dB", int64(size))
	}
	return fmt.Sprintf("%.1f%s", size, units[unit])
}
