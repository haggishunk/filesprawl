package browsing

import (
	"context"
	"testing"
	"time"
	"github.com/haggishunk/filesprawl/internal/consolidation/flags"
)

func TestSessionHandlerCreation(t *testing.T) {
	sm := NewSessionManager(1, nil)
	pfm := NewPersistentFlagManager(1, nil)
	sh := NewSessionHandler(sm, pfm)

	if sh.session != sm {
		t.Errorf("expected session to be set")
	}
	if sh.flagManager != pfm {
		t.Errorf("expected flag manager to be set")
	}
	if sh.summaryRenderer == nil {
		t.Errorf("expected summary renderer to be initialized")
	}
}

func TestHasUnsavedChanges(t *testing.T) {
	sm := NewSessionManager(1, nil)
	pfm := NewPersistentFlagManager(1, nil)
	sh := NewSessionHandler(sm, pfm)

	if sh.HasUnsavedChanges() {
		t.Errorf("expected no unsaved changes initially")
	}

	sh.MarkChanged()
	if !sh.HasUnsavedChanges() {
		t.Errorf("expected unsaved changes after marking changed")
	}
}

func TestShouldAutoSave(t *testing.T) {
	sm := NewSessionManager(1, nil)
	pfm := NewPersistentFlagManager(1, nil)
	sh := NewSessionHandler(sm, pfm)

	// Set short interval for testing
	sh.autoSaveInterval = 10 * time.Millisecond

	if sh.ShouldAutoSave() {
		t.Errorf("expected not to auto-save initially")
	}

	time.Sleep(20 * time.Millisecond)
	if !sh.ShouldAutoSave() {
		t.Errorf("expected to auto-save after interval")
	}
}

func TestSaveProgress(t *testing.T) {
	sm := NewSessionManager(1, nil)
	pfm := NewPersistentFlagManager(1, nil)
	sh := NewSessionHandler(sm, pfm)

	ctx := context.Background()
	err := sh.SaveProgress(ctx)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if sh.HasUnsavedChanges() {
		t.Errorf("expected unsaved changes to be cleared")
	}
}

func TestQuitWithoutConsolidation(t *testing.T) {
	sm := NewSessionManager(1, nil)
	pfm := NewPersistentFlagManager(1, nil)
	sh := NewSessionHandler(sm, pfm)

	// Add some flags
	pfm.GetOrCreateStateMachine(10).TransitionKeep()
	sm.flags[10] = &flags.Flag{FileID: 10, FlagType: flags.FlagTypeKeep}

	ctx := context.Background()
	output, err := sh.QuitWithoutConsolidation(ctx)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if output == "" {
		t.Errorf("expected non-empty quit message")
	}
	if !contains(output, "saved") {
		t.Errorf("expected mention of saved flags")
	}
}

func TestRenderQuitPrompt(t *testing.T) {
	sm := NewSessionManager(1, nil)
	pfm := NewPersistentFlagManager(1, nil)
	sh := NewSessionHandler(sm, pfm)

	prompt := sh.RenderQuitPrompt()

	if prompt == "" {
		t.Errorf("expected non-empty quit prompt")
	}
	if !contains(prompt, "QUIT") {
		t.Errorf("expected QUIT in prompt")
	}
}

func TestGetSessionSummary(t *testing.T) {
	sm := NewSessionManager(1, nil)
	pfm := NewPersistentFlagManager(1, nil)
	sh := NewSessionHandler(sm, pfm)

	// Add test data
	pfm.GetOrCreateStateMachine(1).TransitionKeep()
	sm.flags[1] = &flags.Flag{FileID: 1, FlagType: flags.FlagTypeKeep}

	summary := sh.GetSessionSummary()

	if summary == "" {
		t.Errorf("expected non-empty summary")
	}
	if !contains(summary, "SESSION SUMMARY") {
		t.Errorf("expected title in summary")
	}
	if !contains(summary, "KEEP") {
		t.Errorf("expected flag counts in summary")
	}
}

func TestGetResumeInstructions(t *testing.T) {
	sm := NewSessionManager(1, nil)
	pfm := NewPersistentFlagManager(1, nil)
	sh := NewSessionHandler(sm, pfm)

	instructions := sh.GetResumeInstructions()

	if instructions == "" {
		t.Errorf("expected non-empty instructions")
	}
	if !contains(instructions, "filesprawl consolidate") {
		t.Errorf("expected resume command in instructions")
	}
}
