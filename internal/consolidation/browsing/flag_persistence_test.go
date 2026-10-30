package browsing

import (
	"context"
	"testing"
	"github.com/haggishunk/filesprawl/internal/consolidation/flags"
)

func TestPersistentFlagManagerCreation(t *testing.T) {
	pfm := NewPersistentFlagManager(1, nil)

	if pfm.userID != 1 {
		t.Errorf("expected userID 1, got %d", pfm.userID)
	}
	if len(pfm.stateMachines) != 0 {
		t.Errorf("expected no state machines initially")
	}
}

func TestGetOrCreateStateMachine(t *testing.T) {
	pfm := NewPersistentFlagManager(1, nil)

	fsm1 := pfm.GetOrCreateStateMachine(100)
	if fsm1 == nil {
		t.Errorf("expected state machine to be created")
	}

	fsm2 := pfm.GetOrCreateStateMachine(100)
	if fsm1 != fsm2 {
		t.Errorf("expected same state machine on second call")
	}
}

func TestCountByType(t *testing.T) {
	pfm := NewPersistentFlagManager(1, nil)
	ctx := context.Background()

	fsm1 := pfm.GetOrCreateStateMachine(1)
	fsm1.TransitionKeep()

	fsm2 := pfm.GetOrCreateStateMachine(2)
	fsm2.TransitionKeep()

	fsm3 := pfm.GetOrCreateStateMachine(3)
	fsm3.TransitionDelete(nil)

	counts := pfm.CountByType()

	if counts[flags.FlagTypeKeep] != 2 {
		t.Errorf("expected 2 KEEP flags, got %d", counts[flags.FlagTypeKeep])
	}
	if counts[flags.FlagTypeDelete] != 1 {
		t.Errorf("expected 1 DELETE flag, got %d", counts[flags.FlagTypeDelete])
	}
}

func TestGetFlaggedFiles(t *testing.T) {
	pfm := NewPersistentFlagManager(1, nil)

	fsm := pfm.GetOrCreateStateMachine(99)
	fsm.TransitionKeep()

	flagged := pfm.GetFlaggedFiles()

	if len(flagged) != 1 {
		t.Errorf("expected 1 flagged file, got %d", len(flagged))
	}

	if flag, ok := flagged[99]; ok {
		if flag.FlagType != flags.FlagTypeKeep {
			t.Errorf("expected KEEP flag type")
		}
	} else {
		t.Errorf("expected file 99 in flagged files")
	}
}

func TestPersistKeepTransition(t *testing.T) {
	pfm := NewPersistentFlagManager(1, nil)
	ctx := context.Background()

	err := pfm.PersistKeep(ctx, 200)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	fsm := pfm.GetOrCreateStateMachine(200)
	if fsm.GetCurrentState().FlagType != flags.FlagTypeKeep {
		t.Errorf("expected KEEP flag after persist")
	}
}

func TestPersistArchiveTransition(t *testing.T) {
	pfm := NewPersistentFlagManager(1, nil)
	ctx := context.Background()

	err := pfm.PersistArchive(ctx, 201, "s3-backup", nil)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	fsm := pfm.GetOrCreateStateMachine(201)
	state := fsm.GetCurrentState()
	if state.FlagType != flags.FlagTypeArchive {
		t.Errorf("expected ARCHIVE flag")
	}
	if state.ArchiveDest == nil || *state.ArchiveDest != "s3-backup" {
		t.Errorf("expected archive destination to be set")
	}
}

func TestPersistDeleteTransition(t *testing.T) {
	pfm := NewPersistentFlagManager(1, nil)
	ctx := context.Background()

	reason := "duplicate"
	err := pfm.PersistDelete(ctx, 202, &reason)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	fsm := pfm.GetOrCreateStateMachine(202)
	if fsm.GetCurrentState().FlagType != flags.FlagTypeDelete {
		t.Errorf("expected DELETE flag")
	}
}

func TestPersistSkipTransition(t *testing.T) {
	pfm := NewPersistentFlagManager(1, nil)
	ctx := context.Background()

	fsm := pfm.GetOrCreateStateMachine(203)
	fsm.TransitionKeep()

	err := pfm.PersistSkip(ctx, 203)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if fsm.IsFlagged() {
		t.Errorf("expected file to be unflagged after skip")
	}
}
