package modals

import "fmt"

// QuitModal confirms saving and exiting
type QuitModal struct {
	title            string
	message          string
	keepCount        int
	archiveCount     int
	deleteCount      int
}

// NewQuitModal creates a quit confirmation modal
func NewQuitModal(keep, archive, delete int) *QuitModal {
	return &QuitModal{
		title:        "Save and Exit?",
		message:      "All flags have been saved to the database.",
		keepCount:    keep,
		archiveCount: archive,
		deleteCount:  delete,
	}
}

// Execute runs the quit modal
func (qm *QuitModal) Execute() ModalResult {
	fmt.Printf("\n%s\n%s\n", qm.title, qm.message)
	fmt.Printf("Summary:\n")
	fmt.Printf("  🟢 Keep:    %d files\n", qm.keepCount)
	fmt.Printf("  📦 Archive: %d files\n", qm.archiveCount)
	fmt.Printf("  🗑️  Delete:   %d files\n", qm.deleteCount)
	fmt.Printf("[Exit]/[Cancel]: ")
	
	var response string
	fmt.Scanln(&response)
	return ModalResult{Confirmed: response == "y" || response == "yes"}
}

// Cancel closes the modal
func (qm *QuitModal) Cancel() {
	// No-op for console implementation
}

// ErrorModal displays an error message
type ErrorModal struct {
	title    string
	message  string
}

// NewErrorModal creates an error display modal
func NewErrorModal(title, message string) *ErrorModal {
	return &ErrorModal{
		title:   title,
		message: message,
	}
}

// Execute shows the error and waits for acknowledgement
func (em *ErrorModal) Execute() ModalResult {
	fmt.Printf("\n❌ %s\n%s\n", em.title, em.message)
	fmt.Print("Press Enter to continue...")
	fmt.Scanln()
	return ModalResult{Confirmed: true}
}

// Cancel closes the modal
func (em *ErrorModal) Cancel() {
	// No-op for console implementation
}
