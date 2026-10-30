package widgets

// BaseListWidget wraps a list widget for common functionality
// NOTE: Phase 5 TUI integration pending. Stub implementation for now.
type BaseListWidget struct {
	selectedIndex int
	title         string
	rows          []string
	x, y, width, height int // Rectangle dimensions
}

// NewBaseListWidget creates a new base list widget
func NewBaseListWidget(title string) *BaseListWidget {
	return &BaseListWidget{
		title:         title,
		selectedIndex: 0,
		rows:          []string{},
		x:             0,
		y:             0,
		width:         0,
		height:        0,
	}
}

// GetWidget returns the underlying widget
func (w *BaseListWidget) GetWidget() interface{} {
	// Return a non-nil placeholder for testing
	// In Phase 5, this would return the actual gotui List widget
	return w
}

// SetRows updates the list rows
func (w *BaseListWidget) SetRows(rows []string) error {
	w.rows = rows
	return nil
}

// GetRows returns the current rows
func (w *BaseListWidget) GetRows() []string {
	return w.rows
}

// GetSelectedIndex returns the currently selected item index
func (w *BaseListWidget) GetSelectedIndex() int {
	return w.selectedIndex
}

// SetSelectedIndex sets the selected item index
func (w *BaseListWidget) SetSelectedIndex(index int) error {
	if index < 0 || index >= len(w.rows) {
		return nil // Silently ignore invalid indices
	}
	w.selectedIndex = index
	return nil
}

// GetTitle returns the widget title
func (w *BaseListWidget) GetTitle() string {
	return w.title
}

// SetTitle sets the widget title
func (w *BaseListWidget) SetTitle(title string) {
	w.title = title
}

// SetRect sets the drawing rectangle
func (w *BaseListWidget) SetRect(x, y, width, height int) {
	w.x = x
	w.y = y
	w.width = width
	w.height = height
}

// GetRect returns the current drawing rectangle
func (w *BaseListWidget) GetRect() (int, int, int, int) {
	return w.x, w.y, w.width, w.height
}

// ScrollUp scrolls up in the list
func (w *BaseListWidget) ScrollUp() {
	if w.selectedIndex > 0 {
		w.selectedIndex--
	}
}

// ScrollDown scrolls down in the list
func (w *BaseListWidget) ScrollDown() {
	if w.selectedIndex < len(w.rows)-1 {
		w.selectedIndex++
	}
}

// GetSelectedRow returns the selected row text
func (w *BaseListWidget) GetSelectedRow() string {
	if w.selectedIndex >= 0 && w.selectedIndex < len(w.rows) {
		return w.rows[w.selectedIndex]
	}
	return ""
}

// IsEmpty checks if the list is empty
func (w *BaseListWidget) IsEmpty() bool {
	return len(w.rows) == 0
}

// Clear clears all rows
func (w *BaseListWidget) Clear() {
	w.rows = []string{}
	w.selectedIndex = 0
}

// Count returns the number of rows
func (w *BaseListWidget) Count() int {
	return len(w.rows)
}

// SetBorderStyle sets border styling (stub)
func (w *BaseListWidget) SetBorderStyle(_ interface{}) {
	// TODO: Implement full gotui style support
}

// SetTitleStyle sets title styling (stub)
func (w *BaseListWidget) SetTitleStyle(_ interface{}) {
	// TODO: Implement full gotui style support
}
