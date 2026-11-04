package ui

import tea "github.com/charmbracelet/bubbletea"

// Message types for Bubble Tea model updates

// KeyMsg represents a keyboard press
type KeyMsg struct {
	Key rune
}

// ScrollMsg indicates list scrolling
type ScrollMsg struct {
	Direction string // "up" or "down"
}

// SelectionMsg indicates a new item was selected
type SelectionMsg struct {
	Index int
}

// ViewSwitchMsg indicates switch between hash and filename views
type ViewSwitchMsg struct{}

// ActionMsg represents a user action (keep, archive, delete, skip)
type ActionMsg struct {
	Action  string // "keep", "archive", "delete", "skip"
	FileID  int    // ID of the file being acted upon
	GroupID string // ID of the group containing the file
}

// ModalShowMsg shows a modal dialog
type ModalShowMsg struct {
	ModalType string      // "keep", "archive", "delete", "skip", "help", "quit"
	Data      interface{} // Optional data for the modal
}

// ModalDismissMsg hides the current modal
type ModalDismissMsg struct{}

// ModalConfirmMsg indicates modal was confirmed
type ModalConfirmMsg struct {
	Value string // For modals that return values (e.g., archive destination)
}

// UpdateDetailsMsg triggers update of details pane
type UpdateDetailsMsg struct {
	Index int
}

// QuitMsg indicates user wants to exit
type QuitMsg struct{}

// ErrorMsg represents an error
type ErrorMsg struct {
	Err error
}

// Cmd helpers
func KeyCommand(key rune) tea.Cmd {
	return func() tea.Msg {
		return KeyMsg{Key: key}
	}
}

func ScrollCommand(direction string) tea.Cmd {
	return func() tea.Msg {
		return ScrollMsg{Direction: direction}
	}
}

func SelectCommand(index int) tea.Cmd {
	return func() tea.Msg {
		return SelectionMsg{Index: index}
	}
}
