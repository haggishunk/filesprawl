package components

import (
	"fmt"
	"strings"
)

// HelpModal displays keyboard shortcuts and usage instructions
type HelpModal struct {
	Width  int
	Height int
}

// NewHelpModal creates a new help modal
func NewHelpModal(width, height int) *HelpModal {
	return &HelpModal{
		Width:  width,
		Height: height,
	}
}

// Render returns the rendered help modal
func (hm *HelpModal) Render() string {
	var buf strings.Builder

	// Modal background and border
	borderChar := "═"
	sideBorder := "║"

	// Header
	header := fmt.Sprintf("╔%s KEYBOARD SHORTCUTS %s╗", strings.Repeat(borderChar, 1), strings.Repeat(borderChar, 23))
	buf.WriteString(header)
	buf.WriteString("\n")

	// Help content
	lines := []string{
		"NAVIGATION:",
		"  j/k         Scroll up/down in focused pane",
		"  Tab         Switch between List and Table panes",
		"",
		"TABLE PANE (file list):",
		"  K            Flag file as KEEP (preserve copy)",
		"  a            Flag file as ARCHIVE (backup)",
		"  d            Flag file as DELETE (removal)",
		"  s            Skip/clear all flags on file",
		"  v            Toggle wide view (show filepath)",
		"",
		"LIST PANE (hash groups):",
		"  j/k          Navigate to different hash groups",
		"",
		"GENERAL:",
		"  ?            Show this help",
		"  q            Quit application",
		"  ctrl+c       Force quit",
		"",
		"SAFETY:",
		"  Cannot delete ALL copies of a file",
		"  Must keep at least 1 copy per remote/filename",
	}

	// Calculate content width
	contentWidth := hm.Width - 4 // Account for borders and padding
	if contentWidth < 40 {
		contentWidth = 40
	}

	// Render content lines
	for _, line := range lines {
		// Truncate or pad line to fit
		if len(line) > contentWidth {
			line = line[:contentWidth]
		}
		paddedLine := fmt.Sprintf("%-*s", contentWidth, line)
		fmt.Fprintf(&buf, "%s %s %s\n", sideBorder, paddedLine, sideBorder)
	}

	// Footer
	footer := fmt.Sprintf("╚%s Press any key to close %s╝", strings.Repeat(borderChar, 1), strings.Repeat(borderChar, 22))
	buf.WriteString(footer)
	buf.WriteString("\n")

	return buf.String()
}
