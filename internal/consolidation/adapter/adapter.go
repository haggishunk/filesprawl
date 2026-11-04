package adapter

import (
	"strings"

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
			Id:        file.ID,
			Name:      file.Name,
			Path:      file.Path,
			MimeType:  file.MimeType,
			Size:      file.Size,
			ScanRoot:  file.ScanRoot,
			Persisted: true,
		}

		hash := &object.Hash{
			Hash:      hashValue,
			Type:      hashType,
			Persisted: true,
		}

		// For local files, compute the full canonical path
		// DB stores: RemoteName = scan root, Path = relative path
		// UI expects: RemoteName = full canonical path
		remoteName := file.RemoteName
		if file.ScanRoot != "" {
			// Local file: construct full canonical path
			if file.Path != "" && !strings.HasPrefix(file.Path, "/") {
				// Path is relative, prepend remote name
				remoteName = file.RemoteName + "/" + file.Path
			} else if file.Path != "" {
				// Path is absolute (shouldn't happen, but handle it)
				remoteName = file.Path
			}
		}
		// For rclone files, RemoteName is already correct

		fv := browsing.FileVersion{
			FileID:     file.ID,
			FileInfo:   meta,
			Hash:       hash,
			Locations:  locations,
			Flag:       nil, // File-level flags are not set during initial conversion
			RemoteName: remoteName,
			Hostname:   file.Hostname,
		}
		versions = append(versions, fv)
	}
	return versions
}

// extractLocations builds Location objects for a file based on distinct locations in the group.
// Groups all files by location (hostname/remoteName combination) and counts files per location.
func extractLocations(allFiles []analysis.FileInfo, targetFile analysis.FileInfo) []browsing.Location {
	if len(allFiles) == 0 {
		return []browsing.Location{}
	}

	// Build map of distinct locations and their file counts
	locationMap := make(map[string]*browsing.Location)
	locOrder := []string{} // Track insertion order

	for _, file := range allFiles {
		locKey := getLocationKey(file)
		if _, exists := locationMap[locKey]; !exists {
			_ = append(locOrder, locKey) // locOrder not used in current implementation
		}
		loc := locationMap[locKey]
		if loc == nil {
			locType := inferLocationType(file)
			loc = &browsing.Location{
				Name:      getLocationName(file),
				Type:      locType,
				FileCount: 0,
				Indicator: browsing.PresenceGreen, // Default, may be updated elsewhere
			}
			locationMap[locKey] = loc
		}
		loc.FileCount++
	}

	// Return locations only for the target file's location
	targetLocKey := getLocationKey(targetFile)
	if loc, ok := locationMap[targetLocKey]; ok {
		return []browsing.Location{*loc}
	}

	return []browsing.Location{}
}

// getLocationKey creates a unique key for a file's location.
func getLocationKey(file analysis.FileInfo) string {
	hostname := file.Hostname
	if hostname == "" {
		hostname = "localhost"
	}
	return hostname + ":" + file.RemoteName
}

// getLocationName formats the location name from hostname and remote name.
func getLocationName(file analysis.FileInfo) string {
	hostname := file.Hostname
	if hostname == "" {
		// Local file
		return file.RemoteName
	}
	// Remote file - combine hostname and remote name
	return hostname + "/" + file.RemoteName
}

// inferLocationType determines if a location is "local" or "remote" based on the FileInfo.
// A location is considered local if the hostname is empty (local filesystem).
func inferLocationType(file analysis.FileInfo) string {
	if file.Hostname == "" {
		return "local"
	}
	return "remote"
}
