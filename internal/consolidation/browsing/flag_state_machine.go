package browsing

import (
	"fmt"
	"github.com/haggishunk/filesprawl/internal/consolidation/flags"
)

// FlagStateMachine manages file flagging state transitions
type FlagStateMachine struct {
	currentState FlagState
}

// FlagState represents the current flagging state of a file
type FlagState struct {
	FileID       int
	FlagType     flags.FlagType
	Reason       *string
	ArchiveDest  *string
	IsModified   bool
}

// NewFlagStateMachine creates a new state machine for a file
func NewFlagStateMachine(fileID int) *FlagStateMachine {
	return &FlagStateMachine{
		currentState: FlagState{
			FileID:     fileID,
			FlagType:   "",
			IsModified: false,
		},
	}
}

// SetFlagFromExisting loads the state machine with an existing flag
func (fsm *FlagStateMachine) SetFlagFromExisting(flag *flags.Flag) {
	if flag != nil {
		fsm.currentState = FlagState{
			FileID:      flag.FileID,
			FlagType:    flag.FlagType,
			Reason:      flag.Reason,
			ArchiveDest: flag.ArchiveDest,
			IsModified:  false,
		}
	}
}

// TransitionKeep transitions the file to KEEP state
func (fsm *FlagStateMachine) TransitionKeep() error {
	if fsm.currentState.FlagType == flags.FlagTypeKeep {
		return nil // Already in this state
	}

	fsm.currentState.FlagType = flags.FlagTypeKeep
	fsm.currentState.Reason = nil
	fsm.currentState.ArchiveDest = nil
	fsm.currentState.IsModified = true

	return nil
}

// TransitionArchive transitions the file to ARCHIVE state
func (fsm *FlagStateMachine) TransitionArchive(destination *string) error {
	if destination == nil || *destination == "" {
		return fmt.Errorf("archive destination required")
	}

	if fsm.currentState.FlagType == flags.FlagTypeArchive &&
		fsm.currentState.ArchiveDest != nil &&
		*fsm.currentState.ArchiveDest == *destination {
		return nil // Already in this state with same destination
	}

	fsm.currentState.FlagType = flags.FlagTypeArchive
	fsm.currentState.ArchiveDest = destination
	fsm.currentState.Reason = nil
	fsm.currentState.IsModified = true

	return nil
}

// TransitionDelete transitions the file to DELETE state
func (fsm *FlagStateMachine) TransitionDelete(reason *string) error {
	// Red/critical files cannot be deleted
	if fsm.isCriticalState() {
		return fmt.Errorf("cannot delete file - critical backup status")
	}

	if fsm.currentState.FlagType == flags.FlagTypeDelete {
		return nil // Already in this state
	}

	fsm.currentState.FlagType = flags.FlagTypeDelete
	fsm.currentState.Reason = reason
	fsm.currentState.ArchiveDest = nil
	fsm.currentState.IsModified = true

	return nil
}

// TransitionSkip removes the flag (no action on file)
func (fsm *FlagStateMachine) TransitionSkip() error {
	fsm.currentState.FlagType = ""
	fsm.currentState.Reason = nil
	fsm.currentState.ArchiveDest = nil
	fsm.currentState.IsModified = true

	return nil
}

// GetCurrentState returns the current flag state
func (fsm *FlagStateMachine) GetCurrentState() FlagState {
	return fsm.currentState
}

// IsFlagged returns true if file has any flag
func (fsm *FlagStateMachine) IsFlagged() bool {
	return fsm.currentState.FlagType != ""
}

// IsModified returns true if state has changed
func (fsm *FlagStateMachine) IsModified() bool {
	return fsm.currentState.IsModified
}

// CanTransitionTo checks if transition is allowed
func (fsm *FlagStateMachine) CanTransitionTo(targetType flags.FlagType) bool {
	if targetType == flags.FlagTypeDelete && fsm.isCriticalState() {
		return false
	}
	return true
}

// isCriticalState returns true if file is in critical backup status
func (fsm *FlagStateMachine) isCriticalState() bool {
	// This would check location presence indicators
	// For now, simple placeholder
	return false
}

// GetTransitionDescription returns human-readable description of current state
func (fsm *FlagStateMachine) GetTransitionDescription() string {
	if !fsm.IsFlagged() {
		return "No action flagged"
	}

	switch fsm.currentState.FlagType {
	case flags.FlagTypeKeep:
		return "Flagged to KEEP"
	case flags.FlagTypeArchive:
		dest := "unknown"
		if fsm.currentState.ArchiveDest != nil {
			dest = *fsm.currentState.ArchiveDest
		}
		return fmt.Sprintf("Flagged to ARCHIVE → %s", dest)
	case flags.FlagTypeDelete:
		return "Flagged to DELETE"
	default:
		return "Unknown state"
	}
}
    
