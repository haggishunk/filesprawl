package widgets

import (
	"github.com/metaspartan/gotui/v5"
	"github.com/metaspartan/gotui/v5/widgets"
)

// BaseListWidget wraps a gotui List for common functionality
type BaseListWidget struct {
	list      *widgets.List
	selectedIndex int
	title     string
	rows      []string
}

// NewBaseListWidget creates a new base list widget
func NewBaseListWidget(title string) *BaseListWidget {
	list := widgets.NewList()
	list.Title = title
	list.SetRect(0, 0, 80, 20)

	return &BaseListWidget{
		list:      list,
		title:     title,
		selectedIndex: 0,
		rows:      []string{},
	}
}

// GetWidget returns the underlying gotui widget
func (w *BaseListWidget) GetWidget() gotui.Drawable {
	return w.list
}

// SetRows updates the list rows
func (w *BaseListWidget) SetRows(rows []string) error {
	w.rows = rows
	w.list.Rows = rows
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
	w.list.SelectedRow = index
	return nil
}

// GetTitle returns the widget title
func (w *BaseListWidget) GetTitle() string {
	return w.title
}

// SetTitle sets the widget title
func (w *BaseListWidget) SetTitle(title string) {
	w.title = title
	w.list.Title = title
}

// SetRect sets the drawing rectangle
func (w *BaseListWidget) SetRect(x, y, width, height int) {
	w.list.SetRect(x, y, width, height)
}

// GetRect returns the current drawing rectangle
func (w *BaseListWidget) GetRect() (int, int, int, int) {
	return w.list.GetRect()
}

// ScrollUp scrolls up in the list
func (w *BaseListWidget) ScrollUp() {
	if w.selectedIndex > 0 {
		w.selectedIndex--
		w.list.SelectedRow = w.selectedIndex
	}
}

// ScrollDown scrolls down in the list
func (w *BaseListWidget) ScrollDown() {
	if w.selectedIndex < len(w.rows)-1 {
		w.selectedIndex++
		w.list.SelectedRow = w.selectedIndex
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
	w.list.Rows = []string{}
	w.selectedIndex = 0
}

// Count returns the number of rows
func (w *BaseListWidget) Count() int {
	return len(w.rows)
}

// SetBorderStyle sets border styling
func (w *BaseListWidget) SetBorderStyle(fg gotui.Color) {
	w.list.BorderStyle = gotui.NewStyle().Foreground(fg)
}

// SetTitleStyle sets title styling
func (w *BaseListWidget) SetTitleStyle(fg gotui.Color) {
	w.list.TitleStyle = gotui.NewStyle().Foreground(fg)
}
