package browsing

import (
	"fmt"
)

// PresenceAnalyzer analyzes backup replication status
type PresenceAnalyzer struct {
	// MinimumLocationsForGreen defines how many locations required for green status
	MinimumLocationsForGreen int
	// MinimumLocationsForYellow defines how many locations required for yellow status
	MinimumLocationsForYellow int
}

// NewPresenceAnalyzer creates a new analyzer with default thresholds
func NewPresenceAnalyzer() *PresenceAnalyzer {
	return &PresenceAnalyzer{
		MinimumLocationsForGreen:  3,
		MinimumLocationsForYellow: 2,
	}
}

// AnalyzeFileVersion returns presence indicator for a file version
func (pa *PresenceAnalyzer) AnalyzeFileVersion(version FileVersion) PresenceIndicator {
	locationCount := len(version.Locations)

	if locationCount == 0 {
		return PresenceRed // No backup at all
	}
	if locationCount == 1 {
		return PresenceRed // Only one location, critical risk
	}
	if locationCount < pa.MinimumLocationsForYellow {
		return PresenceYellow // Limited replication
	}
	if locationCount < pa.MinimumLocationsForGreen {
		return PresenceYellow // Getting closer to safe
	}

	return PresenceGreen // Well replicated
}

// AnalyzeDuplicateGroup returns presence indicator for entire group
func (pa *PresenceAnalyzer) AnalyzeDuplicateGroup(group DuplicateGroup) PresenceIndicator {
	// Group is as safe as its least replicated copy
	worstIndicator := PresenceGreen

	for _, version := range group.Versions {
		indicator := pa.AnalyzeFileVersion(version)
		// Red < Yellow < Green in terms of severity
		if indicator == PresenceRed {
			return PresenceRed
		}
		if indicator == PresenceYellow && worstIndicator == PresenceGreen {
			worstIndicator = PresenceYellow
		}
	}

	return worstIndicator
}

// AnalyzeLocation returns presence indicator and description for a location
func (pa *PresenceAnalyzer) AnalyzeLocation(location Location, totalVersions int) (PresenceIndicator, string) {
	if location.FileCount == 0 {
		return PresenceRed, "No files in this location"
	}

	if location.FileCount == 1 && totalVersions > 1 {
		return PresenceRed, "ONLY copy! No backup elsewhere"
	}

	if location.FileCount < totalVersions {
		return PresenceYellow, fmt.Sprintf("Limited backup (%d/%d)", location.FileCount, totalVersions)
	}

	return PresenceGreen, fmt.Sprintf("Complete backup (%d copies)", location.FileCount)
}

// GetIndicatorDescription returns human-readable description of indicator
func GetIndicatorDescription(indicator PresenceIndicator) string {
	switch indicator {
	case PresenceGreen:
		return "Safe - Well replicated across multiple locations"
	case PresenceYellow:
		return "Caution - Limited replication, backup recommended"
	case PresenceRed:
		return "Critical - Single location or no backup, high risk"
	default:
		return "Unknown"
	}
}

// GetLocationPresenceDisplay formats location with indicator and details
func FormatLocationPresence(location Location, indicator PresenceIndicator, description string) string {
	return fmt.Sprintf("%s %s (%s) — %s",
		indicator,
		location.Name,
		location.Type,
		description,
	)
}

// SummarizeGroupPresence returns a summary of all locations in a group with indicators
func (pa *PresenceAnalyzer) SummarizeGroupPresence(group DuplicateGroup) string {
	if len(group.Versions) == 0 {
		return "No versions found"
	}

	output := ""
	seenLocations := make(map[string]Location)

	// Collect all unique locations
	for _, version := range group.Versions {
		for _, loc := range version.Locations {
			if _, seen := seenLocations[loc.Name]; !seen {
				seenLocations[loc.Name] = loc
			}
		}
	}

	// Analyze each location
	for locName, loc := range seenLocations {
		indicator, desc := pa.AnalyzeLocation(loc, len(group.Versions))
		output += FormatLocationPresence(loc, indicator, desc) + "\n"
	}

	return output
}
