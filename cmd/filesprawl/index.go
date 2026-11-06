package main

import (
	"fmt"

	"github.com/haggishunk/filesprawl/internal/cli"
	"github.com/haggishunk/filesprawl/internal/operation"
	"github.com/haggishunk/filesprawl/internal/rclone"
	"github.com/haggishunk/filesprawl/internal/remote"
	"github.com/spf13/cobra"
)

var (
	indexRemote          string
	indexPath            string
	indexOnlyMissingDirs bool
)

var indexCmd = &cobra.Command{
	Use:   "index",
	Short: "Index files from a remote storage location",
	Long: `Index all files from the specified remote, computing hashes and
persisting metadata to the database.

With --only-missing-dirs, only scans directories that haven't been indexed yet,
enabling efficient incremental scanning.`,
	Example: `  # Index entire remote
  filesprawl index --remote media

  # Index specific path within remote
  filesprawl index --remote media --path code/flux

  # Index only missing directories (incremental)
  filesprawl index --remote media --only-missing-dirs`,
	RunE: runIndexCmd,
}

func init() {
	indexCmd.Flags().StringVarP(&indexRemote, "remote", "r", "", "remote name (required)")
	indexCmd.Flags().StringVarP(&indexPath, "path", "p", "", "path within remote")
	indexCmd.Flags().BoolVar(&indexOnlyMissingDirs, "only-missing-dirs", false,
		"Only scan directories that haven't been indexed yet (incremental mode)")
	_ = indexCmd.MarkFlagRequired("remote")

	rootCmd.AddCommand(indexCmd)
}

func runIndexCmd(cmd *cobra.Command, _ []string) error {
	ctx := cmd.Context()

	repo, cleanup, err := cli.OpenRepository(ctx)
	if err != nil {
		return err
	}
	defer cleanup()

	remoteName := normalizeRemoteName(indexRemote)
	if remoteName == "" {
		return fmt.Errorf("remote name is required")
	}

	rem := remote.NewRcloneRemote(remoteName, indexPath)

	// If using --only-missing-dirs, fetch the set of already-scanned directories
	var scannedDirs map[string]bool
	if indexOnlyMissingDirs {
		// Ensure location exists
		if err := repo.EnsureLocation(ctx, &rem.Location); err != nil {
			return fmt.Errorf("failed to ensure location: %w", err)
		}

		scannedDirs, err = repo.GetScannedDirectories(ctx, rem.Location.Id)
		if err != nil {
			return fmt.Errorf("failed to get scanned directories: %w", err)
		}
		fmt.Printf("Found %d previously scanned directories\n", len(scannedDirs))
	}

	scn := operation.NewScanner(
		operation.WithRemote(&rem),
		operation.WithRepo(repo),
		operation.WithScanRoot(indexPath),
		operation.WithScannerOnlyMissingDirs(indexOnlyMissingDirs),
		operation.WithScannerScannedDirs(scannedDirs),
	)

	lo := rclone.NewListOption(rclone.ListOptionFilesOnly())
	lc := rclone.NewListConfig(toRcloneFS(rem.Location.Name), indexPath, &lo)
	if err := operation.Scan(ctx, &scn, lc); err != nil {
		return err
	}

	mode := "full"
	if indexOnlyMissingDirs {
		mode = "incremental"
	}
	_, err = fmt.Fprintf(cmd.OutOrStdout(), "indexed remote %s from location %s (%s scan)\n", rem.Name, rem.Location.Name, mode)
	if err != nil {
		return fmt.Errorf("failed to write index result: %w", err)
	}

	return nil
}
