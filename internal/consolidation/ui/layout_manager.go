package ui

// LayoutManager calculates widget rectangles for 3-pane TUI layout
type LayoutManager struct {
	maxX         int
	maxY         int
	leftPaneW    int // ~40% for list
	statusH      int // 3 lines for status
	dividerX     int // right border of left pane
	detailsPaneX int // start of right pane
}

// NewLayoutManager creates a layout manager for terminal dimensions
func NewLayoutManager(maxX, maxY int) *LayoutManager {
	if maxX < 80 || maxY < 24 {
		// Minimum terminal size: 80x24
		maxX, maxY = 80, 24
	}

	lm := &LayoutManager{
		maxX:      maxX,
		maxY:      maxY,
		leftPaneW: maxX * 40 / 100, // 40% width
		statusH:   3,               // Status bar height
	}

	// Adjust for odd terminal widths
	if lm.leftPaneW < 20 {
		lm.leftPaneW = 20
	}
	if lm.leftPaneW > maxX-20 {
		lm.leftPaneW = maxX - 20
	}

	lm.dividerX = lm.leftPaneW
	lm.detailsPaneX = lm.dividerX + 1

	return lm
}

// GetListPaneRect returns coordinates for left pane (list widget)
// Format: (x, y, width, height)
func (lm *LayoutManager) GetListPaneRect() (int, int, int, int) {
	contentH := lm.maxY - lm.statusH - 1
	return 0, 0, lm.leftPaneW, contentH
}

// GetDetailsPaneRect returns coordinates for right pane (details widget)
func (lm *LayoutManager) GetDetailsPaneRect() (int, int, int, int) {
	contentH := lm.maxY - lm.statusH - 1
	detailsW := lm.maxX - lm.detailsPaneX
	return lm.detailsPaneX, 0, detailsW, contentH
}

// GetStatusBarRect returns coordinates for bottom status bar
func (lm *LayoutManager) GetStatusBarRect() (int, int, int, int) {
	statusY := lm.maxY - lm.statusH
	return 0, statusY, lm.maxX, lm.statusH
}

// GetTerminalSize returns current terminal dimensions
func (lm *LayoutManager) GetTerminalSize() (int, int) {
	return lm.maxX, lm.maxY
}

// UpdateSize recalculates layout for new terminal dimensions
func (lm *LayoutManager) UpdateSize(maxX, maxY int) {
	if maxX < 80 || maxY < 24 {
		maxX, maxY = 80, 24
	}
	lm.maxX = maxX
	lm.maxY = maxY
	lm.leftPaneW = maxX * 40 / 100
	if lm.leftPaneW < 20 {
		lm.leftPaneW = 20
	}
	if lm.leftPaneW > maxX-20 {
		lm.leftPaneW = maxX - 20
	}
	lm.dividerX = lm.leftPaneW
	lm.detailsPaneX = lm.dividerX + 1
}
