package ui

import (
	"context"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/haggishunk/filesprawl/internal/consolidation/browsing"
	"github.com/haggishunk/filesprawl/internal/consolidation/flags"
	"github.com/haggishunk/filesprawl/internal/consolidation/ui/screens"
	"github.com/haggishunk/filesprawl/internal/consolidation/users"
)

// ConsolidationApp orchestrates the interactive consolidation workflow
type ConsolidationApp struct {
	ctx            context.Context
	userRepo       users.Repository
	flagRepo       flags.Repository
	currentUser    *users.User
	sessionManager *browsing.SessionManager
	flagManager    *browsing.PersistentFlagManager
	version        string
	reviewMode     bool
	remotes        []string
}

// NewConsolidationApp creates a new consolidation app
func NewConsolidationApp(ctx context.Context, userRepo users.Repository, flagRepo flags.Repository, version string, reviewMode bool, remotes []string) (*ConsolidationApp, error) {
	return &ConsolidationApp{
		ctx:        ctx,
		userRepo:   userRepo,
		flagRepo:   flagRepo,
		version:    version,
		reviewMode: reviewMode,
		remotes:    remotes,
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
	fmt.Printf("🔍 Loading %d duplicate group(s) into session\n", len(duplicateGroups))
	ca.sessionManager.LoadDuplicates(duplicateGroups)

	// Step 4: Interactive browsing with Bubble Tea
	if err := ca.runBrowsingUI(duplicateGroups); err != nil {
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

// runBrowsingUI runs the interactive browsing interface using Bubble Tea
func (ca *ConsolidationApp) runBrowsingUI(duplicateGroups []browsing.DuplicateGroup) error {
	if ca.ctx == nil || ca.sessionManager == nil || ca.flagManager == nil {
		return fmt.Errorf("app not fully initialized for browsing")
	}

	// Run Bubble Tea TUI
	if err := ca.runBubbleTeaTUI(duplicateGroups); err != nil {
		return fmt.Errorf("TUI failed: %w", err)
	}

	// After browsing, show summary and ask for final confirmation
	summaryScreen := screens.NewSummaryScreen(ca.sessionManager)
	confirmed, err := summaryScreen.Show(ca.ctx)
	if err != nil {
		return fmt.Errorf("summary screen failed: %w", err)
	}

	if !confirmed {
		fmt.Println("Session saved. Flags preserved for next session.")
		return nil
	}

	fmt.Println("Consolidation flags confirmed. Ready for execution phase.")
	return nil
}

// runBubbleTeaTUI runs the Bubble Tea TUI
func (ca *ConsolidationApp) runBubbleTeaTUI(duplicateGroups []browsing.DuplicateGroup) error {
	// Create the Bubble Tea model
	model := NewConsolidationModel(
		ca.ctx,
		ca.sessionManager,
		ca.flagManager,
		duplicateGroups,
		ca.version,
		ca.reviewMode,
		ca.remotes,
	)

	// Run the TUI
	program := tea.NewProgram(model, tea.WithAltScreen())
	_, err := program.Run()
	return err
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
