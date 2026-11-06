package components

import (
	"fmt"
	"strings"
)

// Banner displays header information about the consolidation session
type Banner struct {
	version    string
	groupCount int
	reviewMode bool
	remotes    []string
	width      int
}

// NewBanner creates a new Banner component
func NewBanner(version string, groupCount int, reviewMode bool, remotes []string, width int) *Banner {
	return &Banner{
		version:    version,
		groupCount: groupCount,
		reviewMode: reviewMode,
		remotes:    remotes,
		width:      width,
	}
}

// SetDimensions updates the banner width
func (b *Banner) SetDimensions(width int) {
	b.width = width
}

// Update updates the banner data
func (b *Banner) Update(groupCount int) {
	b.groupCount = groupCount
}

// Render returns the rendered banner
func (b *Banner) Render() string {
	mode := "REVIEW"
	if !b.reviewMode {
		mode = "EXECUTE"
	}

	remoteStr := fmt.Sprintf("%v", b.remotes)
	if len(remoteStr) > 30 {
		remoteStr = remoteStr[:30] + "…"
	}

	content := fmt.Sprintf("📎 Filesprawl v%s | Groups: %d | Mode: %s | Remotes: %s",
		b.version, b.groupCount, mode, remoteStr)

	// Pad and truncate to width
	if len(content) > b.width {
		content = content[:b.width]
	}
	content = fmt.Sprintf("%-*s", b.width, content)

	return content + "\n"
}

// StatusBar displays keyboard shortcuts and flag summary
type StatusBar struct {
	keepCount    int
	archiveCount int
	deleteCount  int
	focusedPane  string // "list" or "table"
	width        int
}

// NewStatusBar creates a new StatusBar component
func NewStatusBar(width int) *StatusBar {
	return &StatusBar{
		keepCount:    0,
		archiveCount: 0,
		deleteCount:  0,
		width:        width,
	}
}

// SetDimensions updates the status bar width
func (sb *StatusBar) SetDimensions(width int) {
	sb.width = width
}

// SetFocusedPane sets which pane has focus
func (sb *StatusBar) SetFocusedPane(pane string) {
	sb.focusedPane = pane
}

// UpdateCounts updates the flag counts
func (sb *StatusBar) UpdateCounts(keep, archive, del int) {
	sb.keepCount = keep
	sb.archiveCount = archive
	sb.deleteCount = del
}

// Render returns the rendered status bar
func (sb *StatusBar) Render() string {
	// Show keyboard shortcuts based on focused pane
	var shortcuts string
	if sb.focusedPane == "list" {
		shortcuts = "[LIST] j/k: nav | Tab: switch | ?: help | q: quit"
	} else {
		shortcuts = "[TABLE] j/k: nav | Tab: switch | w: wide | v: preview | K/a/d/s: flag | ?: help | q: quit"
	}

	// Summary on second line
	summary := fmt.Sprintf("Flags: %d Keep | %d Archive | %d Delete",
		sb.keepCount, sb.archiveCount, sb.deleteCount)

	var result strings.Builder
	result.WriteString(shortcuts)
	result.WriteString("\n")
	result.WriteString(summary)
	result.WriteString("\n")

	return result.String()
}
