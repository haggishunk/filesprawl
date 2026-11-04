package widgets

import (
	"fmt"

	"github.com/haggishunk/filesprawl/internal/consolidation/browsing"
	"github.com/metaspartan/gotui/v5/widgets"
)

// DetailsWidget displays detailed information about a selected DuplicateGroup
type DetailsWidget struct {
	paragraph           *widgets.Paragraph
	analyzer            *browsing.PresenceAnalyzer
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
	output += "\n╔════════════════════════════════════════════════════════════════╗\n"
	output += fmt.Sprintf("║ DUPLICATE GROUP - %s                                         ║\n", truncateStr(group.Identifier, 30))
	output += "╚════════════════════════════════════════════════════════════════╝\n\n"

	output += fmt.Sprintf("📊 Summary: %d copies found", group.Count)
	if group.Flag != nil {
		output += fmt.Sprintf(" │ Status: [%s]", group.Flag.FlagType)
	}
	output += "\n\n"

	// Group versions
	if len(group.Versions) == 0 {
		output += "No versions found\n"
	} else {
		output += fmt.Sprintf("📋 COPIES (%d):\n", len(group.Versions))
		output += "─────────────────────────────────────────────────────────────────\n"

		for vIdx, version := range group.Versions {
			if version.Hash == nil && version.FileInfo == nil {
				continue
			}

			indicator := w.analyzer.AnalyzeFileVersion(version)
			flagStr := ""
			if version.Flag != nil {
				flagStr = fmt.Sprintf(" 🚩 [%s]", version.Flag.FlagType)
			}

			// File name and details
			fileStr := "Unknown"
			if version.FileInfo != nil {
				fileStr = version.FileInfo.Name
			}

			output += fmt.Sprintf("\n%d️⃣  %s %s%s\n", vIdx+1, indicator, fileStr, flagStr)

			// Show detailed information
			if version.FileInfo != nil {
				meta := version.FileInfo
				output += fmt.Sprintf("   📁 Path:     %s\n", meta.Path)
				if meta.Size > 0 {
					output += fmt.Sprintf("   📏 Size:     %s\n", formatSize(meta.Size))
				}
			}

			// Show locations for this version
			if len(version.Locations) > 0 {
				output += "   📍 Locations:\n"
				for _, loc := range version.Locations {
					locType := "📂"
					if loc.Type == "remote" {
						locType = "☁️"
					}
					output += fmt.Sprintf("      %s %s (%s)\n", locType, loc.Name, loc.Type)
				}
			} else {
				output += "   ⚠️  No location information\n"
			}
		}

		output += "\n─────────────────────────────────────────────────────────────────\n"

		// Show locations summary
		output += w.formatLocationsSummary(group)
	}

	output += "\n🎯 Actions: [K]eep | [a]rchive | [d]elete | [s]kip | [v]iew toggle\n"
	return output
}

// formatLocationsSummary shows a summary of all locations across all copies
func (w *DetailsWidget) formatLocationsSummary(group *browsing.DuplicateGroup) string {
	if group == nil || len(group.Versions) == 0 {
		return ""
	}

	output := "\n📍 Location Summary:\n"

	// Collect all unique locations
	locations := make(map[string]int)
	for _, version := range group.Versions {
		for _, loc := range version.Locations {
			locations[loc.Name]++
		}
	}

	// Show which locations have copies
	if len(locations) == 0 {
		output += "   No location information available\n"
	} else {
		for locName, count := range locations {
			indicator := "🟢"
			if count == 1 && group.Count > 1 {
				indicator = "🔴" // Only copy in this location
			} else if count < group.Count {
				indicator = "🟡" // Partial copy
			}
			output += fmt.Sprintf("   %s %s (%d/%d copies)\n", indicator, locName, count, group.Count)
		}
	}

	return output
}

// formatSize formats bytes into human-readable size
func formatSize(bytes int64) string {
	units := []string{"B", "KB", "MB", "GB", "TB"}
	size := float64(bytes)
	unitIdx := 0

	for size >= 1024 && unitIdx < len(units)-1 {
		size /= 1024
		unitIdx++
	}

	return fmt.Sprintf("%.2f %s", size, units[unitIdx])
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
