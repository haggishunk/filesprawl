package widgets

import (
	"github.com/metaspartan/gotui/v5/widgets"
)

// BaseListWidget wraps a gotui List widget for common functionality
type BaseListWidget struct {
	list                *widgets.List
	title               string
	rows                []string
	x, y, width, height int // Rectangle dimensions
}

// NewBaseListWidget creates a new base list widget
func NewBaseListWidget(title string) *BaseListWidget {
	list := widgets.NewList()
	if list != nil {
		list.Title = title
	}

	return &BaseListWidget{
		list:   list,
		title:  title,
		rows:   []string{},
		x:      0,
		y:      0,
		width:  0,
		height: 0,
	}
}

// GetWidget returns the underlying gotui List widget
func (w *BaseListWidget) GetWidget() interface{} {
	if w == nil || w.list == nil {
		return nil
	}
	return w.list
}

// SetRows updates the list rows and syncs with gotui.List
func (w *BaseListWidget) SetRows(rows []string) error {
	if w == nil {
		return nil
	}
	w.rows = rows
	if w.list != nil {
		w.list.Rows = rows
	}
	return nil
}

// GetRows returns the current rows
func (w *BaseListWidget) GetRows() []string {
	if w == nil {
		return []string{}
	}
	return w.rows
}

// GetSelectedIndex returns the currently selected item index
func (w *BaseListWidget) GetSelectedIndex() int {
	if w == nil || w.list == nil {
		return 0
	}
	return w.list.SelectedRow
}

// SetSelectedIndex sets the selected item index
func (w *BaseListWidget) SetSelectedIndex(index int) error {
	if w == nil || w.list == nil {
		return nil
	}
	if index < 0 || index >= len(w.rows) {
		return nil // Silently ignore invalid indices
	}
	w.list.SelectedRow = index
	return nil
}

// GetTitle returns the widget title
func (w *BaseListWidget) GetTitle() string {
	if w == nil {
		return ""
	}
	return w.title
}

// SetTitle sets the widget title
func (w *BaseListWidget) SetTitle(title string) {
	if w == nil {
		return
	}
	w.title = title
	if w.list != nil {
		w.list.Title = title
	}
}

// SetRect sets the drawing rectangle
func (w *BaseListWidget) SetRect(x, y, width, height int) {
	if w == nil {
		return
	}
	w.x = x
	w.y = y
	w.width = width
	w.height = height
	if w.list != nil {
		w.list.SetRect(x, y, width, height)
	}
}

// GetRect returns the current drawing rectangle
func (w *BaseListWidget) GetRect() (int, int, int, int) {
	if w == nil {
		return 0, 0, 0, 0
	}
	return w.x, w.y, w.width, w.height
}

// ScrollUp scrolls up in the list
func (w *BaseListWidget) ScrollUp() {
	if w == nil || w.list == nil {
		return
	}
	if w.list.SelectedRow > 0 {
		w.list.SelectedRow--
	}
}

// ScrollDown scrolls down in the list
func (w *BaseListWidget) ScrollDown() {
	if w == nil || w.list == nil {
		return
	}
	if w.list.SelectedRow < len(w.rows)-1 {
		w.list.SelectedRow++
	}
}

// GetSelectedRow returns the selected row text
func (w *BaseListWidget) GetSelectedRow() string {
	if w == nil || w.list == nil {
		return ""
	}
	idx := w.list.SelectedRow
	if idx >= 0 && idx < len(w.rows) {
		return w.rows[idx]
	}
	return ""
}

// IsEmpty checks if the list is empty
func (w *BaseListWidget) IsEmpty() bool {
	if w == nil {
		return true
	}
	return len(w.rows) == 0
}

// Clear clears all rows
func (w *BaseListWidget) Clear() {
	if w == nil {
		return
	}
	w.rows = []string{}
	if w.list != nil {
		w.list.Rows = []string{}
		w.list.SelectedRow = 0
	}
}

// Count returns the number of rows
func (w *BaseListWidget) Count() int {
	if w == nil {
		return 0
	}
	return len(w.rows)
}

// SetBorderStyle sets border styling
func (w *BaseListWidget) SetBorderStyle(_ interface{}) {
	// gotui List handles border styling through its internal state
	// This is a placeholder for future enhancements
}

// SetTitleStyle sets title styling
func (w *BaseListWidget) SetTitleStyle(_ interface{}) {
	// gotui List handles title styling through its internal state
	// This is a placeholder for future enhancements
}
