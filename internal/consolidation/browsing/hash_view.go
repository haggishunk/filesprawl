package browsing

import (
	"fmt"
	"strings"

	"github.com/haggishunk/filesprawl/internal/consolidation/flags"
)

// HashViewRenderer renders the hash-centric view
type HashViewRenderer struct{}

// NewHashViewRenderer creates a new hash view renderer
func NewHashViewRenderer() *HashViewRenderer {
	return &HashViewRenderer{}
}

// Render renders the hash-centric view
func (r *HashViewRenderer) Render(session *BrowsingSession, pageSize int) string {
	if len(session.DuplicateGroups) == 0 {
		return "No duplicate files found."
	}

	var output strings.Builder
	output.WriteString("\n📊 CONSOLIDATION BROWSER - BY HASH (Content-Centric)\n")
	output.WriteString("════════════════════════════════════════════════════════════════\n\n")

	// Show summary
	flagCounts := countFlags(session.Flags)
	output.WriteString(fmt.Sprintf("Total duplicate groups: %d | Flagged: %d | ",
		len(session.DuplicateGroups),
		len(session.Flags),
	))
	output.WriteString(fmt.Sprintf("Actions: [k]eep [a]rchive [d]elete [TAB]view [Q]uit\n\n"))

	// Render each duplicate group
	for i, group := range session.DuplicateGroups {
		r.renderGroup(&output, group, session.ExpandedItems, i == session.Cursor)
	}

	// Show navigation
	output.WriteString("\n════════════════════════════════════════════════════════════════\n")
	output.WriteString("Summary: ")
	for flagType, count := range flagCounts {
		if count > 0 {
			output.WriteString(fmt.Sprintf("%s=%d ", flagType, count))
		}
	}
	output.WriteString("\n\n")

	return output.String()
}

// renderGroup renders a single duplicate group
func (r *HashViewRenderer) renderGroup(output *strings.Builder, group DuplicateGroup, expanded map[string]bool, highlighted bool) {
	prefix := "└─ "
	if highlighted {
		prefix = ">>> "
	}

	indicator := getIndicator(group)
	status := ""
	if group.Flag != nil {
		status = fmt.Sprintf(" [%s]", group.Flag.FlagType)
	}

	output.WriteString(fmt.Sprintf("%sHash: %s (%d copies)%s %s\n",
		prefix,
		truncateHash(group.Identifier, 8),
		group.Count,
		status,
		indicator,
	))

	// Render expanded details if needed
	if expanded[group.Identifier] {
		r.renderGroupDetails(output, group)
	} else {
		output.WriteString("    [Press ENTER to expand]\n\n")
	}
}

// renderGroupDetails renders detailed info for a group
func (r *HashViewRenderer) renderGroupDetails(output *strings.Builder, group DuplicateGroup) {
	// Group files by location
	locationMap := make(map[string][]FileVersion)
	for _, fv := range group.Versions {
		for _, loc := range fv.Locations {
			locationMap[loc.Name] = append(locationMap[loc.Name], fv)
		}
	}

	// Render each location
	for locName, files := range locationMap {
		indicator := getLocationIndicator(len(files), len(group.Versions))
		output.WriteString(fmt.Sprintf("    %s %s (%s) — %d copy/copies\n",
			indicator,
			locName,
			"local", // TODO: Get actual type from location
			len(files),
		))

		// Show files in this location
		for _, fv := range files {
			flagStr := ""
			if fv.Flag != nil {
				flagStr = fmt.Sprintf(" [%s]", fv.Flag.FlagType)
			}
			output.WriteString(fmt.Sprintf("      • %s @ %s%s\n",
				fv.FileInfo.Name,
				truncatePath(fv.FileInfo.Path, 40),
				flagStr,
			))
		}
	}

	output.WriteString("\n")
}

// getIndicator returns the appropriate indicator for a group
func getIndicator(group DuplicateGroup) string {
	// If all files have multiple copies, green
	// If some have limited copies, yellow
	// If any are single copies, red
	for _, version := range group.Versions {
		if len(version.Locations) <= 1 {
			return "🔴"
		}
		if len(version.Locations) == 2 {
			return "🟡"
		}
	}
	return "🟢"
}

// getLocationIndicator returns indicator for a specific location
func getLocationIndicator(copies, total int) string {
	if copies == 1 && total > 1 {
		return "🟡"
	}
	if copies == total {
		return "🔴"
	}
	return "🟢"
}

// truncateHash truncates a hash to a shorter display form
func truncateHash(hash string, length int) string {
	if len(hash) <= length {
		return hash
	}
	return hash[:length] + "..."
}

// truncatePath truncates a path to fit display
func truncatePath(path string, maxLen int) string {
	if len(path) <= maxLen {
		return path
	}
	// Show end of path if it's too long
	start := len(path) - maxLen + 3
	if start < 0 {
		start = 0
	}
	return "..." + path[start:]
}

// countFlags counts flagged files by type
func countFlags(flagsMap map[int]*flags.Flag) map[string]int {
	counts := make(map[string]int)
	for _, flag := range flagsMap {
		counts[string(flag.FlagType)]++
	}
	return counts
}
