package browsing

import (
	"testing"
)

func TestViewControllerCreation(t *testing.T) {
	sm := NewSessionManager(1, nil)
	vc := NewViewController(sm)

	if vc.session != sm {
		t.Errorf("expected session to be set")
	}
	if vc.hashRenderer == nil {
		t.Errorf("expected hash renderer to be initialized")
	}
	if vc.filenameRenderer == nil {
		t.Errorf("expected filename renderer to be initialized")
	}
}

func TestGetViewName(t *testing.T) {
	sm := NewSessionManager(1, nil)
	vc := NewViewController(sm)

	name := vc.GetViewName()
	if name != "Hash-Centric (Content)" {
		t.Errorf("expected 'Hash-Centric (Content)', got %s", name)
	}

	sm.SwitchView()
	name = vc.GetViewName()
	if name != "Filename-Centric (Version)" {
		t.Errorf("expected 'Filename-Centric (Version)', got %s", name)
	}
}

func TestGetViewDescription(t *testing.T) {
	sm := NewSessionManager(1, nil)
	vc := NewViewController(sm)

	desc := vc.GetViewDescription()
	if desc == "" {
		t.Errorf("expected non-empty description")
	}
	if !contains(desc, "hash") {
		t.Errorf("expected hash in description")
	}
}

func TestSwitchViewAndRender(t *testing.T) {
	sm := NewSessionManager(1, nil)
	vc := NewViewController(sm)

	initialView := vc.session.GetCurrentView()
	output := vc.SwitchViewAndRender()

	if output == "" {
		t.Errorf("expected non-empty output")
	}

	newView := vc.session.GetCurrentView()
	if newView == initialView {
		t.Errorf("expected view to change after switch")
	}
}

func TestGetStatusLine(t *testing.T) {
	sm := NewSessionManager(1, nil)
	vc := NewViewController(sm)

	status := vc.GetStatusLine()
	if status == "" {
		t.Errorf("expected non-empty status line")
	}
	if !contains(status, "View:") {
		t.Errorf("expected View in status line")
	}
	if !contains(status, "Flagged:") {
		t.Errorf("expected Flagged in status line")
	}
}

func TestMoveCursorUp(t *testing.T) {
	sm := NewSessionManager(1, nil)

	// Create test duplicates
	sm.duplicates = []DuplicateGroup{
		{Identifier: "hash1", Count: 1},
		{Identifier: "hash2", Count: 2},
		{Identifier: "hash3", Count: 3},
	}

	sm.SetCursor(2)
	if sm.GetCursor() != 2 {
		t.Errorf("expected cursor at 2")
	}

	sm.MoveCursorUp()
	if sm.GetCursor() != 1 {
		t.Errorf("expected cursor at 1 after moving up")
	}

	sm.MoveCursorUp()
	if sm.GetCursor() != 0 {
		t.Errorf("expected cursor at 0")
	}

	sm.MoveCursorUp()
	if sm.GetCursor() != 0 {
		t.Errorf("expected cursor to stay at 0")
	}
}

func TestMoveCursorDown(t *testing.T) {
	sm := NewSessionManager(1, nil)

	// Create test duplicates
	sm.duplicates = []DuplicateGroup{
		{Identifier: "hash1", Count: 1},
		{Identifier: "hash2", Count: 2},
		{Identifier: "hash3", Count: 3},
	}

	sm.SetCursor(0)
	sm.MoveCursorDown()
	if sm.GetCursor() != 1 {
		t.Errorf("expected cursor at 1")
	}

	sm.MoveCursorDown()
	if sm.GetCursor() != 2 {
		t.Errorf("expected cursor at 2")
	}

	sm.MoveCursorDown()
	if sm.GetCursor() != 2 {
		t.Errorf("expected cursor to stay at 2")
	}
}

func TestHandleKeypressTab(t *testing.T) {
	sm := NewSessionManager(1, nil)
	vc := NewViewController(sm)

	initialView := vc.session.GetCurrentView()
	_, err := vc.HandleKeypress('\t')
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	newView := vc.session.GetCurrentView()
	if newView == initialView {
		t.Errorf("expected view to change on TAB")
	}
}

func TestRenderCurrentView(t *testing.T) {
	sm := NewSessionManager(1, nil)
	vc := NewViewController(sm)

	output := vc.RenderCurrentView()
	if output == "" {
		t.Errorf("expected non-empty render output")
	}
}
