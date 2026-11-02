package main

import (
	"context"
	"fmt"

	"github.com/haggishunk/filesprawl/internal/analysis"
	"github.com/haggishunk/filesprawl/internal/cli"
	"github.com/haggishunk/filesprawl/internal/consolidation/adapter"
	"github.com/haggishunk/filesprawl/internal/consolidation/execution"
	"github.com/haggishunk/filesprawl/internal/consolidation/flags"
	"github.com/haggishunk/filesprawl/internal/consolidation/ui"
	"github.com/haggishunk/filesprawl/internal/consolidation/users"
	"github.com/haggishunk/filesprawl/internal/repository"
	"github.com/spf13/cobra"
)

var (
	consolidateRemote        string
	consolidateReview        bool
	consolidateReviewAll     bool
	consolidateExecute       bool
	consolidateDryRun        bool
	consolidateAcrossRemotes bool
	consolidateUser          string
)

var _ = execution.Executor{} // quiet unused import

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
	consolidateCmd.Flags().StringVar(&consolidateUser, "user", "", "specify user directly (skip authentication)")
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
	if consolidateRemote != "" {
		remoteName := normalizeRemoteName(consolidateRemote)
		// TODO: Implement duplicates within remote
		_ = remoteName
		return fmt.Errorf("consolidate --review for single remote not yet implemented")
	}

	if !consolidateAcrossRemotes {
		return fmt.Errorf("must specify --remote or --across-remotes")
	}

	// Get the underlying ObjectRepository to access the database
	objRepo, ok := repo.(*repository.ObjectRepository)
	if !ok {
		return fmt.Errorf("invalid repository type: expected *ObjectRepository")
	}

	// Get the database interface from the repository
	db := objRepo.GetDatabase()

	// Create user repository using the database interface
	userRepo := users.NewPostgresRepository(db)

	// Create flag repository using the database interface
	flagRepo := flags.NewPostgresRepository(db)

	// Find duplicates across all remotes
	detector := analysis.NewDuplicateDetector(objRepo)
	filterOpts := analysis.FilterOptions{
		HashType: "",
		MinSize:  0,
		Limit:    0,
		Offset:   0,
	}

	duplicateGroups, err := detector.FindAcrossRemotes(ctx, filterOpts)
	if err != nil {
		return fmt.Errorf("failed to find duplicates across remotes: %w", err)
	}

	if len(duplicateGroups) == 0 {
		fmt.Println("No duplicates found across remotes")
		return nil
	}

	// Convert analysis duplicates to browsing duplicates
	browsingGroups := adapter.ToGroupsForBrowsing(duplicateGroups)

	// Create and run consolidation app
	app, err := ui.NewConsolidationApp(ctx, userRepo, flagRepo)
	if err != nil {
		return fmt.Errorf("failed to create consolidation app: %w", err)
	}

	if err := app.Run(browsingGroups); err != nil {
		return fmt.Errorf("consolidation workflow failed: %w", err)
	}

	return nil
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
