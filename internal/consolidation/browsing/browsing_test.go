package browsing

import (
	"testing"
	"github.com/haggishunk/filesprawl/internal/consolidation/flags"
)

func TestViewType(t *testing.T) {
	tests := []struct {
		name string
		view ViewType
		want string
	}{
		{"HASH view", ViewTypeHash, "HASH"},
		{"FILENAME view", ViewTypeFilename, "FILENAME"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.view) != tt.want {
				t.Errorf("ViewType = %s, want %s", tt.view, tt.want)
			}
		})
	}
}

func TestPresenceIndicator(t *testing.T) {
	tests := []struct {
		name      string
		indicator PresenceIndicator
	}{
		{"Green safe", PresenceGreen},
		{"Yellow limited", PresenceYellow},
		{"Red critical", PresenceRed},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.indicator) == "" {
				t.Errorf("PresenceIndicator is empty")
			}
		})
	}
}

func TestSessionManagerCreation(t *testing.T) {
	sm := NewSessionManager(1, nil)

	if sm.userID != 1 {
		t.Errorf("expected userID 1, got %d", sm.userID)
	}
	if sm.currentView != ViewTypeHash {
		t.Errorf("expected default view HASH, got %s", sm.currentView)
	}
	if len(sm.flags) != 0 {
		t.Errorf("expected no flags, got %d", len(sm.flags))
	}
}

func TestSwitchView(t *testing.T) {
	sm := NewSessionManager(1, nil)

	if sm.GetCurrentView() != ViewTypeHash {
		t.Errorf("expected initial view HASH")
	}

	sm.SwitchView()
	if sm.GetCurrentView() != ViewTypeFilename {
		t.Errorf("expected view FILENAME after switch")
	}

	sm.SwitchView()
	if sm.GetCurrentView() != ViewTypeHash {
		t.Errorf("expected view HASH after second switch")
	}
}

func TestExpandItem(t *testing.T) {
	sm := NewSessionManager(1, nil)
	itemID := "hash_abc123"

	if sm.IsExpanded(itemID) {
		t.Errorf("expected item not expanded initially")
	}

	sm.ExpandItem(itemID)
	if !sm.IsExpanded(itemID) {
		t.Errorf("expected item expanded after toggle")
	}

	sm.ExpandItem(itemID)
	if sm.IsExpanded(itemID) {
		t.Errorf("expected item collapsed after second toggle")
	}
}

func TestSessionManagerCountByType(t *testing.T) {
	sm := NewSessionManager(1, nil)

	// Add some mock flags
	sm.flags[1] = &flags.Flag{ID: 1, FlagType: flags.FlagTypeKeep}
	sm.flags[2] = &flags.Flag{ID: 2, FlagType: flags.FlagTypeArchive}
	sm.flags[3] = &flags.Flag{ID: 3, FlagType: flags.FlagTypeArchive}

	counts := sm.CountByType()

	if counts[flags.FlagTypeKeep] != 1 {
		t.Errorf("expected 1 KEEP flag, got %d", counts[flags.FlagTypeKeep])
	}
	if counts[flags.FlagTypeArchive] != 2 {
		t.Errorf("expected 2 ARCHIVE flags, got %d", counts[flags.FlagTypeArchive])
	}
}

func TestDuplicateGroup(t *testing.T) {
	group := DuplicateGroup{
		Identifier: "abc123",
		Count:      5,
		Versions:   []FileVersion{},
		Flag:       nil,
	}

	if group.Identifier != "abc123" {
		t.Errorf("expected identifier abc123")
	}
	if group.Count != 5 {
		t.Errorf("expected count 5, got %d", group.Count)
	}
}
