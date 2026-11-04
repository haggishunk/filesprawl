package browsing

import (
	"testing"

	"github.com/haggishunk/filesprawl/internal/consolidation/flags"
)

// TestBrowsingWorkflow simulates the complete browse session workflow
func TestBrowsingWorkflow(t *testing.T) {
	userID := 1

	// Create session manager
	sessionMgr := NewSessionManager(userID, nil)
	if sessionMgr == nil {
		t.Fatal("failed to create session manager")
	}

	// Create flag manager
	flagMgr := NewPersistentFlagManager(userID, nil)
	if flagMgr == nil {
		t.Fatal("failed to create flag manager")
	}

	// Create presence analyzer
	analyzer := NewPresenceAnalyzer()
	if analyzer == nil {
		t.Fatal("failed to create presence analyzer")
	}

	// Create sample groups
	groups := []DuplicateGroup{
		{
			Identifier: "hash1",
			Count:      2,
			Versions: []FileVersion{
				{FileID: 1},
				{FileID: 2},
			},
		},
	}

	// Load groups into session
	sessionMgr.LoadDuplicates(groups)

	// Verify session is set
	if sessionMgr.GetSession() == nil {
		t.Fatal("session not set")
	}

	// Verify groups loaded
	loaded := sessionMgr.GetSession()
	if len(loaded.DuplicateGroups) != 1 {
		t.Errorf("expected 1 group, got %d", len(loaded.DuplicateGroups))
	}
}

// TestFlagingWorkflow tests the flagging state transitions
func TestFlagingWorkflow(t *testing.T) {
	userID := 1
	fileID := 100

	// Create flag manager
	flagMgr := NewPersistentFlagManager(userID, nil)

	// Get or create state machine for file
	fsm := flagMgr.GetOrCreateStateMachine(fileID)
	if fsm == nil {
		t.Fatal("failed to create state machine")
	}

	// Test KEEP transition
	if !fsm.CanTransitionTo(flags.FlagTypeKeep) {
		t.Error("should allow KEEP transition from NONE")
	}

	_ = fsm.TransitionKeep()
	state := fsm.GetCurrentState()
	if state.FlagType != flags.FlagTypeKeep {
		t.Errorf("expected KEEP, got %s", state.FlagType)
	}

	// Test ARCHIVE transition from KEEP
	dest := "s3-backup"
	if !fsm.CanTransitionTo(flags.FlagTypeArchive) {
		t.Error("should allow ARCHIVE transition")
	}

	_ = fsm.TransitionArchive(&dest)
	state = fsm.GetCurrentState()
	if state.FlagType != flags.FlagTypeArchive {
		t.Errorf("expected ARCHIVE, got %s", state.FlagType)
	}
}

// TestPresenceAnalysis tests location indicator calculation
func TestPresenceAnalysis(t *testing.T) {
	analyzer := NewPresenceAnalyzer()

	// Create sample file version with locations
	version := FileVersion{
		FileID: 1,
		Locations: []Location{
			{ID: 1, Name: "loc1"},
			{ID: 2, Name: "loc2"},
		},
	}

	// Should have a presence indicator
	indicator := analyzer.AnalyzeFileVersion(version)
	if indicator == "" {
		t.Error("expected non-empty indicator")
	}
}

// TestViewTypeSwitching tests hash/filename view types
func TestViewTypeSwitching(t *testing.T) {
	if string(ViewTypeHash) != "HASH" {
		t.Errorf("expected HASH, got %s", ViewTypeHash)
	}
	if string(ViewTypeFilename) != "FILENAME" {
		t.Errorf("expected FILENAME, got %s", ViewTypeFilename)
	}
}
