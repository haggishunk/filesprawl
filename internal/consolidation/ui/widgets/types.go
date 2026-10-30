package widgets

// BrowsingWidget is the interface for browsing views
// NOTE: Phase 5 TUI integration pending. Stub implementation for now.
type BrowsingWidget interface {
	// GetWidget returns the underlying widget for rendering
	GetWidget() interface{}

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
// NOTE: Phase 5 TUI integration pending. Stub implementation for now.
type ModalWidget interface {
	// GetWidget returns the underlying widget
	GetWidget() interface{}

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
	UserID              int
	Username            string
	HashGroupsCount     int
	FilenameGroupsCount int
	FlaggedCount        int
	ShowHash            bool // true for hash view, false for filename view
}

// WidgetFactory creates widgets for consolidation
// NOTE: Phase 5 TUI integration pending. Stub implementation for now.
type WidgetFactory struct {
	// Stub for gotui theme
}

// NewWidgetFactory creates a new widget factory
func NewWidgetFactory() *WidgetFactory {
	return &WidgetFactory{}
}

// CreateTable creates a new Table widget (stub)
func (wf *WidgetFactory) CreateTable(title string) interface{} {
	// TODO: Implement full gotui Table widget
	return nil
}

// CreateList creates a new List widget (stub)
func (wf *WidgetFactory) CreateList(title string) interface{} {
	// TODO: Implement full gotui List widget
	return nil
}

// CreateParagraph creates a new Paragraph widget (stub)
func (wf *WidgetFactory) CreateParagraph(title string) interface{} {
	// TODO: Implement full gotui Paragraph widget
	return nil
}

// CreateModal creates a new Modal dialog (stub)
func (wf *WidgetFactory) CreateModal(title string, width, height int) interface{} {
	// TODO: Implement full gotui Modal widget
	return nil
}

// CreateInput creates a new Input widget (stub)
func (wf *WidgetFactory) CreateInput(label string) interface{} {
	// TODO: Implement full gotui Input widget
	return nil
}

// CreateButton creates a new Button widget (stub)
func (wf *WidgetFactory) CreateButton(label string) interface{} {
	// TODO: Implement full gotui Button widget
	return nil
}
