package screens

import (
	"context"
	"fmt"

	"github.com/haggishunk/filesprawl/internal/consolidation/browsing"
	"github.com/haggishunk/filesprawl/internal/consolidation/ui/modals"
)

// SummaryScreen displays flagging decisions and allows final confirmation
type SummaryScreen struct {
	sessionManager *browsing.SessionManager
	summaryRenderer *browsing.SummaryRenderer
}

// NewSummaryScreen creates a new summary screen
func NewSummaryScreen(sessionMgr *browsing.SessionManager) *SummaryScreen {
	return &SummaryScreen{
		sessionManager: sessionMgr,
		summaryRenderer: browsing.NewSummaryRenderer(),
	}
}

// Show displays the summary and returns user's confirmation
func (ss *SummaryScreen) Show(ctx context.Context) (bool, error) {
	if ss == nil || ss.sessionManager == nil {
		return false, fmt.Errorf("summary screen not initialized")
	}

	session := ss.sessionManager.GetSession()
	if session == nil {
		return false, fmt.Errorf("no session available")
	}

	// Count flags by type
	keepCount := 0
	archiveCount := 0
	deleteCount := 0

	for _, flag := range session.Flags {
		if flag == nil {
			continue
		}
		switch flag.FlagType {
		case "KEEP":
			keepCount++
		case "ARCHIVE":
			archiveCount++
		case "DELETE":
			deleteCount++
		}
	}

	// Display summary using renderer
	summary := ss.summaryRenderer.RenderDecisionSummary(session)
	fmt.Println(summary)

	// Show quit modal for final confirmation
	quitModal := modals.NewQuitModal(keepCount, archiveCount, deleteCount)
	result := quitModal.Execute()

	return result.Confirmed, result.Error
}

// GetSummary returns formatted summary text
func (ss *SummaryScreen) GetSummary() string {
	if ss == nil || ss.sessionManager == nil {
		return "No summary available"
	}
	session := ss.sessionManager.GetSession()
	if session == nil {
		return "No session available"
	}
	return ss.summaryRenderer.RenderDecisionSummary(session)
}
