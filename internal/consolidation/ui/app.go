package ui

import (
	"context"
	"fmt"

	"github.com/haggishunk/filesprawl/internal/consolidation/browsing"
	"github.com/haggishunk/filesprawl/internal/consolidation/flags"
	"github.com/haggishunk/filesprawl/internal/consolidation/ui/screens"
	"github.com/haggishunk/filesprawl/internal/consolidation/ui/widgets"
	"github.com/haggishunk/filesprawl/internal/consolidation/users"
)

// ConsolidationApp orchestrates the interactive consolidation workflow
// NOTE: Phase 5 TUI integration pending. Stub implementation for now.
type ConsolidationApp struct {
	ctx            context.Context
	userRepo       users.Repository
	flagRepo       flags.Repository
	currentUser    *users.User
	sessionManager *browsing.SessionManager
	flagManager    *browsing.PersistentFlagManager
	hashWidget     *widgets.HashCentricWidget
	filenameWidget *widgets.FilenameCentricWidget
	isDualView     bool // true for hash, false for filename
}

// NewConsolidationApp creates a new consolidation app
func NewConsolidationApp(ctx context.Context, userRepo users.Repository, flagRepo flags.Repository) (*ConsolidationApp, error) {
	return &ConsolidationApp{
		ctx:            ctx,
		userRepo:       userRepo,
		flagRepo:       flagRepo,
		hashWidget:     widgets.NewHashCentricWidget(),
		filenameWidget: widgets.NewFilenameCentricWidget(),
		isDualView:     true,
	}, nil
}

// Run executes the full consolidation workflow
func (ca *ConsolidationApp) Run(duplicateGroups []browsing.DuplicateGroup) error {
	// Step 1: Authentication
	authScreen := screens.NewAuthScreen(ca.userRepo)
	user, err := authScreen.Show(ca.ctx)
	if err != nil {
		return fmt.Errorf("authentication failed: %w", err)
	}
	ca.currentUser = user

	// Step 2: Session resume check
	sessionLoader := screens.NewSessionLoaderScreen(ca.flagRepo, user.ID)
	resumeSession, err := sessionLoader.Show(ca.ctx)
	if err != nil {
		return fmt.Errorf("session check failed: %w", err)
	}

	// Step 3: Initialize browsing session
	ca.sessionManager = browsing.NewSessionManager(user.ID, ca.flagRepo)
	ca.flagManager = browsing.NewPersistentFlagManager(user.ID, ca.flagRepo)

	// Load existing flags if resuming
	if resumeSession {
		if err := ca.sessionManager.LoadExistingFlags(ca.ctx); err != nil {
			return fmt.Errorf("failed to load existing flags: %w", err)
		}
		if err := ca.flagManager.LoadExistingFlags(ca.ctx); err != nil {
			return fmt.Errorf("failed to load flag manager: %w", err)
		}
	}

	// Load duplicates
	ca.sessionManager.LoadDuplicates(duplicateGroups)

	// Update widgets
	ca.hashWidget.UpdateData(duplicateGroups)
	ca.filenameWidget.UpdateData(duplicateGroups)

	// Step 4: Interactive browsing
	if err := ca.runBrowsingUI(); err != nil {
		return fmt.Errorf("browsing failed: %w", err)
	}

	// Step 5: Summary and confirmation
	session := ca.sessionManager.GetSession()
	summaryRenderer := browsing.NewSummaryRenderer()
	summary := summaryRenderer.RenderDecisionSummary(session)
	fmt.Println(summary)

	// Step 6: Save all flags (already done via immediate persistence)
	// But can add checkpoint save here if needed

	return nil
}

// runBrowsingUI runs the interactive browsing interface
func (ca *ConsolidationApp) runBrowsingUI() error {
	// For now, render to string output (Phase 5 placeholder)
	// In Phase 6, this would use gotui full rendering
	return nil
}

// HandleViewSwitch switches between hash and filename views
func (ca *ConsolidationApp) HandleViewSwitch() {
	ca.isDualView = !ca.isDualView
}

// GetCurrentWidget returns the currently active widget
func (ca *ConsolidationApp) GetCurrentWidget() widgets.BrowsingWidget {
	if ca.isDualView {
		return ca.hashWidget
	}
	return ca.filenameWidget
}

// FlagCurrentItem flags the currently selected item
func (ca *ConsolidationApp) FlagCurrentItem(flagType flags.FlagType, reason *string, archiveDest *string) error {
	widget := ca.GetCurrentWidget()
	selectedIdx := widget.GetSelectedIndex()

	switch flagType {
	case flags.FlagTypeKeep:
		return ca.flagManager.PersistKeep(ca.ctx, selectedIdx)
	case flags.FlagTypeArchive:
		if archiveDest == nil || *archiveDest == "" {
			return fmt.Errorf("archive destination required")
		}
		return ca.flagManager.PersistArchive(ca.ctx, selectedIdx, *archiveDest, reason)
	case flags.FlagTypeDelete:
		return ca.flagManager.PersistDelete(ca.ctx, selectedIdx, reason)
	default:
		return fmt.Errorf("unknown flag type: %v", flagType)
	}
}

// GetSessionSummary returns a summary of flagged decisions
func (ca *ConsolidationApp) GetSessionSummary() string {
	session := ca.sessionManager.GetSession()
	summaryRenderer := browsing.NewSummaryRenderer()
	return summaryRenderer.RenderDecisionSummary(session)
}

// Close closes the application
func (ca *ConsolidationApp) Close() error {
	return nil
}
