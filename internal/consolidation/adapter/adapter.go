package adapter

import (
	"github.com/haggishunk/filesprawl/internal/analysis"
	"github.com/haggishunk/filesprawl/internal/consolidation/browsing"
	"github.com/haggishunk/filesprawl/internal/object"
)

// ToGroupsForBrowsing converts analysis.DuplicateGroup to browsing.DuplicateGroup.
// Handles nil safety and empty slices appropriately.
func ToGroupsForBrowsing(groups []analysis.DuplicateGroup) []browsing.DuplicateGroup {
	if len(groups) == 0 {
		return []browsing.DuplicateGroup{}
	}

	result := make([]browsing.DuplicateGroup, 0, len(groups))
	for _, group := range groups {
		result = append(result, toGroupForBrowsing(group))
	}
	return result
}

// toGroupForBrowsing converts a single analysis.DuplicateGroup to browsing.DuplicateGroup.
func toGroupForBrowsing(group analysis.DuplicateGroup) browsing.DuplicateGroup {
	versions := toFileVersions(group.Files, group.HashValue, group.HashType)

	return browsing.DuplicateGroup{
		Identifier: group.HashValue,
		Count:      group.Count,
		Versions:   versions,
		Flag:       nil, // Group-level flags are not set during initial conversion
	}
}

// toFileVersions converts analysis.FileInfo to browsing.FileVersion objects.
// Each FileVersion includes location information extracted from the FileInfo.
func toFileVersions(files []analysis.FileInfo, hashValue, hashType string) []browsing.FileVersion {
	if len(files) == 0 {
		return []browsing.FileVersion{}
	}

	versions := make([]browsing.FileVersion, 0, len(files))
	for _, file := range files {
		locations := extractLocations(files, file)

		meta := &object.Meta{
			Id:         file.ID,
			Name:       file.Name,
			Path:       file.Path,
			MimeType:   file.MimeType,
			Size:       file.Size,
			LocationId: file.LocationId,
			Persisted:  true,
		}

		hash := &object.Hash{
			Hash:      hashValue,
			Type:      hashType,
			Persisted: true,
		}

		fv := browsing.FileVersion{
			FileID:       file.ID,
			FileInfo:     meta,
			Hash:         hash,
			Locations:    locations,
			Flag:         nil, // File-level flags are not set during initial conversion
			LocationId:   file.LocationId,
			LocationName: file.LocationName,
		}
		versions = append(versions, fv)
	}
	return versions
}

// extractLocations builds Location objects for a file based on distinct locations in the group.
// Groups all files by location ID and counts files per location.
func extractLocations(allFiles []analysis.FileInfo, targetFile analysis.FileInfo) []browsing.Location {
	if len(allFiles) == 0 {
		return []browsing.Location{}
	}

	// Build map of distinct locations and their file counts
	locationMap := make(map[int]*browsing.Location)

	for _, file := range allFiles {
		locID := file.LocationId
		loc := locationMap[locID]
		if loc == nil {
			loc = &browsing.Location{
				ID:        locID,
				Name:      file.LocationName,
				Type:      file.LocationType,
				FileCount: 0,
				Indicator: browsing.PresenceGreen, // Default, may be updated elsewhere
			}
			locationMap[locID] = loc
		}
		loc.FileCount++
	}

	// Return locations only for the target file's location
	if loc, ok := locationMap[targetFile.LocationId]; ok {
		return []browsing.Location{*loc}
	}

	return []browsing.Location{}
}
