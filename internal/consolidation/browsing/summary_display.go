package browsing

import (
	"fmt"
	"strings"
	"github.com/haggishunk/filesprawl/internal/consolidation/flags"
)

// SummaryRenderer renders consolidation decision summaries
type SummaryRenderer struct {
	presenceAnalyzer *PresenceAnalyzer
}

// NewSummaryRenderer creates a new summary renderer
func NewSummaryRenderer() *SummaryRenderer {
	return &SummaryRenderer{
		presenceAnalyzer: NewPresenceAnalyzer(),
	}
}

// RenderDecisionSummary renders a summary of all flagged decisions
func (sr *SummaryRenderer) RenderDecisionSummary(session *BrowsingSession) string {
	var output strings.Builder

	output.WriteString("\n╔════════════════════════════════════════════════════════════════╗\n")
	output.WriteString("║        CONSOLIDATION DECISION SUMMARY - REVIEW BEFORE ACTION      ║\n")
	output.WriteString("╚════════════════════════════════════════════════════════════════╝\n\n")

	// Count decisions by type
	keepCount, archiveCount, deleteCount := sr.countDecisions(session.Flags)

	// Show summary counts
	output.WriteString(fmt.Sprintf("KEEP:    %3d files | 🟢 Will be preserved\n", keepCount))
	output.WriteString(fmt.Sprintf("ARCHIVE: %3d files | 📦 Will be copied to backup location\n", archiveCount))
	output.WriteString(fmt.Sprintf("DELETE:  %3d files | 🗑️  Will be permanently removed\n", deleteCount))
	output.WriteString(fmt.Sprintf("─────────────────────────────────────\n"))
	output.WriteString(fmt.Sprintf("TOTAL:   %3d files\n\n", keepCount+archiveCount+deleteCount))

	// Show detail breakdown by action
	sr.renderKeepSummary(&output, session.Flags)
	sr.renderArchiveSummary(&output, session.Flags)
	sr.renderDeleteSummary(&output, session.Flags)

	// Show warnings if needed
	sr.renderWarnings(&output, session.Flags)

	output.WriteString("\n╔════════════════════════════════════════════════════════════════╗\n")
	output.WriteString("║ [ENTER] to proceed | [Q] to cancel                               ║\n")
	output.WriteString("╚════════════════════════════════════════════════════════════════╝\n\n")

	return output.String()
}

// renderKeepSummary renders files flagged as KEEP
func (sr *SummaryRenderer) renderKeepSummary(output *strings.Builder, flagMap map[int]*flags.Flag) {
	keepFiles := sr.filterByType(flagMap, flags.FlagTypeKeep)

	if len(keepFiles) == 0 {
		return
	}

	output.WriteString("KEEP - Will be preserved:\n")
	for _, flag := range keepFiles {
		output.WriteString(fmt.Sprintf("  🟢 File ID %d\n", flag.FileID))
		if flag.Reason != nil && *flag.Reason != "" {
			output.WriteString(fmt.Sprintf("     Reason: %s\n", *flag.Reason))
		}
	}
	output.WriteString("\n")
}

// renderArchiveSummary renders files flagged as ARCHIVE
func (sr *SummaryRenderer) renderArchiveSummary(output *strings.Builder, flagMap map[int]*flags.Flag) {
	archiveFiles := sr.filterByType(flagMap, flags.FlagTypeArchive)

	if len(archiveFiles) == 0 {
		return
	}

	output.WriteString("ARCHIVE - Will be copied to backup locations:\n")

	// Group by archive destination
	byDest := make(map[string]int)
	for _, flag := range archiveFiles {
		dest := "unknown"
		if flag.ArchiveDest != nil {
			dest = *flag.ArchiveDest
		}
		byDest[dest]++
	}

	for dest, count := range byDest {
		output.WriteString(fmt.Sprintf("  📦 → %s: %d files\n", dest, count))
	}
	output.WriteString("\n")
}

// renderDeleteSummary renders files flagged as DELETE
func (sr *SummaryRenderer) renderDeleteSummary(output *strings.Builder, flagMap map[int]*flags.Flag) {
	deleteFiles := sr.filterByType(flagMap, flags.FlagTypeDelete)

	if len(deleteFiles) == 0 {
		return
	}

	output.WriteString("DELETE - Will be permanently removed:\n")
	for _, flag := range deleteFiles {
		output.WriteString(fmt.Sprintf("  🗑️  File ID %d\n", flag.FileID))
		if flag.Reason != nil && *flag.Reason != "" {
			output.WriteString(fmt.Sprintf("     Reason: %s\n", *flag.Reason))
		}
	}
	output.WriteString("\n")
}

// renderWarnings renders any warnings about the consolidation
func (sr *SummaryRenderer) renderWarnings(output *strings.Builder, flagMap map[int]*flags.Flag) {
	deleteCount := len(sr.filterByType(flagMap, flags.FlagTypeDelete))

	if deleteCount == 0 {
		return
	}

	output.WriteString("⚠️  WARNINGS:\n")
	output.WriteString(fmt.Sprintf("  • %d files will be permanently deleted\n", deleteCount))
	output.WriteString("  • This action cannot be undone\n")
	output.WriteString("  • Files should be backed up before proceeding\n\n")
}

// countDecisions counts flags by type
func (sr *SummaryRenderer) countDecisions(flagMap map[int]*flags.Flag) (int, int, int) {
	keep := 0
	archive := 0
	delete := 0

	for _, flag := range flagMap {
		switch flag.FlagType {
		case flags.FlagTypeKeep:
			keep++
		case flags.FlagTypeArchive:
			archive++
		case flags.FlagTypeDelete:
			delete++
		}
	}

	return keep, archive, delete
}

// filterByType returns all flags of a specific type
func (sr *SummaryRenderer) filterByType(flagMap map[int]*flags.Flag, flagType flags.FlagType) []*flags.Flag {
	var result []*flags.Flag

	for _, flag := range flagMap {
		if flag.FlagType == flagType {
			result = append(result, flag)
		}
	}

	return result
}

// RenderCompactSummary returns a one-line summary
func (sr *SummaryRenderer) RenderCompactSummary(session *BrowsingSession) string {
	keep, archive, delete := sr.countDecisions(session.Flags)
	return fmt.Sprintf("[K:%d A:%d D:%d]", keep, archive, delete)
}
