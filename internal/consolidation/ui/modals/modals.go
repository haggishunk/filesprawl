package modals

import "fmt"

// ModalResult holds the result of a modal dialog
type ModalResult struct {
	Confirmed bool   // User clicked OK/Yes
	Value     string // Text input value (for archive destination, etc)
	Reason    string // Reason text (for delete, etc)
	Error     error  // Any error during modal execution
}

// Modal is an interface for dialog implementations
type Modal interface {
	// Execute runs the modal and returns result
	Execute() ModalResult
	
	// Cancel closes the modal without action
	Cancel()
}

// KeepModal confirms keeping the current file
type KeepModal struct {
	title   string
	message string
}

// NewKeepModal creates a keep confirmation modal
func NewKeepModal() *KeepModal {
	return &KeepModal{
		title:   "Flag as KEEP?",
		message: "This file will be preserved.",
	}
}

// Execute runs the keep modal (TODO: Phase 4 gotui integration)
func (km *KeepModal) Execute() ModalResult {
	fmt.Printf("\n%s\n%s\n[Confirm]/[Cancel]: ", km.title, km.message)
	var response string
	fmt.Scanln(&response)
	return ModalResult{Confirmed: response == "y" || response == "yes"}
}

// Cancel closes the modal
func (km *KeepModal) Cancel() {
	// No-op for console implementation
}

// ArchiveModal gets archive destination
type ArchiveModal struct {
	title       string
	message     string
	destination string
}

// NewArchiveModal creates an archive destination modal
func NewArchiveModal() *ArchiveModal {
	return &ArchiveModal{
		title:   "Archive File",
		message: "Enter destination path:",
	}
}

// Execute runs the archive modal
func (am *ArchiveModal) Execute() ModalResult {
	fmt.Printf("\n%s\n%s\n", am.title, am.message)
	var dest string
	fmt.Scanln(&dest)
	if dest == "" {
		return ModalResult{Confirmed: false, Error: fmt.Errorf("destination required")}
	}
	return ModalResult{Confirmed: true, Value: dest}
}

// Cancel closes the modal
func (am *ArchiveModal) Cancel() {
	// No-op for console implementation
}

// DeleteModal confirms file deletion
type DeleteModal struct {
	title       string
	isCritical  bool
	message     string
}

// NewDeleteModal creates a delete confirmation modal
func NewDeleteModal(isCritical bool) *DeleteModal {
	msg := "Delete this file?"
	if isCritical {
		msg = "⚠️  WARNING: Only 1 copy exists!\nDelete permanently?"
	}
	return &DeleteModal{
		title:      "Delete File",
		isCritical: isCritical,
		message:    msg,
	}
}

// Execute runs the delete modal
func (dm *DeleteModal) Execute() ModalResult {
	fmt.Printf("\n%s\n%s\n[Confirm]/[Cancel]: ", dm.title, dm.message)
	var response string
	fmt.Scanln(&response)
	return ModalResult{Confirmed: response == "y" || response == "yes"}
}

// Cancel closes the modal
func (dm *DeleteModal) Cancel() {
	// No-op for console implementation
}

// SkipModal confirms removing flag
type SkipModal struct {
	title   string
	message string
}

// NewSkipModal creates a skip confirmation modal
func NewSkipModal() *SkipModal {
	return &SkipModal{
		title:   "Remove Flag?",
		message: "This file will not be flagged for any action.",
	}
}

// Execute runs the skip modal
func (sm *SkipModal) Execute() ModalResult {
	fmt.Printf("\n%s\n%s\n[Confirm]/[Cancel]: ", sm.title, sm.message)
	var response string
	fmt.Scanln(&response)
	return ModalResult{Confirmed: response == "y" || response == "yes"}
}

// Cancel closes the modal
func (sm *SkipModal) Cancel() {
	// No-op for console implementation
}
