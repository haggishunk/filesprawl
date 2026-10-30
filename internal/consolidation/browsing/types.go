package browsing

import (
	"github.com/haggishunk/filesprawl/internal/consolidation/flags"
	"github.com/haggishunk/filesprawl/internal/object"
)

// ViewType represents the current browsing view mode
type ViewType string

const (
	ViewTypeHash     ViewType = "HASH"
	ViewTypeFilename ViewType = "FILENAME"
)

// PresenceIndicator represents the backup replication status of a file/location
type PresenceIndicator string

const (
	// Green: Multiple locations, safe to delete
	PresenceGreen PresenceIndicator = "🟢"
	// Yellow: Limited replication, recommend backup
	PresenceYellow PresenceIndicator = "🟡"
	// Red: Only copy exists, critical - no backup
	PresenceRed PresenceIndicator = "🔴"
)

// Location represents a file location with presence indicator
type Location struct {
	ID        int
	Name      string
	Type      string // "local" or "remote"
	FileCount int
	Indicator PresenceIndicator
}

// FileVersion represents a file with a specific hash
type FileVersion struct {
	FileID   int
	FileInfo *object.Meta
	Hash     *object.Hash
	Locations []Location
	Flag     *flags.Flag
}

// DuplicateGroup represents a group of duplicate files (by hash or filename)
type DuplicateGroup struct {
	Identifier string         // Hash value or filename
	Count      int            // Total copies or versions
	Versions   []FileVersion  // Individual files/variants
	Flag       *flags.Flag    // Group-level flag (if set)
}

// BrowsingSession represents the current browsing state
type BrowsingSession struct {
	UserID       int
	CurrentView  ViewType
	DuplicateGroups []DuplicateGroup
	ExpandedItems map[string]bool // Track which items are expanded
	Flags        map[int]*flags.Flag // fileID -> Flag
	Cursor       int // Current cursor position
}

// ViewData represents data needed for rendering a view
type ViewData struct {
	Groups     []DuplicateGroup
	TotalFiles int
	Flagged    int
	View       ViewType
}

// SummaryData represents consolidation decision summary
type SummaryData struct {
	Keep    []FileVersion
	Archive []FileVersion
	Delete  []FileVersion
	Total   int
}

// Renderer interface for rendering different view types
type Renderer interface {
	// RenderHashView renders the hash-centric view
	RenderHashView(session *BrowsingSession) string

	// RenderFilenameView renders the filename-centric view
	RenderFilenameView(session *BrowsingSession) string

	// RenderSummary renders the decision summary
	RenderSummary(data *SummaryData) string
}

// FlagManager interface for managing flag state
type FlagManager interface {
	// FlagFile sets a flag on a file
	FlagFile(fileID int, flagType flags.FlagType, reason *string) error

	// UnflagFile removes a flag
	UnflagFile(fileID int) error

	// GetFlags returns all flagged files
	GetFlags() map[int]*flags.Flag

	// CountByType returns count of flags by type
	CountByType() map[flags.FlagType]int
}
