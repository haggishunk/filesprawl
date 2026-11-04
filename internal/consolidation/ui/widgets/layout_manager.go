package widgets

// LayoutManager calculates widget rectangles for 4-pane TUI layout
// Left: hash list (40%)
// Right upper: hash-grouped files (30%)
// Right lower: filename-grouped files (30%) - toggleable
// Bottom: status bar (3 lines)
type LayoutManager struct {
	maxX          int
	maxY          int
	leftPaneW     int // ~40% for list
	bannerH       int // 1 line for banner
	statusH       int // 3 lines for status
	rightPaneW    int // ~60% width for right side
	dividerX      int // right border of left pane
	hashPaneX     int // start of hash pane (right side upper)
	hashPaneH     int // height of hash pane
	filenamePaneX int // start of filename pane (right side lower, same as hash)
	filenamePaneY int // y position of filename pane
	filenamePaneH int // height of filename pane
}

// NewLayoutManager creates a layout manager for terminal dimensions
func NewLayoutManager(maxX, maxY int) *LayoutManager {
	if maxX < 80 || maxY < 24 {
		// Minimum terminal size: 80x24
		maxX, maxY = 80, 24
	}

	lm := &LayoutManager{
		maxX:    maxX,
		maxY:    maxY,
		bannerH: 1, // Banner takes 1 line
		statusH: 3, // Status bar height
	}

	// Calculate left pane width (40%)
	lm.leftPaneW = maxX * 40 / 100
	if lm.leftPaneW < 20 {
		lm.leftPaneW = 20
	}
	if lm.leftPaneW > maxX-20 {
		lm.leftPaneW = maxX - 20
	}

	// Calculate right pane dimensions
	lm.dividerX = lm.leftPaneW
	lm.hashPaneX = lm.dividerX + 1
	lm.rightPaneW = maxX - lm.hashPaneX

	// Calculate vertical split of right pane (50/50 or adjust based on height)
	contentH := maxY - lm.bannerH - lm.statusH - 1 // -1 for divider line
	if contentH < 6 {
		contentH = 6 // Minimum height
	}
	lm.hashPaneH = contentH / 2
	lm.filenamePaneH = contentH - lm.hashPaneH - 1 // -1 for divider line
	lm.filenamePaneY = lm.bannerH + lm.hashPaneH + 1

	return lm
}

// GetBannerRect returns coordinates for banner (top of screen)
// Format: (x, y, width, height)
func (lm *LayoutManager) GetBannerRect() (int, int, int, int) {
	return 0, 0, lm.maxX, lm.bannerH
}

// GetListPaneRect returns coordinates for left pane (list widget)
// Format: (x, y, width, height)
func (lm *LayoutManager) GetListPaneRect() (int, int, int, int) {
	contentH := lm.maxY - lm.bannerH - lm.statusH - 1
	return 0, lm.bannerH, lm.leftPaneW, contentH
}

// GetDetailsPaneRect returns coordinates for right pane (details widget) - DEPRECATED
// Use GetHashTablePaneRect for upper right and GetFilenameTablePaneRect for lower right
func (lm *LayoutManager) GetDetailsPaneRect() (int, int, int, int) {
	contentH := lm.maxY - lm.bannerH - lm.statusH - 1
	return lm.hashPaneX, lm.bannerH, lm.rightPaneW, contentH
}

// GetHashTablePaneRect returns coordinates for upper right pane (hash-grouped files)
func (lm *LayoutManager) GetHashTablePaneRect() (int, int, int, int) {
	return lm.hashPaneX, lm.bannerH, lm.rightPaneW, lm.hashPaneH
}

// GetFilenameTablePaneRect returns coordinates for lower right pane (filename-grouped files)
func (lm *LayoutManager) GetFilenameTablePaneRect() (int, int, int, int) {
	return lm.filenamePaneX, lm.filenamePaneY, lm.rightPaneW, lm.filenamePaneH
}

// GetStatusBarRect returns coordinates for status bar
func (lm *LayoutManager) GetStatusBarRect() (int, int, int, int) {
	statusY := lm.maxY - lm.statusH
	return 0, statusY, lm.maxX, lm.statusH
}
