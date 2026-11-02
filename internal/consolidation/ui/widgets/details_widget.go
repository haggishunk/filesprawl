package widgets

import (
	"fmt"

	"github.com/haggishunk/filesprawl/internal/consolidation/browsing"
	"github.com/metaspartan/gotui/v5/widgets"
)

// DetailsWidget displays detailed information about a selected DuplicateGroup
type DetailsWidget struct {
	paragraph *widgets.Paragraph
	analyzer  *browsing.PresenceAnalyzer
	x, y, width, height int // Rectangle dimensions
}

// NewDetailsWidget creates a new details widget
func NewDetailsWidget() *DetailsWidget {
	paragraph := widgets.NewParagraph()
	if paragraph != nil {
		paragraph.Title = "📋 DETAILS"
	}

	return &DetailsWidget{
		paragraph: paragraph,
		analyzer:  browsing.NewPresenceAnalyzer(),
	}
}

// GetWidget returns the underlying gotui Paragraph widget
func (w *DetailsWidget) GetWidget() interface{} {
	if w == nil || w.paragraph == nil {
		return nil
	}
	return w.paragraph
}

// ShowGroup displays details for a duplicate group
func (w *DetailsWidget) ShowGroup(group *browsing.DuplicateGroup) error {
	if w == nil || w.paragraph == nil || group == nil {
		return nil
	}

	detail := w.formatGroupDetails(group)
	w.paragraph.Text = detail
	return nil
}

// formatGroupDetails formats the group information for display
func (w *DetailsWidget) formatGroupDetails(group *browsing.DuplicateGroup) string {
	if group == nil {
		return "No group selected"
	}

	output := ""
	output += fmt.Sprintf("═══════════════════════════════════════\n")
	output += fmt.Sprintf("Identifier: %s\n", group.Identifier)
	output += fmt.Sprintf("Total Copies: %d\n", group.Count)

	if group.Flag != nil {
		output += fmt.Sprintf("Flag Status: [%s]\n", group.Flag.FlagType)
	}
	output += fmt.Sprintf("\n")

	// Group versions
	if len(group.Versions) == 0 {
		output += "No versions found\n"
	} else {
		output += fmt.Sprintf("VERSIONS (%d):\n", len(group.Versions))
		for vIdx, version := range group.Versions {
			if version.Hash == nil && version.FileInfo == nil {
				continue
			}
			indicator := w.analyzer.AnalyzeFileVersion(version)
			hashStr := ""
			if version.Hash != nil {
				hashStr = truncateStr(version.Hash.Hash, 12)
			}
			fileStr := ""
			if version.FileInfo != nil {
				fileStr = version.FileInfo.Name
			}
			output += fmt.Sprintf("  %d. %s %s (Hash: %s)\n", vIdx+1, indicator, fileStr, hashStr)

			// Show locations for this version
			if len(version.Locations) > 0 {
				for _, loc := range version.Locations {
					output += fmt.Sprintf("     → %s (%s)\n", loc.Name, loc.Type)
				}
			}
			if version.Flag != nil {
				output += fmt.Sprintf("     [Flag: %s]\n", version.Flag.FlagType)
			}
		}
	}

	output += fmt.Sprintf("\n═══════════════════════════════════════\n")
	return output
}

// SetRect sets the drawing rectangle
func (w *DetailsWidget) SetRect(x, y, width, height int) {
	if w == nil || w.paragraph == nil {
		return
	}
	w.x = x
	w.y = y
	w.width = width
	w.height = height
	w.paragraph.SetRect(x, y, width, height)
}

// GetRect returns the current drawing rectangle
func (w *DetailsWidget) GetRect() (int, int, int, int) {
	if w == nil {
		return 0, 0, 0, 0
	}
	return w.x, w.y, w.width, w.height
}

// Clear clears the displayed content
func (w *DetailsWidget) Clear() {
	if w == nil || w.paragraph == nil {
		return
	}
	w.paragraph.Text = ""
}

// SetTitle sets the widget title
func (w *DetailsWidget) SetTitle(title string) {
	if w == nil || w.paragraph == nil {
		return
	}
	w.paragraph.Title = title
}

// truncateStr truncates a string to a maximum length
func truncateStr(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}
