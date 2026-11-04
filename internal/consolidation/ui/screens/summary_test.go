package screens

import (
	"context"
	"testing"

	"github.com/haggishunk/filesprawl/internal/consolidation/browsing"
)

func TestNewSummaryScreen(t *testing.T) {
	sessionMgr := browsing.NewSessionManager(1, nil)
	screen := NewSummaryScreen(sessionMgr)

	if screen == nil {
		t.Fatal("expected non-nil SummaryScreen")
	}
	if screen.sessionManager == nil {
		t.Fatal("expected sessionManager to be set")
	}
}

func TestSummaryScreenShow_NoSession(t *testing.T) {
	// Create a summary screen with nil session manager
	screen := &SummaryScreen{
		sessionManager: nil,
	}

	_, err := screen.Show(context.Background())
	if err == nil {
		t.Fatal("expected error when session manager is nil")
	}
}

func TestSummaryScreenGetSummary_NoSession(t *testing.T) {
	screen := &SummaryScreen{
		sessionManager: nil,
	}

	summary := screen.GetSummary()
	if summary == "" {
		t.Fatal("expected non-empty error message")
	}
}

func TestSummaryScreenGetSummary_WithSession(t *testing.T) {
	sessionMgr := browsing.NewSessionManager(1, nil)
	screen := NewSummaryScreen(sessionMgr)

	summary := screen.GetSummary()
	// Should return something (even if empty session)
	if summary == "" {
		t.Errorf("expected non-empty summary")
	}
}

func TestSummaryScreenNilCheck(t *testing.T) {
	var screen *SummaryScreen
	_, err := screen.Show(context.Background())
	if err == nil {
		t.Fatal("expected error for nil screen")
	}
}
