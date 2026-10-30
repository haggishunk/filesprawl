package main

import (
	"context"
	"fmt"

	"github.com/haggishunk/filesprawl/internal/analysis"
	"github.com/haggishunk/filesprawl/internal/cli"
	"github.com/haggishunk/filesprawl/internal/consolidation/execution"
	"github.com/haggishunk/filesprawl/internal/consolidation/flags"
	"github.com/haggishunk/filesprawl/internal/consolidation/ui"
	"github.com/haggishunk/filesprawl/internal/consolidation/users"
	"github.com/spf13/cobra"
)

var (
	consolidateRemote     string
	consolidateReview     bool
	consolidateReviewAll  bool
	consolidateExecute    bool
	consolidateDryRun     bool
	consolidateAcrossRemotes bool
)

var consolidateCmd = &cobra.Command{
	Use:   "consolidate",
	Short: "Interactive consolidation of duplicate files",
	Long: `Consolidate duplicate files interactively with flagging, review,
and execution capabilities. Use --review to start the interactive workflow,
--review-all to see all pending flags, or --execute to run flagged operations.`,
	Example: `  # Start interactive consolidation for duplicates
  filesprawl consolidate --remote media --review

  # Consolidate across all remotes
  filesprawl consolidate --across-remotes --review

  # Review all pending flags from previous sessions
  filesprawl consolidate --review-all

  # Execute flagged operations (dry-run)
  filesprawl consolidate --execute --dry-run

  # Execute flagged operations for real
  filesprawl consolidate --execute`,
	RunE: runConsolidateCmd,
}

func init() {
	consolidateCmd.Flags().StringVarP(&consolidateRemote, "remote", "r", "", "consolidate duplicates within this remote")
	consolidateCmd.Flags().BoolVar(&consolidateAcrossRemotes, "across-remotes", false, "consolidate duplicates across all remotes")
	consolidateCmd.Flags().BoolVar(&consolidateReview, "review", false, "start interactive review workflow")
	consolidateCmd.Flags().BoolVar(&consolidateReviewAll, "review-all", false, "review all pending flags from previous sessions")
	consolidateCmd.Flags().BoolVar(&consolidateExecute, "execute", false, "execute flagged consolidation operations")
	consolidateCmd.Flags().BoolVar(&consolidateDryRun, "dry-run", false, "simulate operations without making changes (use with --execute)")
}

func runConsolidateCmd(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()

	// Validate arguments
	if !consolidateReview && !consolidateReviewAll && !consolidateExecute {
		return fmt.Errorf("must specify one of --review, --review-all, or --execute")
	}

	repo, cleanup, err := cli.OpenRepository(ctx)
	if err != nil {
		return err
	}
	defer cleanup()

	// For --review and --review-all, we need duplicates to work with
	if consolidateReview {
		return runConsolidateReview(ctx, repo)
	}

	if consolidateReviewAll {
		return runConsolidateReviewAll(ctx, repo)
	}

	if consolidateExecute {
		return runConsolidateExecute(ctx, repo)
	}

	return nil
}

// runConsolidateReview runs the interactive consolidation workflow
func runConsolidateReview(ctx context.Context, repo interface{}) error {
	// Get duplicates based on filter
	analyzer := analysis.NewAnalyzer(repo)

	if consolidateRemote != "" {
		remoteName := normalizeRemoteName(consolidateRemote)
		// TODO: Implement duplicates within remote
		_ = remoteName
		return fmt.Errorf("consolidate --review for single remote not yet implemented")
	}

	if consolidateAcrossRemotes {
		// Get duplicates across all remotes
		groups, err := analyzer.FindDuplicatesAcross(ctx)
		if err != nil {
			return fmt.Errorf("failed to find duplicates: %w", err)
		}

		if len(groups) == 0 {
			fmt.Println("No duplicates found")
			return nil
		}

		fmt.Println("consolidate --review: Feature not yet implemented in TUI")
		fmt.Printf("Found %d duplicate groups\n", len(groups))
		return nil
	}

	return fmt.Errorf("must specify --remote or --across-remotes")
}

// runConsolidateReviewAll shows all pending flags from previous sessions
func runConsolidateReviewAll(ctx context.Context, repo interface{}) error {
	// TODO: Implement review-all interface
	fmt.Println("consolidate --review-all: Not yet implemented")
	return nil
}

// runConsolidateExecute runs the executor on all flagged files
func runConsolidateExecute(ctx context.Context, repo interface{}) error {
	fmt.Println("consolidate --execute: Feature not yet fully integrated")
	fmt.Println("This requires:")
	fmt.Println("  - User authentication/selection")
	fmt.Println("  - Loading flagged files from database")
	fmt.Println("  - Executing archive/delete operations via rclone")
	fmt.Println("  - Recording audit trail in consolidation_reviews table")
	return nil
}
