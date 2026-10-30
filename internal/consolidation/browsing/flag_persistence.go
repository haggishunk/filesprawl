package browsing

import (
	"context"
	"fmt"
	"github.com/haggishunk/filesprawl/internal/consolidation/flags"
)

// PersistentFlagManager integrates state machine with database persistence
type PersistentFlagManager struct {
	stateMachines map[int]*FlagStateMachine
	flagRepo      flags.Repository
	userID        int
}

// NewPersistentFlagManager creates a new persistent flag manager
func NewPersistentFlagManager(userID int, flagRepo flags.Repository) *PersistentFlagManager {
	return &PersistentFlagManager{
		stateMachines: make(map[int]*FlagStateMachine),
		flagRepo:      flagRepo,
		userID:        userID,
	}
}

// GetOrCreateStateMachine gets or creates a state machine for a file
func (pfm *PersistentFlagManager) GetOrCreateStateMachine(fileID int) *FlagStateMachine {
	if fsm, ok := pfm.stateMachines[fileID]; ok {
		return fsm
	}

	fsm := NewFlagStateMachine(fileID)
	pfm.stateMachines[fileID] = fsm
	return fsm
}

// LoadExistingFlags loads all existing flags from database
func (pfm *PersistentFlagManager) LoadExistingFlags(ctx context.Context) error {
	if pfm.flagRepo == nil {
		return nil
	}

	userFlags, err := pfm.flagRepo.GetFlagsByUser(ctx, pfm.userID)
	if err != nil {
		return fmt.Errorf("failed to load existing flags: %w", err)
	}

	for _, flag := range userFlags {
		flagCopy := flag
		fsm := pfm.GetOrCreateStateMachine(flag.FileID)
		fsm.SetFlagFromExisting(&flagCopy)
	}

	return nil
}

// PersistKeep flags file as KEEP and writes to database
func (pfm *PersistentFlagManager) PersistKeep(ctx context.Context, fileID int) error {
	fsm := pfm.GetOrCreateStateMachine(fileID)

	if err := fsm.TransitionKeep(); err != nil {
		return err
	}

	return pfm.persistStateToDatabase(ctx, fsm)
}

// PersistArchive flags file as ARCHIVE with destination
func (pfm *PersistentFlagManager) PersistArchive(ctx context.Context, fileID int, destination string, reason *string) error {
	fsm := pfm.GetOrCreateStateMachine(fileID)

	if err := fsm.TransitionArchive(&destination); err != nil {
		return err
	}

	return pfm.persistStateToDatabase(ctx, fsm)
}

// PersistDelete flags file as DELETE with reason
func (pfm *PersistentFlagManager) PersistDelete(ctx context.Context, fileID int, reason *string) error {
	fsm := pfm.GetOrCreateStateMachine(fileID)

	if err := fsm.TransitionDelete(reason); err != nil {
		return err
	}

	return pfm.persistStateToDatabase(ctx, fsm)
}

// PersistSkip removes the flag and deletes from database
func (pfm *PersistentFlagManager) PersistSkip(ctx context.Context, fileID int) error {
	fsm := pfm.GetOrCreateStateMachine(fileID)

	if err := fsm.TransitionSkip(); err != nil {
		return err
	}

	if pfm.flagRepo != nil {
		dbFlags, err := pfm.flagRepo.GetFlagsByFile(ctx, fileID)
		if err == nil && len(dbFlags) > 0 {
			if err := pfm.flagRepo.DeleteFlag(ctx, dbFlags[0].ID); err != nil {
				return fmt.Errorf("failed to delete flag: %w", err)
			}
		}
	}

	return nil
}

// persistStateToDatabase writes state machine state to database
func (pfm *PersistentFlagManager) persistStateToDatabase(ctx context.Context, fsm *FlagStateMachine) error {
	if pfm.flagRepo == nil || !fsm.IsModified() {
		return nil
	}

	state := fsm.GetCurrentState()
	existingFlags, err := pfm.flagRepo.GetFlagsByFile(ctx, state.FileID)
	if err != nil {
		return fmt.Errorf("failed to check existing flag: %w", err)
	}

	if len(existingFlags) > 0 {
		_, err := pfm.flagRepo.UpdateFlag(ctx, existingFlags[0].ID, state.FlagType, state.Reason, state.ArchiveDest)
		if err != nil {
			return fmt.Errorf("failed to update flag: %w", err)
		}
	} else if state.FlagType != "" {
		_, err := pfm.flagRepo.CreateFlag(ctx, state.FileID, pfm.userID, state.FlagType, state.Reason, state.ArchiveDest)
		if err != nil {
			return fmt.Errorf("failed to create flag: %w", err)
		}
	}

	return nil
}

// GetAllStateMachines returns all state machines
func (pfm *PersistentFlagManager) GetAllStateMachines() map[int]*FlagStateMachine {
	return pfm.stateMachines
}

// GetFlaggedFiles returns map of fileID -> Flag for all flagged files
func (pfm *PersistentFlagManager) GetFlaggedFiles() map[int]*flags.Flag {
	result := make(map[int]*flags.Flag)

	for fileID, fsm := range pfm.stateMachines {
		if fsm.IsFlagged() {
			state := fsm.GetCurrentState()
			flag := &flags.Flag{
				FileID:      state.FileID,
				UserID:      pfm.userID,
				FlagType:    state.FlagType,
				Reason:      state.Reason,
				ArchiveDest: state.ArchiveDest,
			}
			result[fileID] = flag
		}
	}

	return result
}

// CountByType returns counts of flags grouped by type
func (pfm *PersistentFlagManager) CountByType() map[flags.FlagType]int {
	counts := make(map[flags.FlagType]int)

	for _, fsm := range pfm.stateMachines {
		if fsm.IsFlagged() {
			counts[fsm.GetCurrentState().FlagType]++
		}
	}

	return counts
}
