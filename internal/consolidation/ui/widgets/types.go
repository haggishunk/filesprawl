package widgets

import (
	"github.com/metaspartan/gotui/v5"
	"github.com/metaspartan/gotui/v5/widgets"
)

// BrowsingWidget is the interface for browsing views
type BrowsingWidget interface {
	// GetWidget returns the underlying gotui widget for rendering
	GetWidget() gotui.Drawable

	// UpdateData updates the widget with new browsing data
	UpdateData(data interface{}) error

	// HandleKeypress handles keyboard input
	HandleKeypress(key rune) error

	// GetSelectedIndex returns the currently selected item index
	GetSelectedIndex() int

	// SetSelectedIndex sets the selected item index
	SetSelectedIndex(index int) error

	// GetTitle returns the widget title
	GetTitle() string

	// SetRect sets the drawing rectangle
	SetRect(x, y, w, h int)
}

// ModalWidget is the interface for modal dialogs
type ModalWidget interface {
	// GetWidget returns the underlying gotui widget
	GetWidget() gotui.Drawable

	// Show displays the modal
	Show()

	// Hide hides the modal
	Hide()

	// IsVisible checks if modal is showing
	IsVisible() bool

	// SetContent sets the modal content text
	SetContent(text string) error

	// HandleKeypress handles keyboard input
	HandleKeypress(key rune) error

	// SetRect sets the drawing rectangle
	SetRect(x, y, w, h int)
}

// AppConfig holds the application configuration
type AppConfig struct {
	UserID        int
	Username      string
	HashGroupsCount int
	FilenameGroupsCount int
	FlaggedCount  int
	ShowHash      bool // true for hash view, false for filename view
}

// WidgetFactory creates gotui-based widgets for consolidation
type WidgetFactory struct {
	theme *gotui.Theme
}

// NewWidgetFactory creates a new widget factory
func NewWidgetFactory() *WidgetFactory {
	theme := &gotui.Theme{
		Block: gotui.NewStyle(),
	}
	return &WidgetFactory{
		theme: theme,
	}
}

// CreateTable creates a new gotui Table
func (wf *WidgetFactory) CreateTable(title string) *widgets.Table {
	table := widgets.NewTable()
	table.Title = title
	table.SetRect(0, 0, 80, 20)
	return table
}

// CreateList creates a new gotui List
func (wf *WidgetFactory) CreateList(title string) *widgets.List {
	list := widgets.NewList()
	list.Title = title
	list.SetRect(0, 0, 80, 20)
	return list
}

// CreateParagraph creates a new gotui Paragraph
func (wf *WidgetFactory) CreateParagraph(title string) *widgets.Paragraph {
	para := widgets.NewParagraph()
	para.Title = title
	para.SetRect(0, 0, 80, 10)
	return para
}

// CreateModal creates a new gotui Modal dialog
func (wf *WidgetFactory) CreateModal(title string, width, height int) *widgets.Modal {
	modal := widgets.NewModal()
	modal.Title = title
	modal.SetRect(0, 0, width, height)
	return modal
}

// CreateInput creates a new gotui Input widget
func (wf *WidgetFactory) CreateInput(label string) *widgets.Input {
	input := widgets.NewInput()
	input.Label = label
	input.SetRect(0, 0, 60, 3)
	return input
}

// CreateButton creates a new gotui Button
func (wf *WidgetFactory) CreateButton(label string) *widgets.Button {
	button := widgets.NewButton(label)
	button.SetRect(0, 0, 20, 1)
	return button
}
