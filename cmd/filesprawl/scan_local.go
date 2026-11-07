package main

import (
	"fmt"
	"os"

	"github.com/haggishunk/filesprawl/internal/cli"
	"github.com/haggishunk/filesprawl/internal/location"
	"github.com/haggishunk/filesprawl/internal/operation"
	"github.com/spf13/cobra"
)

var (
	scanLocalOnlyMissingDirs bool
)

var scanLocalCmd = &cobra.Command{
	Use:   "scan-local",
	Short: "Index files from a local filesystem path",
	Long: `Index all files from the specified local filesystem path, computing hashes and
persisting metadata to the database.

With --only-missing-dirs, only scans directories that haven't been indexed yet,
enabling efficient incremental scanning.`,
	Example: `  # Index entire local directory
  filesprawl scan-local /home/user/documents

  # Index only missing directories (incremental)
  filesprawl scan-local /home/user/documents --only-missing-dirs

  # Index from environment variable
  LOCAL_PATH=/media filesprawl scan-local $LOCAL_PATH`,
	Args: cobra.ExactArgs(1),
	RunE: runScanLocalCmd,
}

func runScanLocalCmd(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()

	if len(args) == 0 {
		return fmt.Errorf("local path is required")
	}

	localPath := args[0]

	repo, cleanup, err := cli.OpenRepository(ctx)
	if err != nil {
		return err
	}
	defer cleanup()

	// If using --only-missing-dirs, fetch the set of already-scanned directories
	var scannedDirs map[string]bool
	if scanLocalOnlyMissingDirs {
		hostname, err := os.Hostname()
		if err != nil {
			return fmt.Errorf("failed to get hostname: %w", err)
		}
		// Create location to get its ID
		loc := location.LocalLocation(hostname)
		if err := repo.EnsureLocation(ctx, &loc); err != nil {
			return fmt.Errorf("failed to ensure location: %w", err)
		}

		scannedDirs, err = repo.GetScannedDirectories(ctx, loc.Id)
		if err != nil {
			return fmt.Errorf("failed to get scanned directories: %w", err)
		}
		fmt.Printf("Found %d previously scanned directories\n", len(scannedDirs))
	}

	ls := operation.NewLocalScanner(
		operation.WithLocalRoot(localPath),
		operation.WithLocalRepo(repo),
		operation.WithOnlyMissingDirs(scanLocalOnlyMissingDirs),
		operation.WithScannedDirs(scannedDirs),
	)

	if err := operation.ScanLocal(ctx, &ls); err != nil {
		return err
	}

	mode := "full"
	if scanLocalOnlyMissingDirs {
		mode = "incremental"
	}
	fmt.Printf("Successfully indexed local path: %s (%s scan)\n", localPath, mode)
	return nil
}

func init() {
	scanLocalCmd.Flags().BoolVar(&scanLocalOnlyMissingDirs, "only-missing-dirs", false,
		"Only scan directories that haven't been indexed yet (incremental mode)")
	rootCmd.AddCommand(scanLocalCmd)
}
