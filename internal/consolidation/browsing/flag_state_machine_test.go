package browsing

import (
	"testing"

	"github.com/haggishunk/filesprawl/internal/consolidation/flags"
)

// Placeholder: Full tests pending proper Phase 5 TUI implementation
// These tests verify the FlagStateMachine's core functionality
func TestFlagStateMachineCreation(t *testing.T) {
	fsm := NewFlagStateMachine(42)

	if fsm.currentState.FileID != 42 {
		t.Errorf("expected FileID 42, got %d", fsm.currentState.FileID)
	}
	if fsm.IsFlagged() {
		t.Errorf("expected not flagged initially")
	}
	if fsm.IsModified() {
		t.Errorf("expected not modified initially")
	}
}

func TestCanTransitionTo(t *testing.T) {
	fsm := NewFlagStateMachine(8)

	// Should allow all transitions for non-critical files
	if !fsm.CanTransitionTo(flags.FlagTypeKeep) {
		t.Errorf("expected KEEP transition to be allowed")
	}
	if !fsm.CanTransitionTo(flags.FlagTypeArchive) {
		t.Errorf("expected ARCHIVE transition to be allowed")
	}
	if !fsm.CanTransitionTo(flags.FlagTypeDelete) {
		t.Errorf("expected DELETE transition to be allowed")
	}
}

