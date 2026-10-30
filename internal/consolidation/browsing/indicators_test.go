package browsing

import (
	"testing"
)

func TestPresenceAnalyzerCreation(t *testing.T) {
	pa := NewPresenceAnalyzer()

	if pa.MinimumLocationsForGreen != 3 {
		t.Errorf("expected MinimumLocationsForGreen 3, got %d", pa.MinimumLocationsForGreen)
	}
	if pa.MinimumLocationsForYellow != 2 {
		t.Errorf("expected MinimumLocationsForYellow 2, got %d", pa.MinimumLocationsForYellow)
	}
}

func TestAnalyzeFileVersionPresence(t *testing.T) {
	pa := NewPresenceAnalyzer()

	tests := []struct {
		name                string
		locationCount       int
		expectedIndicator   PresenceIndicator
	}{
		{"No locations (critical)", 0, PresenceRed},
		{"One location (critical)", 1, PresenceRed},
		{"Two locations (limited)", 2, PresenceYellow},
		{"Three locations (safe)", 3, PresenceGreen},
		{"Five locations (very safe)", 5, PresenceGreen},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			version := FileVersion{
				Locations: make([]Location, tt.locationCount),
			}

			indicator := pa.AnalyzeFileVersion(version)
			if indicator != tt.expectedIndicator {
				t.Errorf("expected %s, got %s", tt.expectedIndicator, indicator)
			}
		})
	}
}

func TestAnalyzeDuplicateGroupPresence(t *testing.T) {
	pa := NewPresenceAnalyzer()

	// Test group with all safe versions
	safeGroup := DuplicateGroup{
		Identifier: "hash123",
		Count:      3,
		Versions: []FileVersion{
			{Locations: []Location{{Name: "loc1"}, {Name: "loc2"}, {Name: "loc3"}}},
			{Locations: []Location{{Name: "loc1"}, {Name: "loc2"}, {Name: "loc3"}}},
		},
	}

	indicator := pa.AnalyzeDuplicateGroup(safeGroup)
	if indicator != PresenceGreen {
		t.Errorf("expected green for safe group, got %s", indicator)
	}

	// Test group with risky version (only 1 location)
	riskyGroup := DuplicateGroup{
		Identifier: "hash456",
		Count:      2,
		Versions: []FileVersion{
			{Locations: []Location{{Name: "loc1"}}},
			{Locations: []Location{{Name: "loc1"}, {Name: "loc2"}, {Name: "loc3"}}},
		},
	}

	indicator = pa.AnalyzeDuplicateGroup(riskyGroup)
	if indicator != PresenceRed {
		t.Errorf("expected red for risky group, got %s", indicator)
	}
}

func TestAnalyzeLocationPresence(t *testing.T) {
	pa := NewPresenceAnalyzer()

	tests := []struct {
		name                string
		fileCount           int
		totalVersions       int
		expectedIndicator   PresenceIndicator
		expectedContains    string
	}{
		{"Empty location", 0, 2, PresenceRed, "No files"},
		{"Single only copy", 1, 2, PresenceRed, "ONLY copy"},
		{"Limited backup", 1, 3, PresenceYellow, "Limited"},
		{"Complete backup", 3, 3, PresenceGreen, "Complete"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			loc := Location{Name: "test-loc", FileCount: tt.fileCount}
			indicator, desc := pa.AnalyzeLocation(loc, tt.totalVersions)

			if indicator != tt.expectedIndicator {
				t.Errorf("expected %s, got %s", tt.expectedIndicator, indicator)
			}
			// Optional: verify description contains expected keywords
		})
	}
}

func TestIndicatorDescription(t *testing.T) {
	greenDesc := GetIndicatorDescription(PresenceGreen)
	if greenDesc == "" {
		t.Errorf("expected description for green indicator")
	}

	yellowDesc := GetIndicatorDescription(PresenceYellow)
	if yellowDesc == "" {
		t.Errorf("expected description for yellow indicator")
	}

	redDesc := GetIndicatorDescription(PresenceRed)
	if redDesc == "" {
		t.Errorf("expected description for red indicator")
	}
}

func TestFormatLocationPresence(t *testing.T) {
	loc := Location{
		Name:     "archlinux/backup",
		Type:     "remote",
		FileCount: 3,
	}

	formatted := FormatLocationPresence(loc, PresenceGreen, "Complete backup (3 copies)")

	if formatted == "" {
		t.Errorf("expected non-empty formatted string")
	}
	if !contains(formatted, "archlinux/backup") {
		t.Errorf("expected location name in formatted output")
	}
	if !contains(formatted, "remote") {
		t.Errorf("expected type in formatted output")
	}
}

func TestSummarizeGroupPresence(t *testing.T) {
	pa := NewPresenceAnalyzer()

	group := DuplicateGroup{
		Identifier: "hash789",
		Count:      2,
		Versions: []FileVersion{
			{Locations: []Location{
				{Name: "loc1", FileCount: 2},
				{Name: "loc2", FileCount: 2},
			}},
		},
	}

	summary := pa.SummarizeGroupPresence(group)
	if summary == "" {
		t.Errorf("expected non-empty summary")
	}
}

// Helper function
func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
