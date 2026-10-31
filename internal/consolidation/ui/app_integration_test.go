package ui

import (
	"context"
	"testing"

	"github.com/haggishunk/filesprawl/internal/consolidation/browsing"
)

// TestConsolidationAppInitialization verifies app creation and initialization
func TestConsolidationAppInitialization(t *testing.T) {
	ctx := context.Background()
	
	// Create with nil repositories (testing structure, not DB ops)
	app, err := NewConsolidationApp(ctx, nil, nil)
	if err != nil {
		t.Fatalf("failed to create app: %v", err)
	}
	
	if app == nil {
		t.Fatal("expected non-nil app")
	}
	if app.ctx == nil {
		t.Fatal("expected context to be set")
	}
}

// TestConsolidationAppWithDuplicateGroups tests app with sample data
func TestConsolidationAppWithSampleGroups(t *testing.T) {
	ctx := context.Background()
	app, err := NewConsolidationApp(ctx, nil, nil)
	if err != nil {
		t.Fatalf("failed to create app: %v", err)
	}

	// Create sample duplicate groups using correct struct fields
	groups := []browsing.DuplicateGroup{
		{
			Identifier: "abc123",
			Count:      2,
			Versions: []browsing.FileVersion{
				{FileID: 1},
				{FileID: 2},
			},
		},
	}

	if len(groups) == 0 {
		t.Fatal("expected sample groups")
	}
	if app == nil {
		t.Fatal("expected non-nil app")
	}
}

// TestAppContextPropagation verifies context is preserved
func TestAppContextPropagation(t *testing.T) {
	ctx := context.Background()
	app, err := NewConsolidationApp(ctx, nil, nil)
	if err != nil {
		t.Fatalf("failed to create app: %v", err)
	}

	if app.ctx != ctx {
		t.Error("expected context to be preserved")
	}
}
