package browsing

import (
	"fmt"
	"strings"
)

// FilenameViewRenderer renders the filename-centric view
type FilenameViewRenderer struct{}

// NewFilenameViewRenderer creates a new filename view renderer
func NewFilenameViewRenderer() *FilenameViewRenderer {
	return &FilenameViewRenderer{}
}

// Render renders the filename-centric view
func (r *FilenameViewRenderer) Render(session *BrowsingSession, pageSize int) string {
	if len(session.DuplicateGroups) == 0 {
		return "No duplicate files found."
	}

	var output strings.Builder
	output.WriteString("\n📋 CONSOLIDATION BROWSER - BY FILENAME (Name-Centric)\n")
	output.WriteString("════════════════════════════════════════════════════════════════\n\n")

	// Show summary
	flagCounts := countFlags(session.Flags)
	output.WriteString(fmt.Sprintf("Total files: %d | Different names: %d | Flagged: %d\n",
		len(session.Flags),
		len(session.DuplicateGroups),
		len(session.Flags),
	))
	output.WriteString(fmt.Sprintf("Actions: [k]eep [a]rchive [d]elete [TAB]view [Q]uit\n\n"))

	// Render each file (filename) with its versions
	for i, group := range session.DuplicateGroups {
		r.renderFilename(&output, group, session.ExpandedItems, i == session.Cursor)
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

// renderFilename renders a filename with all its versions
func (r *FilenameViewRenderer) renderFilename(output *strings.Builder, group DuplicateGroup, expanded map[string]bool, highlighted bool) {
	prefix := "└─ "
	if highlighted {
		prefix = ">>> "
	}

	indicator := getIndicator(group)
	status := ""
	if group.Flag != nil {
		status = fmt.Sprintf(" [%s]", group.Flag.FlagType)
	}

	// Count unique hashes (versions)
	versionCount := len(group.Versions)

	output.WriteString(fmt.Sprintf("%s%s (%d version%s)%s %s\n",
		prefix,
		group.Identifier,
		versionCount,
		pluralize(versionCount),
		status,
		indicator,
	))

	// Render version details if expanded
	if expanded[group.Identifier] {
		r.renderVersionDetails(output, group)
	} else {
		output.WriteString("    [Press ENTER to expand]\n\n")
	}
}

// renderVersionDetails shows all versions of a filename
func (r *FilenameViewRenderer) renderVersionDetails(output *strings.Builder, group DuplicateGroup) {
	for vIdx, version := range group.Versions {
		// Show version header with hash and location count
		versionNum := vIdx + 1
		indicator := getVersionIndicator(version)
		flagStr := ""
		if version.Flag != nil {
			flagStr = fmt.Sprintf(" [%s]", version.Flag.FlagType)
		}

		output.WriteString(fmt.Sprintf("    Version %d (hash: %s, %d location%s) %s%s\n",
			versionNum,
			truncateHash(version.Hash.HashValue, 6),
			len(version.Locations),
			pluralize(len(version.Locations)),
			indicator,
			flagStr,
		))

		// Show locations
		for _, loc := range version.Locations {
			locIndicator := getLocationIndicator(1, len(group.Versions))
			output.WriteString(fmt.Sprintf("      %s %s\n",
				locIndicator,
				loc.Name,
			))
		}
	}
	output.WriteString("\n")
}

// getVersionIndicator returns indicator for a file version
func getVersionIndicator(version FileVersion) string {
	if len(version.Locations) == 0 {
		return "🔴"
	}
	if len(version.Locations) == 1 {
		return "🟡"
	}
	return "🟢"
}

// pluralize returns plural form if needed
func pluralize(count int) string {
	if count == 1 {
		return ""
	}
	return "s"
}
